package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcements2Fa(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcements2FaDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_2fa.test", "two_factor_duration_desktop", "0,12,24,30,9999"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_2fa.test", "restrict_two_factor_channel_google", "true"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_2fa.test", "restrict_two_factor_channel_text", "true"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_2fa.test", "require_two_factor", "true"),
				),
			},
		}})
}

const testAccEnforcements2FaDataSourceConfig = `
data "keeper_enforcements_2fa" "test" {
  two_factor_duration_desktop = "0,12,24,30,9999"
  restrict_two_factor_channel_google = true
  restrict_two_factor_channel_text = true
  require_two_factor = true
}
`
