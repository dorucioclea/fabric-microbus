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

	"github.com/microbus-io/fabric/services/control/controlapi"
)

var (
	_ context.Context
	_ *json.Decoder
	_ *http.Request
	_ time.Duration
	_ *errors.TracedError
	_ *httpx.ResponseRecorder
	_ sub.Option
	_ controlapi.Client
)

// Mock is a mockable version of the control.sys microservice,
// allowing functions, sinks and web handlers to be mocked.
type Mock struct {
	*connector.Connector
	MockPing func(ctx context.Context) (pong int, err error)
	MockConfigRefresh func(ctx context.Context) (err error)
}

// NewMock creates a new mockable version of the microservice.
func NewMock(version int) *Mock {
	svc := &Mock{
		Connector: connector.New("control.sys"),
	}
	svc.SetVersion(version)
	svc.SetDescription(`This microservice is created for the sake of generating the client API for the :888 control subscriptions.
The microservice itself does nothing and should not be included in applications.`)
	svc.SetOnStartup(svc.doOnStartup)
	
	// Functions
	svc.Subscribe(`:888/ping`, svc.doPing)
	svc.Subscribe(`:888/config-refresh`, svc.doConfigRefresh)

	return svc
}

// doOnStartup makes sure that the mock is not executed in a non-dev environment.
func (svc *Mock) doOnStartup(ctx context.Context) (err error) {
	if svc.Deployment() != connector.LOCAL && svc.Deployment() != connector.TESTINGAPP {
		return errors.Newf("mocking disallowed in '%s' deployment", svc.Deployment())
	}
    return nil
}

// doPing handles marshaling for the Ping function.
func (svc *Mock) doPing(w http.ResponseWriter, r *http.Request) error {
	if svc.MockPing == nil {
		return errors.New("mocked endpoint 'Ping' not implemented")
	}
	var i controlapi.PingIn
	var o controlapi.PingOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	o.Pong, err = svc.MockPing(
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

// doConfigRefresh handles marshaling for the ConfigRefresh function.
func (svc *Mock) doConfigRefresh(w http.ResponseWriter, r *http.Request) error {
	if svc.MockConfigRefresh == nil {
		return errors.New("mocked endpoint 'ConfigRefresh' not implemented")
	}
	var i controlapi.ConfigRefreshIn
	var o controlapi.ConfigRefreshOut
	err := httpx.ParseRequestData(r, &i)
	if err!=nil {
		return errors.Trace(err)
	}
	err = svc.MockConfigRefresh(
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