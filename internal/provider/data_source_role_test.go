package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceRole_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigRoleDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_role.root", "role_id", "5299989643267"),
					resource.TestCheckResourceAttr("data.keeper_role.root", "name", "Keeper Administrator"),
				),
			},
		}})
}

const testConfigRoleDataSource = `
data "keeper_role" "root" {
  role_id = 5299989643267
  include_members = true
}
`
