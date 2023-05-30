package tbclient

func (c *Client) getRemoteAccessService(topologyUid string) *singleNestedResourceService[RemoteAccess] {
	return &singleNestedResourceService[RemoteAccess]{
		client:       c,
		resourcePath: "/remote-access",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetRemoteAccess(topologyUid string) (*RemoteAccess, error) {
	return c.getRemoteAccessService(topologyUid).get()
}

func (c *Client) UpdateRemoteAccess(remoteAccess RemoteAccess) (*RemoteAccess, error) {
	return c.getRemoteAccessService(remoteAccess.Topology.Uid).update(remoteAccess.Uid, remoteAccess)
}
