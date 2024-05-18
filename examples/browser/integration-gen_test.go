/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

package browser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	"github.com/microbus-io/fabric/examples/browser/browserapi"
)

var (
	_ bytes.Buffer
	_ context.Context
	_ fmt.Stringer
	_ io.Reader
	_ *http.Request
	_ *url.URL
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
	_ *browserapi.Client
)

var (
	sequence int
)

var (
	// App manages the lifecycle of the microservices used in the test
	App *application.Application
	// Svc is the browser.example microservice being tested
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

// BrowseTestCase assists in asserting against the results of executing Browse.
type BrowseTestCase struct {
	t *testing.T
	res *http.Response
	err error
	dur time.Duration
}

// StatusOK asserts no error and a status code 200.
func (tc *BrowseTestCase) StatusOK() *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Equal(tc.t, tc.res.StatusCode, http.StatusOK)
	}
	return tc
}

// StatusCode asserts no error and a status code.
func (tc *BrowseTestCase) StatusCode(statusCode int) *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Equal(tc.t, tc.res.StatusCode, statusCode)
	}
	return tc
}

// BodyContains asserts no error and that the response contains a string or byte array.
func (tc *BrowseTestCase) BodyContains(bodyContains any) *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		switch v := bodyContains.(type) {
		case []byte:
			assert.True(tc.t, bytes.Contains(body, v), `"%v" does not contain "%v"`, body, v)
		case string:
			assert.True(tc.t, bytes.Contains(body, []byte(v)), `"%s" does not contain "%s"`, string(body), v)
		default:
			vv := fmt.Sprintf("%v", v)
			assert.True(tc.t, bytes.Contains(body, []byte(vv)), `"%s" does not contain "%s"`, string(body), vv)
		}
	}
	return tc
}

// BodyNotContains asserts no error and that the response does not contain a string or byte array.
func (tc *BrowseTestCase) BodyNotContains(bodyNotContains any) *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		switch v := bodyNotContains.(type) {
		case []byte:
			assert.False(tc.t, bytes.Contains(body, v), `"%v" contains "%v"`, body, v)
		case string:
			assert.False(tc.t, bytes.Contains(body, []byte(v)), `"%s" contains "%s"`, string(body), v)
		default:
			vv := fmt.Sprintf("%v", v)
			assert.False(tc.t, bytes.Contains(body, []byte(vv)), `"%s" contains "%s"`, string(body), vv)
		}
	}
	return tc
}

// HeaderContains asserts no error and that the named header contains a string.
func (tc *BrowseTestCase) HeaderContains(headerName string, value string) *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.True(tc.t, strings.Contains(tc.res.Header.Get(headerName), value), `header "%s: %s" does not contain "%s"`, headerName, tc.res.Header.Get(headerName), value)
	}
	return tc
}

// HeaderNotContains asserts no error and that the named header does not contain a string.
func (tc *BrowseTestCase) HeaderNotContains(headerName string, value string) *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.False(tc.t, strings.Contains(tc.res.Header.Get(headerName), value), `header "%s: %s" contains "%s"`, headerName, tc.res.Header.Get(headerName), value)
	}
	return tc
}

// HeaderExists asserts no error and that the named header exists.
func (tc *BrowseTestCase) HeaderExists(headerName string) *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.NotZero(tc.t, len(tc.res.Header.Values(headerName)), `header "%s" does not exist`, headerName)
	}
	return tc
}

// HeaderNotExists asserts no error and that the named header exists.
func (tc *BrowseTestCase) HeaderNotExists(headerName string) *BrowseTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Zero(tc.t, len(tc.res.Header.Values(headerName)), `header "%s" exists`, headerName)
	}
	return tc
}

// Error asserts an error.
func (tc *BrowseTestCase) Error(errContains string) *BrowseTestCase {
	if assert.Error(tc.t, tc.err) {
		assert.Contains(tc.t, tc.err.Error(), errContains)
	}
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *BrowseTestCase) ErrorCode(statusCode int) *BrowseTestCase {
	if assert.Error(tc.t, tc.err) {
		assert.Equal(tc.t, statusCode, errors.Convert(tc.err).StatusCode)
	}
	return tc
}

// NoError asserts no error.
func (tc *BrowseTestCase) NoError() *BrowseTestCase {
	assert.NoError(tc.t, tc.err)
	return tc
}

// CompletedIn checks that the duration of the operation is less than or equal the threshold.
func (tc *BrowseTestCase) CompletedIn(threshold time.Duration) *BrowseTestCase {
	assert.LessOrEqual(tc.t, tc.dur, threshold)
	return tc
}

// Assert asserts using a provided function.
func (tc *BrowseTestCase) Assert(asserter func(t *testing.T, res *http.Response, err error)) *BrowseTestCase {
	asserter(tc.t, tc.res, tc.err)
	return tc
}

// Get returns the result of executing Browse.
func (tc *BrowseTestCase) Get() (res *http.Response, err error) {
	return tc.res, tc.err
}

/*
BrowseGet performs a GET request to the Browse endpoint.

Browser shows a simple address bar and the source code of a URL.

If a URL is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
*/
func BrowseGet(t *testing.T, ctx context.Context, url string) *BrowseTestCase {
	tc := &BrowseTestCase{t: t}
	var err error
	url, err = httpx.ResolveURL(browserapi.URLOfBrowse, url)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	r, err := http.NewRequest(`GET`, url, nil)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	ctx = context.WithValue(ctx, frame.ContextKey, r.Header)
	w := httpx.NewResponseRecorder()
	t0 := time.Now()
	tc.err = utils.CatchPanic(func() error {
		return Svc.Browse(w, r.WithContext(ctx))
	})
	tc.dur = time.Since(t0)
	tc.res = w.Result()
	return tc
}

/*
BrowsePost performs a POST request to the Browse endpoint.

Browser shows a simple address bar and the source code of a URL.

If a URL is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
If the body if of type io.Reader, []byte or string, it is serialized in binary form.
If it is of type url.Values, it is serialized as form data. All other types are serialized as JSON.
If a content type is not explicitly provided, an attempt will be made to derive it from the body.
*/
func BrowsePost(t *testing.T, ctx context.Context, url string, contentType string, body any) *BrowseTestCase {
	tc := &BrowseTestCase{t: t}
	var err error
	url, err = httpx.ResolveURL(browserapi.URLOfBrowse, url)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	r, err := httpx.NewRequest(`POST`, url, body)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	ctx = context.WithValue(ctx, frame.ContextKey, r.Header)
	w := httpx.NewResponseRecorder()
	t0 := time.Now()
	tc.err = utils.CatchPanic(func() error {
		return Svc.Browse(w, r.WithContext(ctx))
	})
	tc.dur = time.Since(t0)
	tc.res = w.Result()
	return tc
}

/*
Browser shows a simple address bar and the source code of a URL.

If a request is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
*/
func Browse(t *testing.T, ctx context.Context, r *http.Request) *BrowseTestCase {
	tc := &BrowseTestCase{t: t}
	var err error
	if r == nil {
		r, err = http.NewRequest(`GET`, "", nil)
		if err != nil {
			tc.err = errors.Trace(err)
			return tc
		}
	}
	u, err := httpx.ResolveURL(browserapi.URLOfBrowse, r.URL.String())
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	r.URL, err = url.Parse(u)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	ctx = context.WithValue(ctx, frame.ContextKey, r.Header)
	w := httpx.NewResponseRecorder()
	t0 := time.Now()
	tc.err = utils.CatchPanic(func() error {
		return Svc.Browse(w, r.WithContext(ctx))
	})
	tc.res = w.Result()
	tc.dur = time.Since(t0)
	return tc
}
