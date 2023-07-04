package tbclient

func (c *Client) getMailServerService(topologyUid string) *resourceService[MailServer, mailServerCollection] {
	return &resourceService[MailServer, mailServerCollection]{
		collectionService[MailServer, mailServerCollection]{
			client:       c,
			resourcePath: "/mail-servers",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllMailServers(topologyUid string) ([]MailServer, error) {
	return c.getMailServerService(topologyUid).getAll()
}

func (c *Client) CreateMailServer(mailServer MailServer) (*MailServer, error) {
	return c.getMailServerService("").create(mailServer)
}

func (c *Client) DeleteMailServer(uid string) error {
	return c.getMailServerService("").delete(uid)
}
