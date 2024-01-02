data "keeper_node" "root" {
  is_root = true
}

resource "keeper_node" "SSO" {
  name                = "SSO"
  parent_id           = data.keeper_node.root.node_id
  restrict_visibility = true
}