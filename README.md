# ChewCrew Core

## Getting Started

* Install Go using [gvm](https://github.com/moovweb/gvm)
```bash
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.3.3
```

* Install chewcrew
```bash
gvm use go1.3.3
git clone git://github.com/wafflehaus/chewcrew.git $GOPATH/src/github.com/wafflehaus/chewcrew
go get github.com/wafflehaus/chewcrew
go install github.com/wafflehaus/chewcrew
chewcrew
```

## Usage

For server configurables, run
```bash
chewcrew -help
```

### Static Assets

Static assets are served on the root path. The absolute path of the static assets directory is available as a server configurable.

### API

API routes are found on the path `/api/*`. Consult the `api.raml` file for API specification.

## Development

The server is highly dependent on the following dependencies, so it would be best to be somewhat familiar with them:
* [gorillatoolkit:context](http://www.gorillatoolkit.org/pkg/context)
  * used to decorate the request with universal facilities like a request logger
* [httprouter](https://github.com/julienschmidt/httprouter)
  * routing
  * uri and query string access

### Before Merging

Before merging, format and test your code.

```bash
go fmt github.com/wafflehaus/chewcrew
go test github.com/wafflehaus/chewcrew
```

### Logging

Server logging is available as `ServerLog`.

Each incoming API request is provided its own child logger, namespaced under the request's id. This child logger is available on the request's context. See gorilla toolkit documentation for more information on [contexts](http://www.gorillatoolkit.org/pkg/context).

Services logging is TBD.

Consult `Logger` for more information on the logger interface. `Logger` provides a wrapper around [go-logging](https://github.com/op/go-logging).

### Errors

Always handle and try to recover from errors! Panicing should be used sparingly (if at all).

Errors originating from core or third party libraries should be wrapped via `NewMaskedError`, `NewMaskedErrorWithContext`, or `NewMaskedErrorWithContextf`. Errors originating within `main` should be initialized via `NewError` or `NewErrorf`. When creating errors via these functions, a call stack (and optional context message) is associated with the error.

When logging an error directly with a `Logger`, the call stack (and optional context message) will be logged alongside the error's root cause.

Custom error types can be defined as so:
```go
type MyCustomError Error
```

Custom error types encourage maintainable error handling when multiple errors could be returned from a function that require different recovery paths.
```go
err := TestFunc()
switch err.(type) {
  MyCustomError:
    // handle MyCustomError
  MyCustomError2:
    // handle MyCustomerError2
}
```
