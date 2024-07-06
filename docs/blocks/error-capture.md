# Error Capture and Propagation

`Microbus` considers `error`s returned from microservice endpoints as exceptions, and as such it emphasizes their detection and disclosure. In addition to being logged as errors, these cross-microservice errors are also counted in Grafana and recorded in [distributed traces](../blocks/distrib-tracing.md) in Jaeger. Consequently, it is ill-advised to use errors to return state that should not raise the alarm.

The point of view of `Microbus` is that errors will happen, they will be unpredictable, they must never bring down the system and they should be observable and easily debuggable. With this in mind, the framework is taking an opinionated "throw and log" approach to standardize the capturing and surfacing of errors. Note that "capturing" does not mean "handling". The latter is left up to the app developer (or user).

## Web Handler Returns Error

The standard `http.HandlerFunc` signature in Go does not return an `error` but rather leaves it to the developer to set the status code (often to `500`) and output an error message to the body of the response. This results in repetitive "log and throw" error handling pattern that app developers may or may not consistently conform to.

```go
func StandardHandler(w http.ResponseWriter, r *http.Request) {
	err := doSomething()
	if err != nil {
		log.LogError(r.Context(), "doing something", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = doSomethingElse()
	if err != nil {
		log.LogError(r.Context(), "doing something else", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
```

`Microbus` takes the approach of extending the default web handler signature to also include an `error` return value. The error is then processed by the framework in a wrapper function that itself conforms to the standard `http.HandlerFunc` signature. With this pattern, the original code is simpler and more Go-like, and any errors are reported and surfaced in a consistent manner. Specifically, error are printed to the body of the response, logged, and metered so that alerts can be triggered.

```go
func UserCodeHandler(w http.ResponseWriter, r *http.Request) error { // Returning an error
	err := doSomething()
	if err != nil {
		return err
	}
	err = doSomethingElse()
	if err != nil {
		return err
	}
	return nil
}

func FrameworkWrapperOfHandler(w http.ResponseWriter, r *http.Request) {
	err := UserCodeHandler(w, r)
	if err != nil {
		// Standard error capture
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(fmt.Sprintf("%+v", err))
		log.LogError(r.Context(), "Handling request", err)
		metrics.IncrementErrorCount(1)
		return
	}
}
```

## Catching Panics

Panics in Go that are not captured terminate the running process. With microservices that are expected to be always on, this is not a desirable outcome. To stop this from happening, the framework captures panics generated by user code and converts them to standard errors, which are then handled in a similar fashion. The wrapper function introduced earlier comes in handy here.

```go
func UserCodeHandler(w http.ResponseWriter, r *http.Request) error {
	panic("omg")
}

func FrameworkWrapperOfHandler(w http.ResponseWriter, r *http.Request) {
	err := utils.CatchPanic(func() error {return UserCodeHandler(w, r)}) // Convert panics to errors
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(fmt.Sprintf("%+v", err))
		log.LogError(r.Context(), "Handling request", err)
		metrics.IncrementErrorCount(1)
		return
	}
}
```

## Stack Trace

Errors in Go are as simple as an error message and do not include a stack trace. In distributed systems that are built by distributed teams, that makes it very difficult to identify the root cause. To address that, the framework is replacing the standard `errors` package with its own implementation of [augmented errors](../structure/errors.md) with stack locations. This is not entirely transparent and app developers must use `errors.Trace` to capture the stack location.

```go
import "github.com/microbus-io/errors"

func UserCodeHandler(w http.ResponseWriter, r *http.Request) error {
	err := doSomething()
	if err != nil {
		return errors.Trace(err) // Capture this line into the error's stack trace
	}
	return nil
}
```

With this in place, error messages make it clear where the error originated from.

```
strconv.ParseInt: parsing "nan": invalid syntax
[400]

- calculator.(*Service).Square
  /src/github.com/microbus-io/fabric/examples/calculator/service.go:75
- connector.(*Connector).Publish
  /src/github.com/microbus-io/fabric/connector/messaging.go:94
- httpingress.(*Service).ServeHTTP
  /src/github.com/microbus-io/fabric/coreservices/httpingress/service.go:124
```

## Status Codes

All microservices are ultimately web servers where it is common practice to return an appropriate HTTP status code along with an error. To facilitate that, `Microbus` allows status codes to be associated with errors.

```go
if obj, ok := m[key]; !ok {
	return errors.Newc(http.StatusNotFound, "record not found") // New error with status code
}
err = doSomething(obj)
if err != nil {
	return errors.Tracec(http.StatusNotFound, err) // Wrap existing error and attach a status code
}
```

The web handler wrapper is extended to write the errors' status code to the HTTP response writer:

```go
func FrameworkWrapperOfHandler(w http.ResponseWriter, r *http.Request) {
	err := utils.CatchPanic(func() error {return UserCodeHandler(w, r)})
	if err != nil {
		w.WriteHeader(errors.StatusCode(err)) // Return the error's status code
		w.Write(fmt.Sprintf("%+v", err))
		log.LogError(r.Context(), "Handling request", err)
		metrics.IncrementErrorCount(1)
		return
	}
}
```

## Propagation Over the Wire

Microservices run in different processes, often on different hardware or even geographies. When a client microservice calls a remote microservice and the latter returns an error, it would be a lot more developer-friendly if the client experienced the error as if it were local.

```go
func (s *Service) MyEndpoint(w http.ResponseWriter, r *http.Request) error {
	response, err := s.Publish(pub.GET("https://another.service/objects")) // Remote call
	if err != nil {
		return errors.Trace(err)
	}
}
```

The `Microbus` framework serializes any error generated by the remote microservice and reconstitutes it on the client side. Errors responses are identified by the special header `Microbus-Op-Code: Err` with the error serialized in the body as JSON, including its stack trace. The status code of the error is reflected in the status code of the HTTP.

```
HTTP/1.1 500 Internal Server Error
Connection: close
Content-Type: application/json
Microbus-From-Host: beta.error.connector
Microbus-From-Id: k2slru4rof
Microbus-Msg-Id: bHh4yWGo
Microbus-Op-Code: Err

{
	"error": "it's really bad",
	"stack": [
		{
			"file": "/src/github.com/microbus-io/fabric/connector/messaging_test.go",
			"function": "connector.TestConnector_Error.func2",
			"line": 343
		},
		{
			"file": "/src/github.com/microbus-io/fabric/connector/messaging.go",
			"function": "connector.(*Connector).onRequest",
			"line": 225
		},
		{
			"file": "/src/github.com/microbus-io/fabric/connector/messaging.go",
			"function": "connector.(*Connector).onRequest",
			"line": 226
		}
	]
}
```

The web handler wrapper function now looks similar to the following:

```go
func FrameworkWrapperOfHandler(w http.ResponseWriter, r *http.Request) {
	err := utils.CatchPanic(func() error {return UserCodeHandler(w, r)})
	if err != nil {
		w.Header().Set("Microbus-Op-Code", "Err") // Mark the response as an error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errors.StatusCode(err))
		json.NewEncoder(w).Encode(err) // Marhsal the error as JSON
		log.LogError(r.Context(), "Handling request", err)
		metrics.IncrementErrorCount(1)
		return
	}
}
```

## Summary

The capturing of errors at the framework level improves system stability, observability and developer experience.

* All errors are logged and metered
* Panics are captured
* Stack traces help identify the root cause
* HTTP status codes can be associated with an error
* Errors are propagated over the wire, including their stack trace and status code