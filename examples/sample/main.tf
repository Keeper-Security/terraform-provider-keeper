terraform {
  required_providers {
    keeper = {
      source = "Keeper-Security/keeper"
    }
  }
}

provider "keeper" {
  // config_type = "commander"
}

data "keeper_node" "root" {
  is_root = true
}

resource "keeper_node" "sso" {
  name      = "SSO"
  parent_id = data.keeper_node.root.node_id
}

resource "keeper_role" "sso_users" {
  name          = "SSO Users"
  node_id       = resource.keeper_node.sso.node_id
  visible_below = true
}

data "keeper_enforcements_login" "enforcements" {
  allow_alternate_passwords = true
}

data "keeper_enforcements_account" "enforcements" {
  disable_onboarding             = true
  disable_setup_tour             = true
  logout_timer_desktop           = 10
  logout_timer_mobile            = 3
  logout_timer_web               = 10
  restrict_email_change          = true
  restrict_import_shared_folders = true
  restrict_persistent_login      = true
  restrict_personal_license      = true
  //stay_logged_in_default         = true
}

data "keeper_enforcements_vault" "enforcements" {
  restrict_breach_watch = true
}

data "keeper_enforcements_2fa" "enforcements" {
  two_factor_duration_mobile = "30_days"
}

data "keeper_enforcements_sharing" "enforcements" {
  restrict_export = true
  restrict_import = true
}

data "keeper_enforcements_keeper_fill" "enforcements" {
  keeper_fill_hover_locks = "disable"
}

data "keeper_enforcements_record_types" "enforcements" {
  restrict_record_types = [
    "address", "bankAccount", "bankCard", "contact", "databaseCredentials",
    "encryptedNotes", "healthInsurance", "membership",
    "passport", "photo", "softwareLicense", "ssnCard"
  ]
}

resource "keeper_role_enforcements" "sso_enf" {
  role_id = resource.keeper_role.sso_users.role_id
  enforcements = {
    login        = data.keeper_enforcements_login.enforcements
    account      = data.keeper_enforcements_account.enforcements
    vault        = data.keeper_enforcements_vault.enforcements
    sharing      = data.keeper_enforcements_sharing.enforcements
    keeper_fill  = data.keeper_enforcements_keeper_fill.enforcements
    two_factor   = data.keeper_enforcements_2fa.enforcements
    record_types = data.keeper_enforcements_record_types.enforcements
  }
}

resource "keeper_team" "sso_everyone" {
  name           = "All SSO Users"
  node_id        = resource.keeper_node.sso.node_id
  restrict_edit  = true
  restrict_share = true
}

data "keeper_users" "all_sso_users" {
  nodes = {
    node_id = resource.keeper_node.sso.node_id
    cascade = true
  }
}

resource "keeper_team_membership" "sso_team_membership" {
  team_uid = resource.keeper_team.sso_everyone.team_uid
  users    = data.keeper_users.all_sso_users.users[*].enterprise_user_id
}

resource "keeper_role_membership" "sso_role_membership" {
  role_id = resource.keeper_role.sso_users.role_id
  users   = []
  teams   = [resource.keeper_team.sso_everyone.team_uid]
}
output "example" {
  value = resource.keeper_node.sso
}