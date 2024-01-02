data "keeper_users" "example_by_node" {
  nodes = {
    node_id = 5299989643266
    cascade = true
  }
}
data "keeper_users" "example_by_status" {
  filter = {
    field = "status"
    value = "inactive"
  }
}
data "keeper_users" "example_by_emails" {
  is_active = false
  emails    = ["pending_user@company.com"]
}