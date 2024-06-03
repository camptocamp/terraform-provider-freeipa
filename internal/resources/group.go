package resources

import (
	"context"
	"errors"

	"github.com/camptocamp/go-freeipa/freeipa"
	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Group struct {
	provider *provider.Provider
}

type GroupModel struct {
	Name        types.String `tfsdk:"cn"`
	Description types.String `tfsdk:"description"`
	GID         types.Int64  `tfsdk:"gidnumber"`
	NonPosix    types.Bool   `tfsdk:"nonposix"`
	External    types.Bool   `tfsdk:"external"`
}

func (r *Group) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *Group) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"cn": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"gidnumber": schema.Int64Attribute{
				Optional: true,
			},
			"nonposix": schema.BoolAttribute{
				Description: "Create as a non-POSIX group",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"external": schema.BoolAttribute{
				Description: "Allow adding external non-IPA members from trusted domains",
				Optional:    true,
			},
		},
	}
}

func (r *Group) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state GroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.GroupAddArgs{
		Cn: plan.Name.ValueString(),
	}

	var gid *int
	if !plan.GID.IsNull() {
		gid = new(int)
		*gid = int(plan.GID.ValueInt64())
	}

	optArgs := &freeipa.GroupAddOptionalArgs{
		Description: plan.Description.ValueStringPointer(),
		Gidnumber:   gid,
		Nonposix:    plan.NonPosix.ValueBoolPointer(),
		External:    plan.External.ValueBoolPointer(),
		All:         freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling GroupAdd", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().GroupAdd(args, optArgs)
	tflog.Trace(ctx, "Called GroupAdd", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Group", "Reason: "+err.Error())
		return
	}

	state = plan

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Group) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.GroupShowArgs{
		Cn: state.Name.ValueString(),
	}
	optArgs := &freeipa.GroupShowOptionalArgs{
		All: freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling GroupShow", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().GroupShow(args, optArgs)
	tflog.Trace(ctx, "Called GroupShow", map[string]any{
		"res": res,
		"err": err,
	})
	if err != nil {
		var freeipaErr *freeipa.Error

		if errors.As(err, &freeipaErr) && freeipaErr.Code == freeipa.NotFoundCode {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Failed to read Group", "Reason: "+err.Error())
		return
	}

	var gid *int64
	if res.Result.Gidnumber != nil {
		gid = new(int64)
		*gid = int64(*res.Result.Gidnumber)
	}
	state.GID = types.Int64PointerValue(gid)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Group) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan GroupModel
	var hasDiff bool

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var gid *int
	if !plan.GID.IsNull() {
		gid = new(int)
		*gid = int(plan.GID.ValueInt64())
	}

	args := &freeipa.GroupModArgs{
		Cn: plan.Name.ValueString(),
	}
	optArgs := &freeipa.GroupModOptionalArgs{
		Description: plan.Description.ValueStringPointer(),
		Gidnumber:   gid,
		External:    plan.External.ValueBoolPointer(),
		All:         freeipa.Bool(true),
	}

	hasDiff = !plan.GID.Equal(state.GID) || !plan.Description.Equal(state.Description) ||
		!plan.External.Equal(state.External)

	if hasDiff {
		tflog.Trace(ctx, "Calling GroupMod", map[string]any{
			"args":     args,
			"opt_args": optArgs,
		})

		res, err := r.provider.Client().GroupMod(args, optArgs)
		tflog.Trace(ctx, "Called GroupMod", map[string]any{
			"res": res,
			"err": err,
		})
		if err != nil {
			resp.Diagnostics.AddError("Failed to update group", "Reason: "+err.Error())
			return
		}
	} else {
		tflog.Debug(ctx, "Updated DNS record has no effective difference", map[string]any{
			"name":        plan.Name.ValueString(),
			"description": plan.Description.ValueString(),
			"gid":         plan.GID.ValueInt64(),
			"external":    plan.External.ValueBool(),
		})
	}

	state = plan

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Group) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.GroupDelArgs{
		Cn: []string{state.Name.ValueString()},
	}

	tflog.Trace(ctx, "Calling GroupDel", map[string]any{
		"args":     args,
		"opt_args": nil,
	})

	res, err := r.provider.Client().GroupDel(args, nil)

	tflog.Trace(ctx, "Called GroupDel", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		var freeipaErr *freeipa.Error
		if errors.As(err, &freeipaErr) && freeipaErr.Code != freeipa.NotFoundCode {
			resp.Diagnostics.AddError("Failed to delete group", "Reason: "+err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Group) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := GroupModel{
		Name: types.StringValue(req.ID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func NewGroup(p *provider.Provider) resource.Resource {
	r := &Group{
		provider: p,
	}

	var _ resource.Resource = r
	var _ resource.ResourceWithImportState = r

	return r
}

func init() {
	resources = append(resources, NewGroup)
}
