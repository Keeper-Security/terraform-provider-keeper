package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourcePrivileges(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccPrivilegesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_nodes", "true"),
					resource.TestCheckResourceAttr("data.kepr_privileges.test", "manage_users", "false"),
				),
			},
		}})
}

const testAccPrivilegesDataSourceConfig = `
data "kepr_privileges" "test" {
  manage_nodes = true
}
`
