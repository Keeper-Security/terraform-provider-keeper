data "keeper_enforcements_login" "example" {
  allow_alternate_passwords                  = true
  restrict_android_fingerprint               = true
  restrict_windows_fingerprint               = true
  restrict_mac_fingerprint                   = true
  restrict_ios_fingerprint                   = true
  master_password_maximum_days_before_change = 20
  master_password_minimum_length             = 20
  master_password_minimum_upper              = 4
  master_password_minimum_special            = 2
  master_password_minimum_digits             = 2
}
