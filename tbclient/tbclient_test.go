// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

type DockerContainer struct {
	client      *client.Client
	ctx         *context.Context
	containerId string
}

type ContractTestSuite struct {
	suite.Suite
	docker   DockerContainer
	tbClient Client
}

// Test Data
var lonTopology = Topology{
	Uid:         "lontopology",
	Name:        "lontopology name",
	Description: "lontopology description",
	Datacenter:  "LON",
	Notes:       "Notes",
	Status:      "DRAFT",
}

var linuxOsFamily = OsFamily{
	Id:   "LINUX",
	Name: "Linux",
}

var e1000NicType = NicType{
	Id:   "VIRTUAL_E1000",
	Name: "VirtualE1000",
}

var script1InventoryHwScript = InventoryHwScript{
	Uid:  "script1",
	Name: "EN-3850-switch-start.xml",
}

var template1InventoryHwScript = InventoryHwScript{
	Uid:  "template1",
	Name: "en-sw-3850-16-v1.txt",
}

var primaryInventoryNetwork = InventoryNetwork{
	Id:     "VLAN-PRIMARY",
	Type:   "ROUTED",
	Subnet: "198.18.128.0 /18",
}

var routedInventoryNetwork = InventoryNetwork{
	Id:     "L3-VLAN-1",
	Type:   "ROUTED",
	Subnet: "198.18.1.0 /24",
}

var routedInventoryNetwork2 = InventoryNetwork{
	Id:     "L3-VLAN-2",
	Type:   "ROUTED",
	Subnet: "198.18.2.0 /24",
}

var routedNetwork = Network{
	Uid:              "routednetwork",
	Name:             "Routed Network",
	Description:      "Routed Network description",
	InventoryNetwork: &InventoryNetwork{Id: routedInventoryNetwork.Id},
	Topology:         &Topology{Uid: lonTopology.Uid},
}

var lonDefaultNetwork = Network{
	Uid:              "lontopologyv1nwk",
	Name:             "Default Network",
	Description:      "Default Network description",
	InventoryNetwork: &primaryInventoryNetwork,
	Topology:         &Topology{Uid: lonTopology.Uid},
}

var inventoryVm = InventoryVm{
	Id:                  "templatevm3",
	Datacenter:          "LON",
	OriginalName:        "na-edge1",
	OriginalDescription: "na-edge1",
	CpuQty:              2,
	MemoryMb:            4096,
	NetworkInterfaces: []InventoryVmNic{
		// TODO - contract needs to specify additional fields
		{
			InventoryNetworkId: "VLAN-PRIMARY",
			Name:               "Network adapter 1",
			IpAddress:          "198.18.133.115",
		},
		{
			InventoryNetworkId: "VLAN-PRIMARY",
			Name:               "Network adapter 2",
		},
		{
			InventoryNetworkId: "L2-VLAN-15",
			Name:               "Network adapter 3",
		},
		{
			InventoryNetworkId: "L2-VLAN-20",
			Name:               "Network adapter 4",
		},
	},
	RemoteAccess: &InventoryVmRemoteAccess{
		RdpAutoLogin: true,
		RdpEnabled:   true,
		SshEnabled:   true,
	},
}

var nestedHypervisor = false

var vm = Vm{
	Uid:              "lonvm1",
	Name:             "Mail Server 1",
	Description:      "Mail Server 1",
	MemoryMb:         6144,
	CpuQty:           3,
	NestedHypervisor: &nestedHypervisor,
	InventoryVmId:    "templatevm2",
	OsFamily:         "LINUX",
	RemoteAccess: &VmRemoteAccess{
		Username:         "user_name",
		Password:         "password",
		VmConsoleEnabled: true,
		DisplayCredentials: &VmRemoteAccessDisplayCredentials{
			Username: "display_user_name",
			Password: "display_password",
		},
		InternalUrls: []VmRemoteAccessInternalUrl{
			{
				Location:    "http://198.18.131.200:8080/myApp/start?foo=bar",
				Description: "My App Launch Page",
			},
		},
	},
	VmNetworkInterfaces: []VmNic{
		{
			Uid:        "lonvm1msnic",
			IpAddress:  "198.18.131.200",
			Name:       "Network adapter 1",
			MacAddress: "00:50:56:00:00:07",
			Type:       "VIRTUAL_E1000",
			InUse:      true,
			Rdp: &VmNicRdp{
				Enabled:   true,
				AutoLogin: true,
			},
			Ssh: &VmNicSsh{Enabled: true},
			Network: &Network{
				// Contract is missing 'description' field
				Uid:              lonDefaultNetwork.Uid,
				Name:             lonDefaultNetwork.Name,
				InventoryNetwork: lonDefaultNetwork.InventoryNetwork,
			},
		},
		{
			Uid:        "lonvm1natnic",
			IpAddress:  "198.18.131.201",
			Name:       "Network adapter 2",
			MacAddress: "00:50:56:00:00:08",
			Type:       "VIRTUAL_E1000",
			Network: &Network{
				// Contract is missing 'description' field
				Uid:              lonDefaultNetwork.Uid,
				Name:             lonDefaultNetwork.Name,
				InventoryNetwork: lonDefaultNetwork.InventoryNetwork,
			},
		},
	},
	AdvancedSettings: &VmAdvancedSettings{
		NameInHypervisor:      "Demo_007-VM",
		BiosUuid:              "61 62 63 64 65 66 67 68-69 6a 6b 6c 6d 6e 6f 70",
		NotStarted:            false,
		AllDisksNonPersistent: false,
	},
	GuestAutomation: &VmGuestAutomation{
		Command:   "cd /var/; sh script.sh",
		DelaySecs: 120,
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var createVm = Vm{
	Uid:              "newvm",
	Name:             "vmtemplate1 name",
	Description:      "vmtemplate1 description",
	MemoryMb:         8192,
	CpuQty:           4,
	NestedHypervisor: &nestedHypervisor,
	InventoryVmId:    "templatevm1",
	RemoteAccess: &VmRemoteAccess{
		Username:         "user_name",
		Password:         "password",
		VmConsoleEnabled: false,
	},
	VmNetworkInterfaces: []VmNic{
		{
			Uid:       "newvmnic1",
			IpAddress: "198.18.131.200",
			Name:      "Network adapter 0",
			Network: &Network{
				// Contract is missing 'description' field
				Uid: lonDefaultNetwork.Uid,
			},
		},
		{
			Uid:       "newvmnic2",
			IpAddress: "198.18.131.201",
			Name:      "Network adapter 1",
			Network: &Network{
				// Contract is missing 'description' field
				Uid: routedNetwork.Uid,
			},
		},
	},
	AdvancedSettings: &VmAdvancedSettings{
		NameInHypervisor:      "Win_10-VM-1",
		BiosUuid:              "61 62 63 64 65 66 67 68-69 6a 6b 6c 6d 6e 6f 70",
		NotStarted:            false,
		AllDisksNonPersistent: false,
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var inventoryHw = InventoryHw{
	Id:                       "86",
	Name:                     "ISA3000-FTD",
	Description:              "ISA3000-FTD",
	PowerControlAvailable:    true,
	HardwareConsoleAvailable: true,
	NetworkInterfaces: []InventoryHwNic{
		{
			Id: "GigabitEthernet1/1",
		},
		{
			Id: "GigabitEthernet1/2",
		},
	},
}

var powerControlEnabled = true
var hardwareConsoleEnabled = true

var hw = Hw{
	Uid:                    "lonhardwareitem2",
	Name:                   "ISA3000-FTD",
	PowerControlEnabled:    &powerControlEnabled,
	HardwareConsoleEnabled: &hardwareConsoleEnabled,
	StartupScript: &InventoryHwScript{
		Uid:  "script1",
		Name: "EN-3850-switch-start.xml",
	},
	CustomScript: &InventoryHwScript{
		Uid:  "script3",
		Name: "custom.xml",
	},
	ShutdownScript: &InventoryHwScript{
		Uid:  "script2",
		Name: "EN-3850-switch-stop.xml",
	},
	TemplateConfigScript: &InventoryHwScript{
		Uid:  "template1",
		Name: "en-sw-3850-16-v1.txt",
	},
	NetworkInterfaces: []HwNic{
		{
			Uid: "hwnetworkiface1",
			Network: Network{
				Uid:  lonDefaultNetwork.Uid,
				Name: lonDefaultNetwork.Name,
			},
			NetworkInterface: InventoryHwNic{
				Id: "GigabitEthernet1/1",
			},
		},
	},
	InventoryHardwareItem: &inventoryHw,
	Topology:              &Topology{Uid: lonTopology.Uid},
}

var inventoryLicense = InventoryLicense{
	Id:       "24",
	Name:     "APIC License",
	Quantity: 3,
}

var license = License{
	Uid:              "license1",
	Quantity:         1,
	InventoryLicense: &inventoryLicense,
	Topology:         &Topology{Uid: lonTopology.Uid},
}

var vmStartOrder = VmStartOrder{
	Uid:     "lontopologyv1sso",
	Ordered: false,
	Positions: []VmStartPosition{
		{
			Position:     1,
			DelaySeconds: 0,
			Vm:           &Vm{Uid: vm.Uid, Name: vm.Name},
		},
		{
			Position:     2,
			DelaySeconds: 10,
			Vm:           &Vm{Uid: "lonvm2", Name: "Collab DB"},
		},
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var vmStopOrder = VmStopOrder{
	Uid:     "lontopologyv1sso",
	Ordered: false,
	Positions: []VmStopPosition{
		{
			Position: 1,
			Vm:       &Vm{Uid: "lonvm2", Name: "Collab DB"},
		},
		{
			Position: 2,
			Vm:       &Vm{Uid: vm.Uid, Name: vm.Name},
		},
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var hwStartOrder = HwStartOrder{
	Uid:     "lontopologyv1sso",
	Ordered: false,
	Positions: []HwStartPosition{
		{
			Position:     1,
			DelaySeconds: 0,
			Hw:           &Hw{Uid: "lonhardwareitem1", Name: "VCube PSTN Services"},
		},
		{
			Position:     2,
			DelaySeconds: 30,
			Hw:           &Hw{Uid: "lonhardwareitem2", Name: "ISA3000-FTD"},
		},
		{
			Position:     3,
			DelaySeconds: 60,
			Hw:           &Hw{Uid: "lonhardwareitem3", Name: "ISA3000-FTD No.2"},
		},
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var remoteAccess = RemoteAccess{
	Uid:                "lontopologyv1rm",
	AnyconnectEnabled:  true,
	EndpointKitEnabled: true,
	Topology:           &Topology{Uid: lonTopology.Uid},
}

var scenarioOption1 = ScenarioOption{
	Uid:          "scenariooption1",
	InternalName: "Version_A",
	DisplayName:  "Version A",
}

var scenarioOption2 = ScenarioOption{
	Uid:          "scenariooption2",
	InternalName: "Version_B",
	DisplayName:  "Version B",
}

var scenario = Scenario{
	Uid:      "lonscenario",
	Question: "Which version of the demo do you want to launch?",
	Enabled:  true,
	Options:  []ScenarioOption{scenarioOption1, scenarioOption2},
	Topology: &Topology{Uid: lonTopology.Uid},
}

func (suite *ContractTestSuite) SetupSuite() {
	suite.docker = startWiremock(suite)
	suite.tbClient = createTbClient(suite)
}

func (suite *ContractTestSuite) TearDownSuite() {
	suite.docker.client.ContainerRemove(*suite.docker.ctx, suite.docker.containerId, types.ContainerRemoveOptions{Force: true})
}

// Error Tests

func (suite *ContractTestSuite) TestErrorHandling() {

	// Given
	badNetwork := routedNetwork
	badNetwork.InventoryNetwork = &routedInventoryNetwork2
	badNetwork.Name = ""

	// When
	_, err := suite.tbClient.CreateNetwork(badNetwork)

	// Then
	suite.ErrorContains(err, "Field 'name' size must be between 1 and 255")
	suite.ErrorContains(err, "Field 'name' must not be blank")
	suite.ErrorContains(err, "[Server Response: 400 Bad Request]")
}

// Asset Tests

func (suite *ContractTestSuite) TestGetAllTopologies() {

	// Given
	expectedTopology := lonTopology

	// When
	topologies, err := suite.tbClient.GetAllTopologies()
	suite.handleError(err)

	// Then
	suite.Equal(9, len(topologies))
	suite.Contains(topologies, expectedTopology)
}

func (suite *ContractTestSuite) TestGetTopology() {

	// When
	actualTopology, err := suite.tbClient.GetTopology(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(lonTopology, *actualTopology)
}

func (suite *ContractTestSuite) TestUpdateTopology() {
	// Given
	expectedTopology := lonTopology
	expectedTopology.Name = "lontopology name updated"

	// When
	actualTopology, err := suite.tbClient.UpdateTopology(expectedTopology)
	suite.handleError(err)

	// Then
	suite.Equal(expectedTopology, *actualTopology)
}

func (suite *ContractTestSuite) TestCreateTopology() {
	// Given
	newTopology := lonTopology
	newTopology.Uid = ""

	// When
	actualTopology, err := suite.tbClient.CreateTopology(newTopology)
	suite.handleError(err)

	// Then
	newTopology.Uid = "newtopology"
	suite.Equal(newTopology, *actualTopology)
}

func (suite *ContractTestSuite) TestDeleteTopology() {

	// When
	err := suite.tbClient.DeleteTopology(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

func (suite *ContractTestSuite) TestGetAllNetworks() {

	// Given
	expectedNetwork := routedNetwork
	expectedNetwork.InventoryNetwork = &routedInventoryNetwork

	// When
	networks, err := suite.tbClient.GetAllNetworks(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(4, len(networks))
	suite.Contains(networks, expectedNetwork)
}

func (suite *ContractTestSuite) TestGetNetwork() {

	// Given
	expectedNetwork := routedNetwork
	expectedNetwork.InventoryNetwork = &routedInventoryNetwork

	// When
	actualNetwork, err := suite.tbClient.GetNetwork(routedNetwork.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(expectedNetwork, *actualNetwork)
}

func (suite *ContractTestSuite) TestUpdateNetwork() {
	// Given
	expectedInventoryNetwork := routedInventoryNetwork2
	expectedInventoryNetwork.Subnet = "198.18.3.0 /24"

	expectedNetwork := routedNetwork
	expectedNetwork.Name = "Updated Routed Network"
	expectedNetwork.InventoryNetwork = &InventoryNetwork{Id: expectedInventoryNetwork.Id}

	// When
	actualNetwork, err := suite.tbClient.UpdateNetwork(expectedNetwork)
	suite.handleError(err)

	// Then
	expectedNetwork.InventoryNetwork = &expectedInventoryNetwork
	suite.Equal(expectedNetwork, *actualNetwork)
}

func (suite *ContractTestSuite) TestCreateNetwork() {
	// Given
	newNetwork := routedNetwork
	newNetwork.Uid = ""
	newNetwork.InventoryNetwork = &InventoryNetwork{Id: routedInventoryNetwork2.Id}

	// When
	actualNetwork, err := suite.tbClient.CreateNetwork(newNetwork)
	suite.handleError(err)

	// Then
	newNetwork.Uid = "newroutednetwork"
	newNetwork.InventoryNetwork = &routedInventoryNetwork2
	suite.Equal(newNetwork, *actualNetwork)
}

func (suite *ContractTestSuite) TestDeleteNetwork() {

	// When
	err := suite.tbClient.DeleteNetwork(routedNetwork.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

func (suite *ContractTestSuite) TestGetAllVms() {

	// When
	vms, err := suite.tbClient.GetAllVms(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(3, len(vms))
	suite.Contains(vms, vm)
}

func (suite *ContractTestSuite) TestGetVm() {

	// Given
	expectedVm := vm
	// Work around contract data inconsistencies
	nic := vm.VmNetworkInterfaces[1]
	nic.InUse = true
	nic.Network.InventoryNetwork = nil
	expectedVm.VmNetworkInterfaces[1] = nic

	// When
	actualVm, err := suite.tbClient.GetVm(vm.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(expectedVm, *actualVm)
}

func (suite *ContractTestSuite) TestUpdateVm() {

	// Given
	expectedVm := vm

	// Change Value
	expectedVm.Name = "New Name"

	// When
	actualVm, err := suite.tbClient.UpdateVm(expectedVm)
	suite.handleError(err)

	// Then

	// Work around contract data inconsistencies
	nic0 := vm.VmNetworkInterfaces[0]
	nic0.Uid = ""
	nic0.Name = ""
	nic0.Network.InventoryNetwork = nil
	expectedVm.VmNetworkInterfaces[0] = nic0

	nic1 := vm.VmNetworkInterfaces[1]
	nic1.Uid = ""
	nic1.Name = ""
	nic1.Network.InventoryNetwork = nil
	nic1.InUse = true
	expectedVm.VmNetworkInterfaces[1] = nic1
	suite.Equal(expectedVm, *actualVm)
}

func (suite *ContractTestSuite) TestCreateVm() {
	// Given
	expectedVm := createVm
	expectedVm.Uid = ""

	// When
	actualVm, err := suite.tbClient.CreateVm(expectedVm)
	suite.handleError(err)

	// Then
	expectedVm.Uid = "newvm"

	// TODO - contract gives unpredictable network names, need to blank them for assert
	for _, nic := range actualVm.VmNetworkInterfaces {
		nic.Network.Name = ""
	}
	suite.Equal(expectedVm, *actualVm)
}

func (suite *ContractTestSuite) TestDeleteVm() {

	// When
	err := suite.tbClient.DeleteVm(vm.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

func (suite *ContractTestSuite) TestGetAllHws() {

	// When
	hws, err := suite.tbClient.GetAllHws(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(4, len(hws))
	suite.Contains(hws, hw)
}

func (suite *ContractTestSuite) TestGetHw() {

	// Given
	expectedHw := hw

	// When
	actualHw, err := suite.tbClient.GetHw(hw.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(expectedHw, *actualHw)
}

func (suite *ContractTestSuite) TestUpdateHw() {

	// Given
	expectedHw := hw
	expectedHw.Name = "New Name"

	newNics := make([]HwNic, 2)
	newNics[0] = expectedHw.NetworkInterfaces[0]
	newNics[1] = HwNic{
		Uid: "newhwnetworkiface",
		Network: Network{
			Uid:  routedNetwork.Uid,
			Name: routedNetwork.Name,
		},
		NetworkInterface: InventoryHwNic{
			Id: "GigabitEthernet1/2",
		},
	}

	expectedHw.NetworkInterfaces = newNics

	// When
	actualHw, err := suite.tbClient.UpdateHw(expectedHw)
	suite.handleError(err)

	// Then
	suite.Equal(expectedHw, *actualHw)
}

func (suite *ContractTestSuite) TestCreateHw() {
	// Given
	// Modify to match contracts
	expectedHw := hw
	expectedHw.Uid = ""
	expectedHw.TemplateConfigScript = nil
	expectedHw.StartupScript = nil
	expectedHw.CustomScript = nil
	expectedHw.ShutdownScript = nil
	expectedHw.NetworkInterfaces = make([]HwNic, 0)

	// When
	actualHw, err := suite.tbClient.CreateHw(expectedHw)
	suite.handleError(err)

	// Then
	expectedHw.Uid = "newhardwareitem"
	expectedHw.Name = "ISA3000-FTD(1)"
	suite.Equal(expectedHw, *actualHw)
}

func (suite *ContractTestSuite) TestDeleteHw() {

	// When
	err := suite.tbClient.DeleteHw(hw.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

// License Tests

func (suite *ContractTestSuite) TestGetAllLicenses() {

	// When
	licenses, err := suite.tbClient.GetAllLicenses(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(2, len(licenses))
	suite.Contains(licenses, license)
}

func (suite *ContractTestSuite) TestGetLicense() {

	// Given
	expectedLicense := license

	// When
	actualLicense, err := suite.tbClient.GetLicense(license.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(expectedLicense, *actualLicense)
}

func (suite *ContractTestSuite) TestUpdateLicense() {

	// Given
	expectedLicense := license
	expectedLicense.Quantity = license.Quantity + 1

	// When
	actualLicense, err := suite.tbClient.UpdateLicense(expectedLicense)
	suite.handleError(err)

	// Then
	suite.Equal(expectedLicense, *actualLicense)
}

func (suite *ContractTestSuite) TestCreateLicense() {
	// Given
	expectedLicense := license
	// Modify to match contracts
	expectedLicense.InventoryLicense = &InventoryLicense{Id: "24"}

	// When
	actualLicense, err := suite.tbClient.CreateLicense(expectedLicense)
	suite.handleError(err)

	// Then
	expectedLicense.Uid = "newlicense"
	expectedLicense.InventoryLicense.Name = "mimic"
	expectedLicense.InventoryLicense.Quantity = 6
	suite.Equal(expectedLicense, *actualLicense)
}

func (suite *ContractTestSuite) TestDeleteLicense() {

	// When
	err := suite.tbClient.DeleteLicense(license.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

// VM Start Order

func (suite *ContractTestSuite) TestGetVmStartOrder() {

	// When
	actualVmStartOrder, err := suite.tbClient.GetVmStartOrder(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(vmStartOrder, *actualVmStartOrder)
}

func (suite *ContractTestSuite) TestUpdateVmStartOrder() {

	// Given
	expectedVmStartOrder := vmStartOrder
	expectedVmStartOrder.Ordered = true
	// Swap the Positions (Slice ordering matters to suite.Equal)
	pos1 := expectedVmStartOrder.Positions[0]
	pos2 := expectedVmStartOrder.Positions[1]
	pos2.Position = 1
	pos1.Position = 2
	expectedVmStartOrder.Positions[0] = pos2
	expectedVmStartOrder.Positions[1] = pos1

	// When
	actualVmStartOrder, err := suite.tbClient.UpdateVmStartOrder(expectedVmStartOrder)
	suite.handleError(err)

	// Then
	suite.Equal(expectedVmStartOrder, *actualVmStartOrder)
}

// VM Stop Order

func (suite *ContractTestSuite) TestGetVmStopOrder() {

	// When
	actualVmStopOrder, err := suite.tbClient.GetVmStopOrder(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(vmStopOrder, *actualVmStopOrder)
}

func (suite *ContractTestSuite) TestUpdateVmStopOrder() {

	// Given
	expectedVmStopOrder := vmStopOrder
	expectedVmStopOrder.Ordered = true
	// Swap the Positions (Slice ordering matters to suite.Equal)
	pos1 := expectedVmStopOrder.Positions[0]
	pos2 := expectedVmStopOrder.Positions[1]
	pos2.Position = 1
	pos1.Position = 2
	expectedVmStopOrder.Positions[0] = pos2
	expectedVmStopOrder.Positions[1] = pos1

	// When
	actualVmStopOrder, err := suite.tbClient.UpdateVmStopOrder(expectedVmStopOrder)
	suite.handleError(err)

	// Then
	suite.Equal(expectedVmStopOrder, *actualVmStopOrder)
}

// Hw Start Order

func (suite *ContractTestSuite) TestGetHwStartOrder() {

	// When
	actualHwStartOrder, err := suite.tbClient.GetHwStartOrder(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(hwStartOrder, *actualHwStartOrder)
}

func (suite *ContractTestSuite) TestUpdateHwStartOrder() {

	// Given
	expectedHwStartOrder := hwStartOrder
	expectedHwStartOrder.Ordered = true

	// When
	actualHwStartOrder, err := suite.tbClient.UpdateHwStartOrder(expectedHwStartOrder)
	suite.handleError(err)

	// Then
	suite.Equal(expectedHwStartOrder, *actualHwStartOrder)
}

// Remote Access

func (suite *ContractTestSuite) TestGetRemoteAccess() {

	// When
	actualRemoteAccess, err := suite.tbClient.GetRemoteAccess(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(remoteAccess, *actualRemoteAccess)
}

func (suite *ContractTestSuite) TestUpdateRemoteAccess() {

	// Given
	expectedRemoteAccess := remoteAccess
	expectedRemoteAccess.EndpointKitEnabled = false

	// When
	actualRemoteAccess, err := suite.tbClient.UpdateRemoteAccess(expectedRemoteAccess)
	suite.handleError(err)

	// Then
	suite.Equal(expectedRemoteAccess, *actualRemoteAccess)
}

// Scenario

func (suite *ContractTestSuite) TestGetScenario() {

	// When
	actualScenario, err := suite.tbClient.GetScenario(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(scenario, *actualScenario)
}

func (suite *ContractTestSuite) TestUpdateScenario() {

	// Given
	expectedScenario := scenario
	expectedScenario.Question = "What would you like today?"

	// When
	actualScenario, err := suite.tbClient.UpdateScenario(expectedScenario)
	suite.handleError(err)

	// Then
	expectedScenario.Uid = "" // Missing in contract stub
	suite.Equal(expectedScenario, *actualScenario)
}

// Inventory Tests

func (suite *ContractTestSuite) TestGetAllOsFamilies() {

	// When
	osFamilies, err := suite.tbClient.GetAllOsFamilies()
	suite.handleError(err)

	// Then
	suite.Equal(4, len(osFamilies))
	suite.Contains(osFamilies, linuxOsFamily)
}

func (suite *ContractTestSuite) TestGetAllNicTypes() {

	// When
	nicTypes, err := suite.tbClient.GetAllNicTypes()
	suite.handleError(err)

	// Then
	suite.Equal(5, len(nicTypes))
	suite.Contains(nicTypes, e1000NicType)
}

func (suite *ContractTestSuite) TestGetAllInventoryHwScripts() {

	// When
	hwScripts, err := suite.tbClient.GetAllInventoryHwScripts(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(4, len(hwScripts))
	suite.Contains(hwScripts, script1InventoryHwScript)
}

func (suite *ContractTestSuite) TestGetAllInventoryHwTemplateConfigs() {

	// When
	configTemplates, err := suite.tbClient.GetAllInventoryHwTemplateConfigs(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(2, len(configTemplates))
	suite.Contains(configTemplates, template1InventoryHwScript)
}

func (suite *ContractTestSuite) TestGetAllInventoryNetworks() {

	// When
	inventoryNetworks, err := suite.tbClient.GetAllInventoryNetworks(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(25, len(inventoryNetworks))
	suite.Contains(inventoryNetworks, routedInventoryNetwork)
}

func (suite *ContractTestSuite) TestGetAllInventoryVms() {

	// TODO - Contract needs to be updated to remove 'filter' as mandatory
	suite.T().Skip("Skipping until contracts are updated")

	// When
	inventoryVms, err := suite.tbClient.GetAllInventoryVms(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(2, len(inventoryVms))
	suite.Contains(inventoryVms, inventoryVm)
}

func (suite *ContractTestSuite) TestGetAllInventoryHws() {

	// When
	inventoryHws, err := suite.tbClient.GetAllInventoryHws(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(5, len(inventoryHws))
	suite.Contains(inventoryHws, inventoryHw)
}

func (suite *ContractTestSuite) TestGetAllInventoryLicenses() {

	// When
	inventoryLicenses, err := suite.tbClient.GetAllInventoryLicenses(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(3, len(inventoryLicenses))
	suite.Contains(inventoryLicenses, inventoryLicense)
}

func TestContractTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}

func startWiremock(suite *ContractTestSuite) DockerContainer {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	suite.handleError(err)
	defer cli.Close()

	imageName := "wiremock/wiremock:2.35.0"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	suite.handleError(err)
	defer out.Close()

	io.Copy(os.Stdout, out)

	_, b, _, _ := runtime.Caller(0)
	mappingDir := filepath.Join(filepath.Dir(b), "test_stubs", "mappings")

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "9876",
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mappingDir,
				Target: "/home/wiremock/mappings",
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
		Cmd: []string{"--global-response-templating"},
	}, hostConfig, nil, nil, "")
	suite.handleError(err)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		suite.FailNow(err.Error())
	}

	return DockerContainer{
		client:      cli,
		ctx:         &ctx,
		containerId: resp.ID,
	}
}

func createTbClient(suite *ContractTestSuite) Client {
	c := Client{
		HostURL:     "http://localhost:9876",
		Token:       "oauthAuthorized",
		DisableGzip: true,
		Debug:       false,
	}

	timeout := time.Now().Add(time.Minute)
	for _, err := c.GetAllTopologies(); err != nil; _, err = c.GetAllTopologies() {
		if time.Now().After(timeout) {
			suite.T().Log("Timeout starting wiremock")
			suite.handleError(err)
		}
		fmt.Println("Waiting another 5s for wiremock...")
		time.Sleep(5 * time.Second)
	}

	return c
}

func (suite *ContractTestSuite) handleError(err error) {
	if err != nil {
		suite.FailNow(err.Error())
	}
}
