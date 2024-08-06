package main

import (
	"context"
	"flag"
	"log"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name descope

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

	err := providerserver.Serve(ctx, provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
