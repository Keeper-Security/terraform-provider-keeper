data "keeper_enforcements_keeper_fill" "example" {
  keeper_fill_auto_suggest  = "enforce"
  keeper_fill_auto_submit   = "disable"
  restrict_auto_submit      = true
  restrict_prompt_to_save   = true
  restrict_prompt_to_change = true
  restrict_prompt_to_fill   = true
  restrict_prompt_to_login  = true
}
