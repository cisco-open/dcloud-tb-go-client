package tbclient

import (
	"errors"
	"fmt"
)

type natRule struct {
	Uid      string    `json:"uid"`
	Topology *Topology `json:"topology"`
	EastWest bool      `json:"eastWest"`
	Target   struct {
		Type       string      `json:"type"`
		TargetItem interface{} `json:"targetItem"`
		Name       string      `json:"name"`
		IpAddress  string      `json:"ipAddress"`
	} `json:"target"`
}

type natRuleCollection struct {
	Data []natRule `json:"natRules"`
}

func (n natRuleCollection) getData() []natRule {
	return n.Data
}

func (c *Client) GetAllIpNatRules(topologyUid string) ([]IpNatRule, error) {

	url := fmt.Sprintf("%s/topologies/%s/nat-rules", c.HostURL, topologyUid)

	resp, err := executeGet(c.createRestClient(), c.Token, url, collection[natRuleCollection]{})
	if err != nil {
		return nil, err
	}

	natRules := resp.Result().(*collection[natRuleCollection]).Embedded.getData()
	ipNatRules := make([]IpNatRule, 0)

	for _, natRule := range natRules {
		if natRule.Target.Type == "IP" {
			ipNatRule := IpNatRule{
				Uid:      natRule.Uid,
				EastWest: natRule.EastWest,
				Target: IpNatTarget{
					IpAddress: natRule.Target.IpAddress,
					Name:      natRule.Target.Name,
				},
				Topology: natRule.Topology,
			}
			ipNatRules = append(ipNatRules, ipNatRule)
		}
	}

	return ipNatRules, nil
}

func (c *Client) GetIpNatRule(uid string) (*IpNatRule, error) {

	url := fmt.Sprintf("%s/nat-rules/%s", c.HostURL, uid)

	resp, err := executeGet(c.createRestClient(), c.Token, url, natRule{})
	if err != nil {
		return nil, err
	}

	natRule := resp.Result().(natRule)
	if natRule.Target.Type != "IP" {
		return nil, errors.New("IP NAT Rule Not Found")
	}

	return &IpNatRule{
		Uid:      natRule.Uid,
		EastWest: natRule.EastWest,
		Target: IpNatTarget{
			IpAddress: natRule.Target.IpAddress,
			Name:      natRule.Target.Name,
		},
		Topology: natRule.Topology,
	}, nil
}

func (c *Client) CreateIpNatRule(ipNatRule IpNatRule) (*IpNatRule, error) {

	rest := c.createRestClient()
	url := fmt.Sprintf("%s/ip-nat-rules", c.HostURL)

	resp, err := executePost(rest, c.Token, url, ipNatRule, IpNatRule{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*IpNatRule), nil
}

func (c *Client) DeleteIpNatRule(uid string) error {

	_, err := c.GetIpNatRule(uid)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/nat-rules/%s", c.HostURL, uid)

	return executeDelete(c.createRestClient(), c.Token, url)
}
