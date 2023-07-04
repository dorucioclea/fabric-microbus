/*
Copyright (c) 2023 Microbus LLC and various contributors

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

	"github.com/microbus-io/fabric/services/configurator/configuratorapi"
)

var (
	_ context.Context
	_ *json.Decoder
	_ *http.Request
	_ time.Duration
	_ *errors.TracedError
	_ *httpx.ResponseRecorder
	_ sub.Option
	_ configuratorapi.Client
)

// Mock is a mockable version of the configurator.sys microservice,
// allowing functions, sinks and web handlers to be mocked.
type Mock struct {
	*connector.Connector
	MockValues func(ctx context.Context, names []string) (values map[string]string, err error)
	MockRefresh func(ctx context.Context) (err error)
	MockSync func(ctx context.Context, timestamp time.Time, values map[string]map[string]string) (err error)
}

// NewMock creates a new mockable version of the microservice.
func NewMock(version int) *Mock {
	svc := &Mock{
		Connector: connector.New("configurator.sys"),
	}
	svc.SetVersion(version)
	svc.SetDescription(`The Configurator is a system microservice that centralizes the dissemination of configuration values to other microservices.`)
	svc.SetOnStartup(svc.doOnStartup)
	
	// Functions
	svc.Subscribe(`:443/values`, svc.doValues)
	svc.Subscribe(`:443/refresh`, svc.doRefresh)
	svc.Subscribe(`:443/sync`, svc.doSync, sub.NoQueue())

	return svc
}

// doOnStartup makes sure that the mock is not executed in a non-dev environment.
func (svc *Mock) doOnStartup(ctx context.Context) (err error) {
	if svc.Deployment() != connector.LOCAL && svc.Deployment() != connector.TESTINGAPP {
		return errors.Newf("mocking disallowed in '%s' deployment", svc.Deployment())
	}
	return nil
}

// doValues handles marshaling for the Values function.
func (svc *Mock) doValues(w http.ResponseWriter, r *http.Request) error {
	if svc.MockValues == nil {
		return errors.New("mocked endpoint 'Values' not implemented")
	}
	var i configuratorapi.ValuesIn
	var o configuratorapi.ValuesOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	o.Values, err = svc.MockValues(
		r.Context(),
		i.Names,
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

// doRefresh handles marshaling for the Refresh function.
func (svc *Mock) doRefresh(w http.ResponseWriter, r *http.Request) error {
	if svc.MockRefresh == nil {
		return errors.New("mocked endpoint 'Refresh' not implemented")
	}
	var i configuratorapi.RefreshIn
	var o configuratorapi.RefreshOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	err = svc.MockRefresh(
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

// doSync handles marshaling for the Sync function.
func (svc *Mock) doSync(w http.ResponseWriter, r *http.Request) error {
	if svc.MockSync == nil {
		return errors.New("mocked endpoint 'Sync' not implemented")
	}
	var i configuratorapi.SyncIn
	var o configuratorapi.SyncOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	err = svc.MockSync(
		r.Context(),
		i.Timestamp,
		i.Values,
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
