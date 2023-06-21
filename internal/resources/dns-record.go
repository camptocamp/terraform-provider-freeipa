package resources

import (
	"context"
	"errors"
	"strings"

	"github.com/camptocamp/terraform-provider-freeipa/internal/provider"
	"github.com/ccin2p3/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type DnsRecord struct {
	provider *provider.Provider
}

type DnsRecordModel struct {
	Name     types.String `tfsdk:"idnsname"`
	ZoneName types.String `tfsdk:"dnszoneidnsname"`
	Class    types.String `tfsdk:"dnsclass"`
	Type     types.String `tfsdk:"type"`
	TTL      types.Int64  `tfsdk:"dnsttl"`
	Records  types.Set    `tfsdk:"records"`
}

func (r *DnsRecord) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *DnsRecord) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"idnsname": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dnszoneidnsname": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dnsclass": schema.StringAttribute{
				Optional:           true,
				DeprecationMessage: "Only “IN” DNS class is supported.",
			},
			"type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dnsttl": schema.Int64Attribute{
				Optional: true,
			},
			"records": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *DnsRecord) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state DnsRecordModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var zone any = plan.ZoneName.ValueString()

	var ttl *int

	if !plan.TTL.IsNull() {
		ttl = new(int)
		*ttl = int(plan.TTL.ValueInt64())
	}

	var records []string

	resp.Diagnostics.Append(plan.Records.ElementsAs(ctx, &records, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.DnsrecordAddArgs{
		Idnsname: plan.Name.ValueString(),
	}

	optArgs := &freeipa.DnsrecordAddOptionalArgs{
		Dnszoneidnsname: &zone,
		Dnsttl:          ttl,
		All:             freeipa.Bool(true),
	}

	switch plan.Type.ValueString() {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Txtrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	tflog.Trace(ctx, "Calling DnsrecordAdd", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().DnsrecordAdd(args, optArgs)

	tflog.Trace(ctx, "Called DnsrecordAdd", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create DNS record", "Reason: "+err.Error())

		return
	}

	state = plan

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DnsRecord) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DnsRecordModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var zone any = state.ZoneName.ValueString()

	args := &freeipa.DnsrecordShowArgs{
		Idnsname: state.Name.ValueString(),
	}

	optArgs := &freeipa.DnsrecordShowOptionalArgs{
		Dnszoneidnsname: &zone,
		All:             freeipa.Bool(true),
	}

	tflog.Trace(ctx, "Calling DnsrecordShow", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().DnsrecordShow(args, optArgs)

	tflog.Trace(ctx, "Called DnsrecordShow", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		var freeipaErr *freeipa.Error

		if errors.As(err, &freeipaErr) && freeipaErr.Code == freeipa.NotFoundCode {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError("Failed to read DNS record", "Reason: "+err.Error())

		return
	}

	var ttl *int64

	if res.Result.Dnsttl != nil {
		ttl = new(int64)
		*ttl = int64(*res.Result.Dnsttl)
	}

	var records *[]string

	switch state.Type.ValueString() {
	case "A":
		records = res.Result.Arecord
	case "AAAA":
		records = res.Result.Aaaarecord
	case "CNAME":
		records = res.Result.Cnamerecord
	case "MX":
		records = res.Result.Mxrecord
	case "NS":
		records = res.Result.Nsrecord
	case "PTR":
		records = res.Result.Ptrrecord
	case "SRV":
		records = res.Result.Srvrecord
	case "TXT":
		records = res.Result.Txtrecord
	case "SSHFP":
		records = res.Result.Sshfprecord
	}

	var diags diag.Diagnostics

	state.TTL = types.Int64PointerValue(ttl)

	if records != nil {
		state.Records, diags = types.SetValueFrom(ctx, types.StringType, *records)

		resp.Diagnostics.Append(diags...)
	} else {
		state.Records = types.SetValueMust(types.StringType, []attr.Value{})
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DnsRecord) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan DnsRecordModel
	var hasDiff bool

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var zone any = plan.ZoneName.ValueString()

	var ttl *int

	if !plan.TTL.IsNull() {
		ttl = new(int)
		*ttl = int(plan.TTL.ValueInt64())
	}

	var records []string

	resp.Diagnostics.Append(plan.Records.ElementsAs(ctx, &records, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := &freeipa.DnsrecordModArgs{
		Idnsname: plan.Name.ValueString(),
	}

	optArgs := &freeipa.DnsrecordModOptionalArgs{
		Dnszoneidnsname: &zone,
		Dnsttl:          ttl,
		All:             freeipa.Bool(true),
	}

	switch plan.Type.ValueString() {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Txtrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	hasDiff = !plan.TTL.Equal(state.TTL) || !plan.Records.Equal(state.Records)

	if hasDiff {
		tflog.Trace(ctx, "Calling DnsrecordMod", map[string]any{
			"args":     args,
			"opt_args": optArgs,
		})

		res, err := r.provider.Client().DnsrecordMod(args, optArgs)

		tflog.Trace(ctx, "Called DnsrecordMod", map[string]any{
			"res": res,
			"err": err,
		})

		if err != nil {
			resp.Diagnostics.AddError("Failed to update DNS record", "Reason: "+err.Error())

			return
		}
	} else {
		tflog.Debug(ctx, "Updated DNS record has no effective difference", map[string]any{
			"name":      plan.Name.ValueString(),
			"zone_name": plan.ZoneName.ValueString(),
			"type":      plan.Type.ValueString(),
		})
	}

	state = plan

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DnsRecord) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DnsRecordModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var zone any = state.ZoneName.ValueString()

	deleteAllRecords := true

	args := &freeipa.DnsrecordDelArgs{
		Idnsname: state.Name.ValueString(),
	}

	optArgs := &freeipa.DnsrecordDelOptionalArgs{
		Dnszoneidnsname: &zone,
		DelAll:          &deleteAllRecords,
	}

	tflog.Trace(ctx, "Calling DnsrecordDel", map[string]any{
		"args":     args,
		"opt_args": optArgs,
	})

	res, err := r.provider.Client().DnsrecordDel(args, optArgs)

	tflog.Trace(ctx, "Called DnsrecordDel", map[string]any{
		"res": res,
		"err": err,
	})

	if err != nil {
		var freeipaErr *freeipa.Error

		if errors.As(err, &freeipaErr) && freeipaErr.Code != freeipa.NotFoundCode {
			resp.Diagnostics.AddError("Failed to delete DNS record", "Reason: "+err.Error())

			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DnsRecord) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := strings.Split(req.ID, "/")

	if len(id) != 3 {
		resp.Diagnostics.AddError("Invalid ID format", "Expected ID format is “<record name>/<zone name>/<record type>”")

		return
	}

	state := DnsRecordModel{
		Name:     types.StringValue(id[0]),
		ZoneName: types.StringValue(id[1]),
		Type:     types.StringValue(id[2]),
		Records:  types.SetNull(types.StringType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DnsRecord) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":              schema.StringAttribute{},
					"idnsname":        schema.StringAttribute{},
					"dnszoneidnsname": schema.StringAttribute{},
					"dnsclass":        schema.StringAttribute{},
					"type":            schema.StringAttribute{},
					"dnsttl":          schema.Int64Attribute{},
					"records": schema.SetAttribute{
						ElementType: types.StringType,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var oldState struct {
					ID       types.String `tfsdk:"id"`
					Name     types.String `tfsdk:"idnsname"`
					ZoneName types.String `tfsdk:"dnszoneidnsname"`
					Class    types.String `tfsdk:"dnsclass"`
					Type     types.String `tfsdk:"type"`
					TTL      types.Int64  `tfsdk:"dnsttl"`
					Records  types.Set    `tfsdk:"records"`
				}

				resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)

				if resp.Diagnostics.HasError() {
					return
				}

				newState := DnsRecordModel{
					Name:     oldState.Name,
					ZoneName: oldState.ZoneName,
					Type:     oldState.Type,
					TTL:      oldState.TTL,
					Records:  oldState.Records,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
			},
		},
	}
}

func NewDnsRecord(p *provider.Provider) resource.Resource {
	r := &DnsRecord{
		provider: p,
	}

	var _ resource.Resource = r
	var _ resource.ResourceWithImportState = r
	var _ resource.ResourceWithUpgradeState = r

	return r
}

func init() {
	resources = append(resources, NewDnsRecord)
}
