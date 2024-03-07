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
	_ datasource.DataSource              = &OwnersDataSource{}
	_ datasource.DataSourceWithConfigure = &OwnersDataSource{}
)

func NewOwnersDataSource() datasource.DataSource {
	return &OwnersDataSource{}
}

type OwnersDataSource struct {
	client *render.Client
}

type OwnersDataSourceModel struct {
	Owners []OwnerDataSourceModel `tfsdk:"owners"`
}

func (d *OwnersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_owners"
}

func (d *OwnersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This endpoint lists all users and teams that your API key has access to. This can be helpful for getting the correct ownerId to use for creating new resources, such as services.",
		Attributes: map[string]schema.Attribute{
			"owners": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the user or team",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the user or team",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email of the user or team",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type. Valid values are `user` or `team`",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *OwnersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OwnersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OwnersDataSourceModel

	owners, err := d.client.GetOwners()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Render Owners",
			err.Error(),
		)
		return
	}

	for _, owner := range *owners {
		ownerState := OwnerDataSourceModel{
			ID:    types.StringValue(owner.Owner.ID),
			Name:  types.StringValue(owner.Owner.Name),
			Email: types.StringValue(owner.Owner.Email),
			Type:  types.StringValue(owner.Owner.Type),
		}
		state.Owners = append(state.Owners, ownerState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
