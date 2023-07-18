// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInventoryTelephonyService(topologyUid string) *collectionService[InventoryTelephonyItem, inventoryTelephonyItemCollection] {
	return &collectionService[InventoryTelephonyItem, inventoryTelephonyItemCollection]{
		client:       c,
		resourcePath: "/inventory-telephony-items",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getTelephonyItemService(topologyUid string) *resourceService[TelephonyItem, telephonyItemCollection] {
	return &resourceService[TelephonyItem, telephonyItemCollection]{
		collectionService[TelephonyItem, telephonyItemCollection]{
			client:       c,
			resourcePath: "/telephony-items",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInventoryTelephonyItems(topologyUid string) ([]InventoryTelephonyItem, error) {
	return c.getInventoryTelephonyService(topologyUid).getAll()
}

func (c *Client) GetAllTelephonyItems(topologyUid string) ([]TelephonyItem, error) {
	return c.getTelephonyItemService(topologyUid).getAll()
}

func (c *Client) CreateTelephonyItem(telephonyItem TelephonyItem) (*TelephonyItem, error) {
	return c.getTelephonyItemService("").create(telephonyItem)
}

func (c *Client) DeleteTelephonyItem(uid string) error {
	return c.getTelephonyItemService("").delete(uid)
}
