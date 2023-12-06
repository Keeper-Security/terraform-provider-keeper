package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccPrivilegesDataSource_Define(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccPrivilegesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_nodes", "true"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_users", "true"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_roles", "false"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_teams", "false"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_reports", "true"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_sso", "false"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "device_approval", "false"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_record_types", "true"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "share_admin", "true"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "run_compliance_reports", "false"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "transfer_account", "false"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_companies", "true"),
				),
			},
		}})
}

const testAccPrivilegesDataSourceConfig = `
data "kepr_privileges" "test" {
  manage_nodes = true
  manage_users = true
  manage_roles = false
  manage_teams = false
  manage_reports = true
  manage_sso = false
  device_approval = false
  manage_record_types = true
  share_admin = true
  run_compliance_reports = false
  transfer_account = false
  manage_companies = true
}
`
