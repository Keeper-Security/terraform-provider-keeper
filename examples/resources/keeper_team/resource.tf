data "keeper_node" "root" {
  is_root = true
}
resource "keeper_team" "example" {
  name           = "SSO Team"
  node_id        = data.keeper_node.root.node_id
  restrict_edit  = true
  restrict_share = true
}