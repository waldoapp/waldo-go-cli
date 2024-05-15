#!/usr/bin/env bash

set -eu -o pipefail

waldo_cli_bin=
waldo_cli_url="${WALDO_CLI_URL:-https://github.com/waldoapp/waldo-go-cli/releases/latest/download}"

waldo_exec1_name="waldo"
waldo_exec2_name="sim_appcenter_build_and_upload.sh"

waldo_asset1_name=
waldo_asset2_name=

waldo_exec_path=

waldo_found_in_path=false
waldo_is_reinstall=false

function authenticate() {
    local _api_token=${TOKEN:-}
    local _auth_status

    if [[ -z $_api_token ]]; then
        _api_token=$(find_api_token)
    fi

    if [[ -n $_api_token ]]; then
        do_authentication $_api_token

        _auth_status=$?
    else
        _auth_status=1
    fi

    if (( $_auth_status != 0 )); then
        echo ""
        echo "You have not yet been authenticated to access Waldo. Some Waldo CLI functionality will be unavailable to you."
        echo ""
        echo "You must first run the following command:"
        echo ""
        echo "    waldo auth <api-token>"
        echo ""
        echo "You can retrieve your API token here: https://app.waldo.com/settings/profile"
    fi
}

function check_curl_command() {
    if ! command -v curl &>/dev/null; then
        fail "No ‘curl’ command found"
    fi
}

function check_destination() {
    waldo_cli_bin=${HOME}/.waldo/bin

    local _cur_path=$(command -v waldo)

    if [[ $waldo_cli_bin/waldo == $_cur_path  ]]; then
        waldo_found_in_path=true
    elif [[ -n $_cur_path ]]; then
        fail "Conflicting Waldo CLI installation found at ‘${_cur_path}’ -- please remove it"
    fi

    if [[ -e $waldo_cli_bin/waldo ]]; then
        echo "Waldo CLI installation found in ‘${waldo_cli_bin}’ -- will re-install"
        echo ""

        waldo_is_reinstall=true
    else
        echo "Waldo CLI will be installed in ‘${waldo_cli_bin}’"
        echo ""
    fi
}

function check_installation() {
    local _line="export PATH=\"$waldo_cli_bin:\$PATH\""

	if [[ $waldo_found_in_path != true ]]; then
		local _startup_file=$(find_startup_file)

		if [[ -n $_startup_file ]]; then
            echo ""
			echo "Updating your PATH in ‘${_startup_file}’ to support ‘waldo’ with the following:"
			echo ""
			echo "    ${_line}"

			echo "${_line}" >> "$_startup_file"
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
		echo "    ${_line}"
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

function do_authentication() {
    "${waldo_exec_path}" auth $1
}

function fail() {
    echo "install.sh: $1" 1>&2
    exit 1
}

function find_api_token() {
    local _profile_path="${HOME}/.waldo/profile.yml"

    if [[ ! -r $_profile_path ]]; then
        return
    fi

    local _regex='user_token:[ ]*([^ ]+)[ ]*'

    while read -r line; do
        if [[ $line =~ $_regex ]]; then
            echo ${BASH_REMATCH[1]}
            return
        fi
    done < "$_profile_path"
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

    install_binary "${waldo_asset1_name}" "${waldo_exec1_name}" || return

    if [[ -n ${APPCENTER_BUILD_ID:-} ]]; then
        install_binary "${waldo_asset2_name}" "${waldo_exec2_name}" || return
    fi

    waldo_exec_path="${waldo_cli_bin}/${waldo_exec1_name}"
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

detect_ci_mode || exit
check_destination || exit
check_curl_command || exit
determine_asset_names || exit
install_binaries || exit
check_installation || exit
authenticate || exit

exit
