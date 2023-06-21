package datasources

import (
	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	dataSources []func(p *provider.Provider) datasource.DataSource
)

func DataSources() []func(p *provider.Provider) datasource.DataSource {
	return dataSources
}
