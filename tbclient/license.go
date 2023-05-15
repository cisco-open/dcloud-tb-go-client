package tbclient

func (c *Client) getInventoryLicenseService(topologyUid string) *collectionService[InventoryLicense, inventoryLicenseCollection] {
	return &collectionService[InventoryLicense, inventoryLicenseCollection]{
		client:       c,
		resourcePath: "/inventory-licenses",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getLicenseService(topologyUid string) *resourceService[License, licenseCollection] {
	return &resourceService[License, licenseCollection]{
		collectionService[License, licenseCollection]{
			client:       c,
			resourcePath: "/licenses",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInventoryLicenses(topologyUid string) ([]InventoryLicense, error) {
	return c.getInventoryLicenseService(topologyUid).getAll()
}

func (c *Client) GetAllLicenses(topologyUid string) ([]License, error) {
	return c.getLicenseService(topologyUid).getAll()
}

func (c *Client) GetLicense(uid string) (*License, error) {
	return c.getLicenseService("").getOne(uid)
}

func (c *Client) CreateLicense(license License) (*License, error) {
	return c.getLicenseService("").create(license)
}

func (c *Client) UpdateLicense(license License) (*License, error) {
	return c.getLicenseService("").update(license.Uid, license)
}

func (c *Client) DeleteLicense(uid string) error {
	return c.getLicenseService("").delete(uid)
}
