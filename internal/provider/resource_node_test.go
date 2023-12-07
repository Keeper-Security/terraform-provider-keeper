package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kepr/internal/model"
	"testing"
)

func TestAccNodeResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigNodeCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kepr_node.SSO", "name", "SSO"),
					resource.TestCheckResourceAttr("kepr_node.SSO", "restrict_visibility", "true"),
				),
			},
		},
	})
}

const testConfigNodeCreate = `
data "kepr_node" "root" {
	is_root = true
}
resource "kepr_node" "SSO" {
	name = "SSO"
	parent_id = data.kepr_node.root.node_id
	restrict_visibility = true
}
`

func TestAccNodeResource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import
			{
				Config:        testConfigNodeImport,
				ImportState:   true,
				ResourceName:  "kepr_node.subnode",
				ImportStateId: "5299989643274",
				ImportStateCheck: model.ComposeAggregateImportStateCheckFunc(
					model.TestCheckImportStateAttr("parent_id", "5299989643266"),
					model.TestCheckImportStateAttr("node_id", "5299989643274"),
					model.TestCheckImportStateAttr("name", "Subnode"),
				),
				ImportStatePersist: true,
			},
			// Update and Read testing
			{
				Config: testConfigNodeImport,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kepr_node.subnode", "name", "Subnodes"),
				),
			},
		},
	})
}

const testConfigNodeImport = `
resource "kepr_node" "subnode" {
	name = "Subnodes"
}
`
