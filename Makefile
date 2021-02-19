.PHONY: all
all: copy-config build fmt lint test

APP=mangindo-feeder
APP_EXECUTABLE="./out/$(APP)"
UNIT_TEST_PACKAGES=$(shell  go list ./...)

# BUILD #

clean:
	GO111MODULE=on go clean
	rm -rf out/

install-linter:
	bin/install-linter

build: clean
	mkdir -p out/
	GO111MODULE=on go build -o $(APP_EXECUTABLE)

# DEV SETUP #

copy-config:
	cp application.yml.sample application.yml

fmt:
	GO111MODULE=on go fmt ./...

vet:
	GO111MODULE=on go vet ./...

lint: install-linter
	GO111MODULE=on ./bin/golangci-lint --config=".golangci.toml" -v run

imports:
	GO111MODULE=on goimports -w -local github.com ./

# TESTS #

install-gotest:
	GO111MODULE=off go get github.com/rakyll/gotest

test: install-gotest
	GO111MODULE=on gotest -p=1 -mod=readonly $(UNIT_TEST_PACKAGES)

test-cov: install-gotest
	GO111MODULE=on gotest -p=1 -mod=readonly -covermode=count -coverprofile=out/coverage.cov $(UNIT_TEST_PACKAGES)

test-ci: copy-config build lint test-cov
	goveralls -coverprofile=out/coverage.cov -service=travis-ci
