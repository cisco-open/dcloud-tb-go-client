package tbclient

func (c *Client) getInventoryNetworkService(topologyUid string) *readService[InventoryNetwork, inventoryNetworkCollection] {
	return &readService[InventoryNetwork, inventoryNetworkCollection]{
		client:       c,
		resourcePath: "/inventory-networks",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetAllInventoryNetworks(topologyUid string) ([]InventoryNetwork, error) {
	return c.getInventoryNetworkService(topologyUid).getAll()
}
