package tbclient

func (c *Client) getInventoryHwService(topologyUid string) *collectionService[InventoryHw, inventoryHwCollection] {
	return &collectionService[InventoryHw, inventoryHwCollection]{
		client:       c,
		resourcePath: "/inventory-hardware-items",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetAllInventoryHws(topologyUid string) ([]InventoryHw, error) {
	return c.getInventoryHwService(topologyUid).getAll()
}
