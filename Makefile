.PHONY: all
all: build fmt vet lint test

APP=mangindo-feeder
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
UNIT_TEST_PACKAGES=$(shell  go list ./... | grep -v "vendor")

APP_EXECUTABLE="./out/$(APP)"

setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint
	go get -u github.com/axw/gocov/gocov

build-deps:
	dep ensure

update-deps:
	dep ensure

compile:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE)

build: build-deps compile fmt

install:
	go install ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@for p in $(UNIT_TEST_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

test:
	ENVIRONMENT=test go test $(UNIT_TEST_PACKAGES) -p=1

test-ci: copy-config build-deps compile fmt test-cov test-cov-report

test-cov:
	gocov test ${ALL_PACKAGES} > coverage.json

test-cov-report:
	@echo "\nGENERATING TEST REPORT."
	gocov report coverage.json

copy-config:
	cp application.yml.sample application.yml
