# Changelog

All notable changes to this project will be documented in this file.

The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this
project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html): TBD, use modules or another vendor system.

## v1.4.0

### Added

- `Iterate` method for lazy loaded iteration by all instances.

## v1.3.1

### Fixed

- Bug: Resolve type as interface causes type reinitialization.

## v1.3.0: A release that doesn't deserve to be called `v2`

### BREAKING CHANGES

- Provide duplications allowed.
- Removed tag `di`. Now all public fields in injectable type will be injected.
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

- [Using `di.WithName()` breaks when having one entry without a `di.Name()`](https://github.com/goava/di/issues/16):

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
  [this](https://github.com/goava/di#optional-parameters).

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
