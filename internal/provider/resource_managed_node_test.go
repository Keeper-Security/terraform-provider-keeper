package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-keeper/internal/model"
	"testing"
)

func TestAccManagedNodeResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigManagedNodeCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("keeper_managed_node.SSO", "role_id", "5299989643267"),
				),
			},
		},
	})
}

const testConfigManagedNodeCreate = `
data "keeper_node" "root" {
	is_root = true
}
resource "keeper_node" "SSO" {
	name = "SSO"
	parent_id = data.keeper_node.root.node_id
	restrict_visibility = true
}
data "keeper_privileges" "full_admin" {
  manage_nodes = true
  manage_users = true
  manage_roles = false
  manage_teams = false
  manage_reports = true
  manage_sso = true
  device_approval = true
  manage_record_types = false
  share_admin = true
  run_compliance_reports = false
  transfer_account = false
  manage_companies = false
}
resource "keeper_managed_node" "SSO" {
    role_id = 5299989643267
	node_id = resource.keeper_node.SSO.node_id
    cascade_node_management = true
    privileges = data.keeper_privileges.full_admin
}
`

func TestAccManagedNodeResource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import
			{
				Config:        testConfigManagedNodeImport,
				ImportState:   true,
				ResourceName:  "keeper_managed_node.root_admin",
				ImportStateId: fmt.Sprintf("%d,%d", 5299989643267, 5299989643266),
				ImportStateCheck: model.ComposeAggregateImportStateCheckFunc(
					model.TestCheckImportStateAttr("role_id", "5299989643267"),
					model.TestCheckImportStateAttr("node_id", "5299989643266"),
					model.TestCheckImportStateAttr("cascade_node_management", "true"),
					model.TestCheckImportStateAttr("privileges.manage_nodes", "true"),
				),
			},
		},
	})
}

const testConfigManagedNodeImport = `
resource "keeper_managed_node" "root_admin" {
}
`
