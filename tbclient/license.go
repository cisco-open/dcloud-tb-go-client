package tbclient

func (c *Client) getInventoryLicenseService(topologyUid string) *collectionService[InventoryLicense, inventoryLicenseCollection] {
	return &collectionService[InventoryLicense, inventoryLicenseCollection]{
		client:       c,
		resourcePath: "/inventory-licenses",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetAllInventoryLicenses(topologyUid string) ([]InventoryLicense, error) {
	return c.getInventoryLicenseService(topologyUid).getAll()
}
