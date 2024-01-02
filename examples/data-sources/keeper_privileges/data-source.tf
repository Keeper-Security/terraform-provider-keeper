data "keeper_privileges" "example" {
  manage_nodes           = true
  manage_users           = true
  manage_roles           = false
  manage_teams           = false
  manage_reports         = true
  manage_sso             = false
  device_approval        = false
  manage_record_types    = true
  share_admin            = true
  run_compliance_reports = false
  transfer_account       = false
  manage_companies       = true
}