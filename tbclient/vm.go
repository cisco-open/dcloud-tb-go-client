// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInventoryVmService(topologyUid string) *collectionService[InventoryVm, inventoryVmCollection] {
	return &collectionService[InventoryVm, inventoryVmCollection]{
		client:       c,
		resourcePath: "/inventory-template-vms",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getVmService(topologyUid string) *resourceService[Vm, vmCollection] {
	return &resourceService[Vm, vmCollection]{
		collectionService[Vm, vmCollection]{
			client:       c,
			resourcePath: "/vms",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInventoryVms(topologyUid string) ([]InventoryVm, error) {
	return c.getInventoryVmService(topologyUid).getAll()
}

func (c *Client) GetAllVms(topologyUid string) ([]Vm, error) {
	return c.getVmService(topologyUid).getAll()
}

func (c *Client) GetVm(uid string) (*Vm, error) {
	return c.getVmService("").getOne(uid)
}

func (c *Client) CreateVm(vm Vm) (*Vm, error) {
	return c.getVmService("").create(vm)
}

func (c *Client) UpdateVm(vm Vm) (*Vm, error) {
	return c.getVmService("").update(vm.Uid, vm)
}

func (c *Client) DeleteVm(vmUid string) error {
	return c.getVmService("").delete(vmUid)
}
