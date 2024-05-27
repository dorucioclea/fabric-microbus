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

	"github.com/microbus-io/fabric/examples/hello/helloapi"
)

var (
	_ context.Context
	_ *json.Decoder
	_ *http.Request
	_ time.Duration
	_ *errors.TracedError
	_ *httpx.ResponseRecorder
	_ sub.Option
	_ helloapi.Client
)

// Mock is a mockable version of the hello.example microservice, allowing functions, event sinks and web handlers to be mocked.
type Mock struct {
	*connector.Connector
	mockHello func(w http.ResponseWriter, r *http.Request) (err error)
	mockEcho func(w http.ResponseWriter, r *http.Request) (err error)
	mockPing func(w http.ResponseWriter, r *http.Request) (err error)
	mockCalculator func(w http.ResponseWriter, r *http.Request) (err error)
	mockBusJPEG func(w http.ResponseWriter, r *http.Request) (err error)
	mockLocalization func(w http.ResponseWriter, r *http.Request) (err error)
	mockRoot func(w http.ResponseWriter, r *http.Request) (err error)
}

// NewMock creates a new mockable version of the microservice.
func NewMock() *Mock {
	svc := &Mock{
		Connector: connector.New("hello.example"),
	}
	svc.SetVersion(7357) // Stands for TEST
	svc.SetDescription(`The Hello microservice demonstrates the various capabilities of a microservice.`)
	svc.SetOnStartup(svc.doOnStartup)

	// Webs
	svc.Subscribe(`*`, `:443/hello`, svc.doHello)
	svc.Subscribe(`*`, `:443/echo`, svc.doEcho)
	svc.Subscribe(`*`, `:443/ping`, svc.doPing)
	svc.Subscribe(`*`, `:443/calculator`, svc.doCalculator)
	svc.Subscribe(`GET`, `:443/bus.jpeg`, svc.doBusJPEG)
	svc.Subscribe(`*`, `:443/localization`, svc.doLocalization)
	svc.Subscribe(`*`, `//root`, svc.doRoot)

	return svc
}

// doOnStartup makes sure that the mock is not executed in a non-dev environment.
func (svc *Mock) doOnStartup(ctx context.Context) (err error) {
	if svc.Deployment() != connector.LOCAL && svc.Deployment() != connector.TESTING {
		return errors.Newf("mocking disallowed in '%s' deployment", svc.Deployment())
	}
	return nil
}

// doHello handles the Hello web handler.
func (svc *Mock) doHello(w http.ResponseWriter, r *http.Request) (err error) {
	if svc.mockHello == nil {
		return errors.New("mocked endpoint 'Hello' not implemented")
	}
	err = svc.mockHello(w, r)
	return errors.Trace(err)
}

// MockHello sets up a mock handler for the Hello web handler.
func (svc *Mock) MockHello(handler func(w http.ResponseWriter, r *http.Request) (err error)) *Mock {
	svc.mockHello = handler
	return svc
}

// doEcho handles the Echo web handler.
func (svc *Mock) doEcho(w http.ResponseWriter, r *http.Request) (err error) {
	if svc.mockEcho == nil {
		return errors.New("mocked endpoint 'Echo' not implemented")
	}
	err = svc.mockEcho(w, r)
	return errors.Trace(err)
}

// MockEcho sets up a mock handler for the Echo web handler.
func (svc *Mock) MockEcho(handler func(w http.ResponseWriter, r *http.Request) (err error)) *Mock {
	svc.mockEcho = handler
	return svc
}

// doPing handles the Ping web handler.
func (svc *Mock) doPing(w http.ResponseWriter, r *http.Request) (err error) {
	if svc.mockPing == nil {
		return errors.New("mocked endpoint 'Ping' not implemented")
	}
	err = svc.mockPing(w, r)
	return errors.Trace(err)
}

// MockPing sets up a mock handler for the Ping web handler.
func (svc *Mock) MockPing(handler func(w http.ResponseWriter, r *http.Request) (err error)) *Mock {
	svc.mockPing = handler
	return svc
}

// doCalculator handles the Calculator web handler.
func (svc *Mock) doCalculator(w http.ResponseWriter, r *http.Request) (err error) {
	if svc.mockCalculator == nil {
		return errors.New("mocked endpoint 'Calculator' not implemented")
	}
	err = svc.mockCalculator(w, r)
	return errors.Trace(err)
}

// MockCalculator sets up a mock handler for the Calculator web handler.
func (svc *Mock) MockCalculator(handler func(w http.ResponseWriter, r *http.Request) (err error)) *Mock {
	svc.mockCalculator = handler
	return svc
}

// doBusJPEG handles the BusJPEG web handler.
func (svc *Mock) doBusJPEG(w http.ResponseWriter, r *http.Request) (err error) {
	if svc.mockBusJPEG == nil {
		return errors.New("mocked endpoint 'BusJPEG' not implemented")
	}
	err = svc.mockBusJPEG(w, r)
	return errors.Trace(err)
}

// MockBusJPEG sets up a mock handler for the BusJPEG web handler.
func (svc *Mock) MockBusJPEG(handler func(w http.ResponseWriter, r *http.Request) (err error)) *Mock {
	svc.mockBusJPEG = handler
	return svc
}

// doLocalization handles the Localization web handler.
func (svc *Mock) doLocalization(w http.ResponseWriter, r *http.Request) (err error) {
	if svc.mockLocalization == nil {
		return errors.New("mocked endpoint 'Localization' not implemented")
	}
	err = svc.mockLocalization(w, r)
	return errors.Trace(err)
}

// MockLocalization sets up a mock handler for the Localization web handler.
func (svc *Mock) MockLocalization(handler func(w http.ResponseWriter, r *http.Request) (err error)) *Mock {
	svc.mockLocalization = handler
	return svc
}

// doRoot handles the Root web handler.
func (svc *Mock) doRoot(w http.ResponseWriter, r *http.Request) (err error) {
	if svc.mockRoot == nil {
		return errors.New("mocked endpoint 'Root' not implemented")
	}
	err = svc.mockRoot(w, r)
	return errors.Trace(err)
}

// MockRoot sets up a mock handler for the Root web handler.
func (svc *Mock) MockRoot(handler func(w http.ResponseWriter, r *http.Request) (err error)) *Mock {
	svc.mockRoot = handler
	return svc
}
