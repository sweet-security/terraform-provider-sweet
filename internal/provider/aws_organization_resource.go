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

var _ resource.Resource = &AwsOrganizationResource{}
var _ resource.ResourceWithImportState = &AwsOrganizationResource{}

func NewAwsOrganizationResource() resource.Resource {
	return &AwsOrganizationResource{}
}

type AwsOrganizationResource struct {
	client *sweet.ApiClient
}

type AwsOrganizationResourceModel struct {
	AccountId            types.String `tfsdk:"account_id"`
	RoleArn              types.String `tfsdk:"role_arn"`
	RoleNameParameterArn types.String `tfsdk:"role_name_parameter_arn"`
	ExternalId           types.String `tfsdk:"external_id"`
	Regions              types.List   `tfsdk:"regions"`
}

func (r *AwsOrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_organization"
}

func (r *AwsOrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Aws Organization resource",

		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Aws Root Account Id",
				Required:            true,
			},
			"role_arn": schema.StringAttribute{
				MarkdownDescription: "Aws Role Arn",
				Required:            true,
			},
			"role_name_parameter_arn": schema.StringAttribute{
				MarkdownDescription: "Aws Role Name Parameter Arn",
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

func (r *AwsOrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AwsOrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AwsOrganizationResourceModel

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

	_, err := r.client.AddAwsOrganization(&sweet.AwsOrganization{
		AccountId:            data.AccountId.ValueString(),
		RoleArn:              data.RoleArn.ValueString(),
		RoleNameParameterArn: data.RoleNameParameterArn.ValueString(),
		ExternalId:           data.ExternalId.ValueString(),
		Regions:              regions,
	})
	if err != nil {
		resp.Diagnostics.AddError("Cannot add organization", fmt.Sprintf("Cannot add organization: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsOrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AwsOrganizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationRes, err := r.client.GetAwsOrganization(data.AccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot read organization", fmt.Sprintf("Cannot read organization: %s", err))
	}
	data.AccountId = types.StringValue(organizationRes.AccountId)
	data.RoleArn = types.StringValue(organizationRes.RoleArn)
	data.RoleNameParameterArn = types.StringValue(organizationRes.RoleNameParameterArn)
	data.ExternalId = types.StringValue(organizationRes.ExternalId)
	regionsList, diags := types.ListValueFrom(ctx, types.StringType, organizationRes.Regions)
	if diags.HasError() {
		resp.Diagnostics.AddError("Cannot read regions", fmt.Sprintf("Cannot read regions: %s", organizationRes.Regions))
		return
	}
	data.Regions = regionsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsOrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AwsOrganizationResourceModel

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
	_, err := r.client.UpdateAwsOrganization(&sweet.AwsOrganization{
		AccountId:            data.AccountId.ValueString(),
		RoleArn:              data.RoleArn.ValueString(),
		RoleNameParameterArn: data.RoleNameParameterArn.ValueString(),
		ExternalId:           data.ExternalId.ValueString(),
		Regions:              regions,
	})
	if err != nil {
		resp.Diagnostics.AddError("Cannot update organization", fmt.Sprintf("Cannot update organization: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsOrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AwsOrganizationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAwsOrganization(data.AccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot delete organization", fmt.Sprintf("Cannot delete organization: %s", err))
		return
	}
}

func (r *AwsOrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("account_id"), req, resp)
}
