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
	_ datasource.DataSource              = &RegistryCredentialsDataSource{}
	_ datasource.DataSourceWithConfigure = &RegistryCredentialsDataSource{}
)

func NewRegistryCredentialsDataSource() datasource.DataSource {
	return &RegistryCredentialsDataSource{}
}

type RegistryCredentialsDataSource struct {
	client *render.Client
}

type RegistryCredentialsDataSourceModel struct {
	RegistryCredentials []RegistryCredentialDataSourceModel `tfsdk:"registrycredentials"`
}

func (d *RegistryCredentialsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registrycredentials"
}

func (d *RegistryCredentialsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a list of registry credentials.",
		Attributes: map[string]schema.Attribute{
			"registrycredentials": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
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
				},
			},
		},
	}
}

func (d *RegistryCredentialsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RegistryCredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RegistryCredentialsDataSourceModel

	registryCredentials, err := d.client.GetRegistryCredentials()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Render RegistryCredentials",
			err.Error(),
		)
		return
	}

	for _, registryCredential := range *registryCredentials {
		registryCredentialState := RegistryCredentialDataSourceModel{
			ID:       types.StringValue(registryCredential.ID),
			Name:     types.StringValue(registryCredential.Name),
			Registry: types.StringValue(registryCredential.Registry),
			Username: types.StringValue(registryCredential.Username),
		}
		state.RegistryCredentials = append(state.RegistryCredentials, registryCredentialState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
