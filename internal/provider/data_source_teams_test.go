package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceTeams_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigTeamsDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kepr_teams.by_node", "teams.#", "1"),
					resource.TestCheckResourceAttr("data.kepr_teams.by_node", "teams.0.team_uid", "MWaZlKLGNa585bX6sCui3g"),
					resource.TestCheckResourceAttr("data.kepr_teams.by_node", "teams.0.node_id", "5299989643266"),
					resource.TestCheckResourceAttr("data.kepr_teams.by_node", "teams.0.name", "Everyone"),
					resource.TestCheckResourceAttr("data.kepr_teams.by_name", "teams.#", "1"),
					resource.TestCheckResourceAttr("data.kepr_teams.by_node", "teams.0.team_uid", "MWaZlKLGNa585bX6sCui3g"),
				),
			},
		}})
}

const testConfigTeamsDataSource = `
data "kepr_teams" "by_node" {
  nodes = {
    node_id = 5299989643266
    cascade = true
  }
}
data "kepr_teams" "by_name" {
  filter = {
    field = "name"
    value = "everyone"
  }
}
`
