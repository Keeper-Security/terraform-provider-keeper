package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcements(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements.test", "account.logout_timer_desktop", "40"),
				),
			},
		}})
}

const testAccEnforcementsDataSourceConfig = `
data "keeper_enforcements" "test" {
  	account = {
		logout_timer_desktop = 40
	}
}
`
