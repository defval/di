# Changelog

All notable changes to this project will be documented in this file.

The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this
project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html): TBD, use
modules or another vendor system.

## v1.12.0

### Changed

- Update package name

## v1.11.1

### Fixed

- [@chirino](https://github.com/chirino): fix: you could not use a
struct to for both di and json marshalling ([#41](https://github.com/defval/di/pull/41)).

## v1.11.0

### Added

- `di.Decorate` function that applies function after type resolve.

### Changed

- [@chirino](https://github.com/chirino): Prefer using "di:" field tags
  to control injection options to avoid conflicting with tags used by
  other libraries ([#38](https://github.com/defval/di/pull/38)).

## v1.10.0

### Added

- [@chirino](https://github.com/chirino): Container nesting. See
  `AddParent()` function ([#35](https://github.com/defval/di/pull/35)).
- An experimental feature: Instance decoration with `di.Decorate()`.

### Fixed

- [@chirino](https://github.com/chirino): Calling `Resolve()` on a
  `di.Injectable` would overwrite the skip fields
  ([#34](https://github.com/defval/di/pull/34)).

## v1.9.0

### Added

- `container.ProvideValue()` function.

## v1.8.0

### Added

- `container.Apply()` function.

## v1.7.1

### Fixed

- Style and coverage fixes.

## v1.7.0

### Added

- Added embed fields support.

### Fixed

- `di.Inject` now works with structs and pointers correctly.

## v1.6.3

### Fixed

- Fix `optional` fields resolving.

## v1.6.2

### Fixed

- Fix `di.As()` with several interfaces.

## v1.6.1

### Fixed

- Removed debug print.
- Documentation fixes.

## v1.6.0

### Changed

- Changed logging interface. See `di.SetTracer()`.

### Fixed

- Some documentation and test updates.

## v1.5.0

### Added

- Add error to `Has()`.

### Fixed

- `Has()` returns false if container could not build instance.

### Changed

- The supported version of go >1.13.

## v1.4.1

### Fixed

- Fix field injection into interface implementations.

## v1.4.0

### Added

- `Iterate` method for lazy loaded iteration by all instances.

## v1.3.1

### Fixed

- Bug: Resolve type as interface causes type reinitialization.

## v1.3.0: A release that doesn't deserve to be called `v2`

### BREAKING CHANGES

- Provide duplications allowed.
- Removed tag `di`. Now all public fields in injectable type will be
  injected.
- Resolving node without tags, now returns all nodes of this type.
- Now, `di:"type_name"` is a `name:"type_name"`.
- Removed `di.Prototype()`: bad practice.

### Added

- Tagging that allows specifying key value identity for types.
- `skip:"true"` field tag option, that skips field providing.

### Fixed

- A bit of bad code

## v1.2.1

### Fixed

- [Using `di.WithName()` breaks when having one entry without a `di.Name()`](https://github.com/defval/di/issues/16):

## v1.2.0

### Added

- Any type can be automatically resolved as a group.
- The container exposes itself by default.
- The only named type in the group will be resolved without a name.
- Dependency graph can be edited in the runtime (but you need to be
  careful with this).

## v1.1.0

### BREAKING CHANGES

- Changed `di.Parameter` to `di.Inject`.
- Remove `optional` support from `di` tag.
- Add `optional` tag. See
  [this](https://github.com/defval/di#optional-parameters).

### Added

- Support injection into constructor result struct via `di.Inject`.

## v1.0.2

### Added

- Location of `di.Provide()`, `di.Invoke()`, `di.Resolve()` in error.

### Fixed

- Fix: `di.As()` with nil causes panic.

## v1.0.1

### Fixed

- `container.Provide` could not be called after container compilation
  now.
- Improve error messages.


## v1.0.0

Initial release.
