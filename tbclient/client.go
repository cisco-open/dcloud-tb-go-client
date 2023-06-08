// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// HostURL Defaulted to local service instance
const HostURL = "http://localhost:8080"
const etagHeader = "ETag"
const ifMatch = "IF-MATCH"

type Client struct {
	HostURL     string
	Token       string
	Debug       bool
	UserAgent   string
	DisableGzip bool
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

type resourceService[R any, RC embeddedData[R]] struct {
	collectionService[R, RC]
}

type collectionService[R any, RC embeddedData[R]] struct {
	client         *Client
	resourcePath   string
	topologyUid    string
	requestHeaders map[string]string
}

func (s *collectionService[R, RC]) getAll() ([]R, error) {

	topologyPath := getTopologyPath(s.topologyUid)

	url := fmt.Sprintf("%s%s%s", s.client.HostURL, topologyPath, s.resourcePath)

	resp, err := executeGet(s.createRestClient(), s.client.Token, url, collection[RC]{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*collection[RC]).Embedded.getData(), nil
}

func (s *resourceService[R, RC]) getOne(uid string) (*R, error) {

	url := singleResourceUrl(s.client.HostURL, s.resourcePath, uid)

	resp, err := executeGet(s.createRestClient(), s.client.Token, url, new(R))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*R), nil
}

func (s *resourceService[R, RC]) create(resource R) (*R, error) {

	url := fmt.Sprintf("%s%s", s.client.HostURL, s.resourcePath)

	resp, err := executePost(s.createRestClient(), s.client.Token, url, resource, new(R))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}

	return resp.Result().(*R), nil
}

func (s *resourceService[R, RC]) update(uid string, resource R) (*R, error) {

	rest := s.createRestClient()

	url := singleResourceUrl(s.client.HostURL, s.resourcePath, uid)

	current, err := executeGet(rest, s.client.Token, url, new(R))

	if err != nil {
		return nil, err
	}

	etag := current.Header().Get(etagHeader)

	updated, err := executePut(rest, s.client.Token, url, resource, new(R), etag)

	if err != nil {
		return nil, err
	}

	return updated.Result().(*R), nil
}

func (s *resourceService[R, RC]) delete(uid string) error {

	url := singleResourceUrl(s.client.HostURL, s.resourcePath, uid)

	err := executeDelete(s.createRestClient(), s.client.Token, url)
	if err != nil {
		return err
	}

	return nil
}

func singleResourceUrl(hostUrl, resourcePath, uid string) string {
	singlePath := fmt.Sprintf("%s/%s", resourcePath, uid)
	url := fmt.Sprintf("%s%s", hostUrl, singlePath)
	return url
}

type singleNestedResourceService[R any] struct {
	client         *Client
	resourcePath   string
	topologyUid    string
	putUsingGetUrl bool
}

func (s *singleNestedResourceService[R]) get() (*R, error) {

	url := fmt.Sprintf("%s/topologies/%s%s", s.client.HostURL, s.topologyUid, s.resourcePath)

	resp, err := executeGet(s.client.createRestClient(), s.client.Token, url, new(R))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*R), nil
}

func (s *singleNestedResourceService[R]) update(uid string, resource R) (*R, error) {

	rest := s.client.createRestClient()
	getUrl := fmt.Sprintf("%s/topologies/%s%s", s.client.HostURL, s.topologyUid, s.resourcePath)

	current, err := executeGet(s.client.createRestClient(), s.client.Token, getUrl, new(R))
	if err != nil {
		return nil, err
	}

	etag := current.Header().Get(etagHeader)

	var putUrl string
	if s.putUsingGetUrl {
		putUrl = getUrl
	} else {
		putUrl = fmt.Sprintf("%s%s/%s", s.client.HostURL, s.resourcePath, uid)
	}

	updated, err := executePut(rest, s.client.Token, putUrl, resource, new(R), etag)
	if err != nil {
		return nil, err
	}

	return updated.Result().(*R), nil
}

func (s *collectionService[R, RC]) createRestClient() *resty.Client {
	rest := s.client.createRestClient()
	rest.SetHeaders(s.requestHeaders)
	return rest
}

func (c *Client) createRestClient() *resty.Client {
	rest := resty.New()
	rest.SetDebug(c.Debug)
	if c.UserAgent != "" {
		rest.SetHeader("User-Agent", c.UserAgent)
	}
	if c.DisableGzip {
		rest.SetHeader("Accept-Encoding", "")
	}
	return rest
}

func getTopologyPath(topologyUid string) string {
	if topologyUid != "" {
		return fmt.Sprintf("/topologies/%s", topologyUid)
	} else {
		return ""
	}
}

func executeGet(rest *resty.Client, authToken string, url string, result interface{}) (*resty.Response, error) {
	resp, err := rest.R().
		SetResult(result).
		SetError([]vndError{}).
		SetAuthToken(authToken).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}
	return resp, nil
}

func executePut(rest *resty.Client, authToken string, url string, body interface{}, result interface{}, etag string) (*resty.Response, error) {
	resp, err := rest.R().
		SetHeader(ifMatch, etag).
		SetBody(body).
		SetResult(result).
		SetError([]vndError{}).
		SetAuthToken(authToken).
		Put(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}
	return resp, nil
}

func executePost(rest *resty.Client, authToken string, url string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := rest.R().
		SetBody(body).
		SetResult(result).
		SetError([]vndError{}).
		SetAuthToken(authToken).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, generateError(*resp)
	}
	return resp, nil
}

func executeDelete(rest *resty.Client, authToken string, url string) error {
	resp, err := rest.R().
		SetError([]vndError{}).
		SetAuthToken(authToken).
		Delete(url)

	if err != nil {
		return err
	}

	if resp.IsError() && resp.StatusCode() != 404 {
		return generateError(*resp)
	}

	return nil
}

var emptyHeaders = map[string]string{}
