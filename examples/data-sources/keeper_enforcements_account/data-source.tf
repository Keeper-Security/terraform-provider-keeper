data "keeper_enforcements_account" "example" {
  logout_timer_desktop      = 30
  logout_timer_mobile       = 5
  logout_timer_web          = 15
  disable_onboarding        = true
  restrict_offline_access   = true
  require_device_approval   = true
  restrict_persistent_login = true
  minimum_pbkdf2_iterations = 1000000
  disable_setup_tour        = true
}
