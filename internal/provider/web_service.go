package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sonlir/render-client-go"
)

var (
	_ resource.Resource                = &WebService{}
	_ resource.ResourceWithConfigure   = &WebService{}
	_ resource.ResourceWithImportState = &WebService{}
)

func NewWebService() resource.Resource {
	return &WebService{}
}

type WebService struct {
	client *render.Client
}

type WebServiceModel struct {
	AutoDeploy     types.String          `tfsdk:"auto_deploy"`
	Branch         types.String          `tfsdk:"branch"`
	BuildFilter    *BuildFilter          `tfsdk:"build_filter"`
	EnvVars        []EnvironmentVariable `tfsdk:"environment_variables"`
	ID             types.String          `tfsdk:"id"`
	Image          *Image                `tfsdk:"image"`
	Name           types.String          `tfsdk:"name"`
	OwnerID        types.String          `tfsdk:"owner_id"`
	Repo           types.String          `tfsdk:"repo"`
	RootDir        types.String          `tfsdk:"root_dir"`
	SecretFiles    []SecretFiles         `tfsdk:"secret_files"`
	ServiceDetails *WebServiceDetails    `tfsdk:"service_details"`
	Type           types.String          `tfsdk:"type"`
	CreateAt       types.String          `tfsdk:"created_at"`
	ImagePath      types.String          `tfsdk:"image_path"`
	NotifyOnFail   types.String          `tfsdk:"notify_on_fail"`
	Slug           types.String          `tfsdk:"slug"`
	Suspended      types.String          `tfsdk:"suspended"`
	Suspenders     []types.String        `tfsdk:"suspenders"`
	UpdatedAt      types.String          `tfsdk:"updated_at"`
}

type WebServiceDetails struct {
	Autoscaling                *Autoscaling              `tfsdk:"autoscaling"`
	Disk                       *Disk                     `tfsdk:"disk"`
	Env                        types.String              `tfsdk:"env"`
	DockerDetails              *DockerDetails            `tfsdk:"docker_details"`
	NativeEnvironmentDetails   *NativeEnvironmentDetails `tfsdk:"native_environment_details"`
	HealthCheckPath            types.String              `tfsdk:"health_check_path"`
	NumInstances               types.Int64               `tfsdk:"num_instances"`
	Plan                       types.String              `tfsdk:"plan"`
	PullRequestPreviewsEnabled types.String              `tfsdk:"pull_request_previews_enabled"`
	Region                     types.String              `tfsdk:"region"`
	OpenPorts                  []OpenPort                `tfsdk:"open_ports"`
	ParentServer               *ParentServer             `tfsdk:"parent_server"`
	URL                        types.String              `tfsdk:"url"`
}

func (r *WebService) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_service"
}

func (r *WebService) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a new Render Web service owned by you or a team you belong to.\n~> **Note:** You can't create free-tier services with the Render API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service",
				Required:            true,
			},
			"owner_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the owner of the service",
				Required:            true,
			},
			"repo": schema.StringAttribute{
				MarkdownDescription: "The git repository of the service",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"auto_deploy": schema.StringAttribute{
				MarkdownDescription: "Whether the service is set to auto-deploy. Valid values are `yes` or `no`. Default: `yes`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "The branch of the service. If left empty, this will fall back to the default branch of the repository",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"image": schema.SingleNestedAttribute{
				MarkdownDescription: "The image used for this server",
				Optional:            true,
				Default:             nil,
				PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"owner_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner for this image. This should match the owner of the service as well as the owner of any specified registry credential.",
						Required:            true,
					},
					"registry_credential_id": schema.StringAttribute{
						MarkdownDescription: "Optional reference to the registry credential passed to the image repository to retrieve this image.",
						Optional:            true,
					},
					"image_path": schema.StringAttribute{
						MarkdownDescription: "Path to the image used for this server e.g `docker.io/library/nginx:latest`.",
						Required:            true,
					},
				},
			},
			"build_filter": schema.SingleNestedAttribute{
				MarkdownDescription: "The build filter for this service",
				Optional:            true,
				Default:             nil,
				Attributes: map[string]schema.Attribute{
					"paths": schema.ListAttribute{
						ElementType:   types.StringType,
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
					},
					"ignored_paths": schema.ListAttribute{
						ElementType:   types.StringType,
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
					},
				},
			},
			"root_dir": schema.StringAttribute{
				MarkdownDescription: "The root directory of the service",
				Computed:            true,
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_details": schema.SingleNestedAttribute{
				MarkdownDescription: "The service details for the service",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"autoscaling": schema.SingleNestedAttribute{
						MarkdownDescription: "The autoscaling for the service",
						Optional:            true,
						Default:             nil,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether autoscaling is enabled.",
								Optional:            true,
							},
							"min": schema.Int64Attribute{
								MarkdownDescription: "The minimum number of instances.",
								Optional:            true,
							},
							"max": schema.Int64Attribute{
								MarkdownDescription: "The maximum number of instances.",
								Optional:            true,
							},
							"criteria": schema.SingleNestedAttribute{
								MarkdownDescription: "The autoscaling criteria for the service",
								Required:            true,
								Attributes: map[string]schema.Attribute{
									"cpu": schema.SingleNestedAttribute{
										MarkdownDescription: "The CPU autoscaling criteria for the service",
										Required:            true,
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												MarkdownDescription: "Whether CPU autoscaling is enabled.",
												Optional:            true,
											},
											"percentage": schema.Int64Attribute{
												MarkdownDescription: "Determines when your service will be scaled. If the average resource utilization is significantly above/below the target, we will increase/decrease the number of instances.",
												Optional:            true,
											},
										},
									},
									"memory": schema.SingleNestedAttribute{
										MarkdownDescription: "The memory autoscaling criteria for the service",
										Required:            true,
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												MarkdownDescription: "Whether memory autoscaling is enabled.",
												Optional:            true,
											},
											"percentage": schema.Int64Attribute{
												MarkdownDescription: "Determines when your service will be scaled. If the average resource utilization is significantly above/below the target, we will increase/decrease the number of instances.",
												Optional:            true,
											},
										},
									},
								},
							},
						},
					},
					"pull_request_previews_enabled": schema.StringAttribute{
						MarkdownDescription: "Whether pull request previews are enabled. Valid values are `yes` or `no`. Default: `no`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"disk": schema.SingleNestedAttribute{
						MarkdownDescription: "The disk for the service",
						Optional:            true,
						Default:             nil,
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								MarkdownDescription: "The name of the disk",
								Optional:            true,
								Computed:            true,
								PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
							},
							"size_gb": schema.Int64Attribute{
								MarkdownDescription: "The size of the disk in GB. Default: `1`.",
								Optional:            true,
								PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
							},
							"mount_path": schema.StringAttribute{
								MarkdownDescription: "The mount path of the disk.",
								Optional:            true,
								PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
							},
							"id": schema.StringAttribute{
								MarkdownDescription: "The ID of the disk",
								Computed:            true,
								PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
							},
						},
					},
					"env": schema.StringAttribute{
						MarkdownDescription: "Environment (runtime). Valid values are `node`, `python`, `ruby`, `go`, `elixir`, `image`, `rust`, `docker`.",
						Required:            true,
					},
					"native_environment_details": schema.SingleNestedAttribute{
						MarkdownDescription: "The environment specific details for the service",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Attributes: map[string]schema.Attribute{
							"pre_deploy_command": schema.StringAttribute{
								MarkdownDescription: "The ID of the registry credential for the service",
								Optional:            true,
							},
							"build_command": schema.StringAttribute{
								MarkdownDescription: "The build command for the service",
								Required:            true,
							},
							"start_command": schema.StringAttribute{
								MarkdownDescription: "The start command for the service",
								Required:            true,
							},
						},
					},
					"docker_details": schema.SingleNestedAttribute{
						MarkdownDescription: "The environment specific details for the service",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
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
								MarkdownDescription: "The dockerfile path for the service.",
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
					"health_check_path": schema.StringAttribute{
						MarkdownDescription: "The health check path for the service",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"num_instances": schema.Int64Attribute{
						MarkdownDescription: "The number of instances for the service. Default: `1`.",
						Required:            true,
					},
					"plan": schema.StringAttribute{
						MarkdownDescription: "The plan for the service. Valid values are `starter`, `starter_plus`, `standard`, `standard_plus`, `pro`, `pro_plus`, `pro_max`, `pro_ultra`. Default: `starter`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"region": schema.StringAttribute{
						MarkdownDescription: "The region for the service. Valid values are `oregon` `frankfurt` . Defaults to `oregon`.",
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"open_ports": schema.ListNestedAttribute{
						MarkdownDescription: "The open ports for the service",
						Computed:            true,
						PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
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
						Optional:            true,
						Default:             nil,
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
					"url": schema.StringAttribute{
						MarkdownDescription: "The URL for the service",
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
				},
			},
			"secret_files": schema.ListNestedAttribute{
				MarkdownDescription: "The secret files for the service",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the secret file",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"content": schema.StringAttribute{
							MarkdownDescription: "The content of the secret file",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
					},
				},
			},
			"environment_variables": schema.ListNestedAttribute{
				MarkdownDescription: "The environment variables for the service",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "The key of the environment variable",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value of the environment variable",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
					},
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the service. Valid values are `web_service`, `static_site`, `cron_job`, `background_worker`, `private_service`.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the service was created",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the service was last updated",
				Computed:            true,
			},
			"image_path": schema.StringAttribute{MarkdownDescription: "The image path for the service",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"notify_on_fail": schema.StringAttribute{
				MarkdownDescription: "Whether to notify on fail. Valid values are `default`, `notify` or `ignore`.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the service",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"suspended": schema.StringAttribute{
				MarkdownDescription: "Whether the service is suspended. Valid values are `suspended` or `not_suspended`.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"suspenders": schema.ListAttribute{
				MarkdownDescription: "The suspenders of the service",
				ElementType:         types.StringType,
				Computed:            true,
				PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *WebService) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*render.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *WebService) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebServiceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := makeWebServiceData(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Render web service", err.Error())
		return
	}

	service, err := r.client.CreateService(*data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Render web service",
			"Could not create web service, unexpected error: "+err.Error(),
		)
		return
	}

	makeWebServiceModel(&plan, service)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WebService) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebServiceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, err := r.client.GetService(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get Render web service: "+state.ID.ValueString(),
			err.Error(),
		)
		return
	}

	makeWebServiceModel(&state, service)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WebService) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state WebServiceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID

	data, err := makeWebServiceData(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Render web service", err.Error())
		return
	}

	service, err := r.client.UpdateService(plan.ID.ValueString(), *data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Render web service",
			"Could not update web service ID: "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	makeWebServiceModel(&plan, service)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WebService) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteService(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Render web service",
			"Could not delete web service ID: "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *WebService) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func makeWebServiceModel(state *WebServiceModel, service *render.Service) {
	var webServiceDetails WebServiceDetails
	state.AutoDeploy = types.StringValue(service.AutoDeploy)
	state.Branch = types.StringValue(service.Branch)
	if service.BuildFilter != nil {
		state.BuildFilter.Paths = []types.String{}
		for _, path := range service.BuildFilter.Paths {
			state.BuildFilter.Paths = append(state.BuildFilter.Paths, types.StringValue(path))
		}
		state.BuildFilter.IgnoredPaths = []types.String{}
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
	state.RootDir = types.StringValue(service.RootDir)
	state.Slug = types.StringValue(service.Slug)
	state.Suspended = types.StringValue(service.Suspended)
	state.Suspenders = []types.String{}
	for _, suspender := range service.Suspenders {
		state.Suspenders = append(state.Suspenders, types.StringValue(suspender))
	}

	state.Type = types.StringValue(service.Type)
	state.UpdatedAt = types.StringValue(service.UpdatedAt)

	webServiceDetails.NumInstances = types.Int64Value(service.ServiceDetails.NumInstances)
	webServiceDetails.Env = types.StringValue(service.ServiceDetails.Env)
	webServiceDetails.HealthCheckPath = types.StringValue(service.ServiceDetails.HealthCheckPath)
	webServiceDetails.Plan = types.StringValue(service.ServiceDetails.Plan)
	webServiceDetails.PullRequestPreviewsEnabled = types.StringValue(service.ServiceDetails.PullRequestPreviewsEnabled)
	webServiceDetails.Region = types.StringValue(service.ServiceDetails.Region)
	webServiceDetails.URL = types.StringValue(service.ServiceDetails.URL)
	webServiceDetails.OpenPorts = []OpenPort{}
	for _, openPort := range service.ServiceDetails.OpenPorts {
		webServiceDetails.OpenPorts = append(webServiceDetails.OpenPorts, OpenPort{
			Port:     types.Int64Value(openPort.Port),
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
			ID:   types.StringValue(service.ServiceDetails.Disk.Id),
			Name: types.StringValue(service.ServiceDetails.Disk.Name),
		}
	}
	if service.ServiceDetails.Autoscaling != nil {
		webServiceDetails.Autoscaling = &Autoscaling{
			Enabled: types.BoolValue(service.ServiceDetails.Autoscaling.Enabled),
			Min:     types.Int64Value(service.ServiceDetails.Autoscaling.Min),
			Max:     types.Int64Value(service.ServiceDetails.Autoscaling.Max),
			Criteria: AutoscalingCriteria{
				CPU: AutoscalingCriteriaObject{
					Enabled:    types.BoolValue(service.ServiceDetails.Autoscaling.Criteria.CPU.Enabled),
					Percentage: types.Int64Value(service.ServiceDetails.Autoscaling.Criteria.CPU.Percentage),
				},
				Memory: AutoscalingCriteriaObject{
					Enabled:    types.BoolValue(service.ServiceDetails.Autoscaling.Criteria.Memory.Enabled),
					Percentage: types.Int64Value(service.ServiceDetails.Autoscaling.Criteria.Memory.Percentage),
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
	state.SecretFiles = []SecretFiles{}
	for i := len(service.SecretFiles) - 1; i >= 0; i-- {
		state.SecretFiles = append(state.SecretFiles, SecretFiles{
			Name:     types.StringValue(service.SecretFiles[i].Name),
			Contents: types.StringValue(service.SecretFiles[i].Contents),
		})
	}

	state.ServiceDetails = &webServiceDetails
}

func makeWebServiceData(plan *WebServiceModel) (*render.Service, error) {
	webService := render.Service{}
	webServiceDetails := plan.ServiceDetails

	if webServiceDetails.Env.ValueString() != "docker" &&
		webServiceDetails.Env.ValueString() != "image" &&
		webServiceDetails.Env.ValueString() != "rust" &&
		webServiceDetails.Env.ValueString() != "go" &&
		webServiceDetails.Env.ValueString() != "elixir" &&
		webServiceDetails.Env.ValueString() != "node" &&
		webServiceDetails.Env.ValueString() != "python" &&
		webServiceDetails.Env.ValueString() != "ruby" {
		return nil, fmt.Errorf("invalid environment: %s", fmt.Sprintf("the environment value must be one of docker, image, rust, go, elixir, node, python, or ruby, got: %s", webServiceDetails.Env.ValueString()))
	}

	webServiceDetailsData := render.ServiceDetails{
		PullRequestPreviewsEnabled: webServiceDetails.PullRequestPreviewsEnabled.ValueString(),
		HealthCheckPath:            webServiceDetails.HealthCheckPath.ValueString(),
		NumInstances:               webServiceDetails.NumInstances.ValueInt64(),
		Plan:                       webServiceDetails.Plan.ValueString(),
		Region:                     webServiceDetails.Region.ValueString(),
		Env:                        webServiceDetails.Env.ValueString(),
	}

	if webServiceDetails.DockerDetails != nil {
		if webServiceDetails.NativeEnvironmentDetails != nil {
			webServiceDetailsData.EnvSpecificDetails = &render.EnvSpecificDetails{
				PreDeployCommand: webServiceDetails.NativeEnvironmentDetails.PreDeployCommand.ValueStringPointer(),
				BuildCommand:     webServiceDetails.NativeEnvironmentDetails.BuildCommand.ValueStringPointer(),
				StartCommand:     webServiceDetails.NativeEnvironmentDetails.StartCommand.ValueStringPointer(),
			}
		}
	}

	if webServiceDetails.Autoscaling != nil {
		cpuAutoscalingCriteriaObject := render.AutoscalingCriteriaObject{
			Enabled:    webServiceDetails.Autoscaling.Criteria.CPU.Enabled.ValueBool(),
			Percentage: webServiceDetails.Autoscaling.Criteria.CPU.Percentage.ValueInt64(),
		}

		memoryAutoscalingCriteriaObject := render.AutoscalingCriteriaObject{
			Enabled:    webServiceDetails.Autoscaling.Criteria.Memory.Enabled.ValueBool(),
			Percentage: webServiceDetails.Autoscaling.Criteria.Memory.Percentage.ValueInt64(),
		}

		autoscalingCriteria := render.AutoscalingCriteria{
			CPU:    &cpuAutoscalingCriteriaObject,
			Memory: &memoryAutoscalingCriteriaObject,
		}

		autoscaling := render.Autoscaling{
			Enabled:  webServiceDetails.Autoscaling.Enabled.ValueBool(),
			Min:      webServiceDetails.Autoscaling.Min.ValueInt64(),
			Max:      webServiceDetails.Autoscaling.Max.ValueInt64(),
			Criteria: &autoscalingCriteria,
		}

		webServiceDetailsData.Autoscaling = &autoscaling
	}

	secretFiles := []render.SecretFiles{}
	for _, secretFile := range plan.SecretFiles {
		secretFiles = append(secretFiles, render.SecretFiles{
			Name:     secretFile.Name.ValueString(),
			Contents: secretFile.Contents.ValueString(),
		})
	}

	envVars := []render.EnvironmentVariable{}
	for _, envVar := range plan.EnvVars {
		envVars = append(envVars, render.EnvironmentVariable{
			Key:   envVar.Key.ValueString(),
			Value: envVar.Value.ValueString(),
		})
	}

	buildFilter := render.BuildFilter{}
	paths := []string{}
	ignoredPaths := []string{}
	if plan.BuildFilter != nil {
		for _, path := range plan.BuildFilter.Paths {
			paths = append(paths, path.ValueString())
		}
		for _, ignoredPath := range plan.BuildFilter.IgnoredPaths {
			ignoredPaths = append(ignoredPaths, ignoredPath.ValueString())
		}
		buildFilter = render.BuildFilter{
			Paths:        paths,
			IgnoredPaths: ignoredPaths,
		}
	}

	if plan.Image != nil {
		webService.Image = &render.Image{
			OwnerId:              plan.Image.OwnerID.ValueString(),
			RegistryCredentialId: plan.Image.RegistryCredentialId.ValueString(),
			ImagePath:            plan.Image.ImagePath.ValueString(),
		}
	}
	webService.Name = plan.Name.ValueString()
	webService.OwnerID = plan.OwnerID.ValueString()
	webService.Repo = plan.Repo.ValueString()
	webService.AutoDeploy = plan.AutoDeploy.ValueString()
	webService.Branch = plan.Branch.ValueString()
	webService.RootDir = plan.RootDir.ValueString()
	webService.ServiceDetails = webServiceDetailsData
	webService.SecretFiles = secretFiles
	webService.EnvVars = envVars
	webService.BuildFilter = &buildFilter
	webService.Type = "web_service"

	return &webService, nil
}
