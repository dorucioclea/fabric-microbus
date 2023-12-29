/*
Copyright (c) 2023 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package configurator implements the configurator.sys microservice.

The Configurator is a system microservice that centralizes the dissemination of configuration values to other microservices.
*/
package configurator

import (
	"context"
	"net/http"
	"time"

	"github.com/microbus-io/fabric/connector"
	"github.com/microbus-io/fabric/errors"

	"github.com/microbus-io/fabric/services/configurator/intermediate"
	"github.com/microbus-io/fabric/services/configurator/configuratorapi"
)

var (
	_ context.Context
	_ *http.Request
	_ time.Duration
	_ connector.Service
	_ *errors.TracedError
	_ *configuratorapi.Client
)

// HostName is the default host name of the microservice: configurator.sys.
const HostName = "configurator.sys"

// EndpointURLs contains the fully-qualified URLs to the microservice's endpoints.
var EndpointURLs = configuratorapi.EndpointURLs

// NewService creates a new configurator.sys microservice.
func NewService() connector.Service {
	s := &Service{}
	s.Intermediate = intermediate.NewService(s, Version)
	return s
}

// Mock is a mockable version of the configurator.sys microservice,
// allowing functions, sinks and web handlers to be mocked.
type Mock = intermediate.Mock

// New creates a new mockable version of the microservice.
func NewMock() *Mock {
	return intermediate.NewMock(Version)
}

// Config initializers
var (
	_ int
)
