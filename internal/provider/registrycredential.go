// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sonlir/render-client-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RegistryCredential{}
	_ resource.ResourceWithConfigure   = &RegistryCredential{}
	_ resource.ResourceWithImportState = &RegistryCredential{}
)

func NewRegistryCredential() resource.Resource {
	return &RegistryCredential{}
}

type RegistryCredential struct {
	client *render.Client
}

// RegistryCredentialModel describes the data source data model.
type RegistryCredentialModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Registry  types.String `tfsdk:"registry"`
	Username  types.String `tfsdk:"username"`
	AuthToken types.String `tfsdk:"auth_token"`
	OwnerId   types.String `tfsdk:"owner_id"`
}

func (r *RegistryCredential) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registrycredential"
}

func (r *RegistryCredential) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Create registry credential",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this credential",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Descriptive name for this credential",
				Required:            true,
			},
			"registry": schema.StringAttribute{
				MarkdownDescription: "The registry to use this credential with. Valid values are GITHUB, GITLAB, DOCKER.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username associated with the credential",
				Required:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "The auth token associated with the credential",
				Required:            true,
				Sensitive:           true,
			},
			"owner_id": schema.StringAttribute{
				MarkdownDescription: "The owner id associated with the credential",
				Required:            true,
			},
		},
	}
}

func (r *RegistryCredential) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *RegistryCredential) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan RegistryCredentialModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Registry.ValueString() != "GITHUB" && plan.Registry.ValueString() != "GITLAB" && plan.Registry.ValueString() != "DOCKER" {
		resp.Diagnostics.AddError(
			"Invalid registry",
			fmt.Sprintf("The registry value must be one of GITHUB, GITLAB, or DOCKER, got: %s", plan.Registry.ValueString()),
		)
		return
	}

	// Generate API request body from plan
	data := render.RegistryCredentialData{
		Name:      plan.Name.ValueString(),
		Registry:  plan.Registry.ValueString(),
		Username:  plan.Username.ValueString(),
		AuthToken: plan.AuthToken.ValueString(),
		OwnerId:   plan.OwnerId.ValueString(),
	}

	registryCredential, err := r.client.CreateRegistryCredential(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Render registry credential",
			"Could not create registry credential, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(registryCredential.ID)
	plan.Name = types.StringValue(registryCredential.Name)
	plan.Registry = types.StringValue(registryCredential.Registry)
	plan.Username = types.StringValue(registryCredential.Username)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *RegistryCredential) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RegistryCredentialModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Render data into the model
	registryCredential, err := r.client.GetRegistryCredential(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get Render registry credential: "+state.ID.ValueString(),
			err.Error(),
		)
		return
	}

	// Write Render data into the model
	state.ID = types.StringValue(registryCredential.ID)
	state.Name = types.StringValue(registryCredential.Name)
	state.Registry = types.StringValue(registryCredential.Registry)
	state.Username = types.StringValue(registryCredential.Username)

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RegistryCredential) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state RegistryCredentialModel
	diags := req.Plan.Get(ctx, &plan)
	_ = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Registry.ValueString() != "GITHUB" && plan.Registry.ValueString() != "GITLAB" && plan.Registry.ValueString() != "DOCKER" {
		resp.Diagnostics.AddError(
			"Invalid registry",
			fmt.Sprintf("The registry value must be one of GITHUB, GITLAB, or DOCKER, got: %s", plan.Registry.ValueString()),
		)
		return
	}

	// Generate API request body from plan
	data := render.RegistryCredentialData{
		Name:      plan.Name.ValueString(),
		Registry:  plan.Registry.ValueString(),
		Username:  plan.Username.ValueString(),
		AuthToken: plan.AuthToken.ValueString(),
		OwnerId:   plan.OwnerId.ValueString(),
	}

	registryCredential, err := r.client.UpdateRegistryCredential(state.ID.ValueString(), data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Render registry credential",
			"Could not update registry credential ID: "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(registryCredential.ID)
	plan.Name = types.StringValue(registryCredential.Name)
	plan.Registry = types.StringValue(registryCredential.Registry)
	plan.Username = types.StringValue(registryCredential.Username)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RegistryCredential) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RegistryCredentialModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRegistryCredential(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Render registry credential",
			"Could not delete registry credential ID: "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *RegistryCredential) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
