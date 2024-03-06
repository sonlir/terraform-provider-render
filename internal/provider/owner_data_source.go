// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sonlir/render-client-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &OwnerDataSource{}
	_ datasource.DataSourceWithConfigure = &OwnerDataSource{}
)

func NewOwnerDataSource() datasource.DataSource {
	return &OwnerDataSource{}
}

type OwnerDataSource struct {
	client *render.Client
}

// OwnerDataSourceModel describes the data source data model.
type OwnerDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
	Type  types.String `tfsdk:"type"`
}

func (d *OwnerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_owner"
}

func (d *OwnerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This endpoint gets information for a specific user or team that your API key has permission to access, based on ownerId.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user or team",
				Required:            true,
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
	}
}

func (d *OwnerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *OwnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OwnerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Render data into the model
	owner, err := d.client.GetOwner(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Render Owner: "+state.ID.ValueString(),
			err.Error(),
		)
		return
	}

	// Write Render data into the model
	state.ID = types.StringValue(owner.ID)
	state.Name = types.StringValue(owner.Name)
	state.Email = types.StringValue(owner.Email)
	state.Type = types.StringValue(owner.Type)

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
