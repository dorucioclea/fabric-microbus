/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package openapiportalapi implements the public API of the openapiportal.core microservice,
including clients and data structures.

The OpenAPI microservice lists links to the OpenAPI endpoint of all microservices that provide one
on the requested port.
*/
package openapiportalapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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
	_ io.Reader
	_ *http.Request
	_ *url.URL
	_ strings.Reader
	_ time.Duration
	_ *errors.TracedError
	_ *httpx.BodyReader
	_ pub.Option
	_ sub.Option
)

// Hostname is the default hostname of the microservice: openapiportal.core.
const Hostname = "openapiportal.core"

// Fully-qualified URLs of the microservice's endpoints.
var (
	URLOfList = httpx.JoinHostAndPath(Hostname, `//openapi:0`)
)

// Client is an interface to calling the endpoints of the openapiportal.core microservice.
// This simple version is for unicast calls.
type Client struct {
	svc  service.Publisher
	host string
}

// NewClient creates a new unicast client to the openapiportal.core microservice.
func NewClient(caller service.Publisher) *Client {
	return &Client{
		svc:  caller,
		host: "openapiportal.core",
	}
}

// ForHost replaces the default hostname of this client.
func (_c *Client) ForHost(host string) *Client {
	_c.host = host
	return _c
}

// MulticastClient is an interface to calling the endpoints of the openapiportal.core microservice.
// This advanced version is for multicast calls.
type MulticastClient struct {
	svc  service.Publisher
	host string
}

// NewMulticastClient creates a new multicast client to the openapiportal.core microservice.
func NewMulticastClient(caller service.Publisher) *MulticastClient {
	return &MulticastClient{
		svc:  caller,
		host: "openapiportal.core",
	}
}

// ForHost replaces the default hostname of this client.
func (_c *MulticastClient) ForHost(host string) *MulticastClient {
	_c.host = host
	return _c
}

// errChan returns a response channel with a single error response.
func (_c *MulticastClient) errChan(err error) <-chan *pub.Response {
	ch := make(chan *pub.Response, 1)
	ch <- pub.NewErrorResponse(err)
	close(ch)
	return ch
}

/*
List_Get performs a GET request to the List endpoint.

List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.

If a URL is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
*/
func (_c *Client) List_Get(ctx context.Context, url string) (res *http.Response, err error) {
	url, err = httpx.ResolveURL(URLOfList, url)
	if err != nil {
		return nil, errors.Trace(err)
	}
	url, err = httpx.FillPathArguments(url)
	if err != nil {
		return nil, errors.Trace(err)
	}
	res, err = _c.svc.Request(ctx, pub.Method("GET"), pub.URL(url))
	if err != nil {
		return nil, err // No trace
	}
	return res, err
}

/*
List_Get performs a GET request to the List endpoint.

List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.

If a URL is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
*/
func (_c *MulticastClient) List_Get(ctx context.Context, url string) <-chan *pub.Response {
	var err error
	url, err = httpx.ResolveURL(URLOfList, url)
	if err != nil {
		return _c.errChan(errors.Trace(err))
	}
	url, err = httpx.FillPathArguments(url)
	if err != nil {
		return _c.errChan(errors.Trace(err))
	}
	return _c.svc.Publish(ctx, pub.Method("GET"), pub.URL(url))
}

/*
List_Post performs a POST request to the List endpoint.

List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.

If a URL is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
If the body if of type io.Reader, []byte or string, it is serialized in binary form.
If it is of type url.Values, it is serialized as form data. All other types are serialized as JSON.
If a content type is not explicitly provided, an attempt will be made to derive it from the body.
*/
func (_c *Client) List_Post(ctx context.Context, url string, contentType string, body any) (res *http.Response, err error) {
	url, err = httpx.ResolveURL(URLOfList, url)
	if err != nil {
		return nil, errors.Trace(err)
	}
	url, err = httpx.FillPathArguments(url)
	if err != nil {
		return nil, errors.Trace(err)
	}
	res, err = _c.svc.Request(ctx, pub.Method("POST"), pub.URL(url), pub.ContentType(contentType), pub.Body(body))
	if err != nil {
		return nil, err // No trace
	}
	return res, err
}

/*
List_Post performs a POST request to the List endpoint.

List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.

If a URL is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
If the body if of type io.Reader, []byte or string, it is serialized in binary form.
If it is of type url.Values, it is serialized as form data. All other types are serialized as JSON.
If a content type is not explicitly provided, an attempt will be made to derive it from the body.
*/
func (_c *MulticastClient) List_Post(ctx context.Context, url string, contentType string, body any) <-chan *pub.Response {
	var err error
	url, err = httpx.ResolveURL(URLOfList, url)
	if err != nil {
		return _c.errChan(errors.Trace(err))
	}
	url, err = httpx.FillPathArguments(url)
	if err != nil {
		return _c.errChan(errors.Trace(err))
	}
	return _c.svc.Publish(ctx, pub.Method("POST"), pub.URL(url), pub.ContentType(contentType), pub.Body(body))
}

/*
List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.

If a request is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
*/
func (_c *Client) List(r *http.Request) (res *http.Response, err error) {
	if r == nil {
		r, err = http.NewRequest(`GET`, "", nil)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	url, err := httpx.ResolveURL(URLOfList, r.URL.String())
	if err != nil {
		return nil, errors.Trace(err)
	}
	url, err = httpx.FillPathArguments(url)
	if err != nil {
		return nil, errors.Trace(err)
	}
	res, err = _c.svc.Request(r.Context(), pub.Method(r.Method), pub.URL(url), pub.CopyHeaders(r.Header), pub.Body(r.Body))
	if err != nil {
		return nil, err // No trace
	}
	return res, err
}

/*
List displays links to the OpenAPI endpoint of all microservices that provide one on the request's port.

If a request is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
*/
func (_c *MulticastClient) List(ctx context.Context, r *http.Request) <-chan *pub.Response {
	var err error
	if r == nil {
		r, err = http.NewRequest(`GET`, "", nil)
		if err != nil {
			return _c.errChan(errors.Trace(err))
		}
	}
	url, err := httpx.ResolveURL(URLOfList, r.URL.String())
	if err != nil {
		return _c.errChan(errors.Trace(err))
	}
	url, err = httpx.FillPathArguments(url)
	if err != nil {
		return _c.errChan(errors.Trace(err))
	}
	return _c.svc.Publish(ctx, pub.Method(r.Method), pub.URL(url), pub.CopyHeaders(r.Header), pub.Body(r.Body))
}
