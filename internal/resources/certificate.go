package resources

import (
	"context"
	"errors"
	"strconv"

	"github.com/camptocamp/go-freeipa/freeipa"
	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Certificate struct {
	provider *provider.Provider
}

type CertificateModel struct {
	Principal    types.String `tfsdk:"principal"`
	CSR          types.String `tfsdk:"csr"`
	SerialNumber types.Int64  `tfsdk:"serial_number"`
}

func (r *Certificate) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *Certificate) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"principal": schema.StringAttribute{
				Description: "Principal for this certificate (e.g. HTTP/test.example.com)",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"csr": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"serial_number": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *Certificate) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CertificateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.CertRequestArgs{
		Principal: plan.Principal.ValueString(),
		Csr:       plan.CSR.ValueString(),
	}
	optArgs := &freeipa.CertRequestOptionalArgs{
		All: freeipa.Bool(true),
		Raw: freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling CertRequest", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().CertRequest(args, optArgs)
	tflog.Trace(ctx, "Called CertRequest", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Certificate", "Reason: "+err.Error())
		return
	}

	state = plan
	state.SerialNumber = types.Int64Value(res.Result.(int64))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Certificate) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CertificateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.CertShowArgs{
		SerialNumber: int(state.SerialNumber.ValueInt64()),
	}
	optArgs := &freeipa.CertShowOptionalArgs{
		All: freeipa.Bool(true),
		Raw: freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling CertShow", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().CertShow(args, optArgs)
	tflog.Trace(ctx, "Called CertShow", map[string]any{
		"res": res,
		"err": err,
	})
	if err != nil {
		var freeipaErr *freeipa.Error

		if errors.As(err, &freeipaErr) && freeipaErr.Code == freeipa.NotFoundCode {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Failed to read Certificate", "Reason: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Certificate) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan CertificateModel
	var hasDiff bool

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasDiff = !plan.Principal.Equal(state.Principal) || !plan.CSR.Equal(state.CSR)
	if !hasDiff {
		tflog.Debug(ctx, "Updated certificate has no effective difference", map[string]any{
			"principal": plan.Principal.ValueString(),
			"csr":       plan.CSR.ValueString(),
		})
		return
	}

	state = plan
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Certificate) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CertificateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.CertRevokeArgs{
		SerialNumber: int(state.SerialNumber.ValueInt64()),
	}

	tflog.Trace(ctx, "Calling CertRevoke", map[string]any{
		"args":     args,
		"opt_args": nil,
	})

	res, err := r.provider.Client().CertRevoke(args, nil)

	tflog.Trace(ctx, "Called CertRevoke", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		var freeipaErr *freeipa.Error
		if errors.As(err, &freeipaErr) && freeipaErr.Code != freeipa.NotFoundCode {
			resp.Diagnostics.AddError("Failed to delete certificate", "Reason: "+err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Certificate) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to import state", err.Error())
		return
	}

	state := CertificateModel{
		SerialNumber: types.Int64Value(int64(id)),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func NewCertificate(p *provider.Provider) resource.Resource {
	r := &Certificate{
		provider: p,
	}

	var _ resource.Resource = r
	var _ resource.ResourceWithImportState = r

	return r
}

func init() {
	resources = append(resources, NewCertificate)
}
