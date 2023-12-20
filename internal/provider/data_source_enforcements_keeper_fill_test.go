package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcementsKeeperFill(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsKeeperFillDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_keeper_fill.test", "keeper_fill_auto_suggest", "enforce"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_keeper_fill.test", "keeper_fill_auto_submit", "disable"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_keeper_fill.test", "restrict_auto_submit", "true"),
				),
			},
		}})
}

const testAccEnforcementsKeeperFillDataSourceConfig = `
data "keeper_enforcements_keeper_fill" "test" {
  keeper_fill_auto_suggest = "enforce"
  keeper_fill_auto_submit = "disable"
  restrict_auto_submit = true
}
`
