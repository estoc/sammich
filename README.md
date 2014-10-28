# chewcrew server

The backend REST API for chewcrew

## Getting Started
* Get the code
```bash
git clone git://github.com/wafflehaus/chewcrew.git $GOPATH/src/github.com/wafflehaus/chewcrew
```

* Get dependencies
```bash
$GOPATH/src/github.com/wafflehaus/chewcrew/bin/getimports.sh
```
* Compile the code and start the server!
```bash
go install github.com/wafflehaus/chewcrew && chewcrew
```

## Installing Go
Use [gvm](https://github.com/moovweb/gvm), or install Golang (sample script below for Ubuntu 64bit). Installing from gvm is advised due to the simplicity.

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
go install github.com/wafflehaus/chewcrew && $GOPATH/bin/chewcrew -help
```

### Static Assets

Static assets are served on the root path. The absolute path of the static assets directory is available as a server configurable.

### API

API routes are found on the path `/api/*`. Consult the `api.raml` file for API specification.

## Running Tests

```bash
go test github.com/wafflehaus/chewcrew
```

## Documentation

Hosted documentation for master is found [here](http://godoc.org/github.com/wafflehaus/chewcrew).

Or, you can host your documentation locally:
```bash
godoc -http=:6060
```
Open [http://localhost:6060/pkg/github.com/wafflehaus/chewcrew](http://localhost:6060/pkg/github.com/wafflehaus/chewcrew) in your browser.

## Before Merging

Before merging, it is strongly advised that you `go fmt` your fork.

```bash
go fmt github.com/wafflehaus/chewcrew
```

## Development Notes

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
