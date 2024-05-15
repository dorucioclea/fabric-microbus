/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package openapiportalapi implements the public API of the openapiportal.sys microservice,
including clients and data structures.

The OpenAPI microservice lists links to the OpenAPI endpoint of all microservices that provide one
on the requested port.
*/
package openapiportalapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/pub"
	"github.com/microbus-io/fabric/service"
	"github.com/microbus-io/fabric/sub"
)

var (
	_ context.Context
	_ *json.Decoder
	_ *http.Request
	_ strings.Reader
	_ time.Duration
	_ *errors.TracedError
	_ *httpx.BodyReader
	_ pub.Option
	_ sub.Option
)

// HostName is the default host name of the microservice: openapiportal.sys.
const HostName = "openapiportal.sys"

// Fully-qualified URLs of the microservice's endpoints.
var (
	URLOfList = httpx.JoinHostAndPath(HostName, "//openapi:*")
)

// Client is an interface to calling the endpoints of the openapiportal.sys microservice.
// This simple version is for unicast calls.
type Client struct {
	svc  service.Publisher
	host string
}

// NewClient creates a new unicast client to the openapiportal.sys microservice.
func NewClient(caller service.Publisher) *Client {
	return &Client{
		svc:  caller,
		host: "openapiportal.sys",
	}
}

// ForHost replaces the default host name of this client.
func (_c *Client) ForHost(host string) *Client {
	_c.host = host
	return _c
}

// MulticastClient is an interface to calling the endpoints of the openapiportal.sys microservice.
// This advanced version is for multicast calls.
type MulticastClient struct {
	svc  service.Publisher
	host string
}

// NewMulticastClient creates a new multicast client to the openapiportal.sys microservice.
func NewMulticastClient(caller service.Publisher) *MulticastClient {
	return &MulticastClient{
		svc:  caller,
		host: "openapiportal.sys",
	}
}

// ForHost replaces the default host name of this client.
func (_c *MulticastClient) ForHost(host string) *MulticastClient {
	_c.host = host
	return _c
}

/*
List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.
*/
func (_c *Client) List(ctx context.Context, options ...pub.Option) (res *http.Response, err error) {
	method := `*`
	if method == "*" {
		method = "GET"
	}
	opts := []pub.Option{
		pub.Method(method),
		pub.URL(httpx.JoinHostAndPath(_c.host, `//openapi:*`)),
	}
	opts = append(opts, options...)
	res, err = _c.svc.Request(ctx, opts...)
	if err != nil {
		return nil, err // No trace
	}
	return res, err
}

/*
List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.
*/
func (_c *MulticastClient) List(ctx context.Context, options ...pub.Option) <-chan *pub.Response {
	method := `*`
	if method == "*" {
		method = "GET"
	}
	opts := []pub.Option{
		pub.Method(method),
		pub.URL(httpx.JoinHostAndPath(_c.host, `//openapi:*`)),
	}
	opts = append(opts, options...)
	return _c.svc.Publish(ctx, opts...)
}