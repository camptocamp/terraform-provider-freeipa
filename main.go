package main

import (
	"github.com/fiveai/terraform-provider-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: freeipa.Provider,
	})
}
