data "keeper_node" "root" {
  is_root = true
}
resource "keeper_node" "sso_node" {
  name      = "SSO Node"
  parent_id = data.keeper_node.root.node_id
}
data "keeper_users" "all_sso_users" {
  nodes = {
    node_id = resource.keeper_node.sso_node.node_id
    cascade = true
  }
  is_active = true
}
resource "keeper_team" "contractors" {
  name    = "Contractors"
  node_id = data.keeper_node.root.node_id
}
resource "keeper_role" "sso_user_role" {
  node_id = resource.keeper_node.sso_node.node_id
  name    = "SSO User Role"
}
resource "keeper_role_membership" "example" {
  role_id = resource.keeper_role.sso_user_role.role_id
  users   = data.keeper_users.all_sso_users.users[*].enterprise_user_id
  teams   = [resource.keeper_team.contractors.team_uid]
}