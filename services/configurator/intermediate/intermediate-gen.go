/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package intermediate serves as the foundation of the configurator.sys microservice.

The Configurator is a system microservice that centralizes the dissemination of configuration values to other microservices.
*/
package intermediate

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/microbus-io/fabric/cfg"
	"github.com/microbus-io/fabric/connector"
	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/frame"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/log"
	"github.com/microbus-io/fabric/openapi"
	"github.com/microbus-io/fabric/shardedsql"
	"github.com/microbus-io/fabric/sub"

	"gopkg.in/yaml.v3"

	"github.com/microbus-io/fabric/services/configurator/resources"
	"github.com/microbus-io/fabric/services/configurator/configuratorapi"
)

var (
	_ context.Context
	_ *embed.FS
	_ *json.Decoder
	_ fmt.Stringer
	_ *http.Request
	_ filepath.WalkFunc
	_ strconv.NumError
	_ strings.Reader
	_ time.Duration
	_ cfg.Option
	_ *errors.TracedError
	_ frame.Frame
	_ *httpx.ResponseRecorder
	_ *log.Field
	_ *openapi.Service
	_ *shardedsql.DB
	_ sub.Option
	_ yaml.Encoder
	_ configuratorapi.Client
)

// ToDo defines the interface that the microservice must implement.
// The intermediate delegates handling to this interface.
type ToDo interface {
	OnStartup(ctx context.Context) (err error)
	OnShutdown(ctx context.Context) (err error)
	Values(ctx context.Context, names []string) (values map[string]string, err error)
	Refresh(ctx context.Context) (err error)
	Sync(ctx context.Context, timestamp time.Time, values map[string]map[string]string) (err error)
	PeriodicRefresh(ctx context.Context) (err error)
}

// Intermediate extends and customizes the generic base connector.
// Code generated microservices then extend the intermediate.
type Intermediate struct {
	*connector.Connector
	impl ToDo
}

// NewService creates a new intermediate service.
func NewService(impl ToDo, version int) *Intermediate {
	svc := &Intermediate{
		Connector: connector.New("configurator.sys"),
		impl: impl,
	}
	svc.SetVersion(version)
	svc.SetDescription(`The Configurator is a system microservice that centralizes the dissemination of configuration values to other microservices.`)

	// Lifecycle
	svc.SetOnStartup(svc.impl.OnStartup)
	svc.SetOnShutdown(svc.impl.OnShutdown)
	
	// OpenAPI
	svc.Subscribe(`:443/openapi.yaml`, svc.doOpenAPI)	

	// Functions
	svc.Subscribe(`:443/values`, svc.doValues)
	svc.Subscribe(`:443/refresh`, svc.doRefresh)
	svc.Subscribe(`:443/sync`, svc.doSync, sub.NoQueue())
	
	// Tickers
	intervalPeriodicRefresh, _ := time.ParseDuration("20m0s")
	svc.StartTicker("PeriodicRefresh", intervalPeriodicRefresh, svc.impl.PeriodicRefresh)

	// Resources file system
	svc.SetResFS(resources.FS)

	return svc
}

// doOpenAPI renders the OpenAPI document of the microservice.
func (svc *Intermediate) doOpenAPI(w http.ResponseWriter, r *http.Request) error {
	oapiSvc := openapi.Service{
		ServiceName: svc.HostName(),
		Description: svc.Description(),
		Version:     svc.Version(),
		Endpoints:   []*openapi.Endpoint{},
		RemoteURI:   frame.Of(r).XForwardedFullURL(),
	}

	if len(oapiSvc.Endpoints) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	err := yaml.NewEncoder(w).Encode(&oapiSvc)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// doOnConfigChanged is called when the config of the microservice changes.
func (svc *Intermediate) doOnConfigChanged(ctx context.Context, changed func(string) bool) (err error) {
	return nil
}

// doValues handles marshaling for the Values function.
func (svc *Intermediate) doValues(w http.ResponseWriter, r *http.Request) error {
	var i configuratorapi.ValuesIn
	var o configuratorapi.ValuesOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.Values, err = svc.impl.Values(
		r.Context(),
		i.Names,
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

// doRefresh handles marshaling for the Refresh function.
func (svc *Intermediate) doRefresh(w http.ResponseWriter, r *http.Request) error {
	var i configuratorapi.RefreshIn
	var o configuratorapi.RefreshOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	err = svc.impl.Refresh(
		r.Context(),
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

// doSync handles marshaling for the Sync function.
func (svc *Intermediate) doSync(w http.ResponseWriter, r *http.Request) error {
	var i configuratorapi.SyncIn
	var o configuratorapi.SyncOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	err = svc.impl.Sync(
		r.Context(),
		i.Timestamp,
		i.Values,
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
