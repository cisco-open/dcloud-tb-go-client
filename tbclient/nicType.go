// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getNicTypeService() *collectionService[NicType, nicTypeCollection] {
	return &collectionService[NicType, nicTypeCollection]{
		client:       c,
		resourcePath: "/vm-network-interface-types",
	}
}

func (c *Client) GetAllNicTypes() ([]NicType, error) {
	return c.getNicTypeService().getAll()
}
