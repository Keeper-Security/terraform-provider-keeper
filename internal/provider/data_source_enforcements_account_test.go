package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcementsAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsAccountDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_account.test", "minimum_pbkdf2_iterations", "1000000"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_account.test", "restrict_account_recovery", "true"),

					resource.TestCheckResourceAttr("data.keeper_enforcements_account.test", "logout_timer_desktop", "30"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_account.test", "logout_timer_mobile", "3"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_account.test", "logout_timer_web", "15"),
				),
			},
		}})
}

const testAccEnforcementsAccountDataSourceConfig = `
data "keeper_enforcements_account" "test" {
  restrict_account_recovery = true
  minimum_pbkdf2_iterations = 1000000
  logout_timer_desktop = 30
  logout_timer_mobile = 3
  logout_timer_web = 15
}
`
