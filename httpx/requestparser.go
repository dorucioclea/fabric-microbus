/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/microbus-io/fabric/errors"
)

// ParseRequestBody parses the body of an incoming request and populates the fields of a data object.
// It supports JSON and URL-encoded form data content types.
// Use json tags to designate the name of the argument to map to each field.
func ParseRequestBody(r *http.Request, data any) error {
	// Parse JSON in the body
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return errors.Trace(err)
		}
	}
	// Parse form in body
	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			return errors.Trace(err)
		}
		err = DecodeDeepObject(r.PostForm, data)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// ParseRequestData parses the body and query arguments of an incoming request
// and populates the fields of a data object.
// Use json tags to designate the name of the argument to map to each field.
// An argument name can be hierarchical using either notation "a[b][c]" or "a.b.c",
// in which case it is read into the corresponding nested field.
func ParseRequestData(r *http.Request, data any) error {
	err := ParseRequestBody(r, data)
	if err == nil {
		err = DecodeDeepObject(r.URL.Query(), data)
	}
	return errors.Trace(err)
}
