package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDataSourceNodes_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigNodesDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.keeper_nodes.by_subnode", "nodes.#", "2"),
					resource.TestCheckResourceAttr("data.keeper_nodes.by_subnode", "nodes.0.node_id", "5299989643266"),
					resource.TestCheckResourceAttr("data.keeper_nodes.by_subnode", "nodes.1.node_id", "5299989643274"),
					resource.TestCheckResourceAttr("data.keeper_nodes.by_parent_id", "nodes.#", "1"),
					resource.TestCheckResourceAttr("data.keeper_nodes.by_parent_id", "nodes.0.node_id", "5299989643274"),
					resource.TestCheckResourceAttr("data.keeper_nodes.by_visibility", "nodes.#", "1"),
					resource.TestCheckResourceAttr("data.keeper_nodes.by_visibility", "nodes.0.node_id", "5299989643274"),
				),
			},
		}})
}

const testConfigNodesDataSource = `
data "keeper_node" "root" {
  is_root = true
}

data "keeper_nodes" "by_subnode" {
  subnodes = {
    include_parent = true
    node_id = data.keeper_node.root.node_id
    cascade = true
  }
}
data "keeper_nodes" "by_parent_id" {
  filter = {
    field = "parent_id"
    value = "5299989643266"
  }
}
data "keeper_nodes" "by_visibility" {
  filter = {
    field = "restrict_visibility"
    value = "true"
  }
}
`
