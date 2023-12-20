package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-keeper/internal/model"
	"testing"
)

func TestAccRoleResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigRoleCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("keeper_role.sso_user", "name", "SSO User"),
					resource.TestCheckResourceAttr("keeper_role.sso_user", "visible_below", "true"),
				),
			},
		},
	})
}

const testConfigRoleCreate = `
data "keeper_node" "root" {
	is_root = true
}
resource "keeper_role" "sso_user" {
	name = "SSO User"
	node_id = data.keeper_node.root.node_id
	visible_below = true
}
`

func TestAccRoleResource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import
			{
				Config:        testConfigRoleImport,
				ImportState:   true,
				ResourceName:  "keeper_role.admin",
				ImportStateId: "5299989643267",
				ImportStateCheck: model.ComposeAggregateImportStateCheckFunc(
					model.TestCheckImportStateAttr("role_id", "5299989643267"),
					model.TestCheckImportStateAttr("node_id", "5299989643266"),
					model.TestCheckImportStateAttr("name", "Keeper Administrator"),
				),
				ImportStatePersist: true,
			},
		},
	})
}

const testConfigRoleImport = `
resource "keeper_role" "admin" {
	name = "Keeper Administrator"
}
`
