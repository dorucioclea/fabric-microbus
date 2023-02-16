/*
Copyright 2023 Microbus LLC and various contributors

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

	"github.com/microbus-io/fabric/cb"
	"github.com/microbus-io/fabric/cfg"
	"github.com/microbus-io/fabric/connector"
	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/log"
	"github.com/microbus-io/fabric/shardedsql"
	"github.com/microbus-io/fabric/sub"
	"github.com/microbus-io/fabric/utils"

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
	_ cb.Option
	_ cfg.Option
	_ *errors.TracedError
	_ *httpx.ResponseRecorder
	_ *log.Field
	_ *shardedsql.DB
	_ sub.Option
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

	// Functions
	svc.Subscribe(`:443/values`, svc.doValues)
	svc.Subscribe(`:443/refresh`, svc.doRefresh)
	svc.Subscribe(`:443/sync`, svc.doSync, sub.NoQueue())
	
	// Tickers
	intervalPeriodicRefresh, _ := time.ParseDuration("20m0s")
	svc.StartTicker("PeriodicRefresh", intervalPeriodicRefresh, svc.impl.PeriodicRefresh)

	return svc
}

// Resources is the in-memory file system of the embedded resources.
func (svc *Intermediate) Resources() utils.ResourceLoader {
	return utils.ResourceLoader{FS: resources.FS}
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
	if err!=nil {
		return errors.Trace(err)
	}
	o.Values, err = svc.impl.Values(
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
func (svc *Intermediate) doRefresh(w http.ResponseWriter, r *http.Request) error {
	var i configuratorapi.RefreshIn
	var o configuratorapi.RefreshOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	err = svc.impl.Refresh(
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
func (svc *Intermediate) doSync(w http.ResponseWriter, r *http.Request) error {
	var i configuratorapi.SyncIn
	var o configuratorapi.SyncOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	err = svc.impl.Sync(
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
