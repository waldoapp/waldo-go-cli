#!/usr/bin/env bash

set -eu -o pipefail

waldo_cli_bin="${WALDO_CLI_BIN:-}"
waldo_cli_url="${WALDO_CLI_URL:-https://github.com/waldoapp/waldo-go-cli/releases/latest/download}"

waldo_exec1_name="waldo"
waldo_exec2_name="sim_appcenter_build_and_upload.sh"

waldo_asset1_name=
waldo_asset2_name=

waldo_found_in_path=false
waldo_is_reinstall=false

function check_curl_command() {
    if ! command -v curl >/dev/null; then
        fail "No ‘curl’ command found"
    fi
}

function check_destination() {
    local _cur_path=$(command -v waldo)

    if [[ -z $waldo_cli_bin ]]; then
    	if [[ -n $_cur_path ]]; then
			waldo_cli_bin=$(dirname $_cur_path)
		else
			waldo_cli_bin=${HOME}/.waldo/bin
		fi
	fi

    if [[ $waldo_cli_bin/waldo == $_cur_path  ]]; then
        waldo_found_in_path=true
    fi

    if [[ -e $waldo_cli_bin/waldo ]]; then
        echo "Waldo CLI installation detected in ‘${waldo_cli_bin}’ -- will re-install"
        echo ""

        waldo_is_reinstall=true
    else
        echo "Waldo CLI will be installed in ‘${waldo_cli_bin}’"
        echo ""
    fi
}

function check_installation() {
	if [[ $waldo_found_in_path != true ]]; then
		local _startup_file=$(find_startup_file)

		if [[ -n $_startup_file ]]; then
			echo "Updating your PATH in ‘${_startup_file}’ to support ‘waldo’"

			echo 'export PATH=$PATH:'"$waldo_cli_bin" >> "$_startup_file"
		fi
	fi

    if [[ $waldo_is_reinstall == true ]]; then
        echo ""
        echo "Waldo CLI successfully re-installed!"
    else
        echo ""
        echo "Waldo CLI successfully installed!"
    fi

	if [[ $waldo_found_in_path != true ]]; then
        echo ""
        echo "Please open a new terminal window OR run the following in your current one:"
        echo ""
        echo "    export PATH=\"$waldo_cli_bin:\$PATH\""
        echo ""
        echo "Then run the following command:"
        echo ""
        echo "    waldo help"
    fi
}

function detect_ci_mode() {
    local _ci_mode=${CI:-}

    if [[ $_ci_mode == true || $_ci_mode == 1 ]]; then
        fail "CI environment detected -- please use ‘install-waldo.sh’ instead"
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
    echo "install.sh: $1" 1>&2
    exit 1
}

function find_startup_file() {
    local _bash_login="${HOME}/.bash_login"
    local _bash_profile="${HOME}/.bash_profile"
    local _bashrc="${HOME}/.bashrc"
    local _profile="${HOME}/.profile"

    local _zshenv="${ZDOTDIR:-${HOME}}/.zshenv"
    local _zshlogin="${ZDOTDIR:-${HOME}}/.zshlogin"
    local _zshprofile="${ZDOTDIR:-${HOME}}/.zshprofile"
    local _zshrc="${ZDOTDIR:-${HOME}}/.zshrc"

    #
    # These files should be tested in the order that zsh or bash would load them:
    #
    if [[ -f $_zshenv ]]; then
        echo "$_zshenv"
    elif [[ -f $_zshprofile ]]; then
        echo "$_zshprofile"
    elif [[ -f $_zshrc ]]; then
        echo "$_zshrc"
    elif [[ -f $_zshlogin ]]; then
        echo "$_zshlogin"
    elif [[ -f $_bash_profile ]]; then
        echo "$_bash_profile"
    elif [[ -f $_bashrc ]]; then
        echo "$_bashrc"
    elif [[ -f $_bash_login ]]; then
        echo "$_bash_login"
    elif [[ -f $_profile ]]; then
        echo "$_profile"
    else
        echo ""
    fi
}

function install_binaries() {
    mkdir -p "$waldo_cli_bin"

    local _mkdir_status=$?

    if (( $_mkdir_status != 0 )); then
        fail "Unable to create directory ‘${waldo_cli_bin}’"
    fi

    if [[ ! -w $waldo_cli_bin ]]; then
        fail "No write access to ‘${waldo_cli_bin}’"
    fi

    install_binary "${waldo_asset1_name}" "${waldo_exec1_name}"

    if [[ -n ${APPCENTER_BUILD_ID:-} ]]; then
        install_binary "${waldo_asset2_name}" "${waldo_exec2_name}"
    fi
}

function install_binary() {
    local _asset_url="${waldo_cli_url}/${1}"
    local _exec_path="${waldo_cli_bin}/${2}"

    curl --fail             \
         --location         \
         --progress-bar     \
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

detect_ci_mode
check_destination
check_curl_command
determine_asset_names
install_binaries
check_installation

exit
