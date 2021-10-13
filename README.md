# Waldo Go CLI

[![License](https://img.shields.io/badge/license-MIT-000000.svg?style=flat)][license]
![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS-lightgrey.svg?style=flat)

## About Waldo

[Waldo](https://www.waldo.io) provides fast, reliable, and maintainable tests
for the most critical flows in your app. Waldo Go CLI is a command-line tool
which allows you to upload an iOS or Android build to Waldo for processing.

Simply download the appropriate `waldo` executable for the [latest
release][release] and install it into `/usr/local/bin`. You can verify that you
have installed it correctly with the `which waldo` and `waldo --help` commands.

If you ever need to uninstall Waldo Go CLI, simply delete the executable from
`/usr/local/bin`.

## Usage

To get started, first obtain an upload token from Waldo for your app. This is
used to authenticate with the Waldo backend on each call.

Build a new IPA or APK for your app and specify the path to it (along with the
Waldo upload token) on the `waldo` command invocation:

```bash
$ waldo /path/to/YourApp.ipa --upload_token 0123456789abcdef0123456789abcdef
```

Make sure you replace the fake upload token value shown above with the real
value for your Waldo application.

You can also use an environment variable to provide the Waldo upload token to
Waldo Go CLI:

```bash
$ export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef
$ waldo /path/to/YourApp.ipa
```

[license]:  https://github.com/waldoapp/waldo-go-cli/blob/master/LICENSE
[release]:  https://github.com/waldoapp/waldo-go-cli/releases
