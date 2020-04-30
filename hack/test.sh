#!/bin/bash
set -e

debug() {
    msg=$1
    if ! [ -z ${DEBUG} ]; then
        echo "[DEBUG] ${msg}"
    fi
}

if echo "$@" | grep "acceptance" 2>&1 >/dev/null; then
    debug "Building issuez..."

    temp_out_file=$(mktemp)
    make ./issuez 2>&1 >"${temp_out_file}"
    debug && cat "${temp_out_file}"
    rm -f "${temp_out_file}"

    debug "Build artifact is ready."

    export TEST_JIRA_PM_EXE_PATH="${PWD}/issuez"
fi

runner="richgo"
if ! [ -z "${HACK_TEST_RUNNER}" ]; then
    runner="${HACK_TEST_RUNNER}"
fi
debug "Executing tests using runner '${runner}':"

extra_args=""
if ! [ -z "${HACK_TEST_EXTRA_ARGS}" ]; then
    extra_args="${HACK_TEST_EXTRA_ARGS}"
fi
cmd="${runner} test ${extra_args} $@"
debug "cmd='${cmd}'"

eval "${cmd}"
