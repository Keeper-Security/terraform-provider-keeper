package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceUser_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUserDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_user.by_username", "enterprise_user_id", "5299989643284"),
					resource.TestCheckResourceAttr("data.kepr_user.by_username", "teams.#", "1"),
					resource.TestCheckResourceAttr("data.kepr_user.by_username", "roles.#", "1"),
				),
			},
		}})
}

const testConfigUserDataSource = `
data "kepr_user" "by_username" {
  username = "user@company.com"
  include_roles = true
  include_teams = true
}
`
