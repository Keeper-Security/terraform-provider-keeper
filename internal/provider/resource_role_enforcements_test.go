package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strconv"
	"terraform-provider-keeper/internal/model"
	"testing"
)

func TestAccRoleEnforcementsResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssignEnforcements,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("keeper_role_enforcements.sso_user_enf", "enforcements.account.logout_timer_desktop", "30"),
					resource.TestCheckResourceAttr("keeper_role_enforcements.sso_user_enf", "enforcements.sharing.restrict_import", "true"),
					resource.TestCheckResourceAttr("keeper_role_enforcements.sso_user_enf", "enforcements.sharing.restrict_export", "true"),
				),
			},
		},
	})
}

const testAccAssignEnforcements = `
data "keeper_node" "root" {
	is_root = true
}
data "keeper_enforcements_account" "limited" {
	restrict_account_recovery = true
	logout_timer_desktop = 30
	logout_timer_mobile = 5
	logout_timer_web = 30
}
data "keeper_enforcements_record_types" "enforcements" {
  restrict_record_types = ["bankAccount", "bankCard"]
}
resource "keeper_role" "sso_user" {
  node_id = data.keeper_node.root.node_id
  name = "SSO User"
}

resource "keeper_role_enforcements" "sso_user_enf" {
  role_id = resource.keeper_role.sso_user.role_id
  enforcements = {
    account = data.keeper_enforcements_account.limited
    sharing = {
      restrict_import = true
      restrict_export = true
    }
    record_types = data.keeper_enforcements_record_types.enforcements
  }
}
`

func TestAccRoleEnforcementsResource_Import(t *testing.T) {
	var roleId = 5299989643267
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import
			{
				Config:        testAccImportAdminEnforcements,
				ImportState:   true,
				ResourceName:  "keeper_role_enforcements.admin",
				ImportStateId: strconv.Itoa(roleId),
				ImportStateCheck: model.ComposeAggregateImportStateCheckFunc(
					model.TestCheckImportStateAttr("enforcements.login.master_password_minimum_length", "20"),
				),
			},
		},
	})
}

const testAccImportAdminEnforcements = `
resource "keeper_role_enforcements" "admin" {
}
`
