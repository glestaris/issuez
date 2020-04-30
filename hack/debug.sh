#!/bin/bash
set -e

HACK_TEST_RUNNER="dlv" exec ./hack/test.sh $@
