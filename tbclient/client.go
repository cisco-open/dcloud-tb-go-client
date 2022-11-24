package tbclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// Default to local servic instance
const HostURL = "http://localhost:8080"
const etag = "ETag"
const ifMatch = "IF-MATCH"

type Client struct {
	HostURL string
	Token   string
}

func NewClient(host, authToken *string) *Client {
	c := Client{
		HostURL: HostURL,
		Token:   *authToken,
	}

	if host != nil {
		c.HostURL = *host
	}

	return &c
}

type service[R any, RC embeddedData[R]] struct {
	readService[R, RC]
}

type readService[R any, RC embeddedData[R]] struct {
	client       *Client
	resourcePath string
	topologyUid  string
}

func (s *readService[R, RC]) getAll() ([]R, error) {

	rest := resty.New()

	topologyPath := getTopologyPath(s.topologyUid)
	url := fmt.Sprintf("%s%s%s", s.client.HostURL, topologyPath, s.resourcePath)

	fmt.Printf("URL: %s\n", url)

	resp, err := rest.R().
		SetResult(collection[RC]{}).
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}
	
	return resp.Result().(*collection[RC]).Embedded.getData(), nil
}

func getTopologyPath(topologyUid string) string {
	if topologyUid != "" {
		return fmt.Sprintf("/topologies/%s", topologyUid)
	} else {
		return ""
	}
}

func (s *service[R, RC]) getOne(uid string) (*R, error) {

	rest := resty.New()

	url := singleResourceUrl(s.client.HostURL, s.resourcePath, uid)
	fmt.Printf("URL: %s\n", url)

	resp, err := rest.R().
		SetResult(new(R)).
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}

	return resp.Result().(*R), nil
}

func (s *service[R, RC]) create(resource R) (*R, error) {

	rest := resty.New()

	url := fmt.Sprintf("%s%s", s.client.HostURL, s.resourcePath)
	fmt.Printf("URL: %s\n", url)

	resp, err := rest.R().
		SetBody(resource).
		SetResult(new(R)).
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}

	return resp.Result().(*R), nil
}

func (s *service[R, RC]) update(uid string, resource R) (*R, error) {

	rest := resty.New()

	url := singleResourceUrl(s.client.HostURL, s.resourcePath, uid)
	fmt.Printf("URL: %s\n", url)

	current, err := rest.R().
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Get(url)

	if err != nil {
		return nil, err
	}

	if current.IsError() {
		return nil, generateError(*current)
	}

	etag := current.Header().Get(etag)

	updated, err := rest.R().
		SetHeader(ifMatch, etag).
		SetBody(resource).
		SetResult(new(R)).
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Put(url)

	if err != nil {
		return nil, err
	}

	if updated.IsError() {
		return nil, generateError(*updated)
	}

	return updated.Result().(*R), nil
}

func (s *service[R, RC]) delete(uid string) error {

	rest := resty.New()

	url := singleResourceUrl(s.client.HostURL, s.resourcePath, uid)
	fmt.Printf("URL: %s\n", url)

	resp, err := rest.R().
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Delete(url)

	if err != nil {
		return err
	}

	if resp.IsError() && resp.StatusCode() != 404 {
		return generateError(*resp)
	}

	return nil
}

func singleResourceUrl(hostUrl, resourcePath, uid string) string {
	singlePath := fmt.Sprintf("%s/%s", resourcePath, uid)
	url := fmt.Sprintf("%s%s", hostUrl, singlePath)
	return url
}
