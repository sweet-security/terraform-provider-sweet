package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sweet-security/go-sweet/sweet"
)

var _ provider.Provider = &SweetProvider{}

type SweetProvider struct {
	version string
}

type SweetProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
	Secret types.String `tfsdk:"secret"`
	Env    types.String `tfsdk:"env"`
	Subenv types.String `tfsdk:"subenv"`
}

func (p *SweetProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sweet"
	resp.Version = p.version
}

func (p *SweetProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Sweet Api Key",
				Required:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Sweet Api Secret",
				Required:            true,
			},
			"env": schema.StringAttribute{
				MarkdownDescription: "Sweet environment to use",
				Optional:            true,
			},
			"subenv": schema.StringAttribute{
				MarkdownDescription: "Sweet sub environment to use",
				Optional:            true,
			},
		},
	}
}

func (p *SweetProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SweetProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	env := "prod"
	if !data.Env.IsNull() {
		env = data.Env.ValueString()
	}
	subenv := "main"
	if !data.Subenv.IsNull() {
		subenv = data.Subenv.ValueString()
	}

	client := sweet.New(
		data.ApiKey.ValueString(),
		data.Secret.ValueString(),
		sweet.WithEnv(env),
		sweet.WithSubenv(subenv),
	)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SweetProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAwsAccountResource,
	}
}

func (p *SweetProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SweetProvider{
			version: version,
		}
	}
}
