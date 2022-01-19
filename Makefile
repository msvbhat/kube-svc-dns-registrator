.PHONY: build clean

GOOS :=
CGO_CXXFLAGS :=
CMAKE_EXTRA_FLAGS :=
UNAME := $(shell uname -s)
ifeq ($(UNAME),Linux)
	GOOS = linux
endif
ifeq ($(UNAME),Darwin)
	GOOS = darwin
endif

build: ## Build the binary file
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=amd64 go build -o app .

clean: ## Remove the compiled binary file
	rm -f app
