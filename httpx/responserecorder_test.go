/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package httpx

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/microbus-io/fabric/frame"
	"github.com/microbus-io/fabric/rand"
	"github.com/stretchr/testify/assert"
)

func TestHttpx_ResponseRecorder(t *testing.T) {
	rr := NewResponseRecorder()

	// Write once
	rr.Header().Set("Foo", "Bar")
	rr.WriteHeader(http.StatusTeapot)

	bin := []byte("Lorem Ipsum")
	n, err := rr.Write(bin)
	assert.NoError(t, err)
	assert.Equal(t, len(bin), n)

	result := rr.Result()
	assert.Equal(t, bin, result.Body.(*BodyReader).Bytes())

	var buf bytes.Buffer
	err = result.Write(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 418 I'm a teapot\r\nContent-Length: 11\r\nFoo: Bar\r\n\r\nLorem Ipsum", buf.String())

	// Write second time
	rr.Header().Set("Foo", "Baz")
	rr.WriteHeader(http.StatusConflict)

	bin2 := []byte(" Dolor Sit Amet")
	n, err = rr.Write(bin2)
	assert.NoError(t, err)
	assert.Equal(t, len(bin2), n)
	bin = append(bin, bin2...)

	result = rr.Result()
	assert.Equal(t, bin, result.Body.(*BodyReader).Bytes())

	buf.Reset()
	err = result.Write(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 409 Conflict\r\nContent-Length: 26\r\nFoo: Baz\r\n\r\nLorem Ipsum Dolor Sit Amet", buf.String())
}

func TestHttpx_FrameOfResponseRecorder(t *testing.T) {
	utilsRecorder := NewResponseRecorder()
	utilsRecorder.Header().Set(frame.HeaderMsgId, "123")
	assert.Equal(t, "123", frame.Of(utilsRecorder).MessageID())
	httpResponse := utilsRecorder.Result()
	assert.Equal(t, "123", frame.Of(httpResponse).MessageID())
}

func TestHttpx_Copy(t *testing.T) {
	payload := rand.AlphaNum64(256 * 1024)

	recorder := NewResponseRecorder()
	recorder.Write([]byte(payload))
	b, err := io.ReadAll(recorder.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, payload, string(b))

	recorder = NewResponseRecorder()
	n, err := io.Copy(recorder, io.LimitReader(strings.NewReader(payload), int64(len(payload))))
	assert.NoError(t, err)
	assert.Equal(t, int(n), len(payload))
	b, err = io.ReadAll(recorder.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, payload, string(b))
}
