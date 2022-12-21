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

var testResources = map[string]interface{}{
	"lontopology": Topology{
		Uid:         "lontopology",
		Name:        "lontopology name",
		Description: "lontopology description",
		Datacenter:  "LON",
		Notes:       "Notes",
		Status:      "DRAFT",
	},
	"linux": OsFamily{
		Id:   "LINUX",
		Name: "Linux",
	},
	"e1000": NicType{
		Id:   "VIRTUAL_E1000",
		Name: "VirtualE1000",
	},
	"script1": InventoryHwScript{
		Uid:  "script1",
		Name: "EN-3850-switch-start.xml",
	},
	"template1": InventoryHwScript{
		Uid:  "template1",
		Name: "en-sw-3850-16-v1.txt",
	},
}

func (suite *ContractTestSuite) SetupSuite() {
	suite.docker, suite.tbClient = startStubrunner()
}

func (suite *ContractTestSuite) TearDownSuite() {
	suite.docker.client.ContainerRemove(*suite.docker.ctx, suite.docker.containerId, types.ContainerRemoveOptions{Force: true})
}

func (suite *ContractTestSuite) TestGetAllTopologies() {

	// Given
	lontopology := testResources["lontopology"].(Topology)
	lontopology.Notes = "" // Get All Topologies contract is missing 'notes' field

	// When
	topologies, err := suite.tbClient.GetAllTopologies()
	suite.handleError(err)

	// Then
	suite.Equal(9, len(topologies))
	suite.Contains(topologies, lontopology)
}

func (suite *ContractTestSuite) TestGetTopology() {
	// Given
	topologyUid := "lontopology"

	// When
	topology, err := suite.tbClient.GetTopology(topologyUid)
	suite.handleError(err)

	// Then
	suite.Equal(testResources[topologyUid], *topology)
}

func (suite *ContractTestSuite) TestUpdateTopology() {
	// Given
	lontopology := testResources["lontopology"].(Topology)
	lontopology.Name = "lontopology name updated"

	// When
	topology, err := suite.tbClient.UpdateTopology(lontopology)
	suite.handleError(err)

	// Then
	suite.Equal(lontopology, *topology)
}

func (suite *ContractTestSuite) TestCreateTopology() {
	// Given
	newTopology := testResources["lontopology"].(Topology)
	newTopology.Uid = ""

	// When
	topology, err := suite.tbClient.CreateTopology(newTopology)
	suite.handleError(err)

	// Then
	newTopology.Uid = "newtopology"
	suite.Equal(newTopology, *topology)
}

func (suite *ContractTestSuite) TestDeleteTopology() {
	// Given
	topologyUid := "lontopology"

	// When
	err := suite.tbClient.DeleteTopology(topologyUid)
	suite.handleError(err)

	// Then
	suite.Nil(err)
}

func (suite *ContractTestSuite) TestGetAllOsFamilies() {

	// When
	osFamilies, err := suite.tbClient.GetAllOsFamilies()
	suite.handleError(err)

	// Then
	suite.Equal(4, len(osFamilies))
	suite.Contains(osFamilies, testResources["linux"])
}

func (suite *ContractTestSuite) TestGetAllNicTypes() {

	// When
	nicTypes, err := suite.tbClient.GetAllNicTypes()
	suite.handleError(err)

	// Then
	suite.Equal(5, len(nicTypes))
	suite.Contains(nicTypes, testResources["e1000"])
}

func (suite *ContractTestSuite) TestGetAllInventoryHwScripts() {

	// When
	hwScripts, err := suite.tbClient.GetAllInventoryHwScripts("lontopology")
	suite.handleError(err)

	// Then
	suite.Equal(4, len(hwScripts))
	suite.Contains(hwScripts, testResources["script1"])
}

func (suite *ContractTestSuite) TestGetAllInventoryHwTemplateConfigs() {

	// When
	configTemplates, err := suite.tbClient.GetAllInventoryHwTemplateConfigs("lontopology")
	suite.handleError(err)

	// Then
	suite.Equal(2, len(configTemplates))
	suite.Contains(configTemplates, testResources["template1"])
}

func TestContractTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}

func startStubrunner() (DockerContainer, Client) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	imageName := "springcloud/spring-cloud-contract-stub-runner:2.2.2.RELEASE"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
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
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
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
			"8083/tcp": struct{}{},
			"9876/tcp": struct{}{},
		},
	}, host_config, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	c := Client{
		HostURL:     "http://localhost:9876",
		Token:       "oauthAuthorized",
		DisableGzip: true,
	}

	timeout := time.Now().Add(time.Minute)
	for _, err := c.GetAllTopologies(); err != nil; _, err = c.GetAllTopologies() {
		if time.Now().After(timeout) {
			panic("Timeout starting stub runner")
		}
		fmt.Println("Waiting for stub runner...")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		panic(err)
	}

	return DockerContainer{
		client:      cli,
		ctx:         &ctx,
		containerId: resp.ID,
	}, c
}

func (suite *ContractTestSuite) handleError(err error) {
	if err != nil {
		suite.FailNow(err.Error())
	}
}
