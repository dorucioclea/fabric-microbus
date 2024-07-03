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
Package intermediate serves as the foundation of the messaging.example microservice.

The Messaging microservice demonstrates service-to-service communication patterns.
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
	"github.com/microbus-io/fabric/service"
	"github.com/microbus-io/fabric/sub"

	"gopkg.in/yaml.v3"

	"github.com/microbus-io/fabric/examples/messaging/resources"
	"github.com/microbus-io/fabric/examples/messaging/messagingapi"
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
	_ service.Service
	_ sub.Option
	_ yaml.Encoder
	_ messagingapi.Client
)

// ToDo defines the interface that the microservice must implement.
// The intermediate delegates handling to this interface.
type ToDo interface {
	OnStartup(ctx context.Context) (err error)
	OnShutdown(ctx context.Context) (err error)
	Home(w http.ResponseWriter, r *http.Request) (err error)
	NoQueue(w http.ResponseWriter, r *http.Request) (err error)
	DefaultQueue(w http.ResponseWriter, r *http.Request) (err error)
	CacheLoad(w http.ResponseWriter, r *http.Request) (err error)
	CacheStore(w http.ResponseWriter, r *http.Request) (err error)
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
		Connector: connector.New("messaging.example"),
		impl: impl,
	}
	svc.SetVersion(version)
	svc.SetDescription(`The Messaging microservice demonstrates service-to-service communication patterns.`)
	
	// Lifecycle
	svc.SetOnStartup(svc.impl.OnStartup)
	svc.SetOnShutdown(svc.impl.OnShutdown)

	// OpenAPI
	svc.Subscribe("GET", `:0/openapi.json`, svc.doOpenAPI)

	// Webs
	svc.Subscribe(`ANY`, `:443/home`, svc.impl.Home)
	svc.Subscribe(`ANY`, `:443/no-queue`, svc.impl.NoQueue, sub.NoQueue())
	svc.Subscribe(`ANY`, `:443/default-queue`, svc.impl.DefaultQueue)
	svc.Subscribe(`ANY`, `:443/cache-load`, svc.impl.CacheLoad)
	svc.Subscribe(`ANY`, `:443/cache-store`, svc.impl.CacheStore)

	// Resources file system
	svc.SetResFS(resources.FS)

	return svc
}

// doOpenAPI renders the OpenAPI document of the microservice.
func (svc *Intermediate) doOpenAPI(w http.ResponseWriter, r *http.Request) error {
	oapiSvc := openapi.Service{
		ServiceName: svc.Hostname(),
		Description: svc.Description(),
		Version:     svc.Version(),
		Endpoints:   []*openapi.Endpoint{},
		RemoteURI:   frame.Of(r).XForwardedFullURL(),
	}
	if r.URL.Port() == "443" || "443" == "0" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `web`,
			Name:        `Home`,
			Method:      `ANY`,
			Path:        `:443/home`,
			Summary:     `Home()`,
			Description: `Home demonstrates making requests using multicast and unicast request/response patterns.`,
			InputArgs: struct {
			}{},
			OutputArgs: struct {
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "0" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `web`,
			Name:        `NoQueue`,
			Method:      `ANY`,
			Path:        `:443/no-queue`,
			Summary:     `NoQueue()`,
			Description: `NoQueue demonstrates how the NoQueue subscription option is used to create
a multicast request/response communication pattern.
All instances of this microservice will respond to each request.`,
			InputArgs: struct {
			}{},
			OutputArgs: struct {
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "0" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `web`,
			Name:        `DefaultQueue`,
			Method:      `ANY`,
			Path:        `:443/default-queue`,
			Summary:     `DefaultQueue()`,
			Description: `DefaultQueue demonstrates how the DefaultQueue subscription option is used to create
a unicast request/response communication pattern.
Only one of the instances of this microservice will respond to each request.`,
			InputArgs: struct {
			}{},
			OutputArgs: struct {
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "0" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `web`,
			Name:        `CacheLoad`,
			Method:      `ANY`,
			Path:        `:443/cache-load`,
			Summary:     `CacheLoad()`,
			Description: `CacheLoad looks up an element in the distributed cache of the microservice.`,
			InputArgs: struct {
			}{},
			OutputArgs: struct {
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "0" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `web`,
			Name:        `CacheStore`,
			Method:      `ANY`,
			Path:        `:443/cache-store`,
			Summary:     `CacheStore()`,
			Description: `CacheStore stores an element in the distributed cache of the microservice.`,
			InputArgs: struct {
			}{},
			OutputArgs: struct {
			}{},
		})
	}

	if len(oapiSvc.Endpoints) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(&oapiSvc)
	return errors.Trace(err)
}

// doOnConfigChanged is called when the config of the microservice changes.
func (svc *Intermediate) doOnConfigChanged(ctx context.Context, changed func(string) bool) (err error) {
	return nil
}
