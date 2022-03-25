#!/usr/bin/env bash

set -eu -o pipefail

waldo_cli_version="1.1.1"

waldo_cli_bin="${WALDO_CLI_BIN:-/usr/local/bin}"
waldo_cli_url="https://github.com/waldoapp/waldo-go-cli/releases/download/${waldo_cli_version}"

waldo_asset1_name=""
waldo_asset2_name="sim_appcenter_build_and_upload.sh"

function check_platform() {
    if [[ -z $(which curl) ]]; then
        fail "No ‘curl’ command found"
    fi
}

function determine_asset_name() {
    local _platform=$(uname -s)
    local _arch=$(uname -m)

    case $_platform in
        Darwin)
            _platform="macos"
            ;;

        Linux)
            _platform="linux"
            ;;

        *)
            fail "Unsupported platform: ${_platform}"
            ;;
    esac

    case $_arch in
        arm64)
            _arch="arm64"
            ;;

        x86_64)
            _arch="x86_64"
            ;;

        *)
            fail "Unsupported architecture: ${_arch}"
            ;;
    esac

    waldo_asset1_name="waldo-${_platform}-${_arch}"
}

function fail() {
    echo "install-waldo.sh: $1" 1>&2
    exit 1
}

function install_binary() {
    mkdir -p "$waldo_cli_bin" || return

    if [[ ! -w $waldo_cli_bin ]]; then
        fail "No write access to ‘${waldo_cli_bin}’"
    fi

    curl --create-dirs --location --show-error --silent \
         "${waldo_cli_url}/${waldo_asset1_name}"        \
         --output "${waldo_cli_bin}/waldo" || return

    chmod +x "${waldo_cli_bin}/waldo" || return

    echo "Installed ‘${waldo_asset1_name}’ as ‘${waldo_cli_bin}/waldo’"

    if [[ -n ${APPCENTER_BUILD_ID:-} ]]; then
        curl --create-dirs --location --show-error --silent \
             "${waldo_cli_url}/${waldo_asset2_name}"        \
             --output "${waldo_cli_bin}/${waldo_asset2_name}" || return

        chmod +x "${waldo_cli_bin}/${waldo_asset2_name}" || return

        echo "Installed ‘${waldo_asset2_name}’ as ‘${waldo_cli_bin}/${waldo_asset2_name}’"
    fi

    return
}

check_platform || exit
determine_asset_name || exit
install_binary || exit

exit
