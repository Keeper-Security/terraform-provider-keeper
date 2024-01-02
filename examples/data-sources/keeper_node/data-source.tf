# reference the root node
data "keeper_node" "root" {
  is_root = true
}

#
data "keeper_node" "subnode" {
  name      = "Subnode"
  parent_id = data.keeper_node.root.node_id
}