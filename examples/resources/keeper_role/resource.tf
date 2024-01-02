data "keeper_node" "root" {
  is_root = true
}

resource "keeper_role" "sso_user_role" {
  name          = "SSO User Role"
  node_id       = data.keeper_node.root.node_id
  visible_below = true
}