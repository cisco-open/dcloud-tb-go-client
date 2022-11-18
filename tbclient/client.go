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
	client         *Client
	collectionPath string
	singlePath     string
}

func (s *service[R, RC]) getAll() ([]R, error) {

	rest := resty.New()

	resp, err := rest.R().
		SetResult(collection[RC]{}).
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Get(fmt.Sprintf(s.collectionPath, s.client.HostURL))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}

	return resp.Result().(*collection[RC]).Embedded.getData(), nil
}

func (s *service[R, RC]) getOne(uid string) (*R, error) {

	rest := resty.New()

	resp, err := rest.R().
		SetResult(new(R)).
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Get(fmt.Sprintf(s.singlePath, s.client.HostURL, uid))

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

	resp, err := rest.R().
		SetBody(resource).
		SetResult(new(R)).
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Post(fmt.Sprintf(s.collectionPath, s.client.HostURL))

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

	current, err := rest.R().
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Get(fmt.Sprintf(s.singlePath, s.client.HostURL, uid))

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
		Put(fmt.Sprintf(s.singlePath, s.client.HostURL, uid))

	if err != nil {
		return nil, err
	}

	if updated.IsError() {
		return nil, generateError(*updated)
	}

	return updated.Result().(*R), nil
}

func (s service[R, RC]) delete(uid string) error {

	rest := resty.New()

	resp, err := rest.R().
		SetError([]vndError{}).
		SetAuthToken(s.client.Token).
		Delete(fmt.Sprintf(s.singlePath, s.client.HostURL, uid))

	if err != nil {
		return err
	}

	if resp.IsError() && resp.StatusCode() != 404 {
		return generateError(*resp)
	}

	return nil
}
