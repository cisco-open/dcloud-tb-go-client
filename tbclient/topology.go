package tbclient

func (c *Client) getTopologyService() *resourceService[Topology, topologyCollection] {
	return &resourceService[Topology, topologyCollection]{
		collectionService[Topology, topologyCollection]{
			client:       c,
			resourcePath: "/topologies",
		},
	}
}

func (c *Client) GetAllTopologies() ([]Topology, error) {
	return c.getTopologyService().getAll()
}

func (c *Client) GetTopology(uid string) (*Topology, error) {
	return c.getTopologyService().getOne(uid)
}

func (c *Client) CreateTopology(topology Topology) (*Topology, error) {
	return c.getTopologyService().create(topology)
}

func (c *Client) UpdateTopology(topology Topology) (*Topology, error) {
	return c.getTopologyService().update(topology.Uid, topology)
}

func (c *Client) DeleteTopology(topologyUid string) error {
	return c.getTopologyService().delete(topologyUid)
}
