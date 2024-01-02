data "keeper_node" "root" {
  is_root = true
}

resource "keeper_role" "sso_user" {
  node_id = data.keeper_node.root.node_id
  name    = "SSO User Role"
}

data "keeper_enforcements_account" "limited" {
  restrict_account_recovery = true
  logout_timer_desktop      = 30
  logout_timer_mobile       = 5
  logout_timer_web          = 30
}

resource "keeper_role_enforcements" "example" {
  role_id = resource.keeper_role.sso_user.role_id
  enforcements = {
    account = data.keeper_enforcements_account.limited
    sharing = {
      restrict_import = true
      restrict_export = true
    }
  }
}