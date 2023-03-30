# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.2.0] - 2023-03-30

### Changed

- Readme updates.
- `internal` package - has no change to library functionality.
- Rename `maxTokens` to `max` when constructing new rate limiter.
- Fine-grained locking for rate limiter.
- Rename internal struct fields.
- Rename package (it was wrongly packaged under ratelimiter)
- Use nanosecnd precision for timestamps.

### Fixed

- Make rate limiter more thread safe. There was potential for race conditons prior to this.

## [0.1.0] - 2023-03-30

### Added

- Token bucket based rate limiter (initial release).
