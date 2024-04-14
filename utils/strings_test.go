/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_ToKebabCase(t *testing.T) {
	t.Parallel()

	testCases := map[string]string{
		"fooBar":     "foo-bar",
		"FooBar":     "foo-bar",
		"fooBAR":     "foo-bar",
		"FooBAR":     "foo-bar",
		"urlEncoder": "url-encoder",
		"URLEncoder": "url-encoder",
		"foobarX":    "foobar-x",
		"a":          "a",
		"A":          "a",
		"HTTP":       "http",
		"":           "",
	}
	for id, expected := range testCases {
		assert.Equal(t, expected, ToKebabCase(id), "%s", id)
	}
}
