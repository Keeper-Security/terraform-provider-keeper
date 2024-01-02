data "keeper_team" "everyone" {
  team_uid = "MWaZlKLGNa585bX6sCui3g"
}

data "keeper_users" "active_users" {
  is_active = true
}

resource "keeper_team_membership" "example" {
  team_uid = data.keeper_team.everyone.team_uid
  users    = data.keeper_users.active_users.users[*].enterprise_user_id
}