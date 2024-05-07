/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

package eventsource

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
	"github.com/microbus-io/fabric/utils"

	"github.com/stretchr/testify/assert"

	"github.com/microbus-io/fabric/examples/eventsource/eventsourceapi"
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
	_ utils.InfiniteChan[int]
	_ assert.TestingT
	_ *eventsourceapi.Client
)

var (
	sequence int
)

var (
	// App manages the lifecycle of the microservices used in the test
	App *application.Application
	// Svc is the eventsource.example microservice being tested
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
	return context.WithValue(context.Background(), frame.ContextKey, http.Header{})
}

// RegisterTestCase assists in asserting against the results of executing Register.
type RegisterTestCase struct {
	_t *testing.T
	_testName string
	allowed bool
	err error
}

// Name sets a name to the test case.
func (tc *RegisterTestCase) Name(testName string) *RegisterTestCase {
	tc._testName = testName
	return tc
}

// Expect asserts no error and exact return values.
func (tc *RegisterTestCase) Expect(allowed bool) *RegisterTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, allowed, tc.allowed)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *RegisterTestCase) Error(errContains string) *RegisterTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *RegisterTestCase) ErrorCode(statusCode int) *RegisterTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *RegisterTestCase) NoError() *RegisterTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *RegisterTestCase) Assert(asserter func(t *testing.T, allowed bool, err error)) *RegisterTestCase {
	tc._t.Run(tc._testName, func(t *testing.T) {
		asserter(t, tc.allowed, tc.err)
	})
	return tc
}

// Get returns the result of executing Register.
func (tc *RegisterTestCase) Get() (allowed bool, err error) {
	return tc.allowed, tc.err
}

// Register executes the function and returns a corresponding test case.
func Register(t *testing.T, ctx context.Context, email string) *RegisterTestCase {
	tc := &RegisterTestCase{_t: t}
	tc.err = utils.CatchPanic(func() error {
		tc.allowed, tc.err = Svc.Register(ctx, email)
		return tc.err
	})
	return tc
}

// OnAllowRegisterTestCase assists in asserting the sink of OnAllowRegister.
type OnAllowRegisterTestCase struct {
	_t *testing.T
	_testName string
    _asserters []func(*testing.T)
    ctx context.Context
	email string
	allow bool
	err error
}

// Name sets a name to the test case.
func (tc *OnAllowRegisterTestCase) Name(testName string) *OnAllowRegisterTestCase {
	tc._testName = testName
	return tc
}

// Expect asserts an exact match for the input arguments of the event sink.
func (tc *OnAllowRegisterTestCase) Expect(email string) *OnAllowRegisterTestCase {
    tc._asserters = append(tc._asserters, func(t *testing.T) {
        assert.Equal(t, email, tc.email)
    })
	return tc
}

// Assert sets a custom function to assert the input args of the event sink.
func (tc *OnAllowRegisterTestCase) Assert(asserter func(t *testing.T, ctx context.Context, email string)) *OnAllowRegisterTestCase {
	tc._asserters = append(tc._asserters, func(t *testing.T) {
		asserter(t, tc.ctx, tc.email)
    })
	return tc
}

// OnAllowRegister sets an event listener and returns a corresponding test case.
func OnAllowRegister(t *testing.T, allow bool, err error) *OnAllowRegisterTestCase {
	tc := &OnAllowRegisterTestCase{
		_t: t,
		allow: allow,
		err: err,
	}
    sequence ++
	con := connector.New(fmt.Sprintf("%s.%d", "OnAllowRegister", sequence))
	eventsourceapi.NewHook(con).OnAllowRegister(func(ctx context.Context, email string) (allow bool, err error) {
		eventsourceapi.NewHook(con).OnAllowRegister(nil)
		tc.ctx = ctx
        tc.email = email
		for _, asserter := range tc._asserters {
			tc._t.Run(tc._testName, asserter)
		}
		return tc.allow, tc.err
	})
	App.Include(con)
	con.Startup()
	return tc
}

// OnRegisteredTestCase assists in asserting the sink of OnRegistered.
type OnRegisteredTestCase struct {
	_t *testing.T
	_testName string
    _asserters []func(*testing.T)
    ctx context.Context
	email string
	err error
}

// Name sets a name to the test case.
func (tc *OnRegisteredTestCase) Name(testName string) *OnRegisteredTestCase {
	tc._testName = testName
	return tc
}

// Expect asserts an exact match for the input arguments of the event sink.
func (tc *OnRegisteredTestCase) Expect(email string) *OnRegisteredTestCase {
    tc._asserters = append(tc._asserters, func(t *testing.T) {
        assert.Equal(t, email, tc.email)
    })
	return tc
}

// Assert sets a custom function to assert the input args of the event sink.
func (tc *OnRegisteredTestCase) Assert(asserter func(t *testing.T, ctx context.Context, email string)) *OnRegisteredTestCase {
	tc._asserters = append(tc._asserters, func(t *testing.T) {
		asserter(t, tc.ctx, tc.email)
    })
	return tc
}

// OnRegistered sets an event listener and returns a corresponding test case.
func OnRegistered(t *testing.T, err error) *OnRegisteredTestCase {
	tc := &OnRegisteredTestCase{
		_t: t,
		err: err,
	}
    sequence ++
	con := connector.New(fmt.Sprintf("%s.%d", "OnRegistered", sequence))
	eventsourceapi.NewHook(con).OnRegistered(func(ctx context.Context, email string) (err error) {
		eventsourceapi.NewHook(con).OnRegistered(nil)
		tc.ctx = ctx
        tc.email = email
		for _, asserter := range tc._asserters {
			tc._t.Run(tc._testName, asserter)
		}
		return tc.err
	})
	App.Include(con)
	con.Startup()
	return tc
}
