# CI Integration

* [Uploading a Build with fastlane](#upload_f)
  * [Uploading an iOS Simulator Build](#upload_f_ios_sim)
  * [Uploading an iOS Device Build](#upload_f_ios_dev)
  * [Uploading an Android Build](#upload_f_android)
* [Uploading a Build with App Center](#upload_ac)
  * [Uploading an iOS Simulator Build](#upload_ac_ios_sim)
  * [Uploading an iOS Device Build](#upload_ac_ios_dev)
  * [Uploading an Android Build](#upload_ac_android)
* [Uploading a Build with Bitrise](#upload_br)
  * [Uploading an iOS Simulator Build](#upload_br_ios_sim)
  * [Uploading an iOS Device Build](#upload_br_ios_dev)
  * [Uploading an Android Build](#upload_br_android)
* [Uploading a Build with buddybuild](#upload_bb)
  * [Uploading an iOS Simulator Build](#upload_bb_ios_sim)
  * [Uploading an iOS Device Build](#upload_bb_ios_dev)
  * [Uploading an Android Build](#upload_bb_android)
* [Uploading a Build with CircleCI](#upload_cc)
  * [Uploading an iOS Simulator Build](#upload_cc_ios_sim)
  * [Uploading an iOS Device Build](#upload_cc_ios_dev)
  * [Uploading an Android Build](#upload_cc_android)
* [Uploading a Build with GitHub Actions](#upload_gha)
  * [Uploading an iOS Simulator Build](#upload_gha_ios_sim)
  * [Uploading an iOS Device Build](#upload_gha_ios_dev)
  * [Uploading an Android Build](#upload_gha_android)
* [Uploading a Build with Travis CI](#upload_tc)
  * [Uploading an iOS Simulator Build](#upload_tc_ios_sim)
  * [Uploading an iOS Device Build](#upload_tc_ios_dev)
  * [Uploading an Android Build](#upload_tc_android)
* [Uploading a Build Manually](#upload_m)
  * [Uploading an iOS Simulator Build](#upload_m_ios_sim)
  * [Uploading an iOS Device Build](#upload_m_ios_dev)
  * [Uploading an Android Build](#upload_m_android)

## <a name="upload_f">Uploading a Build with fastlane</a>

Waldo integration with [fastlane] requires you only to add the `waldo` plugin
to your project:

```bash
$ fastlane add_plugin waldo
```

### <a name="upload_f_ios_sim">Uploading an iOS Simulator Build</a>

Create a new simulator build for your app.

You can use `gym` (aka `build_ios_app`) to build your app provided that you
supply several parameters in order to convince Xcode to _both_ build for the
simulator _and_ not attempt to generate an IPA:

```ruby
gym(configuration: 'Release',
    derived_data_path: '/path/to/derivedData',
    skip_package_ipa: true,
    skip_archive: true,
    destination: 'generic/platform=iOS Simulator')
```

You can then find your app relative to the derived data path:

```ruby
app_path = File.join(derived_data_path,
                     'Build',
                     'Products',
                     'Release-iphonesimulator',
                     'YourApp.app')
```

Regardless of how you create the actual simulator build for your app, the
upload itself is very simple:

```ruby
waldo(upload_token: '0123456789abcdef0123456789abcdef',
      app_path: '/path/to/YourApp.app')
```

> **Note:** You _must_ specify _both_ the Waldo upload token _and_ the path of
> the `.app`.

### <a name="upload_f_ios_dev">Uploading an iOS Device Build</a>

Build a new IPA for your app. If you use `gym` (aka `build_ios_app`) to build
your IPA, `waldo` will automatically find and upload the generated IPA.

```ruby
gym(export_method: 'ad-hoc')                            # or 'development'

waldo(upload_token: '0123456789abcdef0123456789abcdef')
```

> **Note:** You _must_ specify the Waldo upload token.

If you do _not_ use `gym` to build your IPA, you will need to explicitly
specify the IPA path to `waldo`:

```ruby
waldo(upload_token: '0123456789abcdef0123456789abcdef',
      ipa_path: '/path/to/YourApp.ipa')
```

### <a name="upload_f_android">Uploading an Android Build</a>

Build a new APK for your app. If you use `gradle` to build your APK, `waldo`
will automatically find and upload the generated APK.

```ruby
gradle(task: 'assemble',
       build_type: 'Release')

waldo(upload_token: '0123456789abcdef0123456789abcdef')
```

> **Note:** You _must_ specify the Waldo upload token.

If you do _not_ use `gradle` to build your APK, you will need to explicitly
specify the APK path to `waldo`:

```ruby
waldo(upload_token: '0123456789abcdef0123456789abcdef',
      apk_path: '/path/to/YourApp.apk')
```

----------

## <a name="upload_ac">Uploading a Build with App Center</a>

Waldo integration with [App Center] requires you only to add a couple of
[custom build scripts][ac_scripts].

In all cases, add the following to `appcenter-post-clone.sh`:

```bash
export WALDO_CLI_BIN=/usr/local/bin

bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
```

### <a name="upload_ac_ios_sim">Uploading an iOS Simulator Build</a>

_Not supported by the CI._

See [here](SIM_APPCENTER.md) for a usable workaround.

### <a name="upload_ac_ios_dev">Uploading an iOS Device Build</a>

Add the following to `appcenter-post-build.sh`:

```bash
WALDO_CLI_BIN=/usr/local/bin

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

BUILD_PATH=${APPCENTER_OUTPUT_DIRECTORY}/YourApp.ipa

${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

### <a name="upload_ac_android">Uploading an Android Build</a>

Add the following to `appcenter-post-build.sh`:

```bash
WALDO_CLI_BIN=/usr/local/bin

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

BUILD_PATH=${APPCENTER_OUTPUT_DIRECTORY}/YourApp.apk

${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

----------

## <a name="upload_br">Uploading a Build with Bitrise</a>

Waldo integration with [Bitrise] requires you only to add a [`Waldo
Upload`][br_waldo_upload] step to your workflow.

### <a name="upload_br_ios_sim">Uploading an iOS Simulator Build</a>

First, create a new simulator build for your app. When you use the [`Xcode
build for simulator`][br_xcode_build] step to build your app, output variables
are generated that you can then use as input to the [`Waldo
Upload`][br_waldo_upload] step to find and upload the generated app.

```yaml
workflows:
  primary:
    steps:
    #...
    - xcode-build-for-simulator:
        inputs:
        - xcodebuild_options: CODE_SIGNING_ALLOWED=YES
    - waldo-upload:
        inputs:
        - build_path: $BITRISE_APP_DIR_PATH
        - upload_token: $WALDO_UPLOAD_TOKEN     # from your secrets
    #...
```

> **Note:** The value you supply to the `upload_token` input _should_ be
> specified as a “secret” environment variable by going to the **Secrets** tab
> in the Bitrise **Workflow Editor** and assigning your upload token to
> `WALDO_UPLOAD_TOKEN`.

### <a name="upload_br_ios_dev">Uploading an iOS Device Build</a>

First, build a new IPA for your app. When you use the [`Xcode Archive & Export
for iOS`][br_xcode_archive] step to build your IPA, output variables are
generated that you can then use as input to the [`Waldo
Upload`][br_waldo_upload] step to find and upload the generated IPA.

```yaml
workflows:
  primary:
    steps:
    #...
    - xcode-archive:
        inputs:
        - export_method: ad-hoc                 # or development
        - compile_bitcode: 'no'
        - upload_bitcode: 'no'
    - waldo-upload:
        inputs:
        - build_path: $BITRISE_IPA_PATH
        - upload_token: $WALDO_UPLOAD_TOKEN     # from your secrets
    #...
```

> **Note:** The value you supply to the `upload_token` input _should_ be
> specified as a “secret” environment variable by going to the **Secrets** tab
> in the Bitrise **Workflow Editor** and assigning your upload token to
> `WALDO_UPLOAD_TOKEN`.

### <a name="upload_br_android">Uploading an Android Build</a>

First, build a new APK for your app. When you use the [`Android
Build`][br_android_build] step to build your APK, output variables are
generated that you can then use as input to the [`Waldo
Upload`][br_waldo_upload] step to find and upload the generated APK.

```yaml
workflows:
  primary:
    steps:
    #...
    - android-build: {}
    - waldo-upload:
        inputs:
        - build_path: $BITRISE_APK_PATH
        - upload_token: $WALDO_UPLOAD_TOKEN     # from your secrets
    #...
```

> **Note:** The value you supply to the `upload_token` input _should_ be
> specified as a “secret” environment variable by going to the **Secrets** tab
> in the Bitrise **Workflow Editor** and assigning your upload token to
> `WALDO_UPLOAD_TOKEN`.

----------

## <a name="upload_bb">Uploading a Build with buddybuild</a>

Waldo integration with [buddybuild] requires you only to add a couple of
[custom build steps][bb_custom].

In all cases, add the following to `buddybuild_postclone.sh`:

```bash
export WALDO_CLI_BIN=/usr/local/bin

bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
```

### <a name="upload_bb_ios_sim">Uploading an iOS Simulator Build</a>

Add the following to `buddybuild_postbuild.sh`:

```bash
WALDO_CLI_BIN=/usr/local/bin

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

BUILD_PATH=/path/to/YourApp.app

${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

### <a name="upload_bb_ios_dev">Uploading an iOS Device Build</a>

Add the following to `buddybuild_postbuild.sh`:

```bash
WALDO_CLI_BIN=/usr/local/bin

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

cd $BUDDYBUILD_PRODUCT_DIR

${WALDO_CLI_BIN}/waldo "$BUDDYBUILD_IPA_PATH"
```

### <a name="upload_bb_android">Uploading an Android Build</a>

_Not supported by the CI._

----------

## <a name="upload_cc">Uploading a Build with CircleCI</a>

### <a name="upload_cc_ios_sim">Uploading an iOS Simulator Build</a>

Waldo integration with [CircleCI] requires you only to add a couple of steps to
your [configuration][cc_config]:

```yaml
jobs:
  build:
    steps:
      - run:
        name: Download Waldo CLI
        command: |
          bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
        environment:
          WALDO_CLI_BIN: .circleci/waldo

      #...
      #... (generate .app)
      #...

      - run:
        name: Upload build to Waldo
        command: .circleci/waldo "$WALDO_BUILD_PATH"
        environment:
          WALDO_UPLOAD_TOKEN: 0123456789abcdef0123456789abcdef
          WALDO_BUILD_PATH: /path/to/YourApp.app
```

### <a name="upload_cc_ios_dev">Uploading an iOS Device Build</a>

Waldo integration with [CircleCI] requires you only to add a couple of steps to
your [configuration][cc_config]:

```yaml
jobs:
  build:
    steps:
      - run:
        name: Download Waldo CLI
        command: |
          bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
        environment:
          WALDO_CLI_BIN: .circleci/waldo

      #...
      #... (generate .ipa)
      #...

      - run:
        name: Upload build to Waldo
        command: .circleci/waldo "$WALDO_BUILD_PATH"
        environment:
          WALDO_UPLOAD_TOKEN: 0123456789abcdef0123456789abcdef
          WALDO_BUILD_PATH: /path/to/YourApp.ipa
```

### <a name="upload_cc_android">Uploading an Android Build</a>

Waldo integration with [CircleCI] requires you only to add a couple of steps to
your [configuration][cc_config]:

```yaml
jobs:
  build:
    steps:
      - run:
        name: Download Waldo CLI
        command: |
          bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
        environment:
          WALDO_CLI_BIN: .circleci/waldo

      #...
      #... (generate .apk)
      #...

      - run:
        name: Upload build to Waldo
        command: .circleci/waldo "$WALDO_BUILD_PATH"
        environment:
          WALDO_UPLOAD_TOKEN: 0123456789abcdef0123456789abcdef
          WALDO_BUILD_PATH: /path/to/YourApp.apk
```

----------

## <a name="upload_gha">Uploading a Build with GitHub Actions</a>

### <a name="upload_gha_ios_sim">Uploading an iOS Simulator Build</a>

Waldo integration with [GitHub Actions] requires you only to add an extra step
to your [workflow][gha_workflow]:

```yaml
jobs:
  build:
    steps:
      #...
      #... (generate .app)
      #...

      - name: Upload build to Waldo
        env:
          WALDO_BUILD_PATH: /path/to/YourApp.app
          WALDO_CLI_BIN: /usr/local/bin
          WALDO_UPLOAD_TOKEN: 0123456789abcdef0123456789abcdef
        run: |
          if [ ! -e ${WALDO_CLI_BIN}/waldo ]; then
            bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
          fi

          ${WALDO_CLI_BIN}/waldo "$WALDO_BUILD_PATH"
```

> **Important:** If you use the [Checkout V2][gha_checkout] action in your
> workflow, you _must_ set `fetch-depth: 0`. Otherwise, Waldo CLI will _not_ be
> able to correctly identify the git commit and branch associated with your
> build.

### <a name="upload_gha_ios_dev">Uploading an iOS Device Build</a>

Waldo integration with [GitHub Actions] requires you only to add an extra step
to your [workflow][gha_workflow]:

```yaml
jobs:
  build:
    steps:
      #...
      #... (generate .ipa)
      #...

      - name: Upload build to Waldo
        env:
          WALDO_BUILD_PATH: /path/to/YourApp.ipa
          WALDO_CLI_BIN: /usr/local/bin
          WALDO_UPLOAD_TOKEN: 0123456789abcdef0123456789abcdef
        run: |
          if [ ! -e ${WALDO_CLI_BIN}/waldo ]; then
            bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
          fi

          ${WALDO_CLI_BIN}/waldo "$WALDO_BUILD_PATH"
```

> **Important:** If you use the [Checkout V2][gha_checkout] action in your
> workflow, you _must_ set `fetch-depth: 0`. Otherwise, Waldo CLI will _not_ be
> able to correctly identify the git commit and branch associated with your
> build.

### <a name="upload_gha_android">Uploading an Android Build</a>

Waldo integration with [GitHub Actions] requires you only to add an extra step
to your [workflow][gha_workflow]:

```yaml
jobs:
  build:
    steps:
      #...
      #... (generate .apk)
      #...

      - name: Upload build to Waldo
        env:
          WALDO_BUILD_PATH: /path/to/YourApp.apk
          WALDO_CLI_BIN: /usr/local/bin
          WALDO_UPLOAD_TOKEN: 0123456789abcdef0123456789abcdef
        run: |
          if [ ! -e ${WALDO_CLI_BIN}/waldo ]; then
            bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
          fi

          ${WALDO_CLI_BIN}/waldo "$WALDO_BUILD_PATH"
```

> **Important:** If you use the [Checkout V2][gha_checkout] action in your
> workflow, you _must_ set `fetch-depth: 0`. Otherwise, Waldo CLI will _not_ be
> able to correctly identify the git commit and branch associated with your
> build.

----------

## <a name="upload_tc">Uploading a Build with Travis CI</a>

### <a name="upload_tc_ios_sim">Uploading an iOS Simulator Build</a>

Waldo integration with [Travis CI] requires you only to add a few steps to your
[.travis.yml][tc_docs]:

```yaml
env:
  global:
    - WALDO_CLI_BIN=/usr/local/bin
    - WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef
install:
  - bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
script:
  #...
  #... Build your app for simulator with:
  #...
  #...     - xcodebuild [...] -derivedDataPath "$TRAVIS_BUILD_DIR" [...] build
  #...
  - BUILD_PATH="$TRAVIS_BUILD_DIR"/Build/Products/Release-iphonesimulator/YourApp.app
  - ${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

### <a name="upload_tc_ios_dev">Uploading an iOS Device Build</a>

Waldo integration with [Travis CI] requires you only to add a few steps to your
[.travis.yml][tc_docs]:

```yaml
env:
  global:
    - WALDO_CLI_BIN=/usr/local/bin
    - WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef
install:
  - bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
script:
  #...
  #... Build your IPA with:
  #...
  #...     - xcodebuild [...] -archivePath /path/to/YourApp.xcarchive [...] archive
  #...     - xcodebuild -exportArchive [...] -archivePath /path/to/YourApp.xcarchive -exportPath /path/to/export [...]
  #...
  - BUILD_PATH=/path/to/export/YourApp-release.ipa
  - ${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

### <a name="upload_tc_android">Uploading an Android Build</a>

Waldo integration with [Travis CI] requires you only to add a few steps to your
[.travis.yml][tc_docs]:

```yaml
env:
  global:
    - WALDO_CLI_BIN=/usr/local/bin
    - WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef
install:
  - bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
script:
  #...
  #... Build your APK
  #...
  - ${WALDO_CLI_BIN}/waldo "/path/to/YourApp.apk"
```

----------

## <a name="upload_m">Uploading a Build Manually</a>

If you are building outside of CI/CD or in another CI provider, you can also
upload your iOS build manually using Waldo CLI.

### <a name="upload_m_ios_sim">Uploading an iOS Simulator Build</a>

```
export WALDO_CLI_BIN=/usr/local/bin

if [ ! -e ${WALDO_CLI_BIN}/waldo ]; then
  bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
fi

#...
#... (generate .app)
#...

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

BUILD_PATH=/path/to/YourApp.app

${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

### <a name="upload_m_ios_dev">Uploading an iOS Device Build</a>

```
export WALDO_CLI_BIN=/usr/local/bin

if [ ! -e ${WALDO_CLI_BIN}/waldo ]; then
  bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
fi

#...
#... (generate .ipa)
#...

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

BUILD_PATH=/path/to/YourApp.ipa

${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

### <a name="upload_m_android">Uploading an Android Build</a>

```
export WALDO_CLI_BIN=/usr/local/bin

if [ ! -e ${WALDO_CLI_BIN}/waldo ]; then
  bash -c "$(curl -fLs https://github.com/waldoapp/waldo-go-cli/raw/master/install-waldo.sh)"
fi

#...
#... (generate .apk)
#...

export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

BUILD_PATH=/path/to/YourApp.apk

${WALDO_CLI_BIN}/waldo "$BUILD_PATH"
```

[App Center]:       https://appcenter.ms
[Bitrise]:          https://www.bitrise.io
[buddybuild]:       https://www.buddybuild.com
[CircleCI]:         https://circleci.com
[fastlane]:         https://fastlane.tools
[GitHub Actions]:   https://github.com/features/actions
[Travis CI]:        https://travis-ci.com

[ac_scripts]:       https://docs.microsoft.com/en-us/appcenter/build/custom/scripts/
[bb_custom]:        https://docs.buddybuild.com/builds/custom_build_steps.html
[br_android_build]: https://app.bitrise.io/integrations/steps/android-build
[br_waldo_upload]:  https://app.bitrise.io/integrations/steps/waldo-upload
[br_xcode_archive]: https://app.bitrise.io/integrations/steps/xcode-archive
[br_xcode_build]:   https://app.bitrise.io/integrations/steps/xcode-build-for-simulator
[cc_config]:        https://circleci.com/docs/2.0/configuration-reference/
[gha_checkout]:     https://github.com/marketplace/actions/checkout
[gha_workflow]:     https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions
[tc_docs]:          https://docs.travis-ci.com
