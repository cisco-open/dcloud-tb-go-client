// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInventoryDnsAssetService(topologyUid string) *collectionService[InventoryDnsAsset, inventoryDnsAssetCollection] {
	return &collectionService[InventoryDnsAsset, inventoryDnsAssetCollection]{
		client:       c,
		resourcePath: "/inventory-dns-assets",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getInventorySrvProtocolService() *collectionService[InventorySrvProtocol, inventorySrvProtocolCollection] {
	return &collectionService[InventorySrvProtocol, inventorySrvProtocolCollection]{
		client:       c,
		resourcePath: "/srv-protocols",
	}
}

func (c *Client) GetAllInventoryDnsAssets(topologyUid string) ([]InventoryDnsAsset, error) {
	return c.getInventoryDnsAssetService(topologyUid).getAll()
}

func (c *Client) GetAllInventorySrvProtocols() ([]InventorySrvProtocol, error) {
	return c.getInventorySrvProtocolService().getAll()
}
