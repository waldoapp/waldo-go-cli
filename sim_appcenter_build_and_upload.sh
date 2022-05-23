#!/usr/bin/env bash

set -eu -o pipefail

SIM_XCODE_DATA_PATH=/tmp/${SIM_XCODE_SCHEME}-$(uuidgen)
WALDO_CLI_BIN=$(unset CDPATH && cd "${0%/*}" &>/dev/null && pwd)

function cancel_appcenter_build() {
    local _owner_name=${SIM_APPCENTER_OWNER_NAME:-}
    local _api_token=${SIM_APPCENTER_API_TOKEN:-}
    local _app_name=${SIM_APPCENTER_APP_NAME:-}

    [[ -n $_owner_name ]] || return
    [[ -n $_api_token ]] || return
    [[ -n $_app_name ]] || return

    curl --data "{\"status\":\"cancelling\"}"       \
         --header 'Content-Type: application/json'  \
         --header "X-API-Token: $_api_token"        \
         --include                                  \
         --request PATCH                            \
         "https://appcenter.ms/api/v0.1/apps/${_owner_name}/${_app_name}/builds/${APPCENTER_BUILD_ID}"
}

function create_sim_build() {
    local _xcode_project_suffix=${SIM_XCODE_PROJECT##*.}

    if [[ $_xcode_project_suffix == "xcworkspace" ]]; then
        xcodebuild -workspace "$SIM_XCODE_PROJECT"                  \
                   -scheme "$SIM_XCODE_SCHEME"                      \
                   -configuration "$SIM_XCODE_CONFIGURATION"        \
                   -destination 'generic/platform=iOS Simulator'    \
                   -derivedDataPath "$SIM_XCODE_DATA_PATH"          \
                   clean build
    else
        xcodebuild -project "$SIM_XCODE_PROJECT"                    \
                   -scheme "$SIM_XCODE_SCHEME"                      \
                   -configuration "$SIM_XCODE_CONFIGURATION"        \
                   -destination 'generic/platform=iOS Simulator'    \
                   -derivedDataPath "$SIM_XCODE_DATA_PATH"          \
                   clean build
    fi
}

function upload_sim_build() {
    local _build_path="$SIM_XCODE_DATA_PATH"/Build/Products/"$SIM_XCODE_CONFIGURATION"-iphonesimulator/"$SIM_XCODE_APP_NAME"
    local _upload_options=${SIM_WALDO_UPLOAD_OPTIONS:-}

    ${WALDO_CLI_BIN}/waldo upload "$_build_path" $_upload_options
}

create_sim_build || exit
upload_sim_build || exit
cancel_appcenter_build || exit

exit
