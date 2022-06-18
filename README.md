# config
Golang multi-source configuration module

## Contributing

In order to get started please make sure to have golang installed.
Simplest is to use [gvm](https://github.com/moovweb/gvm) and [direnv](https://github.com/direnv/direnv).
Run `$(grep "^go " go.mod | awk '{print $2}')` to check a version of golang used.

Install required version of golang:
`gvm install $(grep "^go " go.mod | awk '{print $2}')`

Install dev dependencies:
`make tools`

Run tests:
```bash
# all tests
make test

# specific test
go test -v ./val/ --run TestValue

# specific test with watch
# more on gow: https://github.com/mitranim/gow
gow test -v ./val/ --run TestValue
```