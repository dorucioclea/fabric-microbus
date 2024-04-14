/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package intermediate serves as the foundation of the calculator.example microservice.

The Calculator microservice performs simple mathematical operations.
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

	"github.com/microbus-io/fabric/examples/calculator/resources"
	"github.com/microbus-io/fabric/examples/calculator/calculatorapi"
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
	_ calculatorapi.Client
)

// ToDo defines the interface that the microservice must implement.
// The intermediate delegates handling to this interface.
type ToDo interface {
	OnStartup(ctx context.Context) (err error)
	OnShutdown(ctx context.Context) (err error)
	Arithmetic(ctx context.Context, x int, op string, y int) (xEcho int, opEcho string, yEcho int, result int, err error)
	Square(ctx context.Context, x int) (xEcho int, result int, err error)
	Distance(ctx context.Context, p1 calculatorapi.Point, p2 calculatorapi.Point) (d float64, err error)
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
		Connector: connector.New("calculator.example"),
		impl: impl,
	}
	svc.SetVersion(version)
	svc.SetDescription(`The Calculator microservice performs simple mathematical operations.`)

	// Lifecycle
	svc.SetOnStartup(svc.impl.OnStartup)
	svc.SetOnShutdown(svc.impl.OnShutdown)	

	// Functions
	svc.Subscribe(`:443/arithmetic`, svc.doArithmetic)
	svc.Subscribe(`:443/square`, svc.doSquare)
	svc.Subscribe(`:443/distance`, svc.doDistance)

	// Metrics
	svc.DefineCounter(
		`calculator_op_count`,
		`OpCount tracks the types of arithmetic calculations.`,
		[]string{"op"},
	)

	// Resources file system
	svc.SetResFS(resources.FS)

	return svc
}

// doOnConfigChanged is called when the config of the microservice changes.
func (svc *Intermediate) doOnConfigChanged(ctx context.Context, changed func(string) bool) (err error) {
	return nil
}

// doArithmetic handles marshaling for the Arithmetic function.
func (svc *Intermediate) doArithmetic(w http.ResponseWriter, r *http.Request) error {
	var i calculatorapi.ArithmeticIn
	var o calculatorapi.ArithmeticOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.XEcho, o.OpEcho, o.YEcho, o.Result, err = svc.impl.Arithmetic(
		r.Context(),
		i.X,
		i.Op,
		i.Y,
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

// doSquare handles marshaling for the Square function.
func (svc *Intermediate) doSquare(w http.ResponseWriter, r *http.Request) error {
	var i calculatorapi.SquareIn
	var o calculatorapi.SquareOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.XEcho, o.Result, err = svc.impl.Square(
		r.Context(),
		i.X,
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

// doDistance handles marshaling for the Distance function.
func (svc *Intermediate) doDistance(w http.ResponseWriter, r *http.Request) error {
	var i calculatorapi.DistanceIn
	var o calculatorapi.DistanceOut
	err := httpx.ParseRequestData(r, &i)
	if err != nil {
		return errors.Trace(err)
	}
	o.D, err = svc.impl.Distance(
		r.Context(),
		i.P1,
		i.P2,
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

/*
IncrementOpCount increments the value of the "calculator_op_count" metric.
OpCount tracks the types of arithmetic calculations.
*/
func (svc *Intermediate) IncrementOpCount(num int, op string) error {
	xnum := float64(num)
	xop := fmt.Sprintf("%v", op)
	return svc.IncrementMetric("calculator_op_count", xnum, xop)
}
