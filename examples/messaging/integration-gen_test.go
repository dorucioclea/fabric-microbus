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

package messaging

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
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/pub"
	"github.com/microbus-io/fabric/shardedsql"
	"github.com/microbus-io/fabric/utils"

	"github.com/stretchr/testify/assert"

	"github.com/microbus-io/fabric/examples/messaging/messagingapi"
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
	_ *httpx.BodyReader
	_ pub.Option
	_ *shardedsql.DB
	_ utils.InfiniteChan[int]
	_ assert.TestingT
	_ *messagingapi.Client
)

var (
	sequence int
)

var (
	// App manages the lifecycle of the microservices used in the test
	App *application.Application
	// Svc is the messaging.example microservice being tested
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

type WebOption func(req *pub.Request) error

// GET sets the method of the request.
func GET() WebOption {
	return WebOption(pub.Method("GET"))
}

// DELETE sets the method of the request.
func DELETE() WebOption {
	return WebOption(pub.Method("DELETE"))
}

// HEAD sets the method of the request.
func HEAD() WebOption {
	return WebOption(pub.Method("HEAD"))
}

// POST sets the method and body of the request.
func POST(body any) WebOption {
	return func(req *pub.Request) error {
		pub.Method("POST")(req)
		return pub.Body(body)(req)
	}
}

// PUT sets the method and body of the request.
func PUT(body any) WebOption {
	return func(req *pub.Request) error {
		pub.Method("PUT")(req)
		return pub.Body(body)(req)
	}
}

// PATCH sets the method and body of the request.
func PATCH(body any) WebOption {
	return func(req *pub.Request) error {
		pub.Method("PATCH")(req)
		return pub.Body(body)(req)
	}
}

// Method sets the method of the request.
func Method(method string) WebOption {
	return WebOption(pub.Method(method))
}

// Header sets the header of the request. It overwrites any previously set value.
func Header(name string, value string) WebOption {
	return WebOption(pub.Header(name, value))
}

// QueryArg adds the query argument to the request.
// The same argument may have multiple values.
func QueryArg(name string, value any) WebOption {
	return WebOption(pub.QueryArg(name, value))
}

// Query adds the escaped query arguments to the request.
// The same argument may have multiple values.
func Query(escapedQueryArgs string) WebOption {
	return WebOption(pub.Query(escapedQueryArgs))
}

// Body sets the body of the request.
// Arguments of type io.Reader, []byte and string are serialized in binary form.
// url.Values is serialized as form data.
// All other types are serialized as JSON.
func Body(body any) WebOption {
	return WebOption(pub.Body(body))
}

// ContentType sets the Content-Type header.
func ContentType(contentType string) WebOption {
	return WebOption(pub.ContentType(contentType))
}

// HomeTestCase assists in asserting against the results of executing Home.
type HomeTestCase struct {
	t *testing.T
	testName string
	res *http.Response
	err error
}

// Name sets a name to the test case.
func (tc *HomeTestCase) Name(testName string) *HomeTestCase {
	tc.testName = testName
	return tc
}

// StatusOK asserts no error and a status code 200.
func (tc *HomeTestCase) StatusOK() *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, http.StatusOK)
		}
	})
	return tc
}

// StatusCode asserts no error and a status code.
func (tc *HomeTestCase) StatusCode(statusCode int) *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, statusCode)
		}
	})
	return tc
}

// BodyContains asserts no error and that the response contains a string or byte array.
func (tc *HomeTestCase) BodyContains(bodyContains any) *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyContains.(type) {
			case []byte:
				assert.True(t, bytes.Contains(body, v), `"%v" does not contain "%v"`, body, v)
			case string:
				assert.True(t, bytes.Contains(body, []byte(v)), `"%s" does not contain "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.True(t, bytes.Contains(body, []byte(vv)), `"%s" does not contain "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// BodyNotContains asserts no error and that the response does not contain a string or byte array.
func (tc *HomeTestCase) BodyNotContains(bodyNotContains any) *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyNotContains.(type) {
			case []byte:
				assert.False(t, bytes.Contains(body, v), `"%v" contains "%v"`, body, v)
			case string:
				assert.False(t, bytes.Contains(body, []byte(v)), `"%s" contains "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.False(t, bytes.Contains(body, []byte(vv)), `"%s" contains "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// HeaderContains asserts no error and that the named header contains a string.
func (tc *HomeTestCase) HeaderContains(headerName string, valueContains string) *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.True(t, strings.Contains(tc.res.Header.Get(headerName), valueContains), `header "%s: %s" does not contain "%s"`, headerName, tc.res.Header.Get(headerName), valueContains)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *HomeTestCase) Error(errContains string) *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *HomeTestCase) ErrorCode(statusCode int) *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *HomeTestCase) NoError() *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *HomeTestCase) Assert(asserter func(t *testing.T, res *http.Response, err error)) *HomeTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		asserter(t, tc.res, tc.err)
	})
	return tc
}

// Get returns the result of executing Home.
func (tc *HomeTestCase) Get() (res *http.Response, err error) {
	return tc.res, tc.err
}

// Home executes the web handler and returns a corresponding test case.
func Home(t *testing.T, ctx context.Context, options ...WebOption) *HomeTestCase {
	tc := &HomeTestCase{t: t}
	pubOptions := []pub.Option{
		pub.URL(httpx.JoinHostAndPath("messaging.example", `:443/home`)),
	}
	for _, opt := range options {
		pubOptions = append(pubOptions, pub.Option(opt))
	}
	req, err := pub.NewRequest(pubOptions...)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		panic(err)
	}
	for name, value := range req.Header {
		httpReq.Header[name] = value
	}
	r := httpReq.WithContext(ctx)
	w := httpx.NewResponseRecorder()
	tc.err = utils.CatchPanic(func () error {
		return Svc.Home(w, r)
	})
	tc.res = w.Result()
	return tc
}

// NoQueueTestCase assists in asserting against the results of executing NoQueue.
type NoQueueTestCase struct {
	t *testing.T
	testName string
	res *http.Response
	err error
}

// Name sets a name to the test case.
func (tc *NoQueueTestCase) Name(testName string) *NoQueueTestCase {
	tc.testName = testName
	return tc
}

// StatusOK asserts no error and a status code 200.
func (tc *NoQueueTestCase) StatusOK() *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, http.StatusOK)
		}
	})
	return tc
}

// StatusCode asserts no error and a status code.
func (tc *NoQueueTestCase) StatusCode(statusCode int) *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, statusCode)
		}
	})
	return tc
}

// BodyContains asserts no error and that the response contains a string or byte array.
func (tc *NoQueueTestCase) BodyContains(bodyContains any) *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyContains.(type) {
			case []byte:
				assert.True(t, bytes.Contains(body, v), `"%v" does not contain "%v"`, body, v)
			case string:
				assert.True(t, bytes.Contains(body, []byte(v)), `"%s" does not contain "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.True(t, bytes.Contains(body, []byte(vv)), `"%s" does not contain "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// BodyNotContains asserts no error and that the response does not contain a string or byte array.
func (tc *NoQueueTestCase) BodyNotContains(bodyNotContains any) *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyNotContains.(type) {
			case []byte:
				assert.False(t, bytes.Contains(body, v), `"%v" contains "%v"`, body, v)
			case string:
				assert.False(t, bytes.Contains(body, []byte(v)), `"%s" contains "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.False(t, bytes.Contains(body, []byte(vv)), `"%s" contains "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// HeaderContains asserts no error and that the named header contains a string.
func (tc *NoQueueTestCase) HeaderContains(headerName string, valueContains string) *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.True(t, strings.Contains(tc.res.Header.Get(headerName), valueContains), `header "%s: %s" does not contain "%s"`, headerName, tc.res.Header.Get(headerName), valueContains)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *NoQueueTestCase) Error(errContains string) *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *NoQueueTestCase) ErrorCode(statusCode int) *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *NoQueueTestCase) NoError() *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *NoQueueTestCase) Assert(asserter func(t *testing.T, res *http.Response, err error)) *NoQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		asserter(t, tc.res, tc.err)
	})
	return tc
}

// Get returns the result of executing NoQueue.
func (tc *NoQueueTestCase) Get() (res *http.Response, err error) {
	return tc.res, tc.err
}

// NoQueue executes the web handler and returns a corresponding test case.
func NoQueue(t *testing.T, ctx context.Context, options ...WebOption) *NoQueueTestCase {
	tc := &NoQueueTestCase{t: t}
	pubOptions := []pub.Option{
		pub.URL(httpx.JoinHostAndPath("messaging.example", `:443/no-queue`)),
	}
	for _, opt := range options {
		pubOptions = append(pubOptions, pub.Option(opt))
	}
	req, err := pub.NewRequest(pubOptions...)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		panic(err)
	}
	for name, value := range req.Header {
		httpReq.Header[name] = value
	}
	r := httpReq.WithContext(ctx)
	w := httpx.NewResponseRecorder()
	tc.err = utils.CatchPanic(func () error {
		return Svc.NoQueue(w, r)
	})
	tc.res = w.Result()
	return tc
}

// DefaultQueueTestCase assists in asserting against the results of executing DefaultQueue.
type DefaultQueueTestCase struct {
	t *testing.T
	testName string
	res *http.Response
	err error
}

// Name sets a name to the test case.
func (tc *DefaultQueueTestCase) Name(testName string) *DefaultQueueTestCase {
	tc.testName = testName
	return tc
}

// StatusOK asserts no error and a status code 200.
func (tc *DefaultQueueTestCase) StatusOK() *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, http.StatusOK)
		}
	})
	return tc
}

// StatusCode asserts no error and a status code.
func (tc *DefaultQueueTestCase) StatusCode(statusCode int) *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, statusCode)
		}
	})
	return tc
}

// BodyContains asserts no error and that the response contains a string or byte array.
func (tc *DefaultQueueTestCase) BodyContains(bodyContains any) *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyContains.(type) {
			case []byte:
				assert.True(t, bytes.Contains(body, v), `"%v" does not contain "%v"`, body, v)
			case string:
				assert.True(t, bytes.Contains(body, []byte(v)), `"%s" does not contain "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.True(t, bytes.Contains(body, []byte(vv)), `"%s" does not contain "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// BodyNotContains asserts no error and that the response does not contain a string or byte array.
func (tc *DefaultQueueTestCase) BodyNotContains(bodyNotContains any) *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyNotContains.(type) {
			case []byte:
				assert.False(t, bytes.Contains(body, v), `"%v" contains "%v"`, body, v)
			case string:
				assert.False(t, bytes.Contains(body, []byte(v)), `"%s" contains "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.False(t, bytes.Contains(body, []byte(vv)), `"%s" contains "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// HeaderContains asserts no error and that the named header contains a string.
func (tc *DefaultQueueTestCase) HeaderContains(headerName string, valueContains string) *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.True(t, strings.Contains(tc.res.Header.Get(headerName), valueContains), `header "%s: %s" does not contain "%s"`, headerName, tc.res.Header.Get(headerName), valueContains)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *DefaultQueueTestCase) Error(errContains string) *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *DefaultQueueTestCase) ErrorCode(statusCode int) *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *DefaultQueueTestCase) NoError() *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *DefaultQueueTestCase) Assert(asserter func(t *testing.T, res *http.Response, err error)) *DefaultQueueTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		asserter(t, tc.res, tc.err)
	})
	return tc
}

// Get returns the result of executing DefaultQueue.
func (tc *DefaultQueueTestCase) Get() (res *http.Response, err error) {
	return tc.res, tc.err
}

// DefaultQueue executes the web handler and returns a corresponding test case.
func DefaultQueue(t *testing.T, ctx context.Context, options ...WebOption) *DefaultQueueTestCase {
	tc := &DefaultQueueTestCase{t: t}
	pubOptions := []pub.Option{
		pub.URL(httpx.JoinHostAndPath("messaging.example", `:443/default-queue`)),
	}
	for _, opt := range options {
		pubOptions = append(pubOptions, pub.Option(opt))
	}
	req, err := pub.NewRequest(pubOptions...)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		panic(err)
	}
	for name, value := range req.Header {
		httpReq.Header[name] = value
	}
	r := httpReq.WithContext(ctx)
	w := httpx.NewResponseRecorder()
	tc.err = utils.CatchPanic(func () error {
		return Svc.DefaultQueue(w, r)
	})
	tc.res = w.Result()
	return tc
}

// CacheLoadTestCase assists in asserting against the results of executing CacheLoad.
type CacheLoadTestCase struct {
	t *testing.T
	testName string
	res *http.Response
	err error
}

// Name sets a name to the test case.
func (tc *CacheLoadTestCase) Name(testName string) *CacheLoadTestCase {
	tc.testName = testName
	return tc
}

// StatusOK asserts no error and a status code 200.
func (tc *CacheLoadTestCase) StatusOK() *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, http.StatusOK)
		}
	})
	return tc
}

// StatusCode asserts no error and a status code.
func (tc *CacheLoadTestCase) StatusCode(statusCode int) *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, statusCode)
		}
	})
	return tc
}

// BodyContains asserts no error and that the response contains a string or byte array.
func (tc *CacheLoadTestCase) BodyContains(bodyContains any) *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyContains.(type) {
			case []byte:
				assert.True(t, bytes.Contains(body, v), `"%v" does not contain "%v"`, body, v)
			case string:
				assert.True(t, bytes.Contains(body, []byte(v)), `"%s" does not contain "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.True(t, bytes.Contains(body, []byte(vv)), `"%s" does not contain "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// BodyNotContains asserts no error and that the response does not contain a string or byte array.
func (tc *CacheLoadTestCase) BodyNotContains(bodyNotContains any) *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyNotContains.(type) {
			case []byte:
				assert.False(t, bytes.Contains(body, v), `"%v" contains "%v"`, body, v)
			case string:
				assert.False(t, bytes.Contains(body, []byte(v)), `"%s" contains "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.False(t, bytes.Contains(body, []byte(vv)), `"%s" contains "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// HeaderContains asserts no error and that the named header contains a string.
func (tc *CacheLoadTestCase) HeaderContains(headerName string, valueContains string) *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.True(t, strings.Contains(tc.res.Header.Get(headerName), valueContains), `header "%s: %s" does not contain "%s"`, headerName, tc.res.Header.Get(headerName), valueContains)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *CacheLoadTestCase) Error(errContains string) *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *CacheLoadTestCase) ErrorCode(statusCode int) *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *CacheLoadTestCase) NoError() *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *CacheLoadTestCase) Assert(asserter func(t *testing.T, res *http.Response, err error)) *CacheLoadTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		asserter(t, tc.res, tc.err)
	})
	return tc
}

// Get returns the result of executing CacheLoad.
func (tc *CacheLoadTestCase) Get() (res *http.Response, err error) {
	return tc.res, tc.err
}

// CacheLoad executes the web handler and returns a corresponding test case.
func CacheLoad(t *testing.T, ctx context.Context, options ...WebOption) *CacheLoadTestCase {
	tc := &CacheLoadTestCase{t: t}
	pubOptions := []pub.Option{
		pub.URL(httpx.JoinHostAndPath("messaging.example", `:443/cache-load`)),
	}
	for _, opt := range options {
		pubOptions = append(pubOptions, pub.Option(opt))
	}
	req, err := pub.NewRequest(pubOptions...)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		panic(err)
	}
	for name, value := range req.Header {
		httpReq.Header[name] = value
	}
	r := httpReq.WithContext(ctx)
	w := httpx.NewResponseRecorder()
	tc.err = utils.CatchPanic(func () error {
		return Svc.CacheLoad(w, r)
	})
	tc.res = w.Result()
	return tc
}

// CacheStoreTestCase assists in asserting against the results of executing CacheStore.
type CacheStoreTestCase struct {
	t *testing.T
	testName string
	res *http.Response
	err error
}

// Name sets a name to the test case.
func (tc *CacheStoreTestCase) Name(testName string) *CacheStoreTestCase {
	tc.testName = testName
	return tc
}

// StatusOK asserts no error and a status code 200.
func (tc *CacheStoreTestCase) StatusOK() *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, http.StatusOK)
		}
	})
	return tc
}

// StatusCode asserts no error and a status code.
func (tc *CacheStoreTestCase) StatusCode(statusCode int) *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.Equal(t, tc.res.StatusCode, statusCode)
		}
	})
	return tc
}

// BodyContains asserts no error and that the response contains a string or byte array.
func (tc *CacheStoreTestCase) BodyContains(bodyContains any) *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyContains.(type) {
			case []byte:
				assert.True(t, bytes.Contains(body, v), `"%v" does not contain "%v"`, body, v)
			case string:
				assert.True(t, bytes.Contains(body, []byte(v)), `"%s" does not contain "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.True(t, bytes.Contains(body, []byte(vv)), `"%s" does not contain "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// BodyNotContains asserts no error and that the response does not contain a string or byte array.
func (tc *CacheStoreTestCase) BodyNotContains(bodyNotContains any) *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			body := tc.res.Body.(*httpx.BodyReader).Bytes()
			switch v := bodyNotContains.(type) {
			case []byte:
				assert.False(t, bytes.Contains(body, v), `"%v" contains "%v"`, body, v)
			case string:
				assert.False(t, bytes.Contains(body, []byte(v)), `"%s" contains "%s"`, string(body), v)
			default:
				vv := fmt.Sprintf("%v", v)
				assert.False(t, bytes.Contains(body, []byte(vv)), `"%s" contains "%s"`, string(body), vv)
			}
		}
	})
	return tc
}

// HeaderContains asserts no error and that the named header contains a string.
func (tc *CacheStoreTestCase) HeaderContains(headerName string, valueContains string) *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.NoError(t, tc.err) {
			assert.True(t, strings.Contains(tc.res.Header.Get(headerName), valueContains), `header "%s: %s" does not contain "%s"`, headerName, tc.res.Header.Get(headerName), valueContains)
		}
	})
	return tc
}

// Error asserts an error.
func (tc *CacheStoreTestCase) Error(errContains string) *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Contains(t, tc.err.Error(), errContains)
		}
	})
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *CacheStoreTestCase) ErrorCode(statusCode int) *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		if assert.Error(t, tc.err) {
			assert.Equal(t, statusCode, errors.Convert(tc.err).StatusCode)
		}
	})
	return tc
}

// NoError asserts no error.
func (tc *CacheStoreTestCase) NoError() *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		assert.NoError(t, tc.err)
	})
	return tc
}

// Assert asserts using a provided function.
func (tc *CacheStoreTestCase) Assert(asserter func(t *testing.T, res *http.Response, err error)) *CacheStoreTestCase {
	tc.t.Run(tc.testName, func(t *testing.T) {
		asserter(t, tc.res, tc.err)
	})
	return tc
}

// Get returns the result of executing CacheStore.
func (tc *CacheStoreTestCase) Get() (res *http.Response, err error) {
	return tc.res, tc.err
}

// CacheStore executes the web handler and returns a corresponding test case.
func CacheStore(t *testing.T, ctx context.Context, options ...WebOption) *CacheStoreTestCase {
	tc := &CacheStoreTestCase{t: t}
	pubOptions := []pub.Option{
		pub.URL(httpx.JoinHostAndPath("messaging.example", `:443/cache-store`)),
	}
	for _, opt := range options {
		pubOptions = append(pubOptions, pub.Option(opt))
	}
	req, err := pub.NewRequest(pubOptions...)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		panic(err)
	}
	for name, value := range req.Header {
		httpReq.Header[name] = value
	}
	r := httpReq.WithContext(ctx)
	w := httpx.NewResponseRecorder()
	tc.err = utils.CatchPanic(func () error {
		return Svc.CacheStore(w, r)
	})
	tc.res = w.Result()
	return tc
}
