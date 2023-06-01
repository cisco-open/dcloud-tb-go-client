// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInventoryHwScriptService(topologyUid string) *collectionService[InventoryHwScript, inventoryHwScriptCollection] {
	return &collectionService[InventoryHwScript, inventoryHwScriptCollection]{
		client:       c,
		resourcePath: "/inventory-hardware-scripts",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getInventoryHwTemplateConfigService(topologyUid string) *collectionService[InventoryHwScript, inventoryHwScriptCollection] {
	return &collectionService[InventoryHwScript, inventoryHwScriptCollection]{
		client:       c,
		resourcePath: "/inventory-hardware-templates",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetAllInventoryHwScripts(topologyUid string) ([]InventoryHwScript, error) {
	return c.getInventoryHwScriptService(topologyUid).getAll()
}

func (c *Client) GetAllInventoryHwTemplateConfigs(topologyUid string) ([]InventoryHwScript, error) {
	return c.getInventoryHwTemplateConfigService(topologyUid).getAll()
}
