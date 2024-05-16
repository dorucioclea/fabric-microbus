/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package intermediate serves as the foundation of the directory.example microservice.

The directory microservice stores personal records in a SQL database.
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

	"github.com/microbus-io/fabric/examples/directory/resources"
	"github.com/microbus-io/fabric/examples/directory/directoryapi"
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
	_ directoryapi.Client
)

// ToDo defines the interface that the microservice must implement.
// The intermediate delegates handling to this interface.
type ToDo interface {
	OnStartup(ctx context.Context) (err error)
	OnShutdown(ctx context.Context) (err error)
	Create(ctx context.Context, person *directoryapi.Person) (created *directoryapi.Person, err error)
	Load(ctx context.Context, key directoryapi.PersonKey) (person *directoryapi.Person, ok bool, err error)
	Delete(ctx context.Context, key directoryapi.PersonKey) (ok bool, err error)
	Update(ctx context.Context, person *directoryapi.Person) (updated *directoryapi.Person, ok bool, err error)
	LoadByEmail(ctx context.Context, email string) (person *directoryapi.Person, ok bool, err error)
	List(ctx context.Context) (keys []directoryapi.PersonKey, err error)
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
		Connector: connector.New("directory.example"),
		impl: impl,
	}
	svc.SetVersion(version)
	svc.SetDescription(`The directory microservice stores personal records in a SQL database.`)
	
	// Lifecycle
	svc.SetOnStartup(svc.impl.OnStartup)
	svc.SetOnShutdown(svc.impl.OnShutdown)

	// Configs
	svc.SetOnConfigChanged(svc.doOnConfigChanged)
	svc.DefineConfig(
		"SQL",
		cfg.Description(`SQL is the connection string to the database.`),
	)

	// OpenAPI
	svc.Subscribe("GET", `:*/openapi.json`, svc.doOpenAPI)	

	// Functions
	svc.Subscribe(`*`, `:443/create`, svc.doCreate)
	svc.Subscribe(`*`, `:443/load`, svc.doLoad)
	svc.Subscribe(`*`, `:443/delete`, svc.doDelete)
	svc.Subscribe(`*`, `:443/update`, svc.doUpdate)
	svc.Subscribe(`*`, `:443/load-by-email`, svc.doLoadByEmail)
	svc.Subscribe(`*`, `:443/list`, svc.doList)

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
	if r.URL.Port() == "443" || "443" == "*" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `function`,
			Name:        `Create`,
			Method:      `*`,
			Path:        `:443/create`,
			Summary:     `Create(person *Person) (created *Person)`,
			Description: `Create registers the person in the directory.`,
			InputArgs: struct {
				Xperson *directoryapi.Person `json:"person"`
			}{},
			OutputArgs: struct {
				Xcreated *directoryapi.Person `json:"created"`
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "*" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `function`,
			Name:        `Load`,
			Method:      `*`,
			Path:        `:443/load`,
			Summary:     `Load(key PersonKey) (person *Person, ok bool)`,
			Description: `Load looks up a person in the directory.`,
			InputArgs: struct {
				Xkey directoryapi.PersonKey `json:"key"`
			}{},
			OutputArgs: struct {
				Xperson *directoryapi.Person `json:"person"`
				Xok bool `json:"ok"`
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "*" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `function`,
			Name:        `Delete`,
			Method:      `*`,
			Path:        `:443/delete`,
			Summary:     `Delete(key PersonKey) (ok bool)`,
			Description: `Delete removes a person from the directory.`,
			InputArgs: struct {
				Xkey directoryapi.PersonKey `json:"key"`
			}{},
			OutputArgs: struct {
				Xok bool `json:"ok"`
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "*" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `function`,
			Name:        `Update`,
			Method:      `*`,
			Path:        `:443/update`,
			Summary:     `Update(person *Person) (updated *Person, ok bool)`,
			Description: `Update updates the person's data in the directory.`,
			InputArgs: struct {
				Xperson *directoryapi.Person `json:"person"`
			}{},
			OutputArgs: struct {
				Xupdated *directoryapi.Person `json:"updated"`
				Xok bool `json:"ok"`
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "*" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `function`,
			Name:        `LoadByEmail`,
			Method:      `*`,
			Path:        `:443/load-by-email`,
			Summary:     `LoadByEmail(email string) (person *Person, ok bool)`,
			Description: `LoadByEmail looks up a person in the directory by their email.`,
			InputArgs: struct {
				Xemail string `json:"email"`
			}{},
			OutputArgs: struct {
				Xperson *directoryapi.Person `json:"person"`
				Xok bool `json:"ok"`
			}{},
		})
	}
	if r.URL.Port() == "443" || "443" == "*" {
		oapiSvc.Endpoints = append(oapiSvc.Endpoints, &openapi.Endpoint{
			Type:        `function`,
			Name:        `List`,
			Method:      `*`,
			Path:        `:443/list`,
			Summary:     `List() (keys []PersonKey)`,
			Description: `List returns the keys of all the persons in the directory.`,
			InputArgs: struct {
			}{},
			OutputArgs: struct {
				Xkeys []directoryapi.PersonKey `json:"keys"`
			}{},
		})
	}

	if len(oapiSvc.Endpoints) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := json.MarshalIndent(&oapiSvc, "", "    ")
	if err != nil {
		return errors.Trace(err)
	}
	_, err = w.Write(b)
	return errors.Trace(err)
}

// doOnConfigChanged is called when the config of the microservice changes.
func (svc *Intermediate) doOnConfigChanged(ctx context.Context, changed func(string) bool) (err error) {
	return nil
}

/*
SQL is the connection string to the database.
*/
func (svc *Intermediate) SQL() (dsn string) {
	_val := svc.Config("SQL")
	return _val
}

/*
SetSQL sets the value of the configuration property.
Settings configs is only enabled in the TESTINGAPP environment where the configurator core microservice is disabled.

SQL is the connection string to the database.
*/
func (svc *Intermediate) SetSQL(dsn string) error {
	return svc.SetConfig("SQL", fmt.Sprintf("%v", dsn))
}

// doCreate handles marshaling for the Create function.
func (svc *Intermediate) doCreate(w http.ResponseWriter, r *http.Request) error {
	var i directoryapi.CreateIn
	var o directoryapi.CreateOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.Created, err = svc.impl.Create(
		r.Context(),
		i.Person,
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

// doLoad handles marshaling for the Load function.
func (svc *Intermediate) doLoad(w http.ResponseWriter, r *http.Request) error {
	var i directoryapi.LoadIn
	var o directoryapi.LoadOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.Person, o.Ok, err = svc.impl.Load(
		r.Context(),
		i.Key,
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

// doDelete handles marshaling for the Delete function.
func (svc *Intermediate) doDelete(w http.ResponseWriter, r *http.Request) error {
	var i directoryapi.DeleteIn
	var o directoryapi.DeleteOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.Ok, err = svc.impl.Delete(
		r.Context(),
		i.Key,
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

// doUpdate handles marshaling for the Update function.
func (svc *Intermediate) doUpdate(w http.ResponseWriter, r *http.Request) error {
	var i directoryapi.UpdateIn
	var o directoryapi.UpdateOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.Updated, o.Ok, err = svc.impl.Update(
		r.Context(),
		i.Person,
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

// doLoadByEmail handles marshaling for the LoadByEmail function.
func (svc *Intermediate) doLoadByEmail(w http.ResponseWriter, r *http.Request) error {
	var i directoryapi.LoadByEmailIn
	var o directoryapi.LoadByEmailOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.Person, o.Ok, err = svc.impl.LoadByEmail(
		r.Context(),
		i.Email,
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

// doList handles marshaling for the List function.
func (svc *Intermediate) doList(w http.ResponseWriter, r *http.Request) error {
	var i directoryapi.ListIn
	var o directoryapi.ListOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.Keys, err = svc.impl.List(
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
