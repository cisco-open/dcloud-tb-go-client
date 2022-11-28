package tbclient

func (c *Client) getInventoryVmService(topologyUid string) *collectionService[InventoryVm, inventoryVmCollection] {
	return &collectionService[InventoryVm, inventoryVmCollection]{
		client:       c,
		resourcePath: "/inventory-template-vms",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetAllInventoryVms(topologyUid string) ([]InventoryVm, error) {
	return c.getInventoryVmService(topologyUid).getAll()
}
