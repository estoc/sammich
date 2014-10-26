# chewcrew server

The backend REST API for chewcrew

## Getting Started
* Install Golang (sample script below for Ubuntu 64bit)
* Run ./bin/getimports.sh
* go build && ./chewcrew

```bash
### Go installation script
# Remove previous go installations
sudo rm -r /usr/local/go
cd /tmp
sudo rm -r go

# Download go (Update the version/OS here)
wget https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz
tar -zxf go1.3.3.linux-amd64.tar.gz

# Move go folder to default folder /usr/local
sudo mv go /usr/local/

# Add to path and .profile (comment out if updating)
export PATH=$PATH:/usr/local/go/bin
echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.profile
```

## Usage

For server configurables, execute the following:

```bash
./chewcrew -help
```

### Static Assets

Static assets are served on the root path. The absolute path of the static assets directory is available as a server configurable.

### API

API routes are found on the path `/api/*`. Consult the `api.raml` file for API specification.

## Development Notes

### Logging

Server logging is available to all of `main` as `serverLog`.

Each incoming API request is provided its own child logger, namespaced under the request's id. This child logger is available on the request's context. See gorilla toolkit documentation for more information on [contexts](http://www.gorillatoolkit.org/pkg/context).

Services logging is TBD.

Consult `log.go` for more information on the logger interface. `log.go` provides a wrapper around [go-logging](https://github.com/op/go-logging).

### Errors

Always handle and try to recover from errors! Panicing should be used sparingly (if at all).

Errors originating from core or third party libraries should be wrapped via `NewMaskedError`, `NewMaskedErrorWithContext`, or `NewMaskedErrorWithContextf`. Errors originating within `main` should be initialized via `NewError` or `NewErrorf`. When creating errors via these functions, a call stack (and optional context message) is associated with the error.

When logging an error directly with `log.go`, the call stack (and optional context message) will be logged alongside the error's root cause.

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

### Before Merging

Before merging, it is strongly advised that you `go fmt` your fork.

```bash
go fmt ~/workspace/chewcrew
```