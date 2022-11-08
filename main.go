package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	p "github.com/hashicorp/terraform-provider-archive/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), p.New, providerserver.ServeOpts{
		Address:         "registry.terraform.io/hashicorp/archive",
		Debug:           debug,
		ProtocolVersion: 5,
	})
	if err != nil {
		log.Fatal(err)
	}
}
