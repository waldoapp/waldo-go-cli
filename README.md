# Waldo CLI (Go Edition)

[![License](https://img.shields.io/badge/license-MIT-000000.svg?style=flat)][license]
![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-lightgrey.svg?style=flat)

## About Waldo

[Waldo](https://www.waldo.io) provides fast, reliable, and maintainable tests
for the most critical flows in your app. Waldo CLI is a command-line tool which
allows you to upload an iOS or Android build to Waldo for processing.

## Installation

To install Waldo CLI, simply download and execute the installer script:

```bash
export WALDO_CLI_BIN=/usr/local/bin     # be sure this location is in $PATH

bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
```

You can verify that you have installed Waldo CLI correctly with the `which
waldo` and `waldo --help` commands.

If you ever need to uninstall Waldo CLI, simply delete the executable from
`$WALDO_CLI_BIN`.

## Usage

First, obtain an upload token from Waldo for your app. This is used to
authenticate with the Waldo backend on each call.

Next, generate a new build for your app. Waldo CLI recognizes the following
file extensions:

- `.app` for iOS simulator builds _only_
- `.ipa` for iOS device builds _only_
- `.apk` for all Android builds (emulator or device)

Finally, specify the path to your new build (along with your Waldo upload
token) on the `waldo` command invocation:

```bash
$ waldo /path/to/YourApp.app --upload_token 0123456789abcdef0123456789abcdef
```

> **Important:** Make sure you replace the fake upload token value shown above
> with the _real_ value for your Waldo app.

You can also use an environment variable to provide the Waldo upload token to
Waldo CLI:

```bash
$ export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef
$ waldo /path/to/YourApp.app
```

[license]:  https://github.com/waldoapp/waldo-go-cli/blob/master/LICENSE
