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
	_ datasource.DataSource              = &WebServicesDataSource{}
	_ datasource.DataSourceWithConfigure = &WebServicesDataSource{}
)

func NewWebServicesDataSource() datasource.DataSource {
	return &WebServicesDataSource{}
}

type WebServicesDataSource struct {
	client *render.Client
}

type WebServicesDataSourceModel struct {
	Name        types.String             `tfsdk:"name"`
	WebServices []ServiceDataSourceModel `tfsdk:"web_services"`
}

func (d *WebServicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_services"
}

func (d *WebServicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Returns a list of Render web services owned by you or a team you belong to.",
		Attributes: map[string]schema.Attribute{
			"web_services": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the service",
							Computed:            true,
						},
						"owner_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the owner of the service",
							Computed:            true,
						},
						"repo": schema.StringAttribute{
							MarkdownDescription: "The git repository of the service",
							Computed:            true,
						},
						"auto_deploy": schema.StringAttribute{
							MarkdownDescription: "Whether the service is set to auto-deploy. Valid values are `yes` or `no`.",
							Computed:            true,
						},
						"branch": schema.StringAttribute{
							MarkdownDescription: "The branch of the service. If left empty, this will fall back to the default branch of the repository",
							Computed:            true,
						},
						"build_filter": schema.SingleNestedAttribute{
							MarkdownDescription: "The build filter for this service",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"paths": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
								"ignored_paths": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"root_dir": schema.StringAttribute{
							MarkdownDescription: "The root directory of the service",
							Computed:            true,
						},
						"environment_variables": schema.ListNestedAttribute{
							MarkdownDescription: "The environment variables for the service",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										MarkdownDescription: "The key of the environment variable",
										Computed:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "The value of the environment variable",
										Computed:            true,
									},
								},
							},
						},
						"service_details": schema.SingleNestedAttribute{
							MarkdownDescription: "The service details for the service",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"autoscaling": schema.SingleNestedAttribute{
									MarkdownDescription: "The autoscaling for the service",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											MarkdownDescription: "Whether autoscaling is enabled.",
											Computed:            true,
										},
										"min": schema.Int64Attribute{
											MarkdownDescription: "The minimum number of instances.",
											Computed:            true,
										},
										"max": schema.Int64Attribute{
											MarkdownDescription: "The maximum number of instances.",
											Computed:            true,
										},
										"criteria": schema.SingleNestedAttribute{
											MarkdownDescription: "The autoscaling criteria for the service",
											Computed:            true,
											Attributes: map[string]schema.Attribute{
												"cpu": schema.SingleNestedAttribute{
													MarkdownDescription: "The CPU autoscaling criteria for the service",
													Computed:            true,
													Attributes: map[string]schema.Attribute{
														"enabled": schema.BoolAttribute{
															MarkdownDescription: "Whether CPU autoscaling is enabled.",
															Computed:            true,
														},
														"percentage": schema.Int64Attribute{
															MarkdownDescription: "Determines when your service will be scaled. If the average resource utilization is significantly above/below the target, we will increase/decrease the number of instances.",
															Computed:            true,
														},
													},
												},
												"memory": schema.SingleNestedAttribute{
													MarkdownDescription: "The memory autoscaling criteria for the service",
													Computed:            true,
													Attributes: map[string]schema.Attribute{
														"enabled": schema.BoolAttribute{
															MarkdownDescription: "Whether memory autoscaling is enabled.",
															Computed:            true,
														},
														"percentage": schema.Int64Attribute{
															MarkdownDescription: "Determines when your service will be scaled. If the average resource utilization is significantly above/below the target, we will increase/decrease the number of instances.",
															Computed:            true,
														},
													},
												},
											},
										},
									},
								},
								"pull_request_previews_enabled": schema.StringAttribute{
									MarkdownDescription: "Whether pull request previews are enabled. Valid values are `yes` or `no`.",
									Computed:            true,
								},
								"disk": schema.SingleNestedAttribute{
									MarkdownDescription: "The disk for the service",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											MarkdownDescription: "The name of the disk",
											Computed:            true,
										},
										"id": schema.StringAttribute{
											MarkdownDescription: "The ID of the disk",
											Computed:            true,
										},
										"mount_path": schema.StringAttribute{
											MarkdownDescription: "The mount path of the disk",
											Computed:            true,
										},
										"size_gb": schema.Int64Attribute{
											MarkdownDescription: "The size of the disk in GB",
											Computed:            true,
										},
									},
								},
								"env": schema.StringAttribute{
									MarkdownDescription: "Environment (runtime)",
									Computed:            true,
								},

								"docker_details": schema.SingleNestedAttribute{
									MarkdownDescription: "The docker details for the service",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"docker_command": schema.StringAttribute{
											MarkdownDescription: "The docker command for the service",
											Computed:            true,
										},
										"docker_context": schema.StringAttribute{
											MarkdownDescription: "The docker context for the service",
											Computed:            true,
										},
										"dockerfile_path": schema.StringAttribute{
											MarkdownDescription: "The dockerfile path for the service",
											Computed:            true,
										},
										"pre_deploy_command": schema.StringAttribute{
											MarkdownDescription: "The pre-deploy command for the service",
											Computed:            true,
										},
										"registry_credential_id": schema.StringAttribute{
											MarkdownDescription: "The registry credential ID for the service",
											Computed:            true,
										},
									},
								},
								"native_environment_details": schema.SingleNestedAttribute{
									MarkdownDescription: "The native environment details for the service",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"build_command": schema.StringAttribute{
											MarkdownDescription: "The build command for the service",
											Computed:            true,
										},
										"start_command": schema.StringAttribute{
											MarkdownDescription: "The start command for the service",
											Computed:            true,
										},
										"pre_deploy_command": schema.StringAttribute{
											MarkdownDescription: "The pre-deploy command for the service",
											Computed:            true,
										},
									},
								},
								"health_check_path": schema.StringAttribute{
									MarkdownDescription: "The health check path for the service",
									Computed:            true,
								},
								"num_instances": schema.Int64Attribute{
									MarkdownDescription: "The number of instances for the service. ",
									Computed:            true,
								},
								"plan": schema.StringAttribute{
									MarkdownDescription: "The plan for the service. Valid values are `starter`, `starter_plus`, `standard`, `standard_plus`, `pro`, `pro_plus`, `pro_max`, `pro_ultra`.",
									Computed:            true,
								},
								"region": schema.StringAttribute{
									MarkdownDescription: "The region for the service. Valid values are `oregon` `frankfurt` . Defaults to `oregon`.",
									Computed:            true,
								},
								"url": schema.StringAttribute{
									MarkdownDescription: "The URL for the service",
									Computed:            true,
								},
								"open_ports": schema.ListNestedAttribute{
									MarkdownDescription: "The open ports for the service",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												MarkdownDescription: "The number of the open port",
												Computed:            true,
											},
											"protocol": schema.StringAttribute{
												MarkdownDescription: "The protocol of the open port",
												Computed:            true,
											},
										},
									},
								},
								"parent_server": schema.SingleNestedAttribute{
									MarkdownDescription: "The parent server for the service",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											MarkdownDescription: "The ID of the parent server",
											Computed:            true,
										},
										"name": schema.StringAttribute{
											MarkdownDescription: "The name of the parent server",
											Computed:            true,
										},
									},
								},
							},
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "The date and time the service was created",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "The date and time the service was last updated",
							Computed:            true,
						},
						"notify_on_fail": schema.StringAttribute{
							MarkdownDescription: "Whether to notify on fail. Valid values are `default`, `notify` or `ignore`.",
							Computed:            true,
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "The slug of the service",
							Computed:            true,
						},
						"suspended": schema.StringAttribute{
							MarkdownDescription: "Whether the service is suspended. Valid values are `suspended` or `not_suspended`.",
							Computed:            true,
						},
						"suspenders": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the service.",
							Computed:            true,
						},
						"image_path": schema.StringAttribute{
							MarkdownDescription: "The image path for the service",
							Computed:            true,
						},
					},
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the web service to filter by.",
				Optional:            true,
			},
		},
	}

}

func (d *WebServicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WebServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state WebServicesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	services, err := d.client.GetServices(&render.GetServicesArgs{Name: state.Name.ValueString(), Type: "web_service"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Render Web Services",
			err.Error(),
		)
		return
	}

	for _, service := range services {
		webService := ServiceDataSourceModel{}
		webService.ID = types.StringValue(service.ID)
		makeWebServiceDataSourceModel(&webService, &service)
		state.WebServices = append(state.WebServices, webService)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
