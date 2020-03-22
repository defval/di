# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!--## [Unrealeased]-->

## v1.0.2

### Added

- Location of `di.Provide()`, `di.Invoke()`, `di.Resolve()` in error.

### Fixed

- Fix: `di.As()` with nil causes panic

## v1.0.1

### Fixed

- `container.Provide` could not be called after container compilation now.
- Improve error messages


## v1.0.0

Initial release.