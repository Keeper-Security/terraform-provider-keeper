data "keeper_roles" "example_by_node_id" {
  nodes = {
    node_id = 5299989643266
    cascade = false
  }
}

data "keeper_roles" "example_by_is_admin" {
  filter = {
    field = "is_admin"
    value = true
  }
}

data "keeper_roles" "example_by_name" {
  filter = {
    field = "name"
    cmp   = "matches"
    value = "keeper*"
  }
}