data "keeper_teams" "example_by_node" {
  nodes = {
    node_id = 5299989643266
    cascade = true
  }
}
data "keeper_teams" "example_by_name" {
  filter = {
    field = "name"
    value = "everyone"
  }
}