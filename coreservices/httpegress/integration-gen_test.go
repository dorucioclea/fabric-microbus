/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

package httpegress

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

	"github.com/andybalholm/cascadia"
	"github.com/microbus-io/fabric/application"
	"github.com/microbus-io/fabric/connector"
	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/frame"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/pub"
	"github.com/microbus-io/fabric/rand"
	"github.com/microbus-io/fabric/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/microbus-io/fabric/coreservices/httpegress/httpegressapi"
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
	_ cascadia.Sel
	_ *connector.Connector
	_ *errors.TracedError
	_ frame.Frame
	_ *httpx.BodyReader
	_ pub.Option
	_ rand.Void
	_ utils.InfiniteChan[int]
	_ assert.TestingT
	_ *html.Node
	_ *httpegressapi.Client
)

var (
	// App manages the lifecycle of the microservices used in the test
	App *application.Application
	// Svc is the http.egress.sys microservice being tested
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
	ctx := context.Background()
	if deadline, ok := t.Deadline(); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		t.Cleanup(cancel)
	}
	ctx = frame.CloneContext(ctx)
	return ctx
}

// MakeRequestTestCase assists in asserting against the results of executing MakeRequest.
type MakeRequestTestCase struct {
	t *testing.T
	dur time.Duration
	res *http.Response
	err error
}

// StatusOK asserts no error and a status code 200.
func (tc *MakeRequestTestCase) StatusOK() *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Equal(tc.t, tc.res.StatusCode, http.StatusOK)
	}
	return tc
}

// StatusCode asserts no error and a status code.
func (tc *MakeRequestTestCase) StatusCode(statusCode int) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Equal(tc.t, tc.res.StatusCode, statusCode)
	}
	return tc
}

// BodyContains asserts no error and that the response contains the string or byte array value.
func (tc *MakeRequestTestCase) BodyContains(value any) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		switch v := value.(type) {
		case []byte:
			assert.True(tc.t, bytes.Contains(body, v), "%v does not contain %v", body, v)
		case string:
			assert.Contains(tc.t, string(body), v)
		default:
			vv := fmt.Sprintf("%v", v)
			assert.Contains(tc.t, string(body), vv)
		}
	}
	return tc
}

// BodyNotContains asserts no error and that the response does not contain the string or byte array value.
func (tc *MakeRequestTestCase) BodyNotContains(value any) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		switch v := value.(type) {
		case []byte:
			assert.False(tc.t, bytes.Contains(body, v), "%v contains %v", body, v)
		case string:
			assert.NotContains(tc.t, string(body), v)
		default:
			vv := fmt.Sprintf("%v", v)
			assert.NotContains(tc.t, string(body), vv)
		}
	}
	return tc
}

// HeaderContains asserts no error and that the named header contains the value.
func (tc *MakeRequestTestCase) HeaderContains(headerName string, value string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Contains(tc.t, tc.res.Header.Get(headerName), value)
	}
	return tc
}

// HeaderNotContains asserts no error and that the named header does not contain a string.
func (tc *MakeRequestTestCase) HeaderNotContains(headerName string, value string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.NotContains(tc.t, tc.res.Header.Get(headerName), value)
	}
	return tc
}

// HeaderEqual asserts no error and that the named header matches the value.
func (tc *MakeRequestTestCase) HeaderEqual(headerName string, value string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Equal(tc.t, value, tc.res.Header.Get(headerName))
	}
	return tc
}

// HeaderNotEqual asserts no error and that the named header does not matche the value.
func (tc *MakeRequestTestCase) HeaderNotEqual(headerName string, value string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.NotEqual(tc.t, value, tc.res.Header.Get(headerName))
	}
	return tc
}

// HeaderExists asserts no error and that the named header exists.
func (tc *MakeRequestTestCase) HeaderExists(headerName string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.NotEmpty(tc.t, tc.res.Header.Values(headerName), "Header %s does not exist", headerName)
	}
	return tc
}

// HeaderNotExists asserts no error and that the named header does not exists.
func (tc *MakeRequestTestCase) HeaderNotExists(headerName string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Empty(tc.t, tc.res.Header.Values(headerName), "Header %s exists", headerName)
	}
	return tc
}

// ContentType asserts no error and that the Content-Type header matches the expected value.
func (tc *MakeRequestTestCase) ContentType(expected string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		assert.Equal(tc.t, expected, tc.res.Header.Get("Content-Type"))
	}
	return tc
}

/*
TagExists asserts no error and that the at least one tag matches the CSS selector query.

Examples:

	TagExists(`TR > TD > A.expandable[href]`)
	TagExists(`DIV#main_panel`)
	TagExists(`TR TD INPUT[name="x"]`)
*/
func (tc *MakeRequestTestCase) TagExists(cssSelectorQuery string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		selector, err := cascadia.Compile(cssSelectorQuery)
		if !assert.NoError(tc.t, err, "Invalid selector %s", cssSelectorQuery) {
			return tc
		}
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		doc, err := html.Parse(bytes.NewReader(body))
		if !assert.NoError(tc.t, err, "Failed to parse HTML") {
			return tc
		}
		matches := selector.MatchAll(doc)
		assert.NotEmpty(tc.t, matches, "Found no tags matching %s", cssSelectorQuery)
	}
	return tc
}

/*
TagNotExists asserts no error and that the no tag matches the CSS selector query.

Example:

	TagNotExists(`TR > TD > A.expandable[href]`)
	TagNotExists(`DIV#main_panel`)
	TagNotExists(`TR TD INPUT[name="x"]`)
*/
func (tc *MakeRequestTestCase) TagNotExists(cssSelectorQuery string) *MakeRequestTestCase {
	if assert.NoError(tc.t, tc.err) {
		selector, err := cascadia.Compile(cssSelectorQuery)
		if !assert.NoError(tc.t, err, "Invalid selector %s", cssSelectorQuery) {
			return tc
		}
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		doc, err := html.Parse(bytes.NewReader(body))
		if !assert.NoError(tc.t, err, "Failed to parse HTML") {
			return tc
		}
		matches := selector.MatchAll(doc)
		assert.Empty(tc.t, matches, "Found %d tag(s) matching %s", len(matches), cssSelectorQuery)
	}
	return tc
}

/*
TagEqual asserts no error and that the at least one of the tags matching the CSS selector query
either contains the exact text itself or has a descendant that does.

Example:

	TagEqual("TR > TD > A.expandable[href]", "Expand")
	TagEqual("DIV#main_panel > SELECT > OPTION", "Red")
*/
func (tc *MakeRequestTestCase) TagEqual(cssSelectorQuery string, value string) *MakeRequestTestCase {
	var textMatches func(n *html.Node) bool
	textMatches = func(n *html.Node) bool {
		for x := n.FirstChild; x != nil; x = x.NextSibling {
			if x.Data == value || textMatches(x) {
				return true
			}
		}
		return false
	}

	if assert.NoError(tc.t, tc.err) {
		selector, err := cascadia.Compile(cssSelectorQuery)
		if !assert.NoError(tc.t, err, "Invalid selector %s", cssSelectorQuery) {
			return tc
		}
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		doc, err := html.Parse(bytes.NewReader(body))
		if !assert.NoError(tc.t, err, "Failed to parse HTML") {
			return tc
		}
		matches := selector.MatchAll(doc)
		if !assert.NotEmpty(tc.t, matches, "Selector %s does not match any tags", cssSelectorQuery) {
			return tc
		}
		if value == "" {
			return tc
		}
		found := false
		for _, match := range matches {
			if textMatches(match) {
				found = true
				break
			}
		}
		assert.True(tc.t, found, "No tag matching %s contains %s", cssSelectorQuery, value)
	}
	return tc
}

/*
TagContains asserts no error and that the at least one of the tags matching the CSS selector query
either contains the text itself or has a descendant that does.

Example:

	TagContains("TR > TD > A.expandable[href]", "Expand")
	TagContains("DIV#main_panel > SELECT > OPTION", "Red")
*/
func (tc *MakeRequestTestCase) TagContains(cssSelectorQuery string, value string) *MakeRequestTestCase {
	var textMatches func(n *html.Node) bool
	textMatches = func(n *html.Node) bool {
		for x := n.FirstChild; x != nil; x = x.NextSibling {
			if strings.Contains(x.Data, value) || textMatches(x) {
				return true
			}
		}
		return false
	}

	if assert.NoError(tc.t, tc.err) {
		selector, err := cascadia.Compile(cssSelectorQuery)
		if !assert.NoError(tc.t, err, "Invalid selector %s", cssSelectorQuery) {
			return tc
		}
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		doc, err := html.Parse(bytes.NewReader(body))
		if !assert.NoError(tc.t, err, "Failed to parse HTML") {
			return tc
		}
		matches := selector.MatchAll(doc)
		if !assert.NotEmpty(tc.t, matches, "Selector %s does not match any tags", cssSelectorQuery) {
			return tc
		}
		if value == "" {
			return tc
		}
		found := false
		for _, match := range matches {
			if textMatches(match) {
				found = true
				break
			}
		}
		assert.True(tc.t, found, "No tag matching %s contains %s", cssSelectorQuery, value)
	}
	return tc
}

/*
TagNotEqual asserts no error and that there is no tag matching the CSS selector that
either contains the exact text itself or has a descendant that does.

Example:

	TagNotEqual("TR > TD > A[href]", "Harry Potter")
	TagNotEqual("DIV#main_panel > SELECT > OPTION", "Red")
*/
func (tc *MakeRequestTestCase) TagNotEqual(cssSelectorQuery string, value string) *MakeRequestTestCase {
	var textMatches func(n *html.Node) bool
	textMatches = func(n *html.Node) bool {
		for x := n.FirstChild; x != nil; x = x.NextSibling {
			if x.Data == value || textMatches(x) {
				return true
			}
		}
		return false
	}

	if assert.NoError(tc.t, tc.err) {
		selector, err := cascadia.Compile(cssSelectorQuery)
		if !assert.NoError(tc.t, err, "Invalid selector %s", cssSelectorQuery) {
			return tc
		}
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		doc, err := html.Parse(bytes.NewReader(body))
		if !assert.NoError(tc.t, err, "Failed to parse HTML") {
			return tc
		}
		matches := selector.MatchAll(doc)
		if len(matches) == 0 {
			return tc
		}
		if !assert.NotEmpty(tc.t, value, "Found tag matching %s", cssSelectorQuery) {
			return tc
		}
		found := false
		for _, match := range matches {
			if textMatches(match) {
				found = true
				break
			}
		}
		assert.False(tc.t, found, "Found tag matching %s that contains %s", cssSelectorQuery, value)
	}
	return tc
}

/*
TagNotContains asserts no error and that there is no tag matching the CSS selector that
either contains the text itself or has a descendant that does.

Example:

	TagNotContains("TR > TD > A[href]", "Harry Potter")
	TagNotContains("DIV#main_panel > SELECT > OPTION", "Red")
*/
func (tc *MakeRequestTestCase) TagNotContains(cssSelectorQuery string, value string) *MakeRequestTestCase {
	var textMatches func(n *html.Node) bool
	textMatches = func(n *html.Node) bool {
		for x := n.FirstChild; x != nil; x = x.NextSibling {
			if strings.Contains(x.Data, value) || textMatches(x) {
				return true
			}
		}
		return false
	}

	if assert.NoError(tc.t, tc.err) {
		selector, err := cascadia.Compile(cssSelectorQuery)
		if !assert.NoError(tc.t, err, "Invalid selector %s", cssSelectorQuery) {
			return tc
		}
		body := tc.res.Body.(*httpx.BodyReader).Bytes()
		doc, err := html.Parse(bytes.NewReader(body))
		if !assert.NoError(tc.t, err, "Failed to parse HTML") {
			return tc
		}
		matches := selector.MatchAll(doc)
		if len(matches) == 0 {
			return tc
		}
		if !assert.NotEmpty(tc.t, value, "Found tag matching %s", cssSelectorQuery) {
			return tc
		}
		found := false
		for _, match := range matches {
			if textMatches(match) {
				found = true
				break
			}
		}
		assert.False(tc.t, found, "Found tag matching %s that contains %s", cssSelectorQuery, value)
	}
	return tc
}

// Error asserts an error.
func (tc *MakeRequestTestCase) Error(errContains string) *MakeRequestTestCase {
	if assert.Error(tc.t, tc.err) {
		assert.Contains(tc.t, tc.err.Error(), errContains)
	}
	return tc
}

// ErrorCode asserts an error by its status code.
func (tc *MakeRequestTestCase) ErrorCode(statusCode int) *MakeRequestTestCase {
	if assert.Error(tc.t, tc.err) {
		assert.Equal(tc.t, statusCode, errors.Convert(tc.err).StatusCode)
	}
	return tc
}

// NoError asserts no error.
func (tc *MakeRequestTestCase) NoError() *MakeRequestTestCase {
	assert.NoError(tc.t, tc.err)
	return tc
}

// CompletedIn checks that the duration of the operation is less than or equal the threshold.
func (tc *MakeRequestTestCase) CompletedIn(threshold time.Duration) *MakeRequestTestCase {
	assert.LessOrEqual(tc.t, tc.dur, threshold)
	return tc
}

// Assert asserts using a provided function.
func (tc *MakeRequestTestCase) Assert(asserter func(t *testing.T, res *http.Response, err error)) *MakeRequestTestCase {
	asserter(tc.t, tc.res, tc.err)
	return tc
}

// Get returns the result of executing MakeRequest.
func (tc *MakeRequestTestCase) Get() (res *http.Response, err error) {
	return tc.res, tc.err
}
/*
MakeRequest proxies a request to a URL and returns the HTTP response, respecting the timeout set in the context.
The proxied request is expected to be posted in the body of the request in binary form (RFC7231).

If a URL is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
If the body if of type io.Reader, []byte or string, it is serialized in binary form.
If it is of type url.Values, it is serialized as form data. All other types are serialized as JSON.
If a content type is not explicitly provided, an attempt will be made to derive it from the body.
*/
func MakeRequest(t *testing.T, ctx context.Context, url string, contentType string, body any) *MakeRequestTestCase {
	tc := &MakeRequestTestCase{t: t}
	var err error
	url, err = httpx.ResolveURL(httpegressapi.URLOfMakeRequest, url)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	r, err := httpx.NewRequest(`POST`, url, nil)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	ctx = frame.CloneContext(ctx)
	r = r.WithContext(ctx)
	r.Header = frame.Of(ctx).Header()
	err = httpx.SetRequestBody(r, body)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httpx.NewResponseRecorder()
	t0 := time.Now()
	tc.err = utils.CatchPanic(func() error {
		return Svc.MakeRequest(w, r)
	})
	tc.dur = time.Since(t0)
	tc.res = w.Result()
	return tc
}

/*
MakeRequestAny performs a customized request to the MakeRequest endpoint.

MakeRequest proxies a request to a URL and returns the HTTP response, respecting the timeout set in the context.
The proxied request is expected to be posted in the body of the request in binary form (RFC7231).

If a request is not provided, it defaults to the URL of the endpoint. Otherwise, it is resolved relative to the URL of the endpoint.
*/
func MakeRequestAny(t *testing.T, r *http.Request) *MakeRequestTestCase {
	tc := &MakeRequestTestCase{t: t}
	var err error
	if r == nil {
		r, err = http.NewRequest(`POST`, "", nil)
		if err != nil {
			tc.err = errors.Trace(err)
			return tc
		}
	}
	u, err := httpx.ResolveURL(httpegressapi.URLOfMakeRequest, r.URL.String())
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	r.URL, err = url.Parse(u)
	if err != nil {
		tc.err = errors.Trace(err)
		return tc
	}
	ctx := frame.ContextWithFrameOf(r.Context(), r.Header)
	r = r.WithContext(ctx)
	r.Header = frame.Of(ctx).Header()
	w := httpx.NewResponseRecorder()
	t0 := time.Now()
	tc.err = utils.CatchPanic(func() error {
		return Svc.MakeRequest(w, r)
	})
	tc.res = w.Result()
	tc.dur = time.Since(t0)
	return tc
}
