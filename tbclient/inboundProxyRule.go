// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

func (c *Client) getInboundProxyRuleService(topologyUid string) *resourceService[InboundProxyRule, inboundProxyRulesCollection] {
	return &resourceService[InboundProxyRule, inboundProxyRulesCollection]{
		collectionService[InboundProxyRule, inboundProxyRulesCollection]{
			client:       c,
			resourcePath: "/inbound-proxy-rules",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInboundProxyRules(topologyUid string) ([]InboundProxyRule, error) {
	return c.getInboundProxyRuleService(topologyUid).getAll()
}

func (c *Client) CreateInboundProxyRule(inboundProxyRule InboundProxyRule) (*InboundProxyRule, error) {
	return c.getInboundProxyRuleService("").create(inboundProxyRule)
}

func (c *Client) DeleteInboundProxyRule(uid string) error {
	return c.getInboundProxyRuleService("").delete(uid)
}
