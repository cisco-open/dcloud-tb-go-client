// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getTopologyService(requestHeaders map[string]string) *resourceService[Topology, topologyCollection] {
	return &resourceService[Topology, topologyCollection]{
		collectionService[Topology, topologyCollection]{
			client:         c,
			resourcePath:   "/topologies",
			requestHeaders: requestHeaders,
		},
	}
}

func (c *Client) GetAllTopologies() ([]Topology, error) {
	return c.getTopologyService(emptyHeaders).getAll()
}

func (c *Client) GetTopology(uid string) (*Topology, error) {
	return c.getTopologyService(emptyHeaders).getOne(uid)
}

func (c *Client) CreateTopology(topology Topology) (*Topology, error) {
	return c.getTopologyService(map[string]string{"x-defaultnetwork": "false"}).create(topology)
}

func (c *Client) UpdateTopology(topology Topology) (*Topology, error) {
	return c.getTopologyService(emptyHeaders).update(topology.Uid, topology)
}

func (c *Client) DeleteTopology(topologyUid string) error {
	return c.getTopologyService(emptyHeaders).delete(topologyUid)
}
