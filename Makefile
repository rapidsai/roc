NAME:=roc
VERSION:="DEV"
OS:=linux
DATE:="$(shell date +%F)"
ARCH:=amd64
SOURCE_FILES:=$$(find . -name '*.go' | grep -v vendor)
#BUILD_FLAGS:=-mod=vendor -a -tags netgo
BUILD_FLAGS:= -a -tags netgo
BINPATH:=$(PWD)/bin

all: build_cli

build_cli:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
		    go build $(BUILD_FLAGS) \
		    -ldflags '-w -extldflags "-static" -X github.com/rapidsai/rapids-ops-cli/internal/build.Version=$(VERSION) -X github.com/rapidsai/rapids-ops-cli/internal/build.Date=$(DATE)' \
		    -o $(BINPATH)/$(NAME) \
		    ./cmd/rapids-ops-cli
	strip $(BINPATH)/$(NAME)

test:
	@go test -v ./pkg/...

fmt:
	@goimports -w $(SOURCE_FILES)
	@go fmt ./...

lint: fmt
	@golangci-lint run ./...

install: ## Build and install locally the binary (dev purpose)
	go install .

coverage:
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	@rm coverage.out

clean:
	@rm -r $(BINPATH)

.PHONY: clean install test coverage lint fmt
