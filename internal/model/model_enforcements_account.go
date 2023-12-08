package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

type EnforcementsAccountDataSourceModel struct {
	RestrictAccountRecovery	types.Bool	`tfsdk:"restrict_account_recovery"`
	ResendEnterpriseInviteInXDays	types.Int64	`tfsdk:"resend_enterprise_invite_in_x_days"`
	RestrictIpAutoapproval	types.Bool	`tfsdk:"restrict_ip_autoapproval"`
	DisallowV2Clients	types.Bool	`tfsdk:"disallow_v2_clients"`
	DisableOnboarding	types.Bool	`tfsdk:"disable_onboarding"`
	RestrictPersonalLicense	types.Bool	`tfsdk:"restrict_personal_license"`
	LogoutTimerDesktop	types.Int64	`tfsdk:"logout_timer_desktop"`
	LogoutTimerMobile	types.Int64	`tfsdk:"logout_timer_mobile"`
	LogoutTimerWeb	types.Int64	`tfsdk:"logout_timer_web"`
	RestrictEmailChange	types.Bool	`tfsdk:"restrict_email_change"`
	SendInviteAtRegistration	types.Bool	`tfsdk:"send_invite_at_registration"`
	RestrictOfflineAccess	types.Bool	`tfsdk:"restrict_offline_access"`
	AutomaticBackupEveryXDays	types.Int64	`tfsdk:"automatic_backup_every_x_days"`
	RequireAccountRecoveryApproval	types.Bool	`tfsdk:"require_account_recovery_approval"`
	RequireDeviceApproval	types.Bool	`tfsdk:"require_device_approval"`
	RestrictPersistentLogin	types.Bool	`tfsdk:"restrict_persistent_login"`
	MaxSessionLoginTime	types.Int64	`tfsdk:"max_session_login_time"`
	MinimumPbkdf2Iterations	types.Int64	`tfsdk:"minimum_pbkdf2_iterations"`
	RequireSecurityKeyPin	types.Bool	`tfsdk:"require_security_key_pin"`
	RestrictImportSharedFolders	types.Bool	`tfsdk:"restrict_import_shared_folders"`
	AllowPamDiscovery	types.Bool	`tfsdk:"allow_pam_discovery"`
	AllowPamRotation	types.Bool	`tfsdk:"allow_pam_rotation"`
	MaximumRecordSize	types.Int64	`tfsdk:"maximum_record_size"`
	RequireSelfDestruct	types.Bool	`tfsdk:"require_self_destruct"`
	StayLoggedInDefault	types.Bool	`tfsdk:"stay_logged_in_default"`
}

func (model *EnforcementsAccountDataSourceModel) FromKeeper(enforcements map[string]string) {
	var intValue int64
	var err error
	intValue, err = strconv.ParseInt(enforcements["resend_enterprise_invite_in_x_days"], 0, 64)
	if err != nil {
		model.ResendEnterpriseInviteInXDays = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["logout_timer_desktop"], 0, 64)
	if err != nil {
		model.LogoutTimerDesktop = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["logout_timer_mobile"], 0, 64)
	if err != nil {
		model.LogoutTimerMobile = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["logout_timer_web"], 0, 64)
	if err != nil {
		model.LogoutTimerWeb = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["automatic_backup_every_x_days"], 0, 64)
	if err != nil {
		model.AutomaticBackupEveryXDays = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["max_session_login_time"], 0, 64)
	if err != nil {
		model.MaxSessionLoginTime = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["minimum_pbkdf2_iterations"], 0, 64)
	if err != nil {
		model.MinimumPbkdf2Iterations = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["maximum_record_size"], 0, 64)
	if err != nil {
		model.MaximumRecordSize = types.Int64Value(intValue)
	}
	model.RestrictAccountRecovery = types.BoolValue(enforcements["restrict_account_recovery"] != "")
	model.RestrictIpAutoapproval = types.BoolValue(enforcements["restrict_ip_autoapproval"] != "")
	model.DisallowV2Clients = types.BoolValue(enforcements["disallow_v2_clients"] != "")
	model.DisableOnboarding = types.BoolValue(enforcements["disable_onboarding"] != "")
	model.RestrictPersonalLicense = types.BoolValue(enforcements["restrict_personal_license"] != "")
	model.RestrictEmailChange = types.BoolValue(enforcements["restrict_email_change"] != "")
	model.SendInviteAtRegistration = types.BoolValue(enforcements["send_invite_at_registration"] != "")
	model.RestrictOfflineAccess = types.BoolValue(enforcements["restrict_offline_access"] != "")
	model.RequireAccountRecoveryApproval = types.BoolValue(enforcements["require_account_recovery_approval"] != "")
	model.RequireDeviceApproval = types.BoolValue(enforcements["require_device_approval"] != "")
	model.RestrictPersistentLogin = types.BoolValue(enforcements["restrict_persistent_login"] != "")
	model.RequireSecurityKeyPin = types.BoolValue(enforcements["require_security_key_pin"] != "")
	model.RestrictImportSharedFolders = types.BoolValue(enforcements["restrict_import_shared_folders"] != "")
	model.AllowPamDiscovery = types.BoolValue(enforcements["allow_pam_discovery"] != "")
	model.AllowPamRotation = types.BoolValue(enforcements["allow_pam_rotation"] != "")
	model.RequireSelfDestruct = types.BoolValue(enforcements["require_self_destruct"] != "")
	model.StayLoggedInDefault = types.BoolValue(enforcements["stay_logged_in_default"] != "")
}