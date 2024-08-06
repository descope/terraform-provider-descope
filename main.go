package main

import (
	"context"
	"flag"
	"log"

	"github.com/descope/terraform-provider-descope/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version string = "dev"
var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	ctx := context.Background()
	opts := providerserver.ServeOpts{
		Address: "descope/descope",
		Debug:   debug,
	}

	err := providerserver.Serve(ctx, provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
