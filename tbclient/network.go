// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInventoryNetworkService(topologyUid string) *collectionService[InventoryNetwork, inventoryNetworkCollection] {
	return &collectionService[InventoryNetwork, inventoryNetworkCollection]{
		client:       c,
		resourcePath: "/inventory-networks",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getNetworkService(topologyUid string) *resourceService[Network, networkCollection] {
	return &resourceService[Network, networkCollection]{
		collectionService[Network, networkCollection]{
			client:       c,
			resourcePath: "/networks",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInventoryNetworks(topologyUid string) ([]InventoryNetwork, error) {
	return c.getInventoryNetworkService(topologyUid).getAll()
}

func (c *Client) GetAllNetworks(topologyUid string) ([]Network, error) {
	return c.getNetworkService(topologyUid).getAll()
}

func (c *Client) GetNetwork(uid string) (*Network, error) {
	return c.getNetworkService("").getOne(uid)
}

func (c *Client) CreateNetwork(network Network) (*Network, error) {
	return c.getNetworkService("").create(network)
}

func (c *Client) UpdateNetwork(network Network) (*Network, error) {
	return c.getNetworkService("").update(network.Uid, network)
}

func (c *Client) DeleteNetwork(networkUid string) error {
	return c.getNetworkService("").delete(networkUid)
}
