package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceTeam_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigTeamDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_team.by_team_uid", "team_uid", "MWaZlKLGNa585bX6sCui3g"),
					resource.TestCheckResourceAttr("data.kepr_team.by_team_uid", "users.#", "1"),
				),
			},
		}})
}

const testConfigTeamDataSource = `
data "kepr_team" "by_team_uid" {
  team_uid = "MWaZlKLGNa585bX6sCui3g"
  include_members = true
}
`
