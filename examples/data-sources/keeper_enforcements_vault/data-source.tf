data "keeper_enforcements_vault" "example" {
  restrict_breach_watch                         = true
  mask_notes                                    = true
  days_before_deleted_records_auto_cleared      = 30
  mask_passwords_while_editing                  = true
  mask_custom_fields                            = true
  restrict_create_folder_to_only_shared_folders = true
  restrict_create_record_to_shared_folders      = true
}
