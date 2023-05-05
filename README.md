# config
Golang multi-source configuration module

## Contributing

In order to get started please make sure to have golang of required version installed.

Recommended is to use [gobrew](https://github.com/kevincobain2000/gobrew). It is also recommended to have [direnv](https://github.com/direnv/direnv) installed.

Run `$(grep "^go " go.mod | awk '{print $2}')` to check a version of golang needed.

Install required version of golang:
`gobrew install $(grep "^go " go.mod | awk '{print $2}')@latest`

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