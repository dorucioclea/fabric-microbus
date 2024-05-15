/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package inbox implements the inbox.sys microservice.

Inbox listens for incoming emails and fires appropriate events.
*/
package inbox

import (
	"context"
	"net/http"
	"time"

	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/service"

	"github.com/microbus-io/fabric/coreservices/inbox/intermediate"
	"github.com/microbus-io/fabric/coreservices/inbox/inboxapi"
)

var (
	_ context.Context
	_ *http.Request
	_ time.Duration
	_ service.Service
	_ *errors.TracedError
	_ *inboxapi.Client
)

// HostName is the default host name of the microservice: inbox.sys.
const HostName = "inbox.sys"

// NewService creates a new inbox.sys microservice.
func NewService() service.Service {
	s := &Service{}
	s.Intermediate = intermediate.NewService(s, Version)
	return s
}

// Mock is a mockable version of the inbox.sys microservice,
// allowing functions, sinks and web handlers to be mocked.
type Mock = intermediate.Mock

// New creates a new mockable version of the microservice.
func NewMock() *Mock {
	return intermediate.NewMock(Version)
}

// Config initializers
var (
	_ int
	// Port initializes the Port config property of the microservice
	Port = intermediate.Port
	// Enabled initializes the Enabled config property of the microservice
	Enabled = intermediate.Enabled
	// MaxSize initializes the MaxSize config property of the microservice
	MaxSize = intermediate.MaxSize
	// MaxClients initializes the MaxClients config property of the microservice
	MaxClients = intermediate.MaxClients
	// Workers initializes the Workers config property of the microservice
	Workers = intermediate.Workers
)