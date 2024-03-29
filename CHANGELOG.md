# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog].

This project adheres to [Semantic Versioning].

## [Unreleased]

## [3.4.0] - 2024-03-21

### Changed

- Improved various error and info messages, and help text.
- Modified `add` verb to support app IDs.
- Enhanced `upload` verb to support user tokens and app IDs.

## [3.3.0] - 2024-03-08

### Added

- Added new `auth` verb to support authorizing with a user token.

### Changed

- Upgraded all dependencies to latest and greatest.

## [3.2.1] - 2023-12-14

### Changed

- Upgraded to use Go 1.21.
- Upgraded all dependencies to latest and greatest.

## [3.2.0] - 2023-07-26

### Added

- Added support for React Native (non-Expo) recipes.

## [3.1.0] - 2023-06-14

### Added

- Added support for building Gradle recipes.
- Added support for building Flutter recipes.

### Changed

- The `add` verb now requests confirmation before adding a new recipe to the
  Waldo configuration.

### Fixed

- Fixed issue with `add` verb automatically initializing a Waldo configuration.

## [3.0.0] - 2023-05-30

Initial public release of enhanced functionality.

### Added

- Added new verbs:
  - `add` — Add a recipe describing how to build and upload a specific variant of the app.
  - `build` — Build app from recipe.
  - `init` — Create an empty Waldo configuration.
  - `list` — List defined recipes.
  - `remove` — Remove a recipe.
  - `sync` — Build app from recipe and then upload to Waldo.
  - `version` — Display version information.

### Changed

- Enhanced `upload` verb to interoperate with `build`.

### Removed

- Removed support for omitting the `upload` verb.

## [2.0.5] - 2023-01-30

### Added

- Added support for network retry.

## [2.0.4] - 2022-09-08

### Fixed

- Fixed issue introduced in prior change.

## [2.0.3] - 2022-09-08

### Changed

- Enhanced `sim_appcenter_build_and_upload.sh` Bash script to allow specifying
  extra options on the underlying `xcodebuild build` invocation.

## [2.0.2] - 2022-05-23

### Changed

- Enhanced `sim_appcenter_build_and_upload.sh` Bash script to allow specifying
  extra options on the underlying `waldo upload` invocation.

## [2.0.1] - 2022-04-27

### Fixed

- Fixed edge case in determining correct wrapper name/version.

## [2.0.0] - 2022-04-26

### Changed

- Major rewrite to enhance stability.

## [1.1.3] - 2022-04-13

### Changed

- Greatly improved git branch inference.

## [1.1.2] - 2022-04-07

### Fixed

- Fixed git branch inference issues.
- Corrected info gathered from GitHub Actions.

## [1.1.1] - 2022-03-25

### Changed

- Enhanced info gathered from CI provider.

## [1.1.0] - 2022-03-08

### Added

- Added new `trigger` verb to support triggering a run of one or more test
  flows for an app.
- Added new `upload` verb to support uploading a build. This verb can be
  omitted for backward compatibility with older scripts.

## [1.0.0] - 2021-11-05

Initial public release.

[Unreleased]:   https://github.com/waldoapp/waldo-go-cli/compare/3.4.0...HEAD
[3.4.0]:		https://github.com/waldoapp/waldo-go-cli/compare/3.3.0...3.4.0
[3.3.0]:		https://github.com/waldoapp/waldo-go-cli/compare/3.2.1...3.3.0
[3.2.1]:		https://github.com/waldoapp/waldo-go-cli/compare/3.2.0...3.2.1
[3.2.0]:		https://github.com/waldoapp/waldo-go-cli/compare/3.1.0...3.2.0
[3.1.0]:		https://github.com/waldoapp/waldo-go-cli/compare/3.0.0...3.1.0
[3.0.0]:		https://github.com/waldoapp/waldo-go-cli/compare/2.0.5...3.0.0
[2.0.5]:        https://github.com/waldoapp/waldo-go-cli/compare/2.0.4...2.0.5
[2.0.4]:        https://github.com/waldoapp/waldo-go-cli/compare/2.0.3...2.0.4
[2.0.3]:        https://github.com/waldoapp/waldo-go-cli/compare/2.0.2...2.0.3
[2.0.2]:        https://github.com/waldoapp/waldo-go-cli/compare/2.0.1...2.0.2
[2.0.1]:        https://github.com/waldoapp/waldo-go-cli/compare/2.0.0...2.0.1
[2.0.0]:        https://github.com/waldoapp/waldo-go-cli/compare/1.1.3...2.0.0
[1.1.3]:        https://github.com/waldoapp/waldo-go-cli/compare/1.1.2...1.1.3
[1.1.3]:        https://github.com/waldoapp/waldo-go-cli/compare/1.1.2...1.1.3
[1.1.2]:        https://github.com/waldoapp/waldo-go-cli/compare/1.1.1...1.1.2
[1.1.1]:        https://github.com/waldoapp/waldo-go-cli/compare/1.1.0...1.1.1
[1.1.0]:        https://github.com/waldoapp/waldo-go-cli/compare/1.0.0...1.1.0
[1.0.0]:        https://github.com/waldoapp/waldo-go-cli/compare/f05ec68...1.0.0

[Keep a Changelog]:     https://keepachangelog.com
[Semantic Versioning]:  https://semver.org
