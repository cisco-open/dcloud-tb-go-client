// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInventoryLicenseService(topologyUid string) *collectionService[InventoryLicense, inventoryLicenseCollection] {
	return &collectionService[InventoryLicense, inventoryLicenseCollection]{
		client:       c,
		resourcePath: "/inventory-licenses",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getLicenseService(topologyUid string) *resourceService[License, licenseCollection] {
	return &resourceService[License, licenseCollection]{
		collectionService[License, licenseCollection]{
			client:       c,
			resourcePath: "/licenses",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInventoryLicenses(topologyUid string) ([]InventoryLicense, error) {
	return c.getInventoryLicenseService(topologyUid).getAll()
}

func (c *Client) GetAllLicenses(topologyUid string) ([]License, error) {
	return c.getLicenseService(topologyUid).getAll()
}

func (c *Client) GetLicense(uid string) (*License, error) {
	return c.getLicenseService("").getOne(uid)
}

func (c *Client) CreateLicense(license License) (*License, error) {
	return c.getLicenseService("").create(license)
}

func (c *Client) UpdateLicense(license License) (*License, error) {
	return c.getLicenseService("").update(license.Uid, license)
}

func (c *Client) DeleteLicense(uid string) error {
	return c.getLicenseService("").delete(uid)
}
