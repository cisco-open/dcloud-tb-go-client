// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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
