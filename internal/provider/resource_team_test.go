package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-keeper/internal/model"
	"testing"
)

func TestAccTeamResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigTeamCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("keeper_team.Readers", "name", "SSO Team"),
					resource.TestCheckResourceAttr("keeper_team.Readers", "node_id", "5299989643266"),
					resource.TestCheckResourceAttr("keeper_team.Readers", "restrict_edit", "true"),
				),
			},
		},
	})
}

const testConfigTeamCreate = `
data "keeper_node" "root" {
	is_root = true
}
resource "keeper_team" "Readers" {
	name = "SSO Team"
	node_id = data.keeper_node.root.node_id
	restrict_edit = true
	restrict_share = true
}
`

func TestAccTeamResource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import
			{
				Config:        testConfigTeamImport,
				ImportState:   true,
				ResourceName:  "keeper_team.all_users",
				ImportStateId: "MWaZlKLGNa585bX6sCui3g",
				ImportStateCheck: model.ComposeAggregateImportStateCheckFunc(
					model.TestCheckImportStateAttr("team_uid", "MWaZlKLGNa585bX6sCui3g"),
					model.TestCheckImportStateAttr("node_id", "5299989643266"),
					model.TestCheckImportStateAttr("name", "Everyone"),
				),
				ImportStatePersist: true,
			},
			// Update and Read testing
			{
				Config: testConfigTeamImport,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("keeper_team.all_users", "name", "Anyone"),
				),
			},
		},
	})
}

const testConfigTeamImport = `
resource "keeper_team" "all_users" {
	name = "Anyone"
}
`
