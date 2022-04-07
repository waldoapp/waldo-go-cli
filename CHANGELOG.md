# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog].

This project adheres to [Semantic Versioning].

## [Unreleased]

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

[Unreleased]:   https://github.com/waldoapp/waldo-go-cli/compare/1.1.2...HEAD
[1.1.2]:        https://github.com/waldoapp/waldo-go-cli/compare/1.1.1...1.1.2
[1.1.1]:        https://github.com/waldoapp/waldo-go-cli/compare/1.1.0...1.1.1
[1.1.0]:        https://github.com/waldoapp/waldo-go-cli/compare/1.0.0...1.1.0
[1.0.0]:        https://github.com/waldoapp/waldo-go-cli/compare/f05ec68...1.0.0

[Keep a Changelog]:     https://keepachangelog.com
[Semantic Versioning]:  https://semver.org
