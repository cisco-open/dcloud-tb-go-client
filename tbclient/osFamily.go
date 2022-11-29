package tbclient

func (c *Client) getOsFamilyService() *collectionService[OsFamily, osFamilyCollection] {
	return &collectionService[OsFamily, osFamilyCollection]{
		client:       c,
		resourcePath: "/os-families",
	}
}

func (c *Client) GetAllOsFamilies() ([]OsFamily, error) {
	return c.getOsFamilyService().getAll()
}
