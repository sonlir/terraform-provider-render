package provider

type EnvironmentVariable struct {
	Key   string `tfsdk:"key"`
	Value string `tfsdk:"value"`
}

type EnvironmentVariablesDataSourceModel struct {
	EnvVar []EnvironmentVariable `tfsdk:"environment_variables"`
}
