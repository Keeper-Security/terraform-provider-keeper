# reference the root node
data "keeper_node" "root" {
  is_root = true
}

# Load all nodes
data "keeper_nodes" "by_subnode" {
  subnodes = {
    include_parent = true
    node_id        = data.keeper_node.root.node_id
    cascade        = true
  }
}

# Get root node direct subnodes
data "keeper_nodes" "by_parent_id" {
  filter = {
    field = "parent_id"
    value = data.keeper_node.root.node_id
  }
}

# Get all nodes that have restrict_visibility property set tot True
data "keeper_nodes" "by_visibility" {
  filter = {
    field = "restrict_visibility"
    value = "true"
  }
}