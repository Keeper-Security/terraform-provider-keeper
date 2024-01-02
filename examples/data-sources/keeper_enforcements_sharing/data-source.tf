data "keeper_enforcements_sharing" "example" {
  require_account_share         = resource.keeper_role.admin_role.role_id
  restrict_link_sharing         = true
  restrict_import               = true
  restrict_export               = true
  restrict_sharing_all_outgoing = true
  restrict_sharing_all_incoming = true
  restrict_create_shared_folder = true
}
