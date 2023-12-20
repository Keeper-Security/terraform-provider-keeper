package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-keeper/internal/model"
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
					resource.TestCheckResourceAttr("keeper_node.SSO", "name", "SSO"),
					resource.TestCheckResourceAttr("keeper_node.SSO", "restrict_visibility", "true"),
				),
			},
		},
	})
}

const testConfigNodeCreate = `
data "keeper_node" "root" {
	is_root = true
}
resource "keeper_node" "SSO" {
	name = "SSO"
	parent_id = data.keeper_node.root.node_id
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
				ResourceName:  "keeper_node.subnode",
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
					resource.TestCheckResourceAttr("keeper_node.subnode", "name", "Subnodes"),
				),
			},
		},
	})
}

const testConfigNodeImport = `
resource "keeper_node" "subnode" {
	name = "Subnodes"
}
`
