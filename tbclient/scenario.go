// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getScenarioService(topologyUid string) *singleNestedResourceService[Scenario] {
	return &singleNestedResourceService[Scenario]{
		client:         c,
		resourcePath:   "/scenario",
		topologyUid:    topologyUid,
		putUsingGetUrl: true,
	}
}

func (c *Client) GetScenario(topologyUid string) (*Scenario, error) {
	return c.getScenarioService(topologyUid).get()
}

func (c *Client) UpdateScenario(Scenario Scenario) (*Scenario, error) {
	return c.getScenarioService(Scenario.Topology.Uid).update(Scenario.Uid, Scenario)
}
