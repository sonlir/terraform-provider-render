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
	_ datasource.DataSource              = &WebServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &WebServiceDataSource{}
)

func NewWebServiceDataSource() datasource.DataSource {
	return &WebServiceDataSource{}
}

type WebServiceDataSource struct {
	client *render.Client
}

type WebServiceDetailsDataSource struct {
	Autoscaling                *Autoscaling              `tfsdk:"autoscaling"`
	Disk                       *Disk                     `tfsdk:"disk"`
	Env                        types.String              `tfsdk:"env"`
	DockerDetails              *DockerDetails            `tfsdk:"docker_details"`
	NativeEnvironmentDetails   *NativeEnvironmentDetails `tfsdk:"native_environment_details"`
	HealthCheckPath            types.String              `tfsdk:"health_check_path"`
	NumInstances               types.Int64               `tfsdk:"num_instances"`
	OpenPorts                  []OpenPort                `tfsdk:"open_ports"`
	ParentServer               *ParentServer             `tfsdk:"parent_server"`
	Plan                       types.String              `tfsdk:"plan"`
	PullRequestPreviewsEnabled types.String              `tfsdk:"pull_request_previews_enabled"`
	Region                     types.String              `tfsdk:"region"`
	URL                        types.String              `tfsdk:"url"`
}

func (d *WebServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_service"
}

func (d *WebServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Returns the details of a single Render Web Service (specified by `id`) that's owned by you or a team you belong to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service",
				Required:            true,
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
								MarkdownDescription: "The ID of the registry credential for the service",
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
	}
}

func (d *WebServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WebServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServiceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, err := d.client.GetService(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Render Web Service: "+state.ID.ValueString(),
			err.Error(),
		)
		return
	}

	makeWebServiceDataSourceModel(&state, service)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func makeWebServiceDataSourceModel(state *ServiceDataSourceModel, service *render.Service) {
	var webServiceDetails WebServiceDetailsDataSource
	state.AutoDeploy = types.StringValue(service.AutoDeploy)
	state.Branch = types.StringValue(service.Branch)
	if service.BuildFilter != nil {
		for _, path := range service.BuildFilter.Paths {
			state.BuildFilter.Paths = append(state.BuildFilter.Paths, types.StringValue(path))
		}
		for _, ignoredPath := range service.BuildFilter.IgnoredPaths {
			state.BuildFilter.IgnoredPaths = append(state.BuildFilter.IgnoredPaths, types.StringValue(ignoredPath))
		}
	}
	state.CreateAt = types.StringValue(service.CreateAt)
	state.ImagePath = types.StringValue(service.ImagePath)
	state.Name = types.StringValue(service.Name)
	state.NotifyOnFail = types.StringValue(service.NotifyOnFail)
	state.OwnerID = types.StringValue(service.OwnerID)
	state.Repo = types.StringValue(service.Repo)
	state.RootDir = types.StringValue(service.RootDir)
	state.Slug = types.StringValue(service.Slug)
	state.Suspended = types.StringValue(service.Suspended)
	for _, suspender := range service.Suspenders {
		state.Suspenders = append(state.Suspenders, types.StringValue(suspender))
	}
	state.Type = types.StringValue(service.Type)
	state.UpdatedAt = types.StringValue(service.UpdatedAt)

	webServiceDetails.NumInstances = types.Int64Value(int64(service.ServiceDetails.NumInstances))
	webServiceDetails.Env = types.StringValue(service.ServiceDetails.Env)
	webServiceDetails.HealthCheckPath = types.StringValue(service.ServiceDetails.HealthCheckPath)
	webServiceDetails.Plan = types.StringValue(service.ServiceDetails.Plan)
	webServiceDetails.PullRequestPreviewsEnabled = types.StringValue(service.ServiceDetails.PullRequestPreviewsEnabled)
	webServiceDetails.Region = types.StringValue(service.ServiceDetails.Region)
	webServiceDetails.URL = types.StringValue(service.ServiceDetails.URL)
	for _, openPort := range service.ServiceDetails.OpenPorts {
		webServiceDetails.OpenPorts = append(webServiceDetails.OpenPorts, OpenPort{
			Port:     types.Int64Value(int64(openPort.Port)),
			Protocol: types.StringValue(openPort.Protocol),
		})
	}
	if service.ServiceDetails.ParentServer != nil {
		webServiceDetails.ParentServer = &ParentServer{
			ID:   types.StringValue(service.ServiceDetails.ParentServer.ID),
			Name: types.StringValue(service.ServiceDetails.ParentServer.Name),
		}
	}
	if service.ServiceDetails.EnvSpecificDetails != nil {
		webServiceDetails.DockerDetails = &DockerDetails{
			DockerCommand:    types.StringPointerValue(service.ServiceDetails.EnvSpecificDetails.DockerCommand),
			DockerContext:    types.StringPointerValue(service.ServiceDetails.EnvSpecificDetails.DockerContext),
			DockerfilePath:   types.StringPointerValue(service.ServiceDetails.EnvSpecificDetails.DockerfilePath),
			PreDeployCommand: types.StringPointerValue(service.ServiceDetails.EnvSpecificDetails.PreDeployCommand),
		}
		if service.ServiceDetails.EnvSpecificDetails.RegistryCredential != nil {
			webServiceDetails.DockerDetails.RegistryCredentialId = types.StringValue(service.ServiceDetails.EnvSpecificDetails.RegistryCredential.ID)
		}
		webServiceDetails.NativeEnvironmentDetails = &NativeEnvironmentDetails{
			PreDeployCommand: types.StringPointerValue(service.ServiceDetails.EnvSpecificDetails.PreDeployCommand),
			BuildCommand:     types.StringPointerValue(service.ServiceDetails.EnvSpecificDetails.BuildCommand),
			StartCommand:     types.StringPointerValue(service.ServiceDetails.EnvSpecificDetails.StartCommand),
		}
	}
	if service.ServiceDetails.Disk != nil {
		webServiceDetails.Disk = &Disk{
			ID:        types.StringValue(service.ServiceDetails.Disk.Id),
			Name:      types.StringValue(service.ServiceDetails.Disk.Name),
			MountPath: types.StringValue(service.ServiceDetails.Disk.MountPath),
			SizeGB:    types.Int64Value(int64(service.ServiceDetails.Disk.SizeGB)),
		}
	}
	if service.ServiceDetails.Autoscaling != nil {
		webServiceDetails.Autoscaling = &Autoscaling{
			Enabled: types.BoolValue(service.ServiceDetails.Autoscaling.Enabled),
			Min:     types.Int64Value(int64(service.ServiceDetails.Autoscaling.Min)),
			Max:     types.Int64Value(int64(service.ServiceDetails.Autoscaling.Max)),
			Criteria: AutoscalingCriteria{
				CPU: AutoscalingCriteriaObject{
					Enabled:    types.BoolValue(service.ServiceDetails.Autoscaling.Criteria.CPU.Enabled),
					Percentage: types.Int64Value(int64(service.ServiceDetails.Autoscaling.Criteria.CPU.Percentage)),
				},
				Memory: AutoscalingCriteriaObject{
					Enabled:    types.BoolValue(service.ServiceDetails.Autoscaling.Criteria.Memory.Enabled),
					Percentage: types.Int64Value(int64(service.ServiceDetails.Autoscaling.Criteria.Memory.Percentage)),
				},
			},
		}
	}
	state.EnvVars = []EnvironmentVariable{}
	for i := len(service.EnvVars) - 1; i >= 0; i-- {
		state.EnvVars = append(state.EnvVars, EnvironmentVariable{
			Key:   types.StringValue(service.EnvVars[i].Key),
			Value: types.StringValue(service.EnvVars[i].Value),
		})
	}

	state.ServiceDetails = webServiceDetails
}
