package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceUsers_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUsersDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_users.by_node", "users.#", "2"),
					resource.TestCheckResourceAttr("data.kepr_users.by_status", "users.#", "1"),
					resource.TestCheckResourceAttr("data.kepr_users.by_emails", "users.#", "1"),
				),
			},
		}})
}

const testConfigUsersDataSource = `
data "kepr_users" "by_node" {
  nodes = {
    node_id = 5299989643266
    cascade = true
  }
}
data "kepr_users" "by_status" {
  filter = {
    field = "status"
    value = "inactive"
  }
}
data "kepr_users" "by_emails" {
  is_active = false
  emails = ["pending_user@company.com"]
}
`
