# Waldo CLI

[![License](https://img.shields.io/badge/license-MIT-000000.svg?style=flat)][license]
![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-lightgrey.svg?style=flat)

## About Waldo

[Waldo](https://www.waldo.com) provides fast, reliable, and maintainable tests for the most critical flows in your app. Waldo CLI is a command-line tool which allows you to interact with Waldo in several useful ways:

- Add one or more recipes describing how to build a particular variant of your app from the command-line that is suitable for running on Waldo, and then actually build the app variant from a recipe. These recipes can be persisted to your repository alongside your app’s source code and thereby shared with your teammates. The following build tools are currently supported:
  - Xcode (see [here](https://docs.waldo.com/docs/exporting-your-build-for-waldo) for details)
  - Gradle
  - Flutter
  - React Native (non-Expo)
  - _…and more to come!_
- Upload an iOS or Android build to Waldo for processing. See [here](https://docs.waldo.com/docs/ios-uploading-your-simulator-build-to-waldo) and [here](https://docs.waldo.com/docs/android-uploading-your-emulator-build-to-waldo) for more details.
- Trigger a run of of one or more test flows for your app. See [here](https://docs.waldo.com/docs/ci-run) for more details.

> **Note:** Only the `upload` and `trigger` commands are fully supported when running in a CI environment.

Type `waldo help` to see all that Waldo CLI can do for you!

## Installation

> **Note:** If you intend to use Waldo CLI from a CI script, please refer to the _next_ section — [Installation for CI](#installation-for-ci) — for instructions.

### Linux and macOS

To install Waldo CLI, simply download and execute the installer script:

```bash
curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install.sh | bash
```

The script attempts to install Waldo CLI to `~/.waldo/bin`.

You can verify that you have installed Waldo CLI correctly with the `which waldo` and `waldo help` commands.

If you ever need to uninstall Waldo CLI, simply delete the binary at `~/.waldo/bin/waldo`.

### Windows

To install Waldo CLI, simply navigate to the [latest release](https://github.com/waldoapp/waldo-go-cli/releases/latest), download the appropriate `waldo` executable for your machine (either `waldo-windows-x86_64.exe` or `waldo-windows-arm64.exe`), and install it as `waldo.exe` to a location known to `%PATH%`.

You can verify that you have installed it correctly with the `waldo help` command.

If you ever need to uninstall Waldo CLI, simply delete the executable from the install location.

## Installation for CI

> **Note:** If you intend to use Waldo CLI interactively, please refer to the _previous_ section — [Installation](#installation) — for instructions.

### Linux and macOS

To install Waldo CLI, simply download and execute the installer script:

```bash
curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh | bash
```

By default, the script installs Waldo CLI to `/usr/local/bin`.

If you wish to install Waldo CLI to a different location, simply define the `WALDO_CLI_BIN` environment variable _before_ invoking the installer script:

```bash
export WALDO_CLI_BIN=/path/to/binary
curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh | bash
```

You can verify that you have installed Waldo CLI correctly with the `which waldo` and `waldo help` commands.

If you ever need to uninstall Waldo CLI, simply delete the executable from the install location.

### Windows

To install Waldo CLI, simply navigate to the [latest release](https://github.com/waldoapp/waldo-go-cli/releases/latest), download the appropriate `waldo` executable for your machine (either `waldo-windows-x86_64.exe` or `waldo-windows-arm64.exe`), and install it as `waldo.exe` to a location known to `%PATH%`.

You can verify that you have installed it correctly with the `waldo help` command.

If you ever need to uninstall Waldo CLI, simply delete the executable from the install location.

[license]:  https://github.com/waldoapp/waldo-go-cli/blob/master/LICENSE
