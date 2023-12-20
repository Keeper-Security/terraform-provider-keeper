package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceEnforcementsVault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnforcementsVaultDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_enforcements_vault.test", "mask_notes", "true"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_vault.test", "mask_passwords_while_editing", "true"),
					resource.TestCheckResourceAttr("data.keeper_enforcements_vault.test", "days_before_deleted_records_auto_cleared", "30"),
				),
			},
		}})
}

const testAccEnforcementsVaultDataSourceConfig = `
data "keeper_enforcements_vault" "test" {
  mask_notes = true
  days_before_deleted_records_auto_cleared = 30
  mask_passwords_while_editing = true
}
`
