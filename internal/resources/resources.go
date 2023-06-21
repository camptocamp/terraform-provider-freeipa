package resources

import (
	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	resources []func(p *provider.Provider) resource.Resource
)

func Resources() []func(p *provider.Provider) resource.Resource {
	return resources
}
