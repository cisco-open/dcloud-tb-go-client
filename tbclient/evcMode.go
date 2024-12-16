// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getEvcModeService() *collectionService[EvcMode, evcModeCollection] {
	return &collectionService[EvcMode, evcModeCollection]{
		client:       c,
		resourcePath: "/evc-modes",
	}
}

func (c *Client) GetAllEvcModes() ([]EvcMode, error) {
	return c.getEvcModeService().getAll()
}
