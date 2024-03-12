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
