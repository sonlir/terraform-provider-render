---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "render_registrycredential Data Source - render"
subcategory: ""
description: |-
  RegistryCredential data source
---

# render_registrycredential (Data Source)

RegistryCredential data source

## Example Usage

```terraform
data "render_registrycredential" "example" {
  id = "rgc-abcdefghijklmnopqest"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Unique identifier for this credential

### Read-Only

- `name` (String) Descriptive name for this credential
- `registry` (String) The registry to use this credential with. Valid values are `GITHUB`, `GITLAB`, `DOCKER`.
- `username` (String) The username associated with the credential
