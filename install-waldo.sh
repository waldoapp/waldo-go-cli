#!/usr/bin/env bash

set -eu -o pipefail

waldo_cli_bin="${WALDO_CLI_BIN:-/usr/local/bin}"
waldo_cli_url="https://github.com/waldoapp/waldo-go-cli/releases/latest/download"

waldo_exec1_name="waldo"
waldo_exec2_name="sim_appcenter_build_and_upload.sh"

waldo_asset1_name=""
waldo_asset2_name=""

function check_platform() {
    if [[ -z $(which curl) ]]; then
        fail "No ‘curl’ command found"
    fi
}

function determine_asset_names() {
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
    waldo_asset2_name="$waldo_exec2_name"
}

function fail() {
    echo "install-waldo.sh: $1" 1>&2
    exit 1
}

function install_binaries() {
    mkdir -p "$waldo_cli_bin" || return

    if [[ ! -w $waldo_cli_bin ]]; then
        fail "No write access to ‘${waldo_cli_bin}’"
    fi

    install_binary "${waldo_asset1_name}" "${waldo_exec1_name}" || return

    if [[ -n ${APPCENTER_BUILD_ID:-} ]]; then
        install_binary "${waldo_asset2_name}" "${waldo_exec2_name}" || return
    fi

    return
}

function install_binary() {
    local _asset_url="${waldo_cli_url}/${1}"
    local _exec_path="${waldo_cli_bin}/${2}"

    curl --fail             \
         --location         \
         --retry 1          \
         --show-error       \
         --silent           \
         "${_asset_url}"    \
         --output "${_exec_path}"

    local _curl_status=$?

    if (( $_curl_status != 0 )); then
        fail "Unable to download ‘${_asset_url}’"
    fi

    chmod +x "${_exec_path}"

    local _chmod_status=$?

    if (( $_chmod_status != 0 )); then
        fail "Unable to install ‘${_exec_path}’"
    fi

    echo "Installed ‘${_asset_url}’ as ‘${_exec_path}’"
}

check_platform || exit
determine_asset_names || exit
install_binaries || exit

exit
