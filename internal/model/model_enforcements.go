package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/api"
	"strconv"
	"strings"
)

var enforcementDescriptions map[string]string

func init() {
	enforcementDescriptions = map[string]string{
		"account":                   "Account-related enforcements",
		"allow_ip_list":             "IP whitelist enforcements",
		"sharing":                   "Sharing enforcements",
		"keeper_fill":               "Keeper Fill enforcements",
		"login":                     "Login-related enforcements",
		"platform":                  "Keeper platform enforcements",
		"record_types":              "Record-type enforcements",
		"two_factor_authentication": "2FA enforcements",
		"vault":                     "Vault-related enforcements",

		"master_password_minimum_length":                "Minimum length required for master password",
		"master_password_minimum_special":               "Minimum # of special characters required for master password",
		"master_password_minimum_upper":                 "Minimum # of upper-case characters required for master password",
		"master_password_minimum_lower":                 "Minimum # of lower-case characters required for master password",
		"master_password_minimum_digits":                "Minimum # of digits required for master password",
		"master_password_restrict_days_before_reuse":    "# of days before master password can be re-used",
		"require_two_factor":                            "Require 2fa for login",
		"master_password_maximum_days_before_change":    "Maximum days before master password change",
		"master_password_expired_as_of":                 "Master password expiration",
		"minimum_pbkdf2_iterations":                     "Minimum PBKDF2 iterations",
		"max_session_login_time":                        "Max session login time",
		"restrict_persistent_login":                     "Restrict persistent login",
		"stay_logged_in_default":                        "Enable staying logged-in by default",
		"restrict_sharing_all":                          "Restrict all sharing",
		"restrict_sharing_enterprise_outgoing":          "Restrict sharing to outside the enterprise",
		"restrict_sharing_enterprise_incoming":          "Restrict sharing from outside the enterprise",
		"restrict_export":                               "Restrict record exports",
		"restrict_import":                               "Restrict record imports",
		"restrict_file_upload":                          "Restrict file uploads",
		"require_account_share":                         "Require account-share",
		"restrict_sharing_incoming_all":                 "Restrict all incoming shares",
		"restrict_sharing_record_with_attachments":      "Restrict sharing records with attachments",
		"restrict_ip_addresses":                         "Restrict IP addresses",
		"require_device_approval":                       "Require device approval",
		"require_account_recovery_approval":             "Require account recovery approval",
		"restrict_vault_ip_addresses":                   "Restrict vault-access from IP addresses",
		"tip_zone_restrict_allowed_ip_ranges":           "Restrict allowed IP ranges for tip zone",
		"restrict_offline_access":                       "Restrict offline access",
		"send_invite_at_registration":                   "Send invite at registration",
		"restrict_email_change":                         "Restrict change of email address",
		"restrict_ios_fingerprint":                      "Restrict iOS fingerprint login",
		"restrict_mac_fingerprint":                      "Restrict MacOS fingerprint login",
		"restrict_android_fingerprint":                  "Restrict Android fingerprint login",
		"restrict_windows_fingerprint":                  "Restrict Windows fingerprint login",
		"logout_timer_web":                              "Automatic-logout delay for web client",
		"logout_timer_mobile":                           "Automatic-logout delay for mobile client",
		"logout_timer_desktop":                          "Automatic-logout delay for desktop client",
		"restrict_web_vault_access":                     "Restrict access to Keeper Vault for web",
		"restrict_extensions_access":                    "Restrict access to Keeper browser extensions",
		"restrict_mobile_access":                        "Restrict access to Keeper Vault for mobile",
		"restrict_desktop_access":                       "Restrict access to Keeper Vault for desktop",
		"restrict_mobile_ios_access":                    "Restrict access to Keeper Vault for iOS",
		"restrict_mobile_android_access":                "Restrict access to Keeper Vault for Android",
		"restrict_mobile_windows_phone_access":          "Restrict access to Keeper Vault for Windows mobile",
		"restrict_desktop_win_access":                   "Restrict access to Keeper Desktop for Windows",
		"restrict_desktop_mac_access":                   "Restrict access to Keeper Desktop for MacOS",
		"restrict_chat_desktop_access":                  "Restrict access to Keeper Chat for desktop",
		"restrict_chat_mobile_access":                   "Restrict access to Keeper Chat for mobile",
		"restrict_commander_access":                     "Restrict access to Keeper Commander",
		"restrict_two_factor_channel_text":              "Restrict use of SMS/text message for 2fa",
		"restrict_two_factor_channel_google":            "Restrict use of Google Authenticator for 2fa",
		"restrict_two_factor_channel_dna":               "Restrict use of KeeperDNA for 2fa",
		"restrict_two_factor_channel_duo":               "Restrict use of DUO for 2fa",
		"restrict_two_factor_channel_rsa":               "Restrict use of RSA SecurID for 2fa",
		"two_factor_duration_web":                       "2fa duration for web client app",
		"two_factor_duration_mobile":                    "2fa duration for mobile client app",
		"two_factor_duration_desktop":                   "2fa duration for desktop client app",
		"restrict_two_factor_channel_security_keys":     "Restrict use of security keys (FIDO2 WebAuthn) for 2fa",
		"two_factor_by_ip":                              "2fa by IP",
		"restrict_domain_access":                        "Restrict access to domain(s)",
		"restrict_domain_create":                        "Restrict creation of new records for domain(s)",
		"restrict_hover_locks":                          "Restrict hover-locks",
		"restrict_prompt_to_login":                      "Restrict prompt to login",
		"restrict_prompt_to_fill":                       "Restrict prompt to fill",
		"restrict_auto_submit":                          "Restrict auto-submit",
		"restrict_prompt_to_save":                       "Restrict prompt to save",
		"restrict_prompt_to_change":                     "Restrict prompt to change",
		"restrict_auto_fill":                            "Restrict auto-fill",
		"restrict_create_folder":                        "Restrict creation of folders",
		"restrict_create_folder_to_only_shared_folders": "Restrict creation of folders to within shared-folders only",
		"restrict_create_identity_payment_records":      "Restrict creation of identity payment records",
		"mask_custom_fields":                            "Mask custom fields",
		"mask_notes":                                    "Mask notes",
		"mask_passwords_while_editing":                  "Mask passwords while editing",
		"generated_password_complexity":                 "Generated password complexity",
		"generated_security_question_complexity":        "Generated security question complexity",
		"days_before_deleted_records_cleared_perm":      "# of days before deleted records are automatically cleared permanently",
		"days_before_deleted_records_auto_cleared":      "# of days before deleted records are automatically cleared",
		"allow_alternate_passwords":                     "Allow alternate passwords",
		"restrict_create_record":                        "Restrict creation of records",
		"restrict_create_record_to_shared_folders":      "Restrict record creation within shared-folders",
		"restrict_sharing_record_to_shared_folders":     "Restrict record sharing within shared-folders",
		"restrict_link_sharing":                         "Restrict link-sharing",
		"restrict_sharing_outside_of_isolated_nodes":    "Restrict sharing outside of isolated nodes",
		"disable_setup_tour":                            "Disable setup-tour",
		"restrict_personal_license":                     "Restrict use of personal license",
		"disable_onboarding":                            "Disable onboarding",
		"disallow_v2_clients":                           "Disallow v2 clients",
		"restrict_ip_autoapproval":                      "Restrict IP auto-approval",
		"send_breach_watch_events":                      "Send BreachWatch events",
		"restrict_breach_watch":                         "Restrict BreachWatch",
		"resend_enterprise_invite_in_x_days":            "Resend enterprise invite in X days",
		"master_password_reentry":                       "Master password re-entry",
		"restrict_account_recovery":                     "Restrict Account Recovery",
		"keeper_fill_hover_locks":                       "Keeper Fill hover locks",
		"keeper_fill_auto_fill":                         "Keeper Fill auto-fill",
		"keeper_fill_auto_submit":                       "Keeper Fill auto-submit",
		"keeper_fill_match_on_subdomain":                "Keeper Fill subdomains to match on",
		"restrict_prompt_to_disable":                    "Restrict prompt to disable Keeper Fill",
		"restrict_http_fill_warning":                    "Restrict HTTP fill warning",
		"restrict_record_types":                         "Restrict record-types",
		"allow_secrets_manager":                         "Allow Keeper Secret Manager access",
		"require_self_destruct":                         "Require self-destruct",
		"keeper_fill_auto_suggest":                      "Keeper auto-fill suggestion",
		"maximum_record_size":                           "Maximum record-size",
		"allow_pam_rotation":                            "Allow PAM rotation",
		"allow_pam_discovery":                           "Allow PAM discovery",
		"restrict_import_shared_folders":                "Restrict shared-folder imports",
		"require_security_key_pin":                      "Require security key PIN",
		"restrict_create_shared_folder":                 "Restrict shared-folder creation",
	}
}

type EnforcementsDataSourceModel struct {
	Account                 *EnforcementsAccountDataSourceModel     `tfsdk:"account"`
	AllowIpList             *EnforcementsAllowIpListDataSourceModel `tfsdk:"allow_ip_list"`
	Sharing                 *EnforcementsSharingDataSourceModel     `tfsdk:"sharing"`
	KeeperFill              *EnforcementsKeeperFillDataSourceModel  `tfsdk:"keeper_fill"`
	Login                   *EnforcementsLoginDataSourceModel       `tfsdk:"login"`
	Platform                *EnforcementsPlatformDataSourceModel    `tfsdk:"platform"`
	RecordTypes             *EnforcementsRecordTypesDataSourceModel `tfsdk:"record_types"`
	TwoFactorAuthentication *Enforcements2faDataSourceModel         `tfsdk:"two_factor_authentication"`
	Vault                   *EnforcementsVaultDataSourceModel       `tfsdk:"vault"`
}

func (em *EnforcementsDataSourceModel) FromKeeper(enforcements map[string]string) {
	if enforcements == nil {
		return
	}

	var account = new(EnforcementsAccountDataSourceModel)
	account.FromKeeper(enforcements)
	if !account.IsBlank() {
		em.Account = account
	}

	var allowIpList = new(EnforcementsAllowIpListDataSourceModel)
	allowIpList.FromKeeper(enforcements)
	if !allowIpList.IsBlank() {
		em.AllowIpList = allowIpList
	}

	var sharing = new(EnforcementsSharingDataSourceModel)
	sharing.FromKeeper(enforcements)
	if !sharing.IsBlank() {
		em.Sharing = sharing
	}

	var keeperFill = new(EnforcementsKeeperFillDataSourceModel)
	keeperFill.FromKeeper(enforcements)
	if !keeperFill.IsBlank() {
		em.KeeperFill = keeperFill
	}

	var login = new(EnforcementsLoginDataSourceModel)
	login.FromKeeper(enforcements)
	if !login.IsBlank() {
		em.Login = login
	}

	var platform = new(EnforcementsPlatformDataSourceModel)
	platform.FromKeeper(enforcements)
	if !platform.IsBlank() {
		em.Platform = platform
	}

	var recordTypes = new(EnforcementsRecordTypesDataSourceModel)
	recordTypes.FromKeeper(enforcements)
	if !recordTypes.IsBlank() {
		em.RecordTypes = recordTypes
	}

	var twoFa = new(Enforcements2faDataSourceModel)
	twoFa.FromKeeper(enforcements)
	if !twoFa.IsBlank() {
		em.TwoFactorAuthentication = twoFa
	}

	var vault = new(EnforcementsVaultDataSourceModel)
	vault.FromKeeper(enforcements)
	if !vault.IsBlank() {
		em.Vault = vault
	}
}

func (em *EnforcementsDataSourceModel) ToKeeper(enforcements map[string]string) {
	if em.Account != nil {
		em.Account.ToKeeper(enforcements)
	}
	if em.AllowIpList != nil {
		em.AllowIpList.ToKeeper(enforcements)
	}
	if em.Sharing != nil {
		em.Sharing.ToKeeper(enforcements)
	}
	if em.KeeperFill != nil {
		em.KeeperFill.ToKeeper(enforcements)
	}
	if em.Login != nil {
		em.Login.ToKeeper(enforcements)
	}
	if em.Platform != nil {
		em.Platform.ToKeeper(enforcements)
	}
	if em.RecordTypes != nil {
		em.RecordTypes.ToKeeper(enforcements)
	}
	if em.TwoFactorAuthentication != nil {
		em.TwoFactorAuthentication.ToKeeper(enforcements)
	}
	if em.Vault != nil {
		em.Vault.ToKeeper(enforcements)
	}
}

func getIpWhitelistValue(property types.List, key string, enforcements map[string]string) {
	if !property.IsNull() {
		var strArrayValue []string
		for _, e := range property.Elements() {
			var s = e.String()
			if len(s) > 0 {
				strArrayValue = append(strArrayValue, s)
			}
		}
		if len(strArrayValue) > 0 {
			enforcements[key] = strings.Join(strArrayValue, ",")
		}
	}
}

func setIpWhitelistValue(property *types.List, key string, enforcements map[string]string) {
	var ok bool
	var strValue string
	var values []attr.Value
	if strValue, ok = enforcements[key]; ok {
		var svs = api.SliceSelect(strings.Split(strValue, ","), func(x string) string {
			return strings.TrimSpace(x)
		})
		values = api.SliceSelect(svs, func(x string) attr.Value {
			return types.StringValue(x)
		})
		if len(values) > 0 {
			*property, _ = types.ListValue(types.StringType, values)
		} else {
			*property = types.ListNull(types.StringType)
		}
	} else {
		*property = types.ListNull(types.StringType)
	}
}

func getInt64Value(property types.Int64, key string, enforcements map[string]string) {
	if !property.IsNull() {
		var intValue = property.ValueInt64()
		if intValue > 0 {
			enforcements[key] = strconv.Itoa(int(intValue))
		}
	}
}

func getBoolValue(property types.Bool, key string, enforcements map[string]string) {
	if !property.IsNull() {
		var boolValue = property.ValueBool()
		if boolValue {
			enforcements[key] = "true"
		}
	}
}

func getStringValue(property types.String, key string, enforcements map[string]string) {
	if !property.IsNull() {
		var strValue = property.ValueString()
		if len(strValue) > 0 {
			enforcements[key] = strValue
		}
	}
}

func setInt64Value(property *types.Int64, key string, enforcements map[string]string) {
	var ok bool
	var strValue string
	var intValue int
	var err error
	if strValue, ok = enforcements[key]; ok {
		if intValue, err = strconv.Atoi(strValue); err == nil {
			*property = types.Int64Value(int64(intValue))
		} else {
			*property = types.Int64Null()
		}
	} else {
		*property = types.Int64Null()
	}
}

func setBoolValue(property *types.Bool, key string, enforcements map[string]string) {
	var ok bool
	var strValue string
	var boolValue bool
	if strValue, ok = enforcements[key]; ok {
		boolValue = strings.ToLower(strValue) == "true" || strValue == "1"
		*property = types.BoolValue(boolValue)
	} else {
		*property = types.BoolNull()
	}
}

func setStringValue(property *types.String, key string, enforcements map[string]string) {
	var ok bool
	var strValue string
	if strValue, ok = enforcements[key]; ok {
		*property = types.StringValue(strValue)
	} else {
		*property = types.StringNull()
	}
}
