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
Package configurator implements the configurator.sys microservice.

The Configurator is a core microservice that centralizes the dissemination of configuration values to other microservices.
*/
package configurator

import (
	"context"
	"net/http"
	"time"

	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/service"

	"github.com/microbus-io/fabric/coreservices/configurator/intermediate"
	"github.com/microbus-io/fabric/coreservices/configurator/configuratorapi"
)

var (
	_ context.Context
	_ *http.Request
	_ time.Duration
	_ service.Service
	_ *errors.TracedError
	_ *configuratorapi.Client
)

// Hostname is the default hostname of the microservice: configurator.sys.
const Hostname = "configurator.sys"

// NewService creates a new configurator.sys microservice.
func NewService() *Service {
	s := &Service{}
	s.Intermediate = intermediate.NewService(s, Version)
	return s
}

// Mock is a mockable version of the configurator.sys microservice, allowing functions, event sinks and web handlers to be mocked.
type Mock = intermediate.Mock

// New creates a new mockable version of the microservice.
func NewMock() *Mock {
	return intermediate.NewMock()
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
