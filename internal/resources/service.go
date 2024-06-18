package resources

import (
	"context"
	"errors"

	"github.com/camptocamp/go-freeipa/freeipa"
	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Service struct {
	provider *provider.Provider
}

type ServiceModel struct {
	KrbHostname   types.String `tfsdk:"krb_hostname"`
	Force         types.Bool   `tfsdk:"force"`
	SkipHostCheck types.Bool   `tfsdk:"skip_host_check"`
}

func (r *Service) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *Service) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"krb_hostname": schema.StringAttribute{
				Description: "Principal name Service principal. Format: <service_type>/<hostname>",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"force": schema.BoolAttribute{
				Description: "Force force principal name even if host not in DNS",
				Optional:    true,
			},
			"skip_host_check": schema.BoolAttribute{
				Description: "Skip host check force service to be created even when host object does not exist to manage it",
				Optional:    true,
			},
		},
	}
}

func (r *Service) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ServiceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.ServiceAddArgs{
		Krbcanonicalname: plan.KrbHostname.ValueString(),
	}
	optArgs := &freeipa.ServiceAddOptionalArgs{
		Force:         plan.Force.ValueBoolPointer(),
		SkipHostCheck: plan.SkipHostCheck.ValueBoolPointer(),
		All:           freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling ServiceAdd", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().ServiceAdd(args, optArgs)
	tflog.Trace(ctx, "Called ServiceAdd", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Service", "Reason: "+err.Error())
		return
	}

	state = plan

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Service) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.ServiceShowArgs{
		Krbcanonicalname: state.KrbHostname.ValueString(),
	}
	optArgs := &freeipa.ServiceShowOptionalArgs{
		All: freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling ServiceShow", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().ServiceShow(args, optArgs)
	tflog.Trace(ctx, "Called ServiceShow", map[string]any{
		"res": res,
		"err": err,
	})
	if err != nil {
		var freeipaErr *freeipa.Error

		if errors.As(err, &freeipaErr) && freeipaErr.Code == freeipa.NotFoundCode {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Failed to read Service", "Reason: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Service) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ServiceModel
	var hasDiff bool

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.ServiceModArgs{
		Krbcanonicalname: plan.KrbHostname.ValueString(),
	}
	optArgs := &freeipa.ServiceModOptionalArgs{
		All: freeipa.Bool(true),
	}

	hasDiff = !plan.KrbHostname.Equal(state.KrbHostname)

	if hasDiff {
		tflog.Trace(ctx, "Calling ServiceMod", map[string]any{
			"args":     args,
			"opt_args": optArgs,
		})

		res, err := r.provider.Client().ServiceMod(args, optArgs)
		tflog.Trace(ctx, "Called ServiceMod", map[string]any{
			"res": res,
			"err": err,
		})
		if err != nil {
			resp.Diagnostics.AddError("Failed to update service", "Reason: "+err.Error())
			return
		}
	} else {
		tflog.Debug(ctx, "Updated service has no effective difference", map[string]any{
			"krb_hostname": plan.KrbHostname.ValueString(),
		})
	}

	state = plan

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Service) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.ServiceDelArgs{
		Krbcanonicalname: []string{state.KrbHostname.ValueString()},
	}

	tflog.Trace(ctx, "Calling ServiceDel", map[string]any{
		"args":     args,
		"opt_args": nil,
	})

	res, err := r.provider.Client().ServiceDel(args, nil)

	tflog.Trace(ctx, "Called ServiceDel", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		var freeipaErr *freeipa.Error
		if errors.As(err, &freeipaErr) && freeipaErr.Code != freeipa.NotFoundCode {
			resp.Diagnostics.AddError("Failed to delete service", "Reason: "+err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Service) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := ServiceModel{
		KrbHostname: types.StringValue(req.ID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func NewService(p *provider.Provider) resource.Resource {
	r := &Service{
		provider: p,
	}

	var _ resource.Resource = r
	var _ resource.ResourceWithImportState = r

	return r
}

func init() {
	resources = append(resources, NewService)
}
