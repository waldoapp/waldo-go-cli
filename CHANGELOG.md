# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog].

This project adheres to [Semantic Versioning].

## [Unreleased]

## [4.0.0] - 2024-05-15

Major rework of supported verbs and functionality.

### Changed

- Modify `upload` verb to accept either a CI token or an API token for use as the upload token. If neither the `--upload_token` option nor the `WALDO_UPLOAD_TOKEN` environment variable is specified, the API token specified by the most recent `waldo auth` invocation (if any) is used. Moreover, if the upload token is an API token, the `--app_id` option is _required_. Conversely, if the upload token is a CI token, the `--app_id` option is _not allowed_.

### Removed

- Remove `add`, `build`, `init`, `list`, `remove`, and `sync` verbs.
- Remove recipe support from `upload` verb.

## [3.5.0] - 2024-04-15

### Added

- Add support for Expo recipes.

### Fixed

- Fix issues with building for Android platform.

## [3.4.0] - 2024-03-21

### Changed

- Improve various error and info messages, and help text.
- Modify `add` verb to support app IDs.
- Enhance `upload` verb to support user tokens and app IDs.

## [3.3.0] - 2024-03-08

### Added

- Add new `auth` verb to support authenticating with a user token.

### Changed

- Upgrade all dependencies to latest and greatest.

## [3.2.1] - 2023-12-14

### Changed

- Upgrade to use Go 1.21.
- Upgrade all dependencies to latest and greatest.

## [3.2.0] - 2023-07-26

### Added

- Add support for React Native (non-Expo) recipes.

## [3.1.0] - 2023-06-14

### Added

- Add support for building Gradle recipes.
- Add support for building Flutter recipes.

### Changed

- Request confirmation on `add` verb before adding a new recipe to the Waldo configuration.

### Fixed

- Fix issue with `add` verb automatically initializing a Waldo configuration.

## [3.0.0] - 2023-05-30

Initial public release of enhanced functionality.

### Added

- Add new verbs:
  - `add` — Add a recipe describing how to build and upload a specific variant of the app.
  - `build` — Build app from recipe.
  - `init` — Create an empty Waldo configuration.
  - `list` — List defined recipes.
  - `remove` — Remove a recipe.
  - `sync` — Build app from recipe and then upload to Waldo.
  - `version` — Display version information.

### Changed

- Enhance `upload` verb to interoperate with `build`.

### Removed

- Remove support for omitting the `upload` verb.

## [2.0.5] - 2023-01-30

### Added

- Add support for network retry.

## [2.0.4] - 2022-09-08

### Fixed

- Fix issue introduced in prior change.

## [2.0.3] - 2022-09-08

### Changed

- Enhance `sim_appcenter_build_and_upload.sh` Bash script to allow specifying extra options on the underlying `xcodebuild build` invocation.

## [2.0.2] - 2022-05-23

### Changed

- Enhance `sim_appcenter_build_and_upload.sh` Bash script to allow specifying extra options on the underlying `waldo upload` invocation.

## [2.0.1] - 2022-04-27

### Fixed

- Fix edge case in determining correct wrapper name/version.

## [2.0.0] - 2022-04-26

### Changed

- Major rewrite to enhance stability.

## [1.1.3] - 2022-04-13

### Changed

- Greatly improve git branch inference.

## [1.1.2] - 2022-04-07

### Fixed

- Fix git branch inference issues.
- Correct info gathered from GitHub Actions.

## [1.1.1] - 2022-03-25

### Changed

- Enhance info gathered from CI provider.

## [1.1.0] - 2022-03-08

### Added

- Add new `trigger` verb to support triggering a run of one or more test flows for an app.
- Add new `upload` verb to support uploading a build. This verb can be omitted for backward compatibility with older scripts.

## [1.0.0] - 2021-11-05

Initial public release.

[Unreleased]:   https://github.com/waldoapp/waldo-go-cli/compare/4.0.0...HEAD
[4.0.0]:		https://github.com/waldoapp/waldo-go-cli/compare/3.5.0...4.0.0
[3.5.0]:		https://github.com/waldoapp/waldo-go-cli/compare/3.4.0...3.5.0
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
