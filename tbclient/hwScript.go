package tbclient

func (c *Client) getInventoryHwScriptService(topologyUid string) *collectionService[InventoryHwScript, inventoryHwScriptollection] {
	return &collectionService[InventoryHwScript, inventoryHwScriptollection]{
		client:       c,
		resourcePath: "/inventory-hardware-scripts",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getInventoryHwTemplateConfigService(topologyUid string) *collectionService[InventoryHwScript, inventoryHwScriptollection] {
	return &collectionService[InventoryHwScript, inventoryHwScriptollection]{
		client:       c,
		resourcePath: "/inventory-hardware-templates",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetAllInventoryHwScripts(topologyUid string) ([]InventoryHwScript, error) {
	return c.getInventoryHwScriptService(topologyUid).getAll()
}

func (c *Client) GetAllInventoryHwTemplateConfigs(topologyUid string) ([]InventoryHwScript, error) {
	return c.getInventoryHwTemplateConfigService(topologyUid).getAll()
}
