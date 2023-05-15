package tbclient

import (
	"fmt"
)

func (c *Client) GetVmStartOrder(topologyUid string) (*VmStartOrder, error) {

	rest := c.createRestClient()
	url := generateVmStartGetUrl(c.HostURL, topologyUid)

	resp, err := executeGet(rest, c.Token, url, VmStartOrder{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*VmStartOrder), nil
}

func (c *Client) UpdateVmStartOrder(vmStartOrder VmStartOrder) (*VmStartOrder, error) {

	rest := c.createRestClient()
	getUrl := generateVmStartGetUrl(c.HostURL, vmStartOrder.Topology.Uid)

	current, err := executeGet(rest, c.Token, getUrl, VmStartOrder{})
	if err != nil {
		return nil, err
	}

	etag := current.Header().Get(etagHeader)

	putUrl := fmt.Sprintf("%s/vm-start-order/%s", c.HostURL, vmStartOrder.Uid)

	updated, err := executePut(rest, c.Token, putUrl, vmStartOrder, VmStartOrder{}, etag)
	if err != nil {
		return nil, err
	}

	return updated.Result().(*VmStartOrder), nil
}

func generateVmStartGetUrl(hostUrl string, topologyUid string) string {
	return fmt.Sprintf("%s/topologies/%s/vm-start-order", hostUrl, topologyUid)
}

func (c *Client) GetVmStopOrder(topologyUid string) (*VmStopOrder, error) {

	rest := c.createRestClient()
	url := generateVmStopGetUrl(c.HostURL, topologyUid)

	resp, err := executeGet(rest, c.Token, url, VmStopOrder{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*VmStopOrder), nil
}

func (c *Client) UpdateVmStopOrder(vmStopOrder VmStopOrder) (*VmStopOrder, error) {

	rest := c.createRestClient()
	getUrl := generateVmStopGetUrl(c.HostURL, vmStopOrder.Topology.Uid)

	current, err := executeGet(rest, c.Token, getUrl, VmStopOrder{})
	if err != nil {
		return nil, err
	}

	etag := current.Header().Get(etagHeader)

	putUrl := fmt.Sprintf("%s/vm-stop-order/%s", c.HostURL, vmStopOrder.Uid)

	updated, err := executePut(rest, c.Token, putUrl, vmStopOrder, VmStopOrder{}, etag)
	if err != nil {
		return nil, err
	}

	return updated.Result().(*VmStopOrder), nil
}

func generateVmStopGetUrl(hostUrl string, topologyUid string) string {
	return fmt.Sprintf("%s/topologies/%s/vm-stop-order", hostUrl, topologyUid)
}

func (c *Client) GetHwStartOrder(topologyUid string) (*HwStartOrder, error) {

	rest := c.createRestClient()
	url := generateHwStartGetUrl(c.HostURL, topologyUid)

	resp, err := executeGet(rest, c.Token, url, HwStartOrder{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*HwStartOrder), nil
}

func (c *Client) UpdateHwStartOrder(hwStartOrder HwStartOrder) (*HwStartOrder, error) {

	rest := c.createRestClient()
	getUrl := generateHwStartGetUrl(c.HostURL, hwStartOrder.Topology.Uid)

	current, err := executeGet(rest, c.Token, getUrl, HwStartOrder{})
	if err != nil {
		return nil, err
	}

	etag := current.Header().Get(etagHeader)

	putUrl := fmt.Sprintf("%s/hardware-start-order/%s", c.HostURL, hwStartOrder.Uid)

	updated, err := executePut(rest, c.Token, putUrl, hwStartOrder, HwStartOrder{}, etag)
	if err != nil {
		return nil, err
	}

	return updated.Result().(*HwStartOrder), nil
}

func generateHwStartGetUrl(hostUrl string, topologyUid string) string {
	return fmt.Sprintf("%s/topologies/%s/hardware-start-order", hostUrl, topologyUid)
}
