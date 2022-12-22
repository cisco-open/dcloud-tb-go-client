package tbclient

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
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

func (suite *ContractTestSuite) SetupSuite() {
	suite.docker = startStubrunner(suite)
	suite.tbClient = createTbClient(suite)
}

func (suite *ContractTestSuite) TearDownSuite() {
	suite.docker.client.ContainerRemove(*suite.docker.ctx, suite.docker.containerId, types.ContainerRemoveOptions{Force: true})
}

// Asset Tests

func (suite *ContractTestSuite) TestGetAllTopologies() {

	// Given
	expectedTopology := lonTopology
	expectedTopology.Notes = "" // Get All Topologies contract is missing 'notes' field

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

func TestContractTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}

func startStubrunner(suite *ContractTestSuite) DockerContainer {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	suite.handleError(err)
	defer cli.Close()

	imageName := "springcloud/spring-cloud-contract-stub-runner:2.2.2.RELEASE"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	suite.handleError(err)
	defer out.Close()

	io.Copy(os.Stdout, out)

	environment := []string{
		"STUBRUNNER_STUBS_MODE=REMOTE",
		"STUBRUNNER_IDS=com.cisco.kapua:kapua-content-studio-api-contracts:1.0-SNAPSHOT:stubs:9876",
		"REPO_WITH_BINARIES_URL=http://engci-maven.cisco.com/artifactory/xse-snapshot/",
		"STUBRUNNER_REPOSITORY_ROOT=http://engci-maven.cisco.com/artifactory/xse-snapshot/",
	}

	host_config := &container.HostConfig{
		PortBindings: nat.PortMap{
			"9876/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "9876",
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   environment,
		ExposedPorts: nat.PortSet{
			"9876/tcp": struct{}{},
		},
	}, host_config, nil, nil, "")
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
	}

	timeout := time.Now().Add(time.Minute)
	for _, err := c.GetAllTopologies(); err != nil; _, err = c.GetAllTopologies() {
		if time.Now().After(timeout) {
			suite.T().Log("Timeout starting stub runner")
			suite.handleError(err)
		}
		fmt.Println("Waiting another 5s for stub runner...")
		time.Sleep(5 * time.Second)
	}

	return c
}

func (suite *ContractTestSuite) handleError(err error) {
	if err != nil {
		suite.FailNow(err.Error())
	}
}
