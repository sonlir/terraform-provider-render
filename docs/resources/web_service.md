---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "render_web_service Resource - render"
subcategory: ""
description: |-
  Creates a new Render Web service owned by you or a team you belong to.
  ~> Note: You can't create free-tier services with the Render API.
---

# render_web_service (Resource)

Creates a new Render Web service owned by you or a team you belong to.
~> **Note:** You can't create free-tier services with the Render API.

## Example Usage

```terraform
# Minimal example of a Render web service
data "render_owner" "example" {
  id = "usr-abcdefghijklmnopqest"
}

resource "render_web_service" "example" {
  owner_id = data.render_owner.example.id
  name     = "render-web-service"
  service_details = {
    env           = "image"
    num_instances = 1
  }
}

# Full example of a Render web service
resource "render_web_service" "example" {
  owner_id = data.render_owner.example.id
  name     = "render-web-service"
  service_details = {
    env           = "image"
    num_instances = 1
    region        = "frankfurt"
    autoscaling = {
      enabled = true
      min     = 1
      max     = 3
      criteria = {
        cpu = {
          enabled    = true
          percentage = 50
        }
        memory = {
          enabled    = true
          percentage = 50
        }
      }
    }
    health_check_path = "/health"
    plan              = "starter"
    image = {
      owner_id   = data.render_owner.example.id
      image_path = "docker.io/library/nginx:latest"
    }
    auto_deploy = "yes"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the service
- `owner_id` (String) The ID of the owner of the service
- `service_details` (Attributes) The service details for the service (see [below for nested schema](#nestedatt--service_details))

### Optional

- `auto_deploy` (String) Whether the service is set to auto-deploy. Valid values are `yes` or `no`. Default: `yes`.
- `branch` (String) The branch of the service. If left empty, this will fall back to the default branch of the repository
- `build_filter` (Attributes) The build filter for this service (see [below for nested schema](#nestedatt--build_filter))
- `environment_variables` (Attributes List) The environment variables for the service (see [below for nested schema](#nestedatt--environment_variables))
- `image` (Attributes) The image used for this server (see [below for nested schema](#nestedatt--image))
- `repo` (String) The git repository of the service
- `root_dir` (String) The root directory of the service
- `secret_files` (Attributes List) The secret files for the service (see [below for nested schema](#nestedatt--secret_files))

### Read-Only

- `created_at` (String) The date and time the service was created
- `id` (String) The ID of the service
- `image_path` (String) The image path for the service
- `notify_on_fail` (String) Whether to notify on fail. Valid values are `default`, `notify` or `ignore`.
- `slug` (String) The slug of the service
- `suspended` (String) Whether the service is suspended. Valid values are `suspended` or `not_suspended`.
- `suspenders` (List of String) The suspenders of the service
- `type` (String) The type of the service. Valid values are `web_service`, `static_site`, `cron_job`, `background_worker`, `private_service`.
- `updated_at` (String) The date and time the service was last updated

<a id="nestedatt--service_details"></a>
### Nested Schema for `service_details`

Required:

- `env` (String) Environment (runtime). Valid values are `node`, `python`, `ruby`, `go`, `elixir`, `image`, `rust`, `docker`.
- `num_instances` (Number) The number of instances for the service. Default: `1`.

Optional:

- `autoscaling` (Attributes) The autoscaling for the service (see [below for nested schema](#nestedatt--service_details--autoscaling))
- `disk` (Attributes) The disk for the service (see [below for nested schema](#nestedatt--service_details--disk))
- `docker_details` (Attributes) The environment specific details for the service (see [below for nested schema](#nestedatt--service_details--docker_details))
- `health_check_path` (String) The health check path for the service
- `native_environment_details` (Attributes) The environment specific details for the service (see [below for nested schema](#nestedatt--service_details--native_environment_details))
- `parent_server` (Attributes) The parent server for the service (see [below for nested schema](#nestedatt--service_details--parent_server))
- `plan` (String) The plan for the service. Valid values are `starter`, `starter_plus`, `standard`, `standard_plus`, `pro`, `pro_plus`, `pro_max`, `pro_ultra`. Default: `starter`.
- `pull_request_previews_enabled` (String) Whether pull request previews are enabled. Valid values are `yes` or `no`. Default: `no`.
- `region` (String) The region for the service. Valid values are `oregon` `frankfurt` . Defaults to `oregon`.

Read-Only:

- `open_ports` (Attributes List) The open ports for the service (see [below for nested schema](#nestedatt--service_details--open_ports))
- `url` (String) The URL for the service

<a id="nestedatt--service_details--autoscaling"></a>
### Nested Schema for `service_details.autoscaling`

Required:

- `criteria` (Attributes) The autoscaling criteria for the service (see [below for nested schema](#nestedatt--service_details--autoscaling--criteria))

Optional:

- `enabled` (Boolean) Whether autoscaling is enabled.
- `max` (Number) The maximum number of instances.
- `min` (Number) The minimum number of instances.

<a id="nestedatt--service_details--autoscaling--criteria"></a>
### Nested Schema for `service_details.autoscaling.criteria`

Required:

- `cpu` (Attributes) The CPU autoscaling criteria for the service (see [below for nested schema](#nestedatt--service_details--autoscaling--criteria--cpu))
- `memory` (Attributes) The memory autoscaling criteria for the service (see [below for nested schema](#nestedatt--service_details--autoscaling--criteria--memory))

<a id="nestedatt--service_details--autoscaling--criteria--cpu"></a>
### Nested Schema for `service_details.autoscaling.criteria.memory`

Optional:

- `enabled` (Boolean) Whether CPU autoscaling is enabled.
- `percentage` (Number) Determines when your service will be scaled. If the average resource utilization is significantly above/below the target, we will increase/decrease the number of instances.


<a id="nestedatt--service_details--autoscaling--criteria--memory"></a>
### Nested Schema for `service_details.autoscaling.criteria.memory`

Optional:

- `enabled` (Boolean) Whether memory autoscaling is enabled.
- `percentage` (Number) Determines when your service will be scaled. If the average resource utilization is significantly above/below the target, we will increase/decrease the number of instances.




<a id="nestedatt--service_details--disk"></a>
### Nested Schema for `service_details.disk`

Optional:

- `mount_path` (String) The mount path of the disk.
- `name` (String) The name of the disk
- `size_gb` (Number) The size of the disk in GB. Default: `1`.

Read-Only:

- `id` (String) The ID of the disk


<a id="nestedatt--service_details--docker_details"></a>
### Nested Schema for `service_details.docker_details`

Read-Only:

- `docker_command` (String) The docker command for the service
- `docker_context` (String) The docker context for the service
- `dockerfile_path` (String) The dockerfile path for the service.
- `pre_deploy_command` (String) The pre-deploy command for the service
- `registry_credential_id` (String) The ID of the registry credential for the service


<a id="nestedatt--service_details--native_environment_details"></a>
### Nested Schema for `service_details.native_environment_details`

Required:

- `build_command` (String) The build command for the service
- `start_command` (String) The start command for the service

Optional:

- `pre_deploy_command` (String) The ID of the registry credential for the service


<a id="nestedatt--service_details--parent_server"></a>
### Nested Schema for `service_details.parent_server`

Read-Only:

- `id` (String) The ID of the parent server
- `name` (String) The name of the parent server


<a id="nestedatt--service_details--open_ports"></a>
### Nested Schema for `service_details.open_ports`

Read-Only:

- `port` (Number) The number of the open port
- `protocol` (String) The protocol of the open port



<a id="nestedatt--build_filter"></a>
### Nested Schema for `build_filter`

Optional:

- `ignored_paths` (List of String)
- `paths` (List of String)


<a id="nestedatt--environment_variables"></a>
### Nested Schema for `environment_variables`

Optional:

- `key` (String) The key of the environment variable
- `value` (String) The value of the environment variable


<a id="nestedatt--image"></a>
### Nested Schema for `image`

Required:

- `image_path` (String) Path to the image used for this server e.g `docker.io/library/nginx:latest`.
- `owner_id` (String) The ID of the owner for this image. This should match the owner of the service as well as the owner of any specified registry credential.

Optional:

- `registry_credential_id` (String) Optional reference to the registry credential passed to the image repository to retrieve this image.


<a id="nestedatt--secret_files"></a>
### Nested Schema for `secret_files`

Optional:

- `content` (String) The content of the secret file
- `name` (String) The name of the secret file

## Import

Import is supported using the following syntax:

```shell
# WebService can be imported by specifying the id.
terraform import render_web_service.example srv-cabcdefghijklmnopqest
```