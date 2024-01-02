package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strconv"
	"terraform-provider-keeper/internal/model"
	"testing"
)

func TestAccRoleMembershipResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssignRoleMembership,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("keeper_role_membership.example", "users.#", "0"),
					resource.TestCheckResourceAttr("keeper_role_membership.example", "teams.#", "1"),
				),
			},
		},
	})
}

const testAccAssignRoleMembership = `
data "keeper_node" "root" {
  is_root = true
}
resource "keeper_node" "sso_node" {
  name = "SSO Node"
  parent_id = data.keeper_node.root.node_id
}
data "keeper_users" "all_sso_users" {
  nodes = {
    node_id = resource.keeper_node.sso_node.node_id
    cascade = true
  }
  is_active = true
}
resource "keeper_team" "contractors" {
  name = "Contractors"
  node_id = data.keeper_node.root.node_id
}
resource "keeper_role" "sso_user_role" {
  node_id = resource.keeper_node.sso_node.node_id
  name = "SSO User Role"
}
resource "keeper_role_membership" "example" {
  role_id = resource.keeper_role.sso_user_role.role_id
  users = data.keeper_users.all_sso_users.users[*].enterprise_user_id
  teams = [resource.keeper_team.contractors.team_uid]
}
`

func TestAccRoleMembershipResource_Import(t *testing.T) {
	var roleId = 5299989643267
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import
			{
				Config:        testAccImportRoleMembership,
				ImportState:   true,
				ResourceName:  "keeper_role_membership.example",
				ImportStateId: strconv.Itoa(roleId),
				ImportStateCheck: model.ComposeAggregateImportStateCheckFunc(
					model.TestCheckImportStateAttr("users.#", "1"),
				),
			},
		},
	})
}

const testAccImportRoleMembership = `
resource "keeper_role_membership" "example" {
}
`
