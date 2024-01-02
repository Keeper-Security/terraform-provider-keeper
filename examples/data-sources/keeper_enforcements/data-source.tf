data "keeper_enforcements" "example" {
  account = {
    logout_timer_desktop      = 40
    restrict_ip_autoapproval  = true
    disable_onboarding        = true
    restrict_offline_access   = true
    restrict_persistent_login = true
    require_device_approval   = true
  }
  sharing = {
    restrict_export               = true
    restrict_import               = true
    restrict_create_shared_folder = true
  }
  keeper_fill = {
    restrict_auto_fill       = true
    keeper_fill_auto_suggest = "disable"
    keeper_fill_auto_fill    = "disable"
  }
  login = {
    master_password_minimum_length = 12
    restrict_windows_fingerprint   = true
  }
  vault = {
    restrict_breach_watch  = true
    restrict_create_folder = true
  }
}
