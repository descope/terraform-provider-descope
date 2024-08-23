.DEFAULT_GOAL := help

.PHONY:  help setup test testacc testcoverage install terragen docs terraformrc ensure-courtney ensure-go
.SILENT: help setup test testacc testcoverage install terragen docs terraformrc ensure-courtney ensure-go

help: Makefile ## this help message
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: install terraformrc ## prepares local environment for running the provider

test: ensure-go ## runs unit tests
	go test -v -timeout 30m $(ARGS) ./...

testacc: ensure-go ## runs acceptance and unit tests
	TF_ACC=1 go test -v -timeout 120m $(ARGS) ./...

testcoverage: ensure-go ensure-courtney ## runs all tests and computes test coverage
	TF_ACC=1 go test -race -timeout 120m -coverpkg=./... -coverprofile=coverage.raw -covermode=atomic $(ARGS) ./...
	cat coverage.raw | grep -v -e "\/tools\/.*\.go\:.*" | grep -v -e ".*\/main\.go\:.*" > coverage.out
	rm -f coverage.raw
	courtney -l coverage.out
	go tool cover -func coverage.out | grep total | awk '{print $$3}'
	go tool cover -html=coverage.out -o coverage.html

install: ensure-go ## installs the descope command line tool to $GOPATH/bin
	mkdir -p "$$GOPATH/bin"
	go install .
	echo The $$'\e[33m'terraform-provider-descope$$'\e[0m' tool has been installed to $$GOPATH/bin

terragen: ensure-go ## runs the terragen tool to generate code and model documentation
	go run tools/terragen/main.go $(ARGS)

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

ensure-courtney:
	if ! command -v courtney &> /dev/null; then \
	    echo \\nInstall the courtney tool with $$'\e[33m'go install github.com/dave/courtney@master$$'\e[0m'\\n ;\
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
