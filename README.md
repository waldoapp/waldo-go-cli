# Waldo CLI (Go Edition)

[![License](https://img.shields.io/badge/license-MIT-000000.svg?style=flat)][license]
![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-lightgrey.svg?style=flat)

## About Waldo

[Waldo](https://www.waldo.io) provides fast, reliable, and maintainable tests
for the most critical flows in your app. Waldo CLI is a command-line tool which
allows you to:

- Upload an iOS or Android build to Waldo for processing. See
  [here](https://docs.waldo.com/docs/ios-uploading-your-simulator-build-to-waldo)
  and
  [here](https://docs.waldo.com/docs/android-uploading-your-emulator-build-to-waldo)
  for more details.
- Trigger a run of of one or more test flows for your app. See
  [here](https://docs.waldo.com/docs/ci-run) for more details.

## Installation

### Linux and macOS

To install Waldo CLI, simply download and execute the installer script:

```bash
export WALDO_CLI_BIN=/usr/local/bin     # be sure this location is in $PATH

bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
```

You can verify that you have installed Waldo CLI correctly with the `which
waldo` and `waldo --help` commands.

If you ever need to uninstall Waldo CLI, simply delete the executable from
`$WALDO_CLI_BIN`.

### Windows

To install Waldo CLI, simply navigate to the latest release on GitHub
(https://github.com/waldoapp/waldo-go-cli/releases/latest), download the
appropriate `waldo` executable for your machine (either
`waldo-windows-x86_64.exe` or `waldo-windows-arm64.exe`), and install it in a
location known to `%PATH%`. You can verify that you have installed it correctly
with the `waldo --help` command.

If you ever need to uninstall Waldo CLI, simply delete the executable from that
location.

[license]:  https://github.com/waldoapp/waldo-go-cli/blob/master/LICENSE
