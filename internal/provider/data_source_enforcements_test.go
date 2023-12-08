package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_enforcements.test", "account.automatic_backup_every_x_days", "4"),
				),
			},
		}})
}

const testAccEnforcementsDataSourceConfig = `
data "kepr_enforcements" "test" {
  	account = {
		automatic_backup_every_x_days = 4
	}
}
`
/*const testAccEnforcementsDataSourceConfig = `
data "kepr_enforcements_account" "test" {
	automatic_backup_every_x_days = 4
}

data "kepr_enforcements" "test" {
  	account = data.kepr_enforcements_account.test
}
`
*/