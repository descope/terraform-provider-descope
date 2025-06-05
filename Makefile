.DEFAULT_GOAL := help

.PHONY:  help dev install test testacc testcoverage testcleanup terragen docs terraformrc lint ensure-linter ensure-gitleaks ensure-descope ensure-courtney ensure-brew ensure-go
.SILENT: help dev install test testacc testcoverage testcleanup terragen docs terraformrc lint ensure-linter ensure-gitleaks ensure-descope ensure-courtney ensure-brew ensure-go

ifneq ($(tests),)
  flags := $(flags) -count 1 -run '$(tests)'
endif

env ?= tools/config.env
ifneq ($(wildcard $(env)),)
  ifeq ($(DESCOPE_PROJECT_ID),)
    export DESCOPE_PROJECT_ID = $(shell cat $(env) | grep DESCOPE_PROJECT_ID | sed 's/^.*=//')
  endif
  ifeq ($(DESCOPE_MANAGEMENT_KEY),)
    export DESCOPE_MANAGEMENT_KEY = $(shell cat $(env) | grep DESCOPE_MANAGEMENT_KEY | sed 's/^.*=//')
  endif
  ifeq ($(DESCOPE_BASE_URL),)
    export DESCOPE_BASE_URL = $(shell cat $(env) | grep DESCOPE_BASE_URL | sed 's/^.*=//')
  endif
endif

help: Makefile ## this help message
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dev: install terraformrc ## prepares development environment for running the provider

install: ensure-go ## installs terraform-provider-descope to $GOPATH/bin
	mkdir -p "$$GOPATH/bin"
	go install .
	echo The $$'\e[33m'terraform-provider-descope$$'\e[0m' tool has been installed to $$GOPATH/bin

test: ensure-go ## runs unit tests
	go test -v -timeout 30m $(flags) ./...

testacc: ensure-go ## runs acceptance and unit tests
	TF_ACC=1 go test -v -timeout 120m $(flags) ./...

testcoverage: ensure-go ensure-courtney ## runs all tests and computes test coverage
	TF_ACC=1 go test -v -race -timeout 120m -coverpkg=./... -coverprofile=coverage.raw -covermode=atomic ./...
	cat coverage.raw | grep -v -e "\/tools\/.*\.go\:.*" | grep -v -e ".*\/main\.go\:.*" > coverage.out
	rm -f coverage.raw
	courtney -l coverage.out
	go tool cover -func coverage.out | grep total | awk '{print $$3}'
	go tool cover -html=coverage.out -o coverage.html

testcleanup: ensure-descope ## cleans up redundant projects after running tests
	descope project list | grep '"name":"testacc-.*' | sed -e 's/.*"id":"\([^"]*\)".*/\1/' | xargs -I {} descope project delete {} --force

terragen: ensure-go ## runs the terragen tool to generate code and model documentation
	go run tools/terragen/main.go $(flags)

docs: ensure-go ## runs tfplugindocs to generate documentation for the registry 
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.4 generate -provider-name descope

terraformrc:
	echo 'provider_installation {'                      > ~/.terraformrc
	echo '  dev_overrides {'                            >> ~/.terraformrc
	echo '    "descope/descope" = "'$$GOPATH'/bin"'     >> ~/.terraformrc
	echo '  }'                                          >> ~/.terraformrc
	echo '  direct {}'                                  >> ~/.terraformrc
	echo '}'                                            >> ~/.terraformrc
	echo The $$'\e[33m'.terraformrc$$'\e[0m' file has been created in $$HOME

lint: ensure-linter ensure-gitleaks ## check for linter and gitleaks failures
	golangci-lint --config .github/actions/ci/lint/golangci.yml run
	gitleaks protect --redact -v -c .github/actions/ci/leaks/gitleaks.toml
	gitleaks detect --redact -v -c .github/actions/ci/leaks/gitleaks.toml

ensure-linter: ensure-go
	TOOL_VERSION=$$(cat .github/actions/ci/lint/action.yml | grep 'version: ' | sed 's/.*version: \(v.*\)$$/\1/'); \
	if ! command -v golangci-lint &> /dev/null; then \
	    echo Installing the $$'\e[33m'golangci-lint$$'\e[0m' tool... ;\
	    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$$TOOL_VERSION ;\
	else \
	    CURRENT_VERSION=$$(golangci-lint --version | sed 's/.*version \(v.*\) built with.*/\1/' ) ;\
	    if [ "$$CURRENT_VERSION" != "$$TOOL_VERSION" ]; then \
	        echo Updating the $$'\e[33m'golangci-lint$$'\e[0m' tool to version $$'\e[33m'$$TOOL_VERSION$$'\e[0m'... ;\
	        go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$$TOOL_VERSION ;\
	    fi ;\
	fi

ensure-gitleaks: ensure-brew
	if ! command -v gitleaks &> /dev/null; then \
	    echo Installing the $$'\e[33m'gitleaks$$'\e[0m' tool... ;\
	    brew install gitleaks ;\
	fi

ensure-descope: ensure-brew
	if ! command -v descope &> /dev/null; then \
	    echo Installing the $$'\e[33m'descope$$'\e[0m' CLI tool... ;\
	    brew install descope ;\
	fi

ensure-courtney: ensure-go
	if ! command -v courtney &> /dev/null; then \
	    echo Installing the $$'\e[33m'courtney$$'\e[0m' tool... ;\
	    go install github.com/dave/courtney@master ;\
	fi

ensure-brew:
	if ! command -v brew &> /dev/null; then \
	    echo \\nInstall the brew tool from $$'\e[33m'https://brew.sh$$'\e[0m'\\n ;\
	    false ;\
	fi

ensure-go:
	if ! command -v go &> /dev/null; then \
	    echo \\nInstall the go compiler from $$'\e[33m'https://go.dev/dl$$'\e[0m'\\n ;\
	    false ;\
	fi
	if [ -z "$$GOPATH" ]; then \
	    echo \\nThe $$'\e[33m'GOPATH$$'\e[0m' environment variable must be defined, see $$'\e[33m'https://go.dev/wiki/GOPATH$$'\e[0m'\\n ;\
	    false ;\
	fi
