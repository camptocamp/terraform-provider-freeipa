package resources

import (
	"context"
	"errors"

	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/ccin2p3/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Host struct {
	provider *provider.Provider
}

type HostModel struct {
	Fqdn           types.String `tfsdk:"fqdn"`
	Description    types.String `tfsdk:"description"`
	Random         types.Bool   `tfsdk:"random"`
	UserPassword   types.String `tfsdk:"userpassword"`
	RandomPassword types.String `tfsdk:"randompassword"`
	Force          types.Bool   `tfsdk:"force"`
}

func (r *Host) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *Host) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"fqdn": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"random": schema.BoolAttribute{
				Optional: true,
			},
			"userpassword": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"randompassword": schema.StringAttribute{
				Sensitive: true,
				Computed:  true,
			},
			"force": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (r *Host) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config HostModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Random.ValueBool() {
		if !config.UserPassword.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid configuration",
				`“userpassword” must not be set when “random” is set to true.`,
			)
		}
	}
}

func (r *Host) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var state, plan HostModel
	var isCreation bool

	if req.Plan.Raw.IsNull() {
		return
	}

	if req.State.Raw.IsNull() {
		isCreation = true
	} else {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Random.ValueBool() {
		if !isCreation && plan.Random.Equal(state.Random) {
			plan.RandomPassword = state.RandomPassword
		}
	} else {
		plan.RandomPassword = types.StringNull()

		if !isCreation && plan.UserPassword.IsNull() {
			if !plan.Random.Equal(state.Random) {
				resp.Diagnostics.AddWarning(
					"Attribute modification considerations",
					`Unsetting (removing or setting to false) “random” attribute will not effectively delete the enrollment password from the host in FreeIPA. It will only be removed from Terraform state.`,
				)
			} else if !plan.UserPassword.Equal(state.UserPassword) {
				resp.Diagnostics.AddWarning(
					"Attribute modification considerations",
					`Removing “userpassword” attribute will not effectively delete the enrollment password from the host in FreeIPA. It will only be removed from Terraform state.`,
				)
			}
		}
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (r *Host) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state HostModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.HostAddArgs{
		Fqdn: plan.Fqdn.ValueString(),
	}

	optArgs := &freeipa.HostAddOptionalArgs{
		Description:  plan.Description.ValueStringPointer(),
		Random:       plan.Random.ValueBoolPointer(),
		Userpassword: plan.UserPassword.ValueStringPointer(),
		Force:        plan.Force.ValueBoolPointer(),
		All:          freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling HostAdd", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().HostAdd(args, optArgs)

	tflog.Trace(ctx, "Called HostAdd", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create host", "Reason: "+err.Error())

		return
	}

	state = plan
	state.RandomPassword = types.StringPointerValue(res.Result.Randompassword)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Host) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state HostModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.HostShowArgs{
		Fqdn: state.Fqdn.ValueString(),
	}

	optArgs := &freeipa.HostShowOptionalArgs{
		All: freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling HostShow", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().HostShow(args, optArgs)

	tflog.Trace(ctx, "Called HostShow", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		var freeipaErr *freeipa.Error

		if errors.As(err, &freeipaErr) && freeipaErr.Code == freeipa.NotFoundCode {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError("Failed to read host", "Reason: "+err.Error())

		return
	}

	state.Description = types.StringPointerValue(res.Result.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Host) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan HostModel
	var hasDiff bool

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.HostModArgs{
		Fqdn: plan.Fqdn.ValueString(),
	}

	optArgs := &freeipa.HostModOptionalArgs{
		Description: plan.Description.ValueStringPointer(),
		All:         freeipa.Bool(true),
	}

	hasDiff = !plan.Description.Equal(state.Description)

	// Do not regenerate a new enrollment password if not requested
	if !plan.Random.Equal(state.Random) && plan.Random.ValueBool() {
		hasDiff = true
		optArgs.Random = plan.Random.ValueBoolPointer()
	}

	// Do not set enrollment password unless a change is requested
	if !plan.UserPassword.Equal(state.UserPassword) && !plan.UserPassword.IsNull() {
		hasDiff = true
		optArgs.Userpassword = plan.UserPassword.ValueStringPointer()
	}

	randomPassword := plan.RandomPassword

	if hasDiff {
		tflog.Trace(ctx, "Calling HostMod", map[string]any{
			"args":     args,
			"opt_args": optArgs,
		})

		res, err := r.provider.Client().HostMod(args, optArgs)

		tflog.Trace(ctx, "Called HostMod", map[string]any{
			"res": res,
			"err": err,
		})

		if err != nil {
			resp.Diagnostics.AddError("Failed to update host", "Reason: "+err.Error())

			return
		}

		if optArgs.Random != nil && *optArgs.Random {
			randomPassword = types.StringPointerValue(res.Result.Randompassword)
		}
	} else {
		tflog.Debug(ctx, "Updated host has no effective difference", map[string]any{
			"fqdn": plan.Fqdn.ValueString(),
		})
	}

	state = plan
	state.RandomPassword = randomPassword

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Host) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state HostModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.HostDelArgs{
		Fqdn: []string{
			state.Fqdn.ValueString(),
		},
	}

	optArgs := &freeipa.HostDelOptionalArgs{}

	tflog.Trace(ctx, "Calling HostDel", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().HostDel(args, optArgs)

	tflog.Trace(ctx, "Called HostDel", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		var freeipaErr *freeipa.Error

		if errors.As(err, &freeipaErr) && freeipaErr.Code != freeipa.NotFoundCode {
			resp.Diagnostics.AddError("Failed to delete host", "Reason: "+err.Error())

			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Host) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := HostModel{
		Fqdn:   types.StringValue(req.ID),
		Random: types.BoolValue(true),
	}

	resp.Diagnostics.AddWarning(
		"Resource import considerations",
		`The enrollment password from the host in FreeIPA cannot be imported. To generate a new random password, unset the “random” attribute and set it again to true.`,
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Host) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":             schema.StringAttribute{},
					"fqdn":           schema.StringAttribute{},
					"description":    schema.StringAttribute{},
					"random":         schema.BoolAttribute{},
					"userpassword":   schema.StringAttribute{},
					"randompassword": schema.StringAttribute{},
					"force":          schema.BoolAttribute{},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var oldState struct {
					ID             types.String `tfsdk:"id"`
					Fqdn           types.String `tfsdk:"fqdn"`
					Description    types.String `tfsdk:"description"`
					Random         types.Bool   `tfsdk:"random"`
					UserPassword   types.String `tfsdk:"userpassword"`
					RandomPassword types.String `tfsdk:"randompassword"`
					Force          types.Bool   `tfsdk:"force"`
				}

				resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)

				if resp.Diagnostics.HasError() {
					return
				}

				newState := HostModel{
					Fqdn:           oldState.Fqdn,
					Description:    oldState.Description,
					Random:         oldState.Random,
					UserPassword:   oldState.UserPassword,
					RandomPassword: oldState.RandomPassword,
					Force:          oldState.Force,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
			},
		},
	}
}

func NewHost(p *provider.Provider) resource.Resource {
	r := &Host{
		provider: p,
	}

	var _ resource.Resource = r
	var _ resource.ResourceWithValidateConfig = r
	var _ resource.ResourceWithModifyPlan = r
	var _ resource.ResourceWithImportState = r
	var _ resource.ResourceWithUpgradeState = r

	return r
}

func init() {
	resources = append(resources, NewHost)
}
