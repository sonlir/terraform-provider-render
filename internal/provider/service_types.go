package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServiceDataSourceModel struct {
	AutoDeploy  types.String `tfsdk:"auto_deploy"`
	Branch      types.String `tfsdk:"branch"`
	BuildFilter *BuildFilter `tfsdk:"build_filter"`
	CreateAt    types.String `tfsdk:"created_at"`
	ID          types.String `tfsdk:"id"`
	//Image        Image        `tfsdk:"image"` Commented because Render REST API does not return these fields
	ImagePath      types.String   `tfsdk:"image_path"`
	Name           types.String   `tfsdk:"name"`
	NotifyOnFail   types.String   `tfsdk:"notify_on_fail"`
	OwnerId        types.String   `tfsdk:"owner_id"`
	Repo           types.String   `tfsdk:"repo"`
	RootDir        types.String   `tfsdk:"root_dir"`
	ServiceDetails interface{}    `tfsdk:"service_details"`
	Slug           types.String   `tfsdk:"slug"`
	Suspended      types.String   `tfsdk:"suspended"`
	Suspenders     []types.String `tfsdk:"suspenders"`
	Type           types.String   `tfsdk:"type"`
	UpdatedAt      types.String   `tfsdk:"updated_at"`
}

type BuildFilter struct {
	Paths        []types.String `tfsdk:"paths"`
	IgnoredPaths []types.String `tfsdk:"ignored_paths"`
}

type Image struct {
	OwnerId              types.String `tfsdk:"owner_id"`
	RegistryCredentialId types.String `tfsdk:"registry_credential_id"`
	ImagePath            types.String `tfsdk:"image_path"`
}

type Route struct {
	Type        types.String `tfsdk:"type"`
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`
}

type SecretFiles struct {
	Name     types.String `tfsdk:"name"`
	Contents types.String `tfsdk:"contents"`
}

type DiskDataSource struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type Autoscaling struct {
	Enabled  types.Bool          `tfsdk:"enabled"`
	Min      types.Int64         `tfsdk:"min"`
	Max      types.Int64         `tfsdk:"max"`
	Criteria AutoscalingCriteria `tfsdk:"criteria"`
}

type AutoscalingCriteria struct {
	CPU    AutoscalingCriteriaObject `tfsdk:"cpu"`
	Memory AutoscalingCriteriaObject `tfsdk:"memory"`
}

type AutoscalingCriteriaObject struct {
	Enabled    types.Bool  `tfsdk:"enabled"`
	Percentage types.Int64 `tfsdk:"percentage"`
}

type OpenPort struct {
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

type ParentServer struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type EnvSpecificDetailsDataSource struct {
	DockerCommand      types.String                      `tfsdk:"docker_command"`
	DockerContext      types.String                      `tfsdk:"docker_context"`
	DockerfilePath     types.String                      `tfsdk:"dockerfile_path"`
	PreDeployCommand   types.String                      `tfsdk:"pre_deploy_command"`
	RegistryCredential RegistryCredentialDataSourceModel `tfsdk:"registry_credential"`
	BuildCommand       types.String                      `tfsdk:"build_command"`
	StartCommand       types.String                      `tfsdk:"start_command"`
}
