data "keeper_enforcements_2fa" "example" {
  two_factor_duration_desktop        = "forever"
  two_factor_duration_mobile         = "30_days"
  restrict_two_factor_channel_google = true
  restrict_two_factor_channel_text   = true
  require_two_factor                 = true
}