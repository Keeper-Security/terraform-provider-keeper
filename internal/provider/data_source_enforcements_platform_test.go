package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcementsPlatform(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsPlatformDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_platform.test", "restrict_commander_access", "true"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_platform.test", "restrict_extensions_access", "true"),
				),
			},
		}})
}

const testAccEnforcementsPlatformDataSourceConfig = `
data "keeper_enforcements_platform" "test" {
  restrict_commander_access = true
  restrict_extensions_access = true
}
`
