# Uploading an iOS Simulator Build with App Center

As mentioned in the [CI Integration](CI_INTEGRATION.md) document, uploading an
iOS simulator build is not officially supported in [App Center]. However, we
have found a usable workaround that accomplishes just that. Multiple Waldo
customers have implemented this solution until such time as App Center provides
official support for simulator builds.

This solution “piggybacks” a simulator build on top of a regular device build.
It requires you only to add a couple of [custom build scripts][build_scripts].

## Step 1

First, add the following to `appcenter-post-clone.sh`:

```bash
export WALDO_CLI_BIN=/usr/local/bin

bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
```

> **Note:** This downloads a special Bash script,
> `sim_appcenter_build_and_upload.sh`, in addition to the `waldo` executable
> binary.

## Step 2

Then, add the following to `appcenter-pre-build.sh` (making sure to supply the
appropriate values to the environment variables):

```bash
WALDO_CLI_BIN=/usr/local/bin

export SIM_XCODE_PROJECT=YourApp.xcodeproj  # or YourApp.xcworkspace
export SIM_XCODE_SCHEME=YourApp
export SIM_XCODE_CONFIGURATION=Release      # or equivalent
export SIM_XCODE_APP_NAME=YourApp.app

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

#
# Uncomment and define the following environment variable if you need to pass
# extra options (for example, `--git_branch` or `--git_commit`) to the
# underlying `waldo upload` invocation:
#
# export SIM_WALDO_UPLOAD_OPTIONS="--git_branch ${YOUR_GIT_BRANCH} --verbose"

#
# Uncomment and define the following three environment variables to disable the
# device build operation from running and display as “Canceled” in the App
# Center dashboard:
#
# export SIM_APPCENTER_API_TOKEN=0123456789abcdef0123456789abcdef01234567
# export SIM_APPCENTER_OWNER_NAME="Owner Name"
# export SIM_APPCENTER_APP_NAME=YourApp

${WALDO_CLI_BIN}/sim_appcenter_build_and_upload.sh || exit

#
# Uncomment the following “exit” line to disable the device build operation
# from running and display as “Failed” in the App Center dashboard:
#
# exit 1
```

> **Note 1:** These additions are executed _pre_-build in contrast to those for
> device builds which are executed _post_-build.

> **Note 2:** You can also set environment variable `SIM_XCODE_PROJECT` to
> `$APPCENTER_XCODE_PROJECT` and `SIM_XCODE_SCHEME` to
> `$APPCENTER_XCODE_SCHEME`.

## Details

The `SIM_APPCENTER_*` environment variables enable the helper script to
_cancel_ the App Center build _after_ the simulator build uploads to Waldo, but
_before_ the regular (and expensive) device build operation starts. If you
choose to _not_ set these environment variables, the simulator build will still
upload the simulator build to Waldo. What happens afterward depends on whether
you explicitly exit the pre-build script with success (`0`) or failure (_not_
`0`).

Thus, there are three options you can choose for the pre-build script:

1. **normal** _(default)_ — Leave the `SIM_APPCENTER_*` environment variables
   and the `exit` line commented out. This will upload a simulator build of
   your app to Waldo and then proceed to build it for device, too. If you have
   also created a post-build script (as described in the [CI
   Integration](CI_INTEGRATION.md) document), the device build will be uploaded
   to Waldo as well. The entire build will display as `Succeeded` in the App
   Center dashboard.

2. **cancel** — Uncomment and define the `SIM_APPCENTER_*` environment
   variables (leave the `exit` line commented out). This will upload a
   simulator build of your app to Waldo and then _cancel_ the device build. The
   entire build will display as `Canceled` in the App Center dashboard.

   The `sim_appcenter_build_and_upload.sh` script calls an App Center API
   endpoint to cancel the device build. Therefore, you _must_ define the
   following environment variables correctly in order to cancel the device
   build:

   - `SIM_APPCENTER_API_TOKEN` — This must be set to a valid API token: either
     a user token or an app token.

     See [Creating an App Center App API token][app_api_token] or [Creating an
     App Center User API token][user_api_token] for details.

   - `SIM_APPCENTER_APP_NAME` — This must be set to the name of your app.

     See [Find owner_name and app_name from an App Center URL][owner_app_names]
     for details.

   - `SIM_APPCENTER_OWNER_NAME` — This must be set to the name of the _owner_
     of your app.

     See [Find owner_name and app_name from an App Center URL][owner_app_names]
     for details.

3. **failure** — Uncomment the `exit` line (leave the `SIM_APPCENTER_*`
   environment variables commented out). This will upload a simulator build of
   your app to Waldo and then _fail_ the device build. The
   entire build will display as `Failed` in the App Center dashboard.

We _strongly_ recommend that you choose option 2 unless you actually desire the
device build operation to be run. In this way, the build displays a `Canceled`
status, and thereby reserves the `Failed` status for true build failures.

[App Center]:   https://appcenter.ms

[app_api_token]:    https://docs.microsoft.com/en-us/appcenter/api-docs/#creating-an-app-center-app-api-token
[build_scripts]:    https://docs.microsoft.com/en-us/appcenter/build/custom/scripts/
[owner_app_names]:  https://docs.microsoft.com/en-us/appcenter/api-docs/#find-owner_name-and-app_name-from-an-app-center-url
[user_api_token]:   https://docs.microsoft.com/en-us/appcenter/api-docs/#creating-an-app-center-user-api-token
