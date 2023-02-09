/**
Copyright (c) 2023 Cisco Systems, Inc. and its affiliates
*/

package tbclient

func (c *Client) getTopologyService(requestHeaders map[string]string) *resourceService[Topology, topologyCollection] {
	return &resourceService[Topology, topologyCollection]{
		collectionService[Topology, topologyCollection]{
			client:         c,
			resourcePath:   "/topologies",
			requestHeaders: requestHeaders,
		},
	}
}

func (c *Client) GetAllTopologies() ([]Topology, error) {
	return c.getTopologyService(emptyHeaders).getAll()
}

func (c *Client) GetTopology(uid string) (*Topology, error) {
	return c.getTopologyService(emptyHeaders).getOne(uid)
}

func (c *Client) CreateTopology(topology Topology) (*Topology, error) {
	return c.getTopologyService(map[string]string{"x-defaultnetwork": "false"}).create(topology)
}

func (c *Client) UpdateTopology(topology Topology) (*Topology, error) {
	return c.getTopologyService(emptyHeaders).update(topology.Uid, topology)
}

func (c *Client) DeleteTopology(topologyUid string) error {
	return c.getTopologyService(emptyHeaders).delete(topologyUid)
}
