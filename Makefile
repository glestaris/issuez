.PHONY: release unit-test integration-test acceptance-test test coverage clean

help:
	@echo "release .......................... Build issuez CLI supported OSs"
	@echo "test-unit ........................ Run unit tests"
	@echo "test-integration ................. Run integration tests"
	@echo "test-acceptance .................. Run acceptance tests"
	@echo "test ............................. Run all tests"
	@echo "coverage ......................... Measure code coverage"
	@echo "clean ............................ Clean build artifacts"

GO_FILES=$(shell find . -path '*.go' -not -name '*_test.go')

release: ./dist/issuez-linux ./dist/issuez-darwin

./dist/issuez-linux: ${GO_FILES}
	GOOS=linux go build -o ./dist/issuez-linux .

./dist/issuez-darwin: ${GO_FILES}
	GOOS=darwin go build -o ./dist/issuez-darwin .

# Testing

GO_UNIT_TESTS=$(shell go list ./... | grep -v acceptance | grep -v integration)

test-unit:
	./hack/test.sh ${GO_UNIT_TESTS}

test-integration:
	./hack/test.sh ./integration

test-acceptance:
	./hack/test.sh ./acceptance

test:
	@echo '== UNIT TESTS =================='
	./hack/test.sh ${GO_UNIT_TESTS}
	@echo

	@echo '== INTEGRATION TESTS ==========='
	./hack/test.sh ./integration
	@echo

	@echo '== ACCEPTANCE TESTS ============'
	./hack/test.sh ./acceptance

coverage: clean
	@echo '== TESTS ========================'
	HACK_TEST_EXTRA_ARGS="-test.coverprofile cp.out" ./hack/test.sh ${GO_UNIT_TESTS}
	@echo

	@echo '== COVERAGE ANALYSIS ==========='
	go tool cover -html=cp.out -o=./coverage.html

# Housekeeping

clean:
	rm -rf ./dist ./cp.out ./coverage.html
	go clean -testcache
