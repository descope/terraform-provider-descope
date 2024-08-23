package main

import (
	"context"
	"flag"
	"log"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version = "dev"
	debug   bool
)

func main() {
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	ctx := context.Background()
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/descope/descope",
		Debug:   debug,
	}

	err := providerserver.Serve(ctx, provider.NewDescopeProvider(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
