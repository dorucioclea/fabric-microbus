/*
Copyright (c) 2023 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package eventsource

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/microbus-io/fabric/examples/eventsink"
	"github.com/microbus-io/fabric/examples/eventsource/eventsourceapi"
)

var (
	_ *testing.T
	_ assert.TestingT
	_ *eventsourceapi.Client
)

// Initialize starts up the testing app.
func Initialize() error {
	// Include all downstream microservices in the testing app
	// Use .With(...) to initialize with appropriate config values
	App.Include(
		Svc,
		eventsink.NewService(),
	)

	err := App.Startup()
	if err != nil {
		return err
	}

	// You may call any of the microservices after the app is started

	return nil
}

// Terminate shuts down the testing app.
func Terminate() error {
	err := App.Shutdown()
	if err != nil {
		return err
	}
	return nil
}

func TestEventsource_Register(t *testing.T) {
	t.Parallel()
	/*
		Register(t, ctx, email).
			Name(testName).
			Expect(allowed).
			NoError().
			Error(errContains).
			Assert(func(t, allowed, err))
	*/
	ctx := Context(t)
	Register(t, ctx, "brian@hotmail.com").Name("decline hotmail.com").Expect(false)
	Register(t, ctx, "brian@example.com").Name("accept example.com").Expect(true)
	Register(t, ctx, "brian@example.com").Name("decline dup").Expect(false)
}

func TestEventsource_OnAllowRegister(t *testing.T) {
	// No parallel
	/*
		OnAllowRegister(t, allow, err).
			Name(testName).
			Expect(email).
			Assert(func(t, ctx, email))
	*/
	ctx := Context(t)
	OnAllowRegister(t, true, nil).
		Expect("barb@example.com")
	Register(t, ctx, "barb@example.com").Expect(true)
	OnAllowRegister(t, false, nil).
		Expect("josh@example.com")
	Register(t, ctx, "josh@example.com").Expect(false)
}

func TestEventsource_OnRegistered(t *testing.T) {
	// No parallel
	/*
		OnRegistered(t, err).
			Name(testName).
			Expect(email).
			Assert(func(t, ctx, email))
	*/
	ctx := Context(t)
	OnRegistered(t, nil).Expect("harry@example.com")
	Register(t, ctx, "harry@example.com").Expect(true)
}
