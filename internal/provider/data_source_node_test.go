package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceNode_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigNodeDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_node.root", "node_id", "5299989643266"),
					resource.TestCheckResourceAttr("data.keeper_node.root", "name", "Keeper TF"),
					resource.TestCheckResourceAttr("data.keeper_node.root", "is_root", "true"),
					resource.TestCheckResourceAttr("data.keeper_node.subnode", "node_id", "5299989643274"),
					resource.TestCheckResourceAttr("data.keeper_node.subnode", "parent_id", "5299989643266"),
					resource.TestCheckResourceAttr("data.keeper_node.subnode", "name", "Subnode"),
				),
			},
		}})
}

const testConfigNodeDataSource = `
data "keeper_node" "root" {
  is_root = true
}
data "keeper_node" "subnode" {
  name = "Subnode"
  parent_id = data.keeper_node.root.node_id
}
`
