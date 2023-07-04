/*
Copyright (c) 2023 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

package calculator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/microbus-io/fabric/application"
	"github.com/microbus-io/fabric/connector"
	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/frame"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/pub"
	"github.com/microbus-io/fabric/shardedsql"
	"github.com/microbus-io/fabric/utils"

	"github.com/stretchr/testify/assert"

	"github.com/microbus-io/fabric/examples/calculator/calculatorapi"
)

var (
	_ bytes.Buffer
	_ context.Context
	_ fmt.Stringer
	_ io.Reader
	_ *http.Request
	_ os.File
	_ time.Time
	_ strings.Builder
	_ *connector.Connector
	_ *errors.TracedError
	_ frame.Frame
	_ *httpx.BodyReader
	_ pub.Option
	_ *shardedsql.DB
	_ utils.InfiniteChan[int]
	_ assert.TestingT
	_ *calculatorapi.Client
)

var (
	sequence int
)

var (
	// App manages the lifecycle of the microservices used in the test
	App *application.Application
	// Svc is the calculator.example microservice being tested
	Svc *Service
)

func TestMain(m *testing.M) {
	var code int

	// Initialize the application
	err := func() error {
		var err error
		App = application.NewTesting()
		Svc = NewService().(*Service)
		err = Initialize()
		if err != nil {
			return err
		}
		err = App.Startup()
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "--- FAIL: %+v\n", err)
		code = 19
	}

	// Run the tests
	if err == nil {
		code = m.Run()
	}

	// Terminate the app
	err = func() error {
		var err error
		var lastErr error
		err = Terminate()
		if err != nil {
			lastErr = err
		}
		err = App.Shutdown()
		if err != nil {
			lastErr = err
		}
		return lastErr
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "--- FAIL: %+v\n", err)
	}

	os.Exit(code)
}

// Context creates a new context for a test.
func Context(t *testing.T) context.Context {
	return context.Background()
}

// ArithmeticTestCase assists in asserting against the results of executing Arithmetic.
type ArithmeticTestCase struct {
	_t *testing.T
	_testName string
	xEcho int
	opEcho string
	yEcho int
	result int
	err error
}

// Name sets a name to the test case.
func (tc *ArithmeticTestCase) Name(testName string) *ArithmeticTestCase {
	tc._testName = testName
	return tc
}

// Expect asserts no error and exact return values.
func (tc *ArithmeticTestCase) Expect(xEcho int, opEcho string, yEcho int, result int) *ArithmeticTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, xEcho, tc.xEcho)
			assert.Equal(t, opEcho, tc.opEcho)
			assert.Equal(t, yEcho, tc.yEcho)
			assert.Equal(t, result, tc.result)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *ArithmeticTestCase) Error(errContains string) *ArithmeticTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *ArithmeticTestCase) ErrorCode(statusCode int) *ArithmeticTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *ArithmeticTestCase) NoError() *ArithmeticTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *ArithmeticTestCase) Assert(asserter func(t *testing.T, xEcho int, opEcho string, yEcho int, result int, err error)) *ArithmeticTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		asserter(t, tc.xEcho, tc.opEcho, tc.yEcho, tc.result, tc.err)
	})
	return tc
}

// Get returns the result of executing Arithmetic.
func (tc *ArithmeticTestCase) Get() (xEcho int, opEcho string, yEcho int, result int, err error) {
	return tc.xEcho, tc.opEcho, tc.yEcho, tc.result, tc.err
}

// Arithmetic executes the function and returns a corresponding test case.
func Arithmetic(t *testing.T, ctx context.Context, x int, op string, y int) *ArithmeticTestCase {
	tc := &ArithmeticTestCase{_t: t}
	tc.err = utils.CatchPanic(func() error {
		tc.xEcho, tc.opEcho, tc.yEcho, tc.result, tc.err = Svc.Arithmetic(ctx, x, op, y)
		return tc.err
	})
	return tc
}

// SquareTestCase assists in asserting against the results of executing Square.
type SquareTestCase struct {
	_t *testing.T
	_testName string
	xEcho int
	result int
	err error
}

// Name sets a name to the test case.
func (tc *SquareTestCase) Name(testName string) *SquareTestCase {
	tc._testName = testName
	return tc
}

// Expect asserts no error and exact return values.
func (tc *SquareTestCase) Expect(xEcho int, result int) *SquareTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, xEcho, tc.xEcho)
			assert.Equal(t, result, tc.result)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *SquareTestCase) Error(errContains string) *SquareTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *SquareTestCase) ErrorCode(statusCode int) *SquareTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *SquareTestCase) NoError() *SquareTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *SquareTestCase) Assert(asserter func(t *testing.T, xEcho int, result int, err error)) *SquareTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		asserter(t, tc.xEcho, tc.result, tc.err)
	})
	return tc
}

// Get returns the result of executing Square.
func (tc *SquareTestCase) Get() (xEcho int, result int, err error) {
	return tc.xEcho, tc.result, tc.err
}

// Square executes the function and returns a corresponding test case.
func Square(t *testing.T, ctx context.Context, x int) *SquareTestCase {
	tc := &SquareTestCase{_t: t}
	tc.err = utils.CatchPanic(func() error {
		tc.xEcho, tc.result, tc.err = Svc.Square(ctx, x)
		return tc.err
	})
	return tc
}

// DistanceTestCase assists in asserting against the results of executing Distance.
type DistanceTestCase struct {
	_t *testing.T
	_testName string
	d float64
	err error
}

// Name sets a name to the test case.
func (tc *DistanceTestCase) Name(testName string) *DistanceTestCase {
	tc._testName = testName
	return tc
}

// Expect asserts no error and exact return values.
func (tc *DistanceTestCase) Expect(d float64) *DistanceTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, d, tc.d)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *DistanceTestCase) Error(errContains string) *DistanceTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *DistanceTestCase) ErrorCode(statusCode int) *DistanceTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *DistanceTestCase) NoError() *DistanceTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *DistanceTestCase) Assert(asserter func(t *testing.T, d float64, err error)) *DistanceTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		asserter(t, tc.d, tc.err)
	})
	return tc
}

// Get returns the result of executing Distance.
func (tc *DistanceTestCase) Get() (d float64, err error) {
	return tc.d, tc.err
}

// Distance executes the function and returns a corresponding test case.
func Distance(t *testing.T, ctx context.Context, p1 calculatorapi.Point, p2 calculatorapi.Point) *DistanceTestCase {
	tc := &DistanceTestCase{_t: t}
	tc.err = utils.CatchPanic(func() error {
		tc.d, tc.err = Svc.Distance(ctx, p1, p2)
		return tc.err
	})
	return tc
}
