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

var westmereEvcMode = EvcMode{
	Id:           "WESTMERE",
	Name:         "Westmere",
	DisplayOrder: 1,
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
	Id:                  "templatevm1",
	Datacenter:          "LON",
	Name:                "vmtemplate1 name",
	Description:         "vmtemplate1 description",
	OriginalName:        "Collab-mssql1",
	OriginalDescription: "Collab-mssql1",
	CpuQty:              4,
	MemoryMb:            8192,
	EvcMode:             evcModeWestmere,
	NetworkInterfaces: []InventoryVmNic{
		// TODO - contract needs to specify additional fields
		{
			InventoryNetworkId: "VLAN-PRIMARY",
			Name:               "Network adapter 0",
			IpAddress:          "198.18.133.115",
			Type:               "VIRTUAL_E1000",
		},
	},
	RemoteAccess: &InventoryVmRemoteAccess{
		RdpAutoLogin: false,
		RdpEnabled:   false,
		SshEnabled:   false,
	},
}

var nestedHypervisor = false
var evcModeHaswell = "HASWELL"
var evcModeWestmere = "WESTMERE"

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
			AssignDhcp: true,
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
			Uid:        "lonvm1natnic1",
			IpAddress:  "198.18.131.211",
			Name:       "Network adapter 11",
			MacAddress: "00:50:56:00:00:11",
			Type:       "VIRTUAL_E1000",
			InUse:      true,
			AssignDhcp: false,
			Network: &Network{
				// Contract is missing 'inventoryNetwork' field
				Uid:  lonDefaultNetwork.Uid,
				Name: lonDefaultNetwork.Name,
			},
		},
		{
			Uid:        "lonvm1natnic2",
			IpAddress:  "198.18.131.212",
			Name:       "Network adapter 12",
			MacAddress: "00:50:56:00:00:12",
			Type:       "VIRTUAL_E1000",
			InUse:      true,
			AssignDhcp: false,
			Network: &Network{
				// Contract is missing 'description' field
				Uid:  lonDefaultNetwork.Uid,
				Name: lonDefaultNetwork.Name,
			},
		},
		{
			Uid:        "lonvm1natnic3",
			IpAddress:  "198.18.131.213",
			Name:       "Network adapter 13",
			MacAddress: "00:50:56:00:00:13",
			Type:       "VIRTUAL_E1000",
			InUse:      true,
			AssignDhcp: false,
			Network: &Network{
				// Contract is missing 'description' field
				Uid:  lonDefaultNetwork.Uid,
				Name: lonDefaultNetwork.Name,
			},
		},
	},
	AdvancedSettings: &VmAdvancedSettings{
		NameInHypervisor:      "Demo_007-VM",
		BiosUuid:              "61 62 63 64 65 66 67 68-69 6a 6b 6c 6d 6e 6f 70",
		NotStarted:            false,
		AllDisksNonPersistent: false,
		EvcMode:               evcModeHaswell,
	},
	DhcpConfig: &VmDhcpConfig{
		DefaultGatewayIp: "198.18.130.2",
	},
	GuestAutomation: &VmGuestAutomation{
		Command:   "cd /var/; sh script.sh",
		DelaySecs: 120,
	},
	ShutdownAutomation: &VmGuestAutomation{
		Command:   "cd /var/; sh shutdownscript.sh",
		DelaySecs: 200,
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
			Uid:        "newvmnic1",
			IpAddress:  "198.18.131.200",
			MacAddress: "00:50:56:10:11:12",
			Name:       "Network adapter 0",
			Network: &Network{
				// Contract is missing 'description' field
				Uid: lonDefaultNetwork.Uid,
			},
		},
		{
			Uid:        "newvmnic2",
			IpAddress:  "198.18.131.201",
			MacAddress: "00:50:56:10:11:13",
			Name:       "Network adapter 1",
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
		EvcMode:               evcModeWestmere,
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

var publicScope = "PUBLIC"

var ipNatRule = IpNatRule{
	Uid:      "lonipnatrule1",
	EastWest: false,
	Scope:    &publicScope,
	Target: IpNatTarget{
		IpAddress: "198.18.131.100",
		Name:      "Some Device",
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var vmNatRule = VmNatRule{
	Uid:      "lonvmpublicnat",
	EastWest: false,
	Scope:    &publicScope,
	Target: VmNatTarget{
		VmNic:     &VmNic{Uid: "lonvm1natnic1"},
		IpAddress: "198.18.131.211",
		Name:      "Mail Server 1",
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var inboundProxyRule = InboundProxyRule{
	Uid:     "loninboundproxy1",
	TcpPort: 433,
	Ssl:     true,
	UrlPath: "/demo/demo101",
	VmNicTarget: &TrafficVmNicTarget{
		Uid:       "lonvm2pxnic",
		IpAddress: "198.18.133.110",
		Vm: &Vm{
			Uid:  "lonvm2",
			Name: "Collab DB",
		},
	},
	Hyperlink: &InboundProxyHyperlink{
		Show: true,
		Text: "Click me...",
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var inventoryDnsAsset = InventoryDnsAsset{
	Id:   "3",
	Name: "Collab Edge v2",
}

var inventorySrvProtocol = InventorySrvProtocol{
	Id:       "TCP",
	Protocol: "_tcp",
}

var externalDnsRecord = ExternalDnsRecord{
	Uid:               "lonextdnsrec1",
	Hostname:          "lonhost",
	ARecord:           "lonhost.<subdomain>.dc-03.com",
	InventoryDnsAsset: &inventoryDnsAsset,
	NatRule:           &ExternalDnsNatRule{Uid: "lonvmpublicnat"},
	SrvRecords: []ExternalDnsSrvRecord{
		{
			Service:  "_sip",
			Protocol: "TCP",
			Port:     5060,
		},
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var mailServer = MailServer{
	Uid:               "lonmailserver",
	InventoryDnsAsset: &inventoryDnsAsset,
	VmNicTarget: &TrafficVmNicTarget{
		Uid:       "lonvm1msnic",
		IpAddress: "198.18.131.200",
		Vm: &Vm{
			Uid:  "lonvm1",
			Name: "Mail Server 1",
		},
	},
	Topology: &Topology{Uid: lonTopology.Uid},
}

var documentation = Documentation{
	Uid:              "lontopology",
	DocumentationUrl: "http://www.google.com",
}

var inventoryTelephonyItem = InventoryTelephonyItem{
	Id:          "1",
	Name:        "PSTN Services",
	Description: "PSTN Services (VCube)",
}

var telephonyItem = TelephonyItem{
	Uid:                    "lontelephonyitem1",
	Name:                   "PSTN Services",
	InventoryTelephonyItem: &inventoryTelephonyItem,
	Topology:               &Topology{Uid: lonTopology.Uid},
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
	suite.Equal(2, len(vms))
	suite.Contains(vms, vm)
}

func (suite *ContractTestSuite) TestGetVm() {

	// Given
	expectedVm := vm

	// When
	actualVm, err := suite.tbClient.GetVm(vm.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(expectedVm, *actualVm)
}

func (suite *ContractTestSuite) TestUpdateVm() {

	// Given
	expectedVm := vm
	expectedVm.DhcpConfig = &VmDhcpConfig{DefaultGatewayIp: "198.18.130.1"} // Match Contract
	expectedVm.DhcpConfig.PrimaryDnsIp = &primaryDnsIp
	expectedVm.DhcpConfig.SecondaryDnsIp = &secondaryDnsIp

	// Change Value
	expectedVm.Name = "New Name"

	// When
	actualVm, err := suite.tbClient.UpdateVm(expectedVm)
	suite.handleError(err)

	// Then

	// Work around contract data inconsistencies
	expectedNics := expectedVm.VmNetworkInterfaces

	for i := range expectedNics {
		expectedNics[i].Uid = ""
		expectedNics[i].Name = ""
	}
	expectedNics[0].Network.InventoryNetwork = nil

	suite.Equal(expectedVm, *actualVm)
}

var primaryDnsIp = "198.18.130.111"
var secondaryDnsIp = "198.18.130.112"

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
	actualVm.VmNetworkInterfaces[1].MacAddress = expectedVm.VmNetworkInterfaces[1].MacAddress // Missing in Contract
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
	expectedHw.TemplateConfigScript = &InventoryHwScript{Uid: hw.TemplateConfigScript.Uid}
	expectedHw.CustomScript = &InventoryHwScript{Uid: hw.CustomScript.Uid}
	expectedHw.ShutdownScript = &InventoryHwScript{Uid: hw.ShutdownScript.Uid}
	expectedHw.StartupScript = &InventoryHwScript{Uid: hw.StartupScript.Uid}
	expectedHw.NetworkInterfaces = []HwNic{
		hw.NetworkInterfaces[0],
		{
			Network: Network{
				Uid:  lonDefaultNetwork.Uid,
				Name: lonDefaultNetwork.Name,
			},
			NetworkInterface: InventoryHwNic{
				Id: "GigabitEthernet1/2",
			},
		},
	}

	// When
	actualHw, err := suite.tbClient.CreateHw(expectedHw)
	suite.handleError(err)

	// Then
	expectedHw.Uid = "newhardwareitem"
	expectedHw.NetworkInterfaces = nil // Match Contract
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

// IP Nat Rule

func (suite *ContractTestSuite) TestGetAllIpNatRules() {

	// When
	ipNatRules, err := suite.tbClient.GetAllIpNatRules(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(1, len(ipNatRules))
	suite.Contains(ipNatRules, ipNatRule)
}

// TODO - "Get One" Tests when Contracts are updated

func (suite *ContractTestSuite) TestCreateIpNatRule() {

	// Given
	expectedIpNatRule := ipNatRule

	// When
	actualIpNatRule, err := suite.tbClient.CreateIpNatRule(expectedIpNatRule)
	suite.handleError(err)

	// Then
	expectedIpNatRule.Uid = "newipnat"
	suite.Equal(expectedIpNatRule, *actualIpNatRule)
}

// ToDo "Delete" Tests when Contracts are updated

// VM Nat Rule

func (suite *ContractTestSuite) TestGetAllVmNatRules() {

	// When
	vmNatRules, err := suite.tbClient.GetAllVmNatRules(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(2, len(vmNatRules))
	expected := vmNatRule
	suite.Contains(vmNatRules, expected)
}

// TODO - "Get One" Tests when Contracts are updated

func (suite *ContractTestSuite) TestCreateVmNatRule() {

	// Given
	expectedVmNatRule := vmNatRule

	// When
	actualVmNatRule, err := suite.tbClient.CreateVmNatRule(expectedVmNatRule)
	suite.handleError(err)

	// Then
	// Use Contract Values
	expectedVmNatRule.Uid = "newvmnat"
	expectedVmNatRule.Target.Name = "Collab DB"
	expectedVmNatRule.Target.IpAddress = "198.18.133.111"
	suite.Equal(expectedVmNatRule, *actualVmNatRule)
}

// ToDo "Delete" Tests when Contracts are updated to include "Get One"

// Inbound Proxy Rule Tests

func (suite *ContractTestSuite) TestGetAllInboundProxyRules() {

	// When
	inboundProxyRules, err := suite.tbClient.GetAllInboundProxyRules(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(2, len(inboundProxyRules))
	suite.Contains(inboundProxyRules, inboundProxyRule)
}

func (suite *ContractTestSuite) TestCreateInboundProxyRule() {
	// Given
	expectedInboundProxyRule := inboundProxyRule

	// When
	actualInboundProxyRule, err := suite.tbClient.CreateInboundProxyRule(expectedInboundProxyRule)
	suite.handleError(err)

	// Then
	// Match Contract Wiremock Stubs
	expectedInboundProxyRule.Uid = "newloninboundproxy"
	expectedInboundProxyRule.VmNicTarget = &TrafficVmNicTarget{
		Uid:       inboundProxyRule.VmNicTarget.Uid,
		IpAddress: "192.168.0.9",
		Vm: &Vm{
			Uid:  "HIahoglmRva1DNeuYAII",
			Name: "MHJOGANCZJKHQFWFVSXD",
		},
	}
	suite.Equal(expectedInboundProxyRule, *actualInboundProxyRule)
}

func (suite *ContractTestSuite) TestDeleteInboundProxyRule() {

	// When
	err := suite.tbClient.DeleteInboundProxyRule(inboundProxyRule.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

// External DNS Record Tests

func (suite *ContractTestSuite) TestGetAllExternalDnsRecords() {

	// When
	externalDnsRecords, err := suite.tbClient.GetAllExternalDnsRecords(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(1, len(externalDnsRecords))
	suite.Contains(externalDnsRecords, externalDnsRecord)
}

func (suite *ContractTestSuite) TestGetExternalDnsRecord() {

	// Given
	expectedExternalDnsRecord := externalDnsRecord

	// When
	actualExternalDnsRecord, err := suite.tbClient.GetExternalDnsRecord(externalDnsRecord.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(expectedExternalDnsRecord, *actualExternalDnsRecord)
}

func (suite *ContractTestSuite) TestUpdateExternalDnsRecord() {

	// TODO - GET Contract needs to be updated to include ETAG
	suite.T().Skip("Skipping until contracts are updated")

	// Given
	expectedExternalDnsRecord := externalDnsRecord
	expectedExternalDnsRecord.ARecord = "newhost.<subdomain>.dc-03.com"
	expectedExternalDnsRecord.NatRule = &ExternalDnsNatRule{Uid: "lonipnatrule1"}

	// When
	actualExternalDnsRecord, err := suite.tbClient.UpdateExternalDnsRecord(expectedExternalDnsRecord)
	suite.handleError(err)

	// Then
	suite.Equal(expectedExternalDnsRecord, *actualExternalDnsRecord)
}

func (suite *ContractTestSuite) TestCreateExternalDnsRecord() {
	// Given
	expectedExternalDnsRecord := externalDnsRecord
	expectedExternalDnsRecord.NatRule = &ExternalDnsNatRule{Uid: "lonipnatrule1"}
	expectedExternalDnsRecord.SrvRecords = nil

	// When
	actualExternalDnsRecord, err := suite.tbClient.CreateExternalDnsRecord(expectedExternalDnsRecord)
	suite.handleError(err)

	// Then
	expectedExternalDnsRecord.Uid = "newextdnsrec"
	suite.Equal(expectedExternalDnsRecord, *actualExternalDnsRecord)
}

func (suite *ContractTestSuite) TestDeleteExternalDnsRecord() {

	// When
	err := suite.tbClient.DeleteExternalDnsRecord(externalDnsRecord.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

func (suite *ContractTestSuite) TestGetAllMailServers() {

	// When
	mailServers, err := suite.tbClient.GetAllMailServers(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(1, len(mailServers))
	suite.Contains(mailServers, mailServer)
}

func (suite *ContractTestSuite) TestCreateMailServer() {

	// Given
	// Match Contract expectations
	expectedMailServer := mailServer
	expectedMailServer.VmNicTarget = &TrafficVmNicTarget{
		Uid:       expectedMailServer.VmNicTarget.Uid,
		IpAddress: "192.168.0.9",
		Vm: &Vm{
			Uid:  "INLc0eWy1do3Xlk5GOdj",
			Name: "UZARRCZGMJSBXUAGLOLJ",
		},
	}
	expectedMailServer.InventoryDnsAsset = &InventoryDnsAsset{
		Id:   "3",
		Name: "CollabEdge_Swiss",
	}
	expectedMailServer.Uid = ""

	// When
	actualMailServer, err := suite.tbClient.CreateMailServer(expectedMailServer)
	suite.handleError(err)

	// Then
	expectedMailServer.Uid = "newmailserver"
	suite.Equal(expectedMailServer, *actualMailServer)
}

func (suite *ContractTestSuite) TestDeleteMailServer() {

	// When
	err := suite.tbClient.DeleteMailServer(mailServer.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

func (suite *ContractTestSuite) TestGetDocumentation() {

	// When
	actualDocumentation, err := suite.tbClient.GetDocumentation(documentation.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(documentation, *actualDocumentation)
}

func (suite *ContractTestSuite) TestUpdateDocumentation() {

	// Given
	expectedDocumentation := documentation
	expectedDocumentation.DocumentationUrl = "https://www.cisco.com"

	// When
	actualDocumentation, err := suite.tbClient.UpdateDocumentation(expectedDocumentation)
	suite.handleError(err)

	// Then
	suite.Equal(expectedDocumentation, *actualDocumentation)
}

func (suite *ContractTestSuite) TestGetAllTelephonyItems() {

	// When
	telephonyItems, err := suite.tbClient.GetAllTelephonyItems(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(1, len(telephonyItems))
	suite.Equal(telephonyItems[0], telephonyItem)
	suite.Contains(telephonyItems, telephonyItem)
}

func (suite *ContractTestSuite) TestCreateTelephonyItem() {

	// Given
	// Match contract
	inputTelephonyItem := TelephonyItem{
		InventoryTelephonyItem: &InventoryTelephonyItem{Id: inventoryTelephonyItem.Id},
		Topology:               &Topology{Uid: "rtptopology"},
	}

	// When
	actualTelephonyItem, err := suite.tbClient.CreateTelephonyItem(inputTelephonyItem)
	suite.handleError(err)

	// Then
	// Match contracts
	expectedTelephonyItem := telephonyItem
	expectedTelephonyItem.Uid = "newtelephonyitem"
	expectedTelephonyItem.Topology = inputTelephonyItem.Topology
	suite.Equal(expectedTelephonyItem, *actualTelephonyItem)
}

func (suite *ContractTestSuite) TestDeleteTelephonyItem() {

	// When
	err := suite.tbClient.DeleteTelephonyItem(telephonyItem.Uid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
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

func (suite *ContractTestSuite) TestGetAllEvcModes() {

	// When
	evcModes, err := suite.tbClient.GetAllEvcModes()
	suite.handleError(err)

	// Then
	suite.Equal(3, len(evcModes))
	suite.Contains(evcModes, westmereEvcMode)
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

	// When
	inventoryVms, err := suite.tbClient.GetAllInventoryVms(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(3, len(inventoryVms))
	suite.Equal(inventoryVms[2], inventoryVm)
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

func (suite *ContractTestSuite) TestGetAllInventoryDnsAssets() {

	// When
	inventoryDnsAssets, err := suite.tbClient.GetAllInventoryDnsAssets(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(2, len(inventoryDnsAssets))
	suite.Contains(inventoryDnsAssets, inventoryDnsAsset)
}

func (suite *ContractTestSuite) TestGetAllInventorySrvProtocols() {

	// When
	inventorySrvProtocols, err := suite.tbClient.GetAllInventorySrvProtocols()
	suite.handleError(err)

	// Then
	suite.Equal(5, len(inventorySrvProtocols))
	suite.Contains(inventorySrvProtocols, inventorySrvProtocol)
}

func (suite *ContractTestSuite) TestGetAllInventoryTelephonyItems() {

	// When
	inventoryTelephonyItems, err := suite.tbClient.GetAllInventoryTelephonyItems(lonTopology.Uid)
	suite.handleError(err)

	// Then
	suite.Equal(1, len(inventoryTelephonyItems))
	suite.Contains(inventoryTelephonyItems, inventoryTelephonyItem)
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
		fmt.Println("Waiting another 2s for wiremock...")
		time.Sleep(2 * time.Second)
	}

	return c
}

func (suite *ContractTestSuite) handleError(err error) {
	if err != nil {
		suite.FailNow(err.Error())
	}
}
