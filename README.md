DI
============
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://pkg.go.dev/github.com/goava/di)
[![Release](https://img.shields.io/github/tag/goava/di.svg?label=release&color=24B898&logo=github&style=for-the-badge)](https://github.com/goava/di/releases/latest)
[![Build Status](https://img.shields.io/travis/goava/di.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/goava/di)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge)](https://goreportcard.com/report/github.com/goava/di)
[![Code Coverage](https://img.shields.io/codecov/c/github/goava/di.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/goava/di)

Dependency injection for Go programming language.

[Tutorial](./docs/tutorial.md) | [Examples](./_examples) | [Advanced features](./docs/advanced.md)

Dependency injection is one form of the broader technique of inversion
of control. It is used to increase modularity of the program and make it
extensible. 

This library helps you to organize responsibilities in your codebase and make it 
easy to combine low-level implementation into high-level behavior without boilerplate.

## Features

- Intuitive auto wiring
- Interface implementations
- Constructor injection
- Optional injection
- Field injection
- Lazy-loading
- Named types
- Grouping
- Cleanup

## Documentation

You can use standard [pkg.go.dev](https://pkg.go.dev/github.com/goava/di) and inline code
comments. If you do not have experience with auto-wiring libraries
as [google/wire](https://github.com/google/wire),
[uber-go/dig](https://github.com/uber-go/dig) or another - start with
[tutorial](./docs/tutorial.md).

## Install

```shell
go get github.com/goava/di
```

## Questions

If you have any questions, feel free to create an issue.