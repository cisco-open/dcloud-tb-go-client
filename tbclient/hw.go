// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInventoryHwService(topologyUid string) *collectionService[InventoryHw, inventoryHwCollection] {
	return &collectionService[InventoryHw, inventoryHwCollection]{
		client:       c,
		resourcePath: "/inventory-hardware-items",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getHwService(topologyUid string) *resourceService[Hw, hwCollection] {
	return &resourceService[Hw, hwCollection]{
		collectionService[Hw, hwCollection]{
			client:       c,
			resourcePath: "/hardware-items",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInventoryHws(topologyUid string) ([]InventoryHw, error) {
	return c.getInventoryHwService(topologyUid).getAll()
}

func (c *Client) GetAllHws(topologyUid string) ([]Hw, error) {
	return c.getHwService(topologyUid).getAll()
}

func (c *Client) GetHw(uid string) (*Hw, error) {
	return c.getHwService("").getOne(uid)
}

func (c *Client) CreateHw(hw Hw) (*Hw, error) {
	return c.getHwService("").create(hw)
}

func (c *Client) UpdateHw(hw Hw) (*Hw, error) {
	return c.getHwService("").update(hw.Uid, hw)
}

func (c *Client) DeleteHw(uid string) error {
	return c.getHwService("").delete(uid)
}
