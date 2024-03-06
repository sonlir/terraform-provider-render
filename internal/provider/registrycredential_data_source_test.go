// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegistryCredentialDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccRegistryCredentialDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.render_registrycredential.test", "id", "rgc-cnim7h7sc6pc73d1op80"),
					resource.TestCheckResourceAttr("data.render_registrycredential.test", "name", "asd"),
					resource.TestCheckResourceAttr("data.render_registrycredential.test", "registry", "DOCKER"),
					resource.TestCheckResourceAttr("data.render_registrycredential.test", "username", "ss"),
				),
			},
		},
	})
}

const testAccRegistryCredentialDataSourceConfig = `
data "render_registrycredential" "test" {
  id = "rgc-cnim7h7sc6pc73d1op80"
}
`
