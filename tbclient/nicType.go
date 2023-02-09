/**
Copyright (c) 2023 Cisco Systems, Inc. and its affiliates
*/

package tbclient

func (c *Client) getNicTypeService() *collectionService[NicType, nicTypeCollection] {
	return &collectionService[NicType, nicTypeCollection]{
		client:       c,
		resourcePath: "/vm-network-interface-types",
	}
}

func (c *Client) GetAllNicTypes() ([]NicType, error) {
	return c.getNicTypeService().getAll()
}
