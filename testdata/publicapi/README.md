# Public API smoke test

This independent module exercises Intel as an external consumer. It is not
included in the repository's main `go test ./...` run. From this directory, run:

```sh
GOWORK=off GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...
```

The relative replacement intentionally tests the local Intel checkout.
