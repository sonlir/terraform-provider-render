package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sonlir/render-client-go"
)

var (
	_ datasource.DataSource              = &RegistryCredentialDataSource{}
	_ datasource.DataSourceWithConfigure = &RegistryCredentialDataSource{}
)

func NewRegistryCredentialDataSource() datasource.DataSource {
	return &RegistryCredentialDataSource{}
}

type RegistryCredentialDataSource struct {
	client *render.Client
}

type RegistryCredentialDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Registry types.String `tfsdk:"registry"`
	Username types.String `tfsdk:"username"`
}

func (d *RegistryCredentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registrycredential"
}

func (d *RegistryCredentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RegistryCredential data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this credential",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Descriptive name for this credential",
				Computed:            true,
			},
			"registry": schema.StringAttribute{
				MarkdownDescription: "The registry to use this credential with. Valid values are `GITHUB`, `GITLAB`, `DOCKER`.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username associated with the credential",
				Computed:            true,
			},
		},
	}
}

func (d *RegistryCredentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*render.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RegistryCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RegistryCredentialDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	registryCredential, err := d.client.GetRegistryCredential(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Render RegistryCredential: "+state.ID.ValueString(),
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(registryCredential.ID)
	state.Name = types.StringValue(registryCredential.Name)
	state.Registry = types.StringValue(registryCredential.Registry)
	state.Username = types.StringValue(registryCredential.Username)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
