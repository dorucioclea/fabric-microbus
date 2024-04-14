/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package inboxapi implements the public API of the inbox.sys microservice,
including clients and data structures.

Inbox listens for incoming emails and fires appropriate events.
*/
package inboxapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/pub"
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

// HostName is the default host name of the microservice: inbox.sys.
const HostName = "inbox.sys"

// Fully-qualified URLs of the microservice's endpoints.
var (
)

// Service is an interface abstraction of a microservice used by the client.
// The connector implements this interface.
type Service interface {
	Request(ctx context.Context, options ...pub.Option) (*http.Response, error)
	Publish(ctx context.Context, options ...pub.Option) <-chan *pub.Response
	Subscribe(path string, handler sub.HTTPHandler, options ...sub.Option) error
	Unsubscribe(path string) error
}

// Client is an interface to calling the endpoints of the inbox.sys microservice.
// This simple version is for unicast calls.
type Client struct {
	svc  Service
	host string
}

// NewClient creates a new unicast client to the inbox.sys microservice.
func NewClient(caller Service) *Client {
	return &Client{
		svc:  caller,
		host: "inbox.sys",
	}
}

// ForHost replaces the default host name of this client.
func (_c *Client) ForHost(host string) *Client {
	_c.host = host
	return _c
}

// MulticastClient is an interface to calling the endpoints of the inbox.sys microservice.
// This advanced version is for multicast calls.
type MulticastClient struct {
	svc  Service
	host string
}

// NewMulticastClient creates a new multicast client to the inbox.sys microservice.
func NewMulticastClient(caller Service) *MulticastClient {
	return &MulticastClient{
		svc:  caller,
		host: "inbox.sys",
	}
}

// ForHost replaces the default host name of this client.
func (_c *MulticastClient) ForHost(host string) *MulticastClient {
	_c.host = host
	return _c
}

// MulticastTrigger is an interface to trigger the events of the inbox.sys microservice.
type MulticastTrigger struct {
	svc  Service
	host string
}

// NewMulticastTrigger creates a new multicast trigger of the inbox.sys microservice.
func NewMulticastTrigger(caller Service) *MulticastTrigger {
	return &MulticastTrigger{
		svc:  caller,
		host: "inbox.sys",
	}
}

// ForHost replaces the default host name of this trigger.
func (_c *MulticastTrigger) ForHost(host string) *MulticastTrigger {
	_c.host = host
	return _c
}

// Hook assists in the subscription to the events of the inbox.sys microservice.
type Hook struct {
	svc  Service
	host string
}

// NewHook creates a new hook to the events of the inbox.sys microservice.
func NewHook(listener Service) *Hook {
	return &Hook{
		svc:  listener,
		host: "inbox.sys",
	}
}

// ForHost replaces the default host name of this hook.
func (_c *Hook) ForHost(host string) *Hook {
	_c.host = host
	return _c
}

// OnInboxSaveMailIn are the input arguments of OnInboxSaveMail.
type OnInboxSaveMailIn struct {
	MailMessage *Email `json:"mailMessage"`
}

// OnInboxSaveMailOut are the return values of OnInboxSaveMail.
type OnInboxSaveMailOut struct {
}

// OnInboxSaveMailResponse is the response to OnInboxSaveMail.
type OnInboxSaveMailResponse struct {
	data OnInboxSaveMailOut
	HTTPResponse *http.Response
	err error
}

// Get retrieves the return values.
func (_out *OnInboxSaveMailResponse) Get() (err error) {
	err = _out.err
	return
}

/*
OnInboxSaveMail is triggered when a new email message is received.
*/
func (_c *MulticastTrigger) OnInboxSaveMail(ctx context.Context, mailMessage *Email, _options ...pub.Option) <-chan *OnInboxSaveMailResponse {
	_in := OnInboxSaveMailIn{
		mailMessage,
	}
	_opts := []pub.Option{
		pub.Method("POST"),
		pub.URL(httpx.JoinHostAndPath(_c.host, `:417/on-inbox-save-mail`)),
		pub.Body(_in),
	}
	_opts = append(_opts, _options...)
	_ch := _c.svc.Publish(ctx, _opts...)

	_res := make(chan *OnInboxSaveMailResponse, cap(_ch))
	go func() {
		for _i := range _ch {
			var _r OnInboxSaveMailResponse
			_httpRes, _err := _i.Get()
			_r.HTTPResponse = _httpRes
			if _err != nil {
				_r.err = _err // No trace
			} else {
				_err = json.NewDecoder(_httpRes.Body).Decode(&(_r.data))
				if _err != nil {
					_r.err = errors.Trace(_err)
				}
			}
			_res <- &_r
		}
		close(_res)
	}()
	return _res
}

/*
OnInboxSaveMail is triggered when a new email message is received.
*/
func (_c *Hook) OnInboxSaveMail(handler func(ctx context.Context, mailMessage *Email) (err error), options ...sub.Option) error {
	doOnInboxSaveMail := func(w http.ResponseWriter, r *http.Request) error {
		var i OnInboxSaveMailIn
		var o OnInboxSaveMailOut
		err := httpx.ParseRequestData(r, &i)
		if err!=nil {
			return errors.Trace(err)
		}
		err = handler(
			r.Context(),
			i.MailMessage,
		)
		if err != nil {
			return err // No trace
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(o)
		if err != nil {
			return errors.Trace(err)
		}
		return nil
	}
	path := httpx.JoinHostAndPath(_c.host, `:417/on-inbox-save-mail`)
	if handler == nil {
		return _c.svc.Unsubscribe(path)
	}
	return _c.svc.Subscribe(path, doOnInboxSaveMail, options...)
}
