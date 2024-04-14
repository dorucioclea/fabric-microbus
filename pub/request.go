/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package pub

import (
	"io"
	"net/http"
	"strings"

	"github.com/microbus-io/fabric/errors"
)

// Request is used to construct an HTTP request that can be sent over the bus.
// Although technically public, it is used internally and should not be constructed by microservices directly.
type Request struct {
	Method    string
	URL       string
	Header    http.Header
	Body      io.Reader
	Multicast bool

	queryArgs string
}

// NewRequest constructs a new request from the provided options.
func NewRequest(options ...Option) (*Request, error) {
	req := &Request{
		Method:    "POST",
		Header:    make(http.Header),
		Multicast: true,
	}
	err := req.Apply(options...)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return req, nil
}

// Apply the provided options to the request, in order.
func (req *Request) Apply(options ...Option) error {
	for _, opt := range options {
		err := opt(req)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// Canonical returns the fully-qualified canonical path of the request, without the query arguments.
func (req *Request) Canonical() string {
	qm := strings.Index(req.URL, "?")
	if qm >= 0 {
		return req.URL[:qm]
	}
	return req.URL
}
