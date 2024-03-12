package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServiceDataSourceModel struct {
	AutoDeploy     types.String          `tfsdk:"auto_deploy"`
	Branch         types.String          `tfsdk:"branch"`
	BuildFilter    *BuildFilter          `tfsdk:"build_filter"`
	CreateAt       types.String          `tfsdk:"created_at"`
	EnvVars        []EnvironmentVariable `tfsdk:"environment_variables"`
	ID             types.String          `tfsdk:"id"`
	ImagePath      types.String          `tfsdk:"image_path"`
	Name           types.String          `tfsdk:"name"`
	NotifyOnFail   types.String          `tfsdk:"notify_on_fail"`
	OwnerID        types.String          `tfsdk:"owner_id"`
	Repo           types.String          `tfsdk:"repo"`
	RootDir        types.String          `tfsdk:"root_dir"`
	ServiceDetails interface{}           `tfsdk:"service_details"`
	Slug           types.String          `tfsdk:"slug"`
	Suspended      types.String          `tfsdk:"suspended"`
	Suspenders     []types.String        `tfsdk:"suspenders"`
	Type           types.String          `tfsdk:"type"`
	UpdatedAt      types.String          `tfsdk:"updated_at"`
}

type BuildFilter struct {
	Paths        []types.String `tfsdk:"paths"`
	IgnoredPaths []types.String `tfsdk:"ignored_paths"`
}

type Image struct {
	OwnerID              types.String `tfsdk:"owner_id"`
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

type DockerDetails struct {
	DockerCommand        types.String `tfsdk:"docker_command"`
	DockerContext        types.String `tfsdk:"docker_context"`
	DockerfilePath       types.String `tfsdk:"dockerfile_path"`
	PreDeployCommand     types.String `tfsdk:"pre_deploy_command"`
	RegistryCredentialId types.String `tfsdk:"registry_credential_id"`
}

type NativeEnvironmentDetails struct {
	BuildCommand     types.String `tfsdk:"build_command"`
	StartCommand     types.String `tfsdk:"start_command"`
	PreDeployCommand types.String `tfsdk:"pre_deploy_command"`
}

type Header struct {
	Path  types.String `tfsdk:"path"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type Disk struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SizeGB    types.Int64  `tfsdk:"size_gb"`
	ID        types.String `tfsdk:"id"`
}

type EnvironmentVariable struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}
