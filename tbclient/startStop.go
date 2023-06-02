// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getVmStartOrderService(topologyUid string) *singleNestedResourceService[VmStartOrder] {
	return &singleNestedResourceService[VmStartOrder]{
		client:       c,
		resourcePath: "/vm-start-order",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetVmStartOrder(topologyUid string) (*VmStartOrder, error) {
	return c.getVmStartOrderService(topologyUid).get()
}

func (c *Client) UpdateVmStartOrder(vmStartOrder VmStartOrder) (*VmStartOrder, error) {
	return c.getVmStartOrderService(vmStartOrder.Topology.Uid).update(vmStartOrder.Uid, vmStartOrder)
}

func (c *Client) getVmStopOrderService(topologyUid string) *singleNestedResourceService[VmStopOrder] {
	return &singleNestedResourceService[VmStopOrder]{
		client:       c,
		resourcePath: "/vm-stop-order",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetVmStopOrder(topologyUid string) (*VmStopOrder, error) {
	return c.getVmStopOrderService(topologyUid).get()
}

func (c *Client) UpdateVmStopOrder(vmStopOrder VmStopOrder) (*VmStopOrder, error) {
	return c.getVmStopOrderService(vmStopOrder.Topology.Uid).update(vmStopOrder.Uid, vmStopOrder)
}

func (c *Client) getHwStartOrderService(topologyUid string) *singleNestedResourceService[HwStartOrder] {
	return &singleNestedResourceService[HwStartOrder]{
		client:       c,
		resourcePath: "/hardware-start-order",
		topologyUid:  topologyUid,
	}
}

func (c *Client) GetHwStartOrder(topologyUid string) (*HwStartOrder, error) {
	return c.getHwStartOrderService(topologyUid).get()
}

func (c *Client) UpdateHwStartOrder(hwStartOrder HwStartOrder) (*HwStartOrder, error) {
	return c.getHwStartOrderService(hwStartOrder.Topology.Uid).update(hwStartOrder.Uid, hwStartOrder)
}
