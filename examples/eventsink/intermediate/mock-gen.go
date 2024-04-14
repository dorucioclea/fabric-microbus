/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

package intermediate

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/microbus-io/fabric/connector"
	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/sub"

	"github.com/microbus-io/fabric/examples/eventsink/eventsinkapi"
	
	eventsourceapi1 "github.com/microbus-io/fabric/examples/eventsource/eventsourceapi"
	eventsourceapi2 "github.com/microbus-io/fabric/examples/eventsource/eventsourceapi"
)

var (
	_ context.Context
	_ *json.Decoder
	_ *http.Request
	_ time.Duration
	_ *errors.TracedError
	_ *httpx.ResponseRecorder
	_ sub.Option
	_ eventsinkapi.Client
)

// Mock is a mockable version of the eventsink.example microservice,
// allowing functions, sinks and web handlers to be mocked.
type Mock struct {
	*connector.Connector
	MockRegistered func(ctx context.Context) (emails []string, err error)
	MockOnAllowRegister func(ctx context.Context, email string) (allow bool, err error)
	MockOnRegistered func(ctx context.Context, email string) (err error)
}

// NewMock creates a new mockable version of the microservice.
func NewMock(version int) *Mock {
	svc := &Mock{
		Connector: connector.New("eventsink.example"),
	}
	svc.SetVersion(version)
	svc.SetDescription(`The event sink microservice handles events that are fired by the event source microservice.`)
	svc.SetOnStartup(svc.doOnStartup)
	
	// Functions
	svc.Subscribe(`:443/registered`, svc.doRegistered)
	
	// Sinks
	eventsourceapi1.NewHook(svc).OnAllowRegister(svc.doOnAllowRegister)
	eventsourceapi2.NewHook(svc).OnRegistered(svc.doOnRegistered)

	return svc
}

// doOnStartup makes sure that the mock is not executed in a non-dev environment.
func (svc *Mock) doOnStartup(ctx context.Context) (err error) {
	if svc.Deployment() != connector.LOCAL && svc.Deployment() != connector.TESTINGAPP {
		return errors.Newf("mocking disallowed in '%s' deployment", svc.Deployment())
	}
	return nil
}

// doRegistered handles marshaling for the Registered function.
func (svc *Mock) doRegistered(w http.ResponseWriter, r *http.Request) error {
	if svc.MockRegistered == nil {
		return errors.New("mocked endpoint 'Registered' not implemented")
	}
	var i eventsinkapi.RegisteredIn
	var o eventsinkapi.RegisteredOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	o.Emails, err = svc.MockRegistered(
		r.Context(),
	)
	if err != nil {
		return errors.Trace(err)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(o)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// doOnAllowRegister handles marshaling for the OnAllowRegister event sink.
func (svc *Mock) doOnAllowRegister(ctx context.Context, email string) (allow bool, err error) {
	if svc.MockOnAllowRegister == nil {
		err = errors.New("mocked endpoint 'OnAllowRegister' not implemented")
		return
	}
	allow, err = svc.MockOnAllowRegister(ctx, email)
	err = errors.Trace(err)
	return
}

// doOnRegistered handles marshaling for the OnRegistered event sink.
func (svc *Mock) doOnRegistered(ctx context.Context, email string) (err error) {
	if svc.MockOnRegistered == nil {
		err = errors.New("mocked endpoint 'OnRegistered' not implemented")
		return
	}
	err = svc.MockOnRegistered(ctx, email)
	err = errors.Trace(err)
	return
}
