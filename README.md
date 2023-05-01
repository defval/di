# DI

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/goava/di)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/goava/di?logo=semver&style=for-the-badge)](https://github.com/goava/di/releases/latest)
[![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/goava/di/go.yml?branch=master&logo=github-actions&style=for-the-badge)](https://github.com/goava/di/actions/workflows/go.yml)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge)](https://goreportcard.com/report/github.com/goava/di)
[![Codecov](https://img.shields.io/codecov/c/github/goava/di?logo=codecov&style=for-the-badge)](https://codecov.io/gh/goava/di)

**DI** is a dependency injection library for the Go programming language.

- [Tutorial](./docs/tutorial.md)
- [Examples](./_examples)
- [Advanced Features](./docs/advanced.md)

Dependency injection is a form of inversion of control that increases modularity and extensibility in your programs.
This library helps you organize responsibilities in your codebase and makes it easy to combine low-level implementations
into high-level behavior without boilerplate.

## Features

- Intuitive auto wiring
- Interface implementations
- Constructor injection
- Optional injection
- Field injection
- Lazy-loading
- Tagging
- Grouping
- Cleanup
- Container Chaining / Scopes

## Documentation

You can use the standard [pkg.go.dev](https://pkg.go.dev/github.com/goava/di) and inline code comments. If you are new
to auto-wiring libraries such as [google/wire](https://github.com/google/wire)
or [uber-go/dig](https://github.com/uber-go/dig), start with the [tutorial](./docs/tutorial.md).

## Installation

```shell
go get github.com/goava/di
```

## Example Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/goava/di"
)

func main() {
	di.SetTracer(&di.StdTracer{})
	// create container
	c, err := di.New(
		di.Provide(NewContext),  // provide application context
		di.Provide(NewServer),   // provide http server
		di.Provide(NewServeMux), // provide http serve mux
		// controllers as []Controller group
		di.Provide(NewOrderController, di.As(new(Controller))),
		di.Provide(NewUserController, di.As(new(Controller))),
	)
	// handle container errors
	if err != nil {
		log.Fatal(err)
	}
	// invoke function
	if err := c.Invoke(StartServer); err != nil {
		log.Fatal(err)
	}
}
```

Full code available [here](./_examples/tutorial/main.go).

## Questions

If you have any questions, feel free to create an issue.
