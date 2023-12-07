package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceRoles_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigDataSourceRoles,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_roles.by_node", "roles.#", "1"),
					resource.TestCheckResourceAttr("data.kepr_roles.by_node", "roles.0.role_id", "5299989643267"),
					resource.TestCheckResourceAttr("data.kepr_roles.by_is_admin", "roles.#", "1"),
					resource.TestCheckResourceAttr("data.kepr_roles.by_is_admin", "roles.0.role_id", "5299989643267"),
					resource.TestCheckResourceAttr("data.kepr_roles.by_name", "roles.#", "1"),
					resource.TestCheckResourceAttr("data.kepr_roles.by_name", "roles.0.role_id", "5299989643267"),
				),
			},
		}})
}

const testConfigDataSourceRoles = `
data "kepr_roles" "by_node" {
  nodes = {
    node_id = 5299989643266
    cascade = false
  }
}
data "kepr_roles" "by_is_admin" {
  filter = {
    field = "is_admin"
    value = true
  }
}
data "kepr_roles" "by_name" {
  filter = {
    field = "name"
    cmp = "matches"
    value = "keeper*"
  }
}
`
