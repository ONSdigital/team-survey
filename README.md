# team-survey

A team cohesion survey generator/manager

## Getting started

### Environment

| Var           | Default         | Example  | Description                                      |
| ------------- | --------------- | -------- | ------------------------------------------------ |
| DISABLE_AUTH  | _&lt;unset&gt;_ | `True`   | Set to `True` to disable authentication checking |

### Running with vanilla go

```sh
cd cmd/webservice
go run ./...
```

## Question sets

To run this service, you need to supply the appropriate question sets. They are omitted from this repo as we do not hold the license to publish them. However, if you have your own access to them you can insert them yourself.

**TODO** Note on what question sets are required and where they are to be inserted

## Tests

Several sets of tests are included:

### Unit Tests

Use the standard `go` testing command sets, e.g.:

```sh
go test -v ./...
```

### UI Testing

Use `cypress` to run the test suite. (**Note**: For these you need to have the `nodejs` based `[npx](https://www.npmjs.com/package/npx)` installed)

The tests can be found under `cmd/webservice/web/test` and run with:

```sh
npx cypress open
```

## License

Copyright (c) 2019 Crown Copyright (Office for National Statistics)

Released under MIT license, see [LICENSE](LICENSE) for details.
