package tbclient

import "fmt"

func (c *Client) GetDocumentation(uid string) (*Documentation, error) {
	url := fmt.Sprintf("%s/documentation/%s", c.HostURL, uid)

	resp, err := executeGet(c.createRestClient(), c.Token, url, Documentation{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*Documentation), nil
}

func (c *Client) UpdateDocumentation(documentation Documentation) (*Documentation, error) {
	url := fmt.Sprintf("%s/documentation/%s", c.HostURL, documentation.Uid)
	rest := c.createRestClient()

	current, err := executeGet(rest, c.Token, url, Documentation{})

	if err != nil {
		return nil, err
	}

	etag := current.Header().Get(etagHeader)

	updated, err := executePut(rest, c.Token, url, documentation, Documentation{}, etag)

	if err != nil {
		return nil, err
	}

	return updated.Result().(*Documentation), nil
}
