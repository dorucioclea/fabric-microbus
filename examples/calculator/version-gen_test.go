/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

// Code generated by Microbus. DO NOT EDIT.

package calculator

import (
	"os"
	"testing"

	"github.com/microbus-io/fabric/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculator_Versioning(t *testing.T) {
	t.Parallel()
	
	hash, err := utils.SourceCodeSHA256(".")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, SourceCodeSHA256, "SourceCodeSHA256 is not up to date")
	}
	buf, err := os.ReadFile("version-gen.go")
	if assert.NoError(t, err) {
		assert.Contains(t, string(buf), hash, "SHA256 in version-gen.go is not up to date")
	}
}
