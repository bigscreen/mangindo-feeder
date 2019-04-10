.PHONY: all
all: build fmt vet lint test

APP=mangindo-feeder
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
UNIT_TEST_PACKAGES=$(shell  go list ./... | grep -v "vendor")

APP_EXECUTABLE="./out/$(APP)"

setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint

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

copy-config:
	cp application.yml.sample application.yml
