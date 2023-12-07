package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kepr/internal/model"
	"testing"
)

func TestAccTeamMembershipResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAllActiveToEveryone,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kepr_team_membership.everyone", "team_uid", "MWaZlKLGNa585bX6sCui3g"),
					resource.TestCheckResourceAttr("kepr_team_membership.everyone", "users.#", "1"),
					resource.TestCheckTypeSetElemAttr("kepr_team_membership.everyone", "users.*", "5299989643284"),
				),
			},
		},
	})
}

const testAccAllActiveToEveryone = `
data "kepr_team" "everyone" {
	team_uid = "MWaZlKLGNa585bX6sCui3g"
}

data "kepr_users" "active_users" {
	is_active = true
}

resource "kepr_team_membership" "everyone" {
  team_uid = data.kepr_team.everyone.team_uid
  users = data.kepr_users.active_users.users[*].enterprise_user_id
}
`

func TestAccTeamMembershipResource_Import(t *testing.T) {
	var teamUid = "MWaZlKLGNa585bX6sCui3g"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import
			{
				Config:        testAccImportEveryone,
				ImportState:   true,
				ResourceName:  "kepr_team_membership.everyone",
				ImportStateId: teamUid,
				ImportStateCheck: model.ComposeAggregateImportStateCheckFunc(
					model.TestCheckImportStateAttr("team_uid", "MWaZlKLGNa585bX6sCui3g"),
					model.TestCheckImportStateAttr("users.#", "0"),
				),
			},
		},
	})
}

const testAccImportEveryone = `
data "kepr_team" "everyone" {
	team_uid = "MWaZlKLGNa585bX6sCui3g"
}

resource "kepr_team_membership" "everyone" {
  team_uid = data.kepr_team.everyone.team_uid
  users = [5299989643285]
}
`
