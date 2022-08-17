NAME:=roc
VERSION:="DEV"
OS:=linux
DATE:="$(shell date +%F)"
ARCH:=amd64
SOURCE_FILES:=$$(find . -name '*.go' | grep -v vendor)
BUILD_FLAGS:=-mod=vendor -a -tags netgo
BINPATH:=$(PWD)/bin
RELEASEPATH:=$(PWD)/release

LINUX_AMD64_BIN_NAME:=$(NAME)-linux-amd64
DARWIN_AMD64_BIN_NAME:=$(NAME)-darwin-amd64
DARWIN_ARM64_BIN_NAME:=$(NAME)-darwin-arm64

all: build_roc_linux

build_roc:
	GO111MODULE=on CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
		    go build $(BUILD_FLAGS) \
		    -ldflags '-w -extldflags "-static" -X github.com/rapidsai/$(NAME)/internal/build.Version=$(VERSION) -X github.com/rapidsai/$(NAME)/internal/build.Date=$(DATE)' \
		    -o $(BINPATH)/$(BINNAME) \
		    ./cmd/$(NAME)

build_roc_linux: BINNAME=$(LINUX_AMD64_BIN_NAME)
build_roc_linux: build_roc
	strip $(BINPATH)/$(BINNAME)

build_roc_darwin_amd64: BINNAME=$(DARWIN_AMD64_BIN_NAME)
build_roc_darwin_amd64: OS="darwin"
build_roc_darwin_amd64: build_roc

build_roc_darwin_arm64: BINNAME=$(DARWIN_ARM64_BIN_NAME)
build_roc_darwin_arm64: OS="darwin"
build_roc_darwin_arm64: ARCH="arm64"
build_roc_darwin_arm64: build_roc

release_zips:
	@mkdir -p $(RELEASEPATH)
	make build_roc_linux && zip -j $(RELEASEPATH)/roc-linux-amd64.zip $(BINPATH)/$(LINUX_AMD64_BIN_NAME)
	make build_roc_darwin_amd64 && zip -j $(RELEASEPATH)/roc-darwin-amd64.zip $(BINPATH)/$(DARWIN_AMD64_BIN_NAME)
	make build_roc_darwin_arm64 && zip -j $(RELEASEPATH)/roc-darwin-arm64.zip $(BINPATH)/$(DARWIN_ARM64_BIN_NAME)

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
	-rm -rf $(BINPATH)
	-rm -rf $(RELEASEPATH)

.PHONY: clean install test coverage lint fmt
