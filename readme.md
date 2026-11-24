# Intel

Intel is a Go reflection library for inspecting types and values through
checked operations. Look up fields and tags, resolve paths on live values,
or compile a caller-defined struct shape into a reusable declaration.

## What it does

- Inspects exact types, method sets, and struct members, including promoted
  fields, tags, and full index paths.
- Narrows captured types and values through checked steps and reaches live
  field or slice storage when the supplied value permits it.
- Recognizes and validates struct members with caller-defined rules and
  compiles definitions into typed declarations.

## Why it exists

A field can be visible on a type but unreachable on a particular value. A
method set can satisfy an interface without proving a call will succeed.
Intel keeps those questions separate, so you can build on Go's reflection
rules without wiring up the same checks yourself.

## Core idea

- Capture a type for structural questions or a value when an instance matters.
- Check `OK()` before using a result returned by a checked operation.
- Use the inspection and capture operations directly, or build on them with
  rules and definitions. Your application decides what the structure means.

## Install

In a Go module using Go 1.27 or newer:

```bash
go get github.com/qlustra/intel@v1.0.0
```

## Quick start

Find a field promoted through an embedded struct and print its name, Go type,
and JSON tag. Save this as `main.go` in your module:

```go
package main

import (
	"fmt"
	"log"

	"github.com/qlustra/intel"
)

type Profile struct {
	Name string `json:"name"`
}

type User struct{ Profile }

func main() {
	user, check := intel.CaptureType[User]().AsStruct()
	if !check.OK() {
		log.Fatal(check.Why())
	}
	field, check := user.Field("Name")
	if !check.OK() {
		log.Fatal(check.Why())
	}
	fmt.Println(field.FieldName(), field.ReflectedType(), field.StructField().Tag.Get("json"))
}
```

```sh
go run main.go
```

```text
Name string name
```

## Project shape

- `capture/`
- `inspect/`
- `rules/`
- `define/`
- `compile/`

## Status

This README targets Intel's first public release, `v1.0.0`.

## Links

- [Usage manual](https://github.com/Qlustra/intel-docs/blob/main/docs/manual/readme.md)
- [Runnable examples](https://github.com/Qlustra/intel-docs/blob/main/docs/examples/readme.md)
- [v1 guarantees](https://github.com/Qlustra/intel-docs/blob/main/docs/guarantees.md)

## License

[Apache License 2.0](LICENSE) · [Notice](NOTICE)

Intel is maintained within the Qlustra Engineering Tooling Lab.
https://github.com/Qlustra
