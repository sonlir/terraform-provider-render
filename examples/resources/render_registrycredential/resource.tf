resource "render_registrycredential" "example" {
  name       = "example"
  registry   = "DOCKER"
  username   = "docker-username"
  auth_token = "pat-token"
  owner_id   = "usr-abcdefghijklmnopqest"
}
