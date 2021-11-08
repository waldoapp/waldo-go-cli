# Waldo CLI (Go Edition)

[![License](https://img.shields.io/badge/license-MIT-000000.svg?style=flat)][license]
![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-lightgrey.svg?style=flat)

## About Waldo

[Waldo](https://www.waldo.io) provides fast, reliable, and maintainable tests
for the most critical flows in your app. Waldo CLI is a command-line tool which
allows you to upload an iOS or Android build to Waldo for processing.

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
(https://github.com/waldoapp/waldo-go-cli/releases/latest) and download the
appropriate `waldo` executable for your machine (either
`waldo-windows-x86_64.exe` or `waldo-windows-arm64.exe`) and install it in a
location known to `%PATH%`. You can verify that you have installed it correctly
with the `waldo --help` command.

If you ever need to uninstall Waldo CLI, simply delete the executable from that
location.

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

You can also use an environment variable to provide the upload token to Waldo
CLI:

```bash
$ export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef
$ waldo /path/to/YourApp.app
```

### Advanced Usage

Whereas only the build path and upload token are _required_ to successfully
upload your build to Waldo, there are a few other _non-required_ options
recognized by Waldo CLI that you may find useful:

- `--variant_name <value>` — This option instructs Waldo CLI to associate an
  arbitrary string (provided by you) with your build.
- `--verbose` — If you specify this option, Waldo CLI prints additional debug
  information. This can shed more light on why your build is failing to upload.
- `--git_commit <value>` and `--git_branch <value>` — If you have `git`
  installed and you are running from the working directory of a git repository,
  Waldo CLI attempts to “infer” the most likely commit SHA and branch name to
  associate with your build. In most cases it works very well. However, some
  CIs make it difficult or impossible for Waldo CLI to deduce this information.
  In such cases, you can directly specify the git information to associate with
  your build using these options.

## Integrating Waldo CLI with Your CI

See [CI_INTEGRATION.md][ci] for full details.

[ci]:       https://github.com/waldoapp/waldo-go-cli/blob/master/CI_INTEGRATION.md
[license]:  https://github.com/waldoapp/waldo-go-cli/blob/master/LICENSE
