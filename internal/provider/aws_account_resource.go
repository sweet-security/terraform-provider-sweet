package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sweet-security/go-sweet/sweet"
)

var _ resource.Resource = &AwsAccountResource{}
var _ resource.ResourceWithImportState = &AwsAccountResource{}

func NewAwsAccountResource() resource.Resource {
	return &AwsAccountResource{}
}

type AwsAccountResource struct {
	client *sweet.ApiClient
}

type AwsAccountResourceModel struct {
	AccountId  types.String `tfsdk:"account_id"`
	RoleArn    types.String `tfsdk:"role_arn"`
	ExternalId types.String `tfsdk:"external_id"`
	Regions    types.List   `tfsdk:"regions"`
}

func (r *AwsAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_account"
}

func (r *AwsAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Aws Account resource",

		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Aws Account Id",
				Required:            true,
			},
			"role_arn": schema.StringAttribute{
				MarkdownDescription: "Aws Role Arn",
				Required:            true,
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "Aws External Id",
				Optional:            true,
			},
			"regions": schema.ListAttribute{
				MarkdownDescription: "Aws Regions",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *AwsAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sweet.ApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AwsAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AwsAccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	regions := make([]string, 0, len(data.Regions.Elements()))
	diag := data.Regions.ElementsAs(ctx, &regions, false)
	if diag.HasError() {
		resp.Diagnostics.AddError("Cannot convert regions", fmt.Sprintf("Cannot convert regions: %s", diag.Errors()))
		return
	}

	_, err := r.client.AddAwsAccount(&sweet.AwsAccount{
		AccountId:  data.AccountId.ValueString(),
		RoleArn:    data.RoleArn.ValueString(),
		ExternalId: data.ExternalId.ValueString(),
		Regions:    regions,
	})
	if err != nil {
		resp.Diagnostics.AddError("Cannot add account", fmt.Sprintf("Cannot add account: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AwsAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AwsAccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	regions := make([]string, 0, len(data.Regions.Elements()))
	diag := data.Regions.ElementsAs(ctx, &regions, false)
	if diag.HasError() {
		resp.Diagnostics.AddError("Cannot convert regions", fmt.Sprintf("Cannot convert regions: %s", diag.Errors()))
		return
	}
	_, err := r.client.UpdateAwsAccount(&sweet.AwsAccount{
		AccountId:  data.AccountId.ValueString(),
		RoleArn:    data.RoleArn.ValueString(),
		ExternalId: data.ExternalId.ValueString(),
		Regions:    regions,
	})
	if err != nil {
		resp.Diagnostics.AddError("Cannot update account", fmt.Sprintf("Cannot update account: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AwsAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAwsAccount(data.AccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot delete account", fmt.Sprintf("Cannot delete account: %s", err))
		return
	}
}

func (r *AwsAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("account_id"), req, resp)
}
