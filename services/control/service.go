/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package control

import (
	"context"
	"net/http"

	"github.com/microbus-io/fabric/errors"

	"github.com/microbus-io/fabric/services/control/intermediate"
)

var (
	_ errors.TracedError
	_ http.Request
)

/*
Service implements the control.sys microservice.

This microservice is created for the sake of generating the client API for the :888 control subscriptions.
The microservice itself does nothing and should not be included in applications.
*/
type Service struct {
	*intermediate.Intermediate // DO NOT REMOVE
}

// OnStartup is called when the microservice is started up.
func (svc *Service) OnStartup(ctx context.Context) (err error) {
	return errors.New("unstartable")
}

// OnShutdown is called when the microservice is shut down.
func (svc *Service) OnShutdown(ctx context.Context) (err error) {
	return nil
}

/*
Ping responds to the message with a pong.
*/
func (svc *Service) Ping(ctx context.Context) (pong int, err error) {
	return 0, nil
}

/*
ConfigRefresh pulls the latest config values from the configurator service.
*/
func (svc *Service) ConfigRefresh(ctx context.Context) (err error) {
	return nil
}

/*
Trace forces exporting the indicated tracing span.
*/
func (svc *Service) Trace(ctx context.Context, id string) (err error) {
	return nil
}
