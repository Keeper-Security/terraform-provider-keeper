---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "keeper_role_enforcements Resource - terraform-provider-keeper"
subcategory: ""
description: |-
  
---

# keeper_role_enforcements (Resource)



## Example Usage

```terraform
data "keeper_node" "root" {
  is_root = true
}

resource "keeper_role" "sso_user" {
  node_id = data.keeper_node.root.node_id
  name    = "SSO User Role"
}

data "keeper_enforcements_account" "limited" {
  restrict_account_recovery = true
  logout_timer_desktop      = 30
  logout_timer_mobile       = 5
  logout_timer_web          = 30
}

resource "keeper_role_enforcements" "example" {
  role_id = resource.keeper_role.sso_user.role_id
  enforcements = {
    account = data.keeper_enforcements_account.limited
    sharing = {
      restrict_import = true
      restrict_export = true
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enforcements` (Attributes) Role Enforcements (see [below for nested schema](#nestedatt--enforcements))
- `role_id` (Number) Role ID

<a id="nestedatt--enforcements"></a>
### Nested Schema for `enforcements`

Optional:

- `account` (Attributes) Account-related enforcements (see [below for nested schema](#nestedatt--enforcements--account))
- `allow_ip_list` (Attributes) IP whitelist enforcements (see [below for nested schema](#nestedatt--enforcements--allow_ip_list))
- `keeper_fill` (Attributes) Keeper Fill enforcements (see [below for nested schema](#nestedatt--enforcements--keeper_fill))
- `login` (Attributes) Login-related enforcements (see [below for nested schema](#nestedatt--enforcements--login))
- `platform` (Attributes) Keeper platform enforcements (see [below for nested schema](#nestedatt--enforcements--platform))
- `record_types` (Attributes) Record-type enforcements (see [below for nested schema](#nestedatt--enforcements--record_types))
- `sharing` (Attributes) Sharing enforcements (see [below for nested schema](#nestedatt--enforcements--sharing))
- `two_factor` (Attributes) 2FA enforcements (see [below for nested schema](#nestedatt--enforcements--two_factor))
- `vault` (Attributes) Vault-related enforcements (see [below for nested schema](#nestedatt--enforcements--vault))

<a id="nestedatt--enforcements--account"></a>
### Nested Schema for `enforcements.account`

Optional:

- `allow_pam_discovery` (Boolean) Allow PAM discovery
- `allow_pam_rotation` (Boolean) Allow PAM rotation
- `allow_secrets_manager` (Boolean) Allow Keeper Secret Manager access
- `disable_onboarding` (Boolean) Disable onboarding
- `disable_setup_tour` (Boolean) Disable setup-tour
- `disallow_v2_clients` (Boolean) Disallow v2 clients
- `logout_timer_desktop` (Number) Automatic-logout delay for desktop client
- `logout_timer_mobile` (Number) Automatic-logout delay for mobile client
- `logout_timer_web` (Number) Automatic-logout delay for web client
- `max_session_login_time` (Number) Max session login time
- `maximum_record_size` (Number) Maximum record-size
- `minimum_pbkdf2_iterations` (Number) Minimum PBKDF2 iterations
- `require_account_recovery_approval` (Boolean) Require account recovery approval
- `require_device_approval` (Boolean) Require device approval
- `require_security_key_pin` (Boolean) Require security key PIN
- `require_self_destruct` (Boolean) Require self-destruct
- `resend_enterprise_invite_in_x_days` (Number) Resend enterprise invite in X days
- `restrict_account_recovery` (Boolean) Restrict Account Recovery
- `restrict_email_change` (Boolean) Restrict change of email address
- `restrict_import_shared_folders` (Boolean) Restrict shared-folder imports
- `restrict_ip_autoapproval` (Boolean) Restrict IP auto-approval
- `restrict_offline_access` (Boolean) Restrict offline access
- `restrict_persistent_login` (Boolean) Restrict persistent login
- `restrict_personal_license` (Boolean) Restrict use of personal license
- `send_invite_at_registration` (Boolean) Send invite at registration
- `stay_logged_in_default` (Boolean) Enable staying logged-in by default


<a id="nestedatt--enforcements--allow_ip_list"></a>
### Nested Schema for `enforcements.allow_ip_list`

Optional:

- `restrict_ip_addresses` (List of String) Restrict IP addresses
- `restrict_vault_ip_addresses` (List of String) Restrict vault-access from IP addresses
- `tip_zone_restrict_allowed_ip_ranges` (List of String) Restrict allowed IP ranges for tip zone


<a id="nestedatt--enforcements--keeper_fill"></a>
### Nested Schema for `enforcements.keeper_fill`

Optional:

- `keeper_fill_auto_fill` (String) Keeper Fill auto-fill
- `keeper_fill_auto_submit` (String) Keeper Fill auto-submit
- `keeper_fill_auto_suggest` (String) Keeper auto-fill suggestion
- `keeper_fill_hover_locks` (String) Keeper Fill hover locks
- `keeper_fill_match_on_subdomain` (String) Keeper Fill subdomains to match on
- `master_password_reentry` (String) Master password re-entry
- `restrict_auto_fill` (Boolean) Restrict auto-fill
- `restrict_auto_submit` (Boolean) Restrict auto-submit
- `restrict_domain_access` (String) Restrict access to domain(s)
- `restrict_domain_create` (String) Restrict creation of new records for domain(s)
- `restrict_hover_locks` (Boolean) Restrict hover-locks
- `restrict_http_fill_warning` (Boolean) Restrict HTTP fill warning
- `restrict_prompt_to_change` (Boolean) Restrict prompt to change
- `restrict_prompt_to_disable` (Boolean) Restrict prompt to disable Keeper Fill
- `restrict_prompt_to_fill` (Boolean) Restrict prompt to fill
- `restrict_prompt_to_login` (Boolean) Restrict prompt to login
- `restrict_prompt_to_save` (Boolean) Restrict prompt to save


<a id="nestedatt--enforcements--login"></a>
### Nested Schema for `enforcements.login`

Optional:

- `allow_alternate_passwords` (Boolean) Allow alternate passwords
- `master_password_expired_as_of` (Number) Master password expiration
- `master_password_maximum_days_before_change` (Number) Maximum days before master password change
- `master_password_minimum_digits` (Number) Minimum # of digits required for master password
- `master_password_minimum_length` (Number) Minimum length required for master password
- `master_password_minimum_lower` (Number) Minimum # of lower-case characters required for master password
- `master_password_minimum_special` (Number) Minimum # of special characters required for master password
- `master_password_minimum_upper` (Number) Minimum # of upper-case characters required for master password
- `master_password_restrict_days_before_reuse` (Number) # of days before master password can be re-used
- `restrict_android_fingerprint` (Boolean) Restrict Android fingerprint login
- `restrict_ios_fingerprint` (Boolean) Restrict iOS fingerprint login
- `restrict_mac_fingerprint` (Boolean) Restrict MacOS fingerprint login
- `restrict_windows_fingerprint` (Boolean) Restrict Windows fingerprint login


<a id="nestedatt--enforcements--platform"></a>
### Nested Schema for `enforcements.platform`

Optional:

- `restrict_chat_desktop_access` (Boolean) Restrict access to Keeper Chat for desktop
- `restrict_chat_mobile_access` (Boolean) Restrict access to Keeper Chat for mobile
- `restrict_commander_access` (Boolean) Restrict access to Keeper Commander
- `restrict_desktop_access` (Boolean) Restrict access to Keeper Vault for desktop
- `restrict_desktop_mac_access` (Boolean) Restrict access to Keeper Desktop for MacOS
- `restrict_desktop_win_access` (Boolean) Restrict access to Keeper Desktop for Windows
- `restrict_extensions_access` (Boolean) Restrict access to Keeper browser extensions
- `restrict_mobile_access` (Boolean) Restrict access to Keeper Vault for mobile
- `restrict_mobile_android_access` (Boolean) Restrict access to Keeper Vault for Android
- `restrict_mobile_ios_access` (Boolean) Restrict access to Keeper Vault for iOS
- `restrict_mobile_windows_phone_access` (Boolean) Restrict access to Keeper Vault for Windows mobile
- `restrict_web_vault_access` (Boolean) Restrict access to Keeper Vault for web


<a id="nestedatt--enforcements--record_types"></a>
### Nested Schema for `enforcements.record_types`

Optional:

- `restrict_record_types` (Set of String) Restrict record-types


<a id="nestedatt--enforcements--sharing"></a>
### Nested Schema for `enforcements.sharing`

Optional:

- `require_account_share` (Number) Require account-share
- `restrict_create_shared_folder` (Boolean) Restrict shared-folder creation
- `restrict_export` (Boolean) Restrict record exports
- `restrict_file_upload` (Boolean) Restrict file uploads
- `restrict_import` (Boolean) Restrict record imports
- `restrict_link_sharing` (Boolean) Restrict link-sharing
- `restrict_sharing_all_incoming` (Boolean) Restrict all incoming sharing
- `restrict_sharing_all_outgoing` (Boolean) Restrict all outgoing sharing
- `restrict_sharing_enterprise_incoming` (Boolean) Restrict sharing from outside the enterprise
- `restrict_sharing_enterprise_outgoing` (Boolean) Restrict sharing to outside the enterprise
- `restrict_sharing_outside_of_isolated_nodes` (Boolean) Restrict sharing outside of isolated nodes
- `restrict_sharing_record_to_shared_folders` (Boolean) Restrict record sharing within shared-folders
- `restrict_sharing_record_with_attachments` (Boolean) Restrict sharing records with attachments


<a id="nestedatt--enforcements--two_factor"></a>
### Nested Schema for `enforcements.two_factor`

Optional:

- `require_two_factor` (Boolean) Require 2fa for login
- `restrict_two_factor_channel_dna` (Boolean) Restrict use of KeeperDNA for 2fa
- `restrict_two_factor_channel_duo` (Boolean) Restrict use of DUO for 2fa
- `restrict_two_factor_channel_google` (Boolean) Restrict use of Google Authenticator for 2fa
- `restrict_two_factor_channel_rsa` (Boolean) Restrict use of RSA SecurID for 2fa
- `restrict_two_factor_channel_security_keys` (Boolean) Restrict use of security keys (FIDO2 WebAuthn) for 2fa
- `restrict_two_factor_channel_text` (Boolean) Restrict use of SMS/text message for 2fa
- `two_factor_duration_desktop` (String) 2fa duration for desktop client app
- `two_factor_duration_mobile` (String) 2fa duration for mobile client app
- `two_factor_duration_web` (String) 2fa duration for web client app


<a id="nestedatt--enforcements--vault"></a>
### Nested Schema for `enforcements.vault`

Optional:

- `days_before_deleted_records_auto_cleared` (Number) # of days before deleted records are automatically cleared
- `days_before_deleted_records_cleared_perm` (Number) # of days before deleted records are automatically cleared permanently
- `generated_password_complexity` (String) Generated password complexity
- `generated_security_question_complexity` (String) Generated security question complexity
- `mask_custom_fields` (Boolean) Mask custom fields
- `mask_notes` (Boolean) Mask notes
- `mask_passwords_while_editing` (Boolean) Mask passwords while editing
- `restrict_breach_watch` (Boolean) Restrict BreachWatch
- `restrict_create_folder` (Boolean) Restrict creation of folders
- `restrict_create_folder_to_only_shared_folders` (Boolean) Restrict creation of folders to within shared-folders only
- `restrict_create_identity_payment_records` (Boolean) Restrict creation of identity payment records
- `restrict_create_record` (Boolean) Restrict creation of records
- `restrict_create_record_to_shared_folders` (Boolean) Restrict record creation within shared-folders
- `send_breach_watch_events` (Boolean) Send BreachWatch events

## Import

Import is supported using the following syntax:

```shell
terraform import keeper_role_enforcements.example 5299989643267
```
