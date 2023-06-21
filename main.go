package main

import (
	"context"
	"flag"
	"log"

	"github.com/camptocamp/terraform-provider-freeipa/internal/datasources"
	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/camptocamp/terraform-provider-freeipa/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like Delve")
	flag.Parse()

	err := providerserver.Serve(
		context.Background(),
		provider.NewFactory(datasources.DataSources(), resources.Resources()),
		providerserver.ServeOpts{
			Address: "registry.terraform.io/camptocamp/freeipa",
			Debug:   debug,
		})

	if err != nil {
		log.Fatal(err.Error())
	}
}
