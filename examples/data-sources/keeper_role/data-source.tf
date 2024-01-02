data "keeper_role" "example_by_id" {
  role_id = 5299989643267
}

data "keeper_role" "example_by_name" {
  name            = "Role Name"
  include_members = true
}