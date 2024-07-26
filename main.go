package main

import (
	"context"
	"flag"
	"log"

	"github.com/camptocamp/terraform-provider-freeipa/freeipa"
	"github.com/camptocamp/terraform-provider-freeipa/internal/datasources"
	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/camptocamp/terraform-provider-freeipa/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name freeipa

func main() {
	var debug bool

	ctx := context.Background()

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like Delve")
	flag.Parse()

	providers := []func() tfprotov5.ProviderServer{
		freeipa.Provider().GRPCProvider, // legacy provider using terraform-sdk-v2
		providerserver.NewProtocol5(provider.NewFactory(datasources.DataSources(), resources.Resources())()), // new provider built using terraform-plugin-framework
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/camptocamp/freeipa",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
