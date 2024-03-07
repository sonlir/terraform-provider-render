package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sonlir/render-client-go"
)

var _ provider.Provider = &RenderProvider{}
var _ provider.ProviderWithFunctions = &RenderProvider{}

type RenderProvider struct {
	version string
}

type RenderProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

func (p *RenderProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "render"
	resp.Version = p.version
}

func (p *RenderProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The Render API key to use for authentication. May also be provided via `RENDER_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *RenderProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config RenderProviderModel

	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Render API KEY",
			"The provider cannot create the Render API client as there is an unknown configuration value for the Render API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the RENDER_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("RENDER_API_KEY")

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Render API KEY",
			"The provider cannot create the Render API client as there is a missing or empty value for the Render API Key. "+
				"Set the api_key value in the configuration or use the RENDER_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := render.NewClient(&apiKey, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Render API Client",
			"An unexpected error occurred when creating the Render API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Render Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RenderProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRegistryCredential,
	}
}

func (p *RenderProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOwnerDataSource,
		NewOwnersDataSource,
		NewRegistryCredentialDataSource,
		NewRegistryCredentialsDataSource,
		NewWebServiceDataSource,
		NewWebServicesDataSource,
	}
}

func (p *RenderProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RenderProvider{
			version: version,
		}
	}
}
