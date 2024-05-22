/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package openapiportal implements the openapiportal.sys microservice.

The OpenAPI microservice lists links to the OpenAPI endpoint of all microservices that provide one
on the requested port.
*/
package openapiportal

import (
	"context"
	"net/http"
	"time"

	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/service"

	"github.com/microbus-io/fabric/coreservices/openapiportal/intermediate"
	"github.com/microbus-io/fabric/coreservices/openapiportal/openapiportalapi"
)

var (
	_ context.Context
	_ *http.Request
	_ time.Duration
	_ service.Service
	_ *errors.TracedError
	_ *openapiportalapi.Client
)

// HostName is the default host name of the microservice: openapiportal.sys.
const HostName = "openapiportal.sys"

// NewService creates a new openapiportal.sys microservice.
func NewService() *Service {
	s := &Service{}
	s.Intermediate = intermediate.NewService(s, Version)
	return s
}

// Mock is a mockable version of the openapiportal.sys microservice,
// allowing functions, sinks and web handlers to be mocked.
type Mock = intermediate.Mock

// New creates a new mockable version of the microservice.
func NewMock() *Mock {
	return intermediate.NewMock(Version)
}

/*
Init enables a single-statement pattern for initializing the microservice.

	svc.Init(func(svc Service) {
		svc.SetGreeting("Hello")
	})
*/
func (svc *Service) Init(initializer func(svc *Service)) *Service {
	initializer(svc)
	return svc
}
