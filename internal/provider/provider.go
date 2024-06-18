package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/camptocamp/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	version = "dev"
)

type Provider struct {
	dataSources []func() datasource.DataSource
	resources   []func() resource.Resource

	client *freeipa.Client
}

type Model struct {
	Host               types.String `tfsdk:"host"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure"`
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "freeipa"
	resp.Version = version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "FreeIPA host to connect to",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username to use for connection",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Description: "Password to use for connection",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true to disable FreeIPA host TLS certificate verification",
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("FREEIPA_HOST")
	username := os.Getenv("FREEIPA_USERNAME")
	password := os.Getenv("FREEIPA_PASSWORD")
	insecureSkipVerify := false

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.InsecureSkipVerify.IsNull() {
		insecureSkipVerify = config.InsecureSkipVerify.ValueBool()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(path.Root("host"), "Missing FreeIPA host",
			`Host is required to establish a connection to FreeIPA.`,
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(path.Root("username"), "Missing FreeIPA username",
			`Username is required to establish a connection to FreeIPA.`,
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Missing FreeIPA password",
			`Password is required to establish a connection to FreeIPA.`,
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var err error

	p.client, err = freeipa.Connect(
		host,
		&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureSkipVerify,
			},
		},
		username,
		password,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to connect to FreeIPA", "Reason: "+err.Error())
		return
	}

	tflog.Info(ctx, "Successfully connected to FreeIPA", map[string]any{
		"host":     host,
		"username": username,
	})
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return p.dataSources
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return p.resources
}

func (p *Provider) Client() *freeipa.Client {
	return p.client
}

func NewFactory(ds []func(p *Provider) datasource.DataSource, rs []func(p *Provider) resource.Resource) func() provider.Provider {
	return func() provider.Provider {
		p := &Provider{}

		p.dataSources = make([]func() datasource.DataSource, len(ds))

		for i, d := range ds {
			d := d

			p.dataSources[i] = func() datasource.DataSource {
				return d(p)
			}
		}

		p.resources = make([]func() resource.Resource, len(rs))

		for i, r := range rs {
			r := r

			p.resources[i] = func() resource.Resource {
				return r(p)
			}
		}

		var _ provider.Provider = p

		return p
	}
}
