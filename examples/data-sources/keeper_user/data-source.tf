data "keeper_user" "example_by_username" {
  username      = "user@company.com"
  include_roles = true
  include_teams = true
}