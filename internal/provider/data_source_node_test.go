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
					resource.TestCheckResourceAttr("data.kepr_node.root", "node_id", "5299989643266"),
					resource.TestCheckResourceAttr("data.kepr_node.root", "name", "Kepr TF"),
					resource.TestCheckResourceAttr("data.kepr_node.root", "is_root", "true"),
					resource.TestCheckResourceAttr("data.kepr_node.subnode", "node_id", "5299989643274"),
					resource.TestCheckResourceAttr("data.kepr_node.subnode", "parent_id", "5299989643266"),
					resource.TestCheckResourceAttr("data.kepr_node.subnode", "name", "Subnode"),
				),
			},
		}})
}

const testConfigNodeDataSource = `
data "kepr_node" "root" {
  is_root = true
}
data "kepr_node" "subnode" {
  name = "Subnode"
  parent_id = data.kepr_node.root.node_id
}
`
