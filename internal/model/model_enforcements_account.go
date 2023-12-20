package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsAccountDataSourceModel struct {
	RestrictAccountRecovery        types.Bool  `tfsdk:"restrict_account_recovery"`
	ResendEnterpriseInviteInXDays  types.Int64 `tfsdk:"resend_enterprise_invite_in_x_days"`
	RestrictIpAutoApproval         types.Bool  `tfsdk:"restrict_ip_autoapproval"`
	DisallowV2Clients              types.Bool  `tfsdk:"disallow_v2_clients"`
	DisableOnboarding              types.Bool  `tfsdk:"disable_onboarding"`
	RestrictPersonalLicense        types.Bool  `tfsdk:"restrict_personal_license"`
	LogoutTimerDesktop             types.Int64 `tfsdk:"logout_timer_desktop"`
	LogoutTimerMobile              types.Int64 `tfsdk:"logout_timer_mobile"`
	LogoutTimerWeb                 types.Int64 `tfsdk:"logout_timer_web"`
	RestrictEmailChange            types.Bool  `tfsdk:"restrict_email_change"`
	SendInviteAtRegistration       types.Bool  `tfsdk:"send_invite_at_registration"`
	RestrictOfflineAccess          types.Bool  `tfsdk:"restrict_offline_access"`
	RequireAccountRecoveryApproval types.Bool  `tfsdk:"require_account_recovery_approval"`
	RequireDeviceApproval          types.Bool  `tfsdk:"require_device_approval"`
	RestrictPersistentLogin        types.Bool  `tfsdk:"restrict_persistent_login"`
	MaxSessionLoginTime            types.Int64 `tfsdk:"max_session_login_time"`
	MinimumPbkdf2Iterations        types.Int64 `tfsdk:"minimum_pbkdf2_iterations"`
	RequireSecurityKeyPin          types.Bool  `tfsdk:"require_security_key_pin"`
	RestrictImportSharedFolders    types.Bool  `tfsdk:"restrict_import_shared_folders"`
	AllowPamDiscovery              types.Bool  `tfsdk:"allow_pam_discovery"`
	AllowPamRotation               types.Bool  `tfsdk:"allow_pam_rotation"`
	MaximumRecordSize              types.Int64 `tfsdk:"maximum_record_size"`
	RequireSelfDestruct            types.Bool  `tfsdk:"require_self_destruct"`
	StayLoggedInDefault            types.Bool  `tfsdk:"stay_logged_in_default"`
	DisableSetupTour               types.Bool  `tfsdk:"disable_setup_tour"`
	AllowSecretsManager            types.Bool  `tfsdk:"allow_secrets_manager"`
}

func (eam *EnforcementsAccountDataSourceModel) ToKeeper(enforcements map[string]string) {
	getInt64Value(eam.ResendEnterpriseInviteInXDays, "resend_enterprise_invite_in_x_days", enforcements)
	getInt64Value(eam.LogoutTimerDesktop, "logout_timer_desktop", enforcements)
	getInt64Value(eam.LogoutTimerMobile, "logout_timer_mobile", enforcements)
	getInt64Value(eam.LogoutTimerWeb, "logout_timer_web", enforcements)
	getInt64Value(eam.MaxSessionLoginTime, "max_session_login_time", enforcements)
	getInt64Value(eam.MinimumPbkdf2Iterations, "minimum_pbkdf2_iterations", enforcements)
	getInt64Value(eam.MaximumRecordSize, "maximum_record_size", enforcements)
	getBoolValue(eam.RestrictAccountRecovery, "restrict_account_recovery", enforcements)
	getBoolValue(eam.RestrictIpAutoApproval, "restrict_ip_autoapproval", enforcements)
	getBoolValue(eam.DisallowV2Clients, "disallow_v2_clients", enforcements)
	getBoolValue(eam.DisableOnboarding, "disable_onboarding", enforcements)
	getBoolValue(eam.RestrictPersonalLicense, "restrict_personal_license", enforcements)
	getBoolValue(eam.RestrictEmailChange, "restrict_email_change", enforcements)
	getBoolValue(eam.SendInviteAtRegistration, "send_invite_at_registration", enforcements)
	getBoolValue(eam.RestrictOfflineAccess, "restrict_offline_access", enforcements)
	getBoolValue(eam.RequireAccountRecoveryApproval, "require_account_recovery_approval", enforcements)
	getBoolValue(eam.RequireDeviceApproval, "require_device_approval", enforcements)
	getBoolValue(eam.RestrictPersistentLogin, "restrict_persistent_login", enforcements)
	getBoolValue(eam.RequireSecurityKeyPin, "require_security_key_pin", enforcements)
	getBoolValue(eam.RestrictImportSharedFolders, "restrict_import_shared_folders", enforcements)
	getBoolValue(eam.AllowPamDiscovery, "allow_pam_discovery", enforcements)
	getBoolValue(eam.AllowPamRotation, "allow_pam_rotation", enforcements)
	getBoolValue(eam.RequireSelfDestruct, "require_self_destruct", enforcements)
	getBoolValue(eam.StayLoggedInDefault, "stay_logged_in_default", enforcements)
	getBoolValue(eam.DisableSetupTour, "disable_setup_tour", enforcements)
	getBoolValue(eam.AllowSecretsManager, "allow_secrets_manager", enforcements)
}

func (eam *EnforcementsAccountDataSourceModel) FromKeeper(enforcements map[string]string) {
	setInt64Value(&eam.ResendEnterpriseInviteInXDays, "resend_enterprise_invite_in_x_days", enforcements)
	setInt64Value(&eam.LogoutTimerDesktop, "logout_timer_desktop", enforcements)
	setInt64Value(&eam.LogoutTimerMobile, "logout_timer_mobile", enforcements)
	setInt64Value(&eam.LogoutTimerWeb, "logout_timer_web", enforcements)
	setInt64Value(&eam.MaxSessionLoginTime, "max_session_login_time", enforcements)
	setInt64Value(&eam.MinimumPbkdf2Iterations, "minimum_pbkdf2_iterations", enforcements)
	setInt64Value(&eam.MaximumRecordSize, "maximum_record_size", enforcements)
	setBoolValue(&eam.RestrictAccountRecovery, "restrict_account_recovery", enforcements)
	setBoolValue(&eam.RestrictIpAutoApproval, "restrict_ip_autoapproval", enforcements)
	setBoolValue(&eam.DisallowV2Clients, "disallow_v2_clients", enforcements)
	setBoolValue(&eam.DisableOnboarding, "disable_onboarding", enforcements)
	setBoolValue(&eam.RestrictPersonalLicense, "restrict_personal_license", enforcements)
	setBoolValue(&eam.RestrictEmailChange, "restrict_email_change", enforcements)
	setBoolValue(&eam.SendInviteAtRegistration, "send_invite_at_registration", enforcements)
	setBoolValue(&eam.RestrictOfflineAccess, "restrict_offline_access", enforcements)
	setBoolValue(&eam.RequireAccountRecoveryApproval, "require_account_recovery_approval", enforcements)
	setBoolValue(&eam.RequireDeviceApproval, "require_device_approval", enforcements)
	setBoolValue(&eam.RestrictPersistentLogin, "restrict_persistent_login", enforcements)
	setBoolValue(&eam.RequireSecurityKeyPin, "require_security_key_pin", enforcements)
	setBoolValue(&eam.RestrictImportSharedFolders, "restrict_import_shared_folders", enforcements)
	setBoolValue(&eam.AllowPamDiscovery, "allow_pam_discovery", enforcements)
	setBoolValue(&eam.AllowPamRotation, "allow_pam_rotation", enforcements)
	setBoolValue(&eam.RequireSelfDestruct, "require_self_destruct", enforcements)
	setBoolValue(&eam.StayLoggedInDefault, "stay_logged_in_default", enforcements)
	setBoolValue(&eam.DisableSetupTour, "disable_setup_tour", enforcements)
	setBoolValue(&eam.AllowSecretsManager, "allow_secrets_manager", enforcements)
}

func (eam *EnforcementsAccountDataSourceModel) IsBlank() bool {
	if !eam.RestrictAccountRecovery.IsNull() {
		return false
	}
	if !eam.ResendEnterpriseInviteInXDays.IsNull() {
		return false
	}
	if !eam.RestrictIpAutoApproval.IsNull() {
		return false
	}
	if !eam.DisallowV2Clients.IsNull() {
		return false
	}
	if !eam.DisableOnboarding.IsNull() {
		return false
	}
	if !eam.RestrictPersonalLicense.IsNull() {
		return false
	}
	if !eam.LogoutTimerDesktop.IsNull() {
		return false
	}
	if !eam.LogoutTimerMobile.IsNull() {
		return false
	}
	if !eam.LogoutTimerWeb.IsNull() {
		return false
	}
	if !eam.RestrictEmailChange.IsNull() {
		return false
	}
	if !eam.SendInviteAtRegistration.IsNull() {
		return false
	}
	if !eam.RestrictOfflineAccess.IsNull() {
		return false
	}
	if !eam.RequireAccountRecoveryApproval.IsNull() {
		return false
	}
	if !eam.RequireDeviceApproval.IsNull() {
		return false
	}
	if !eam.RestrictPersistentLogin.IsNull() {
		return false
	}
	if !eam.MaxSessionLoginTime.IsNull() {
		return false
	}
	if !eam.MinimumPbkdf2Iterations.IsNull() {
		return false
	}
	if !eam.RequireSecurityKeyPin.IsNull() {
		return false
	}
	if !eam.RestrictImportSharedFolders.IsNull() {
		return false
	}
	if !eam.AllowPamDiscovery.IsNull() {
		return false
	}
	if !eam.AllowPamRotation.IsNull() {
		return false
	}
	if !eam.MaximumRecordSize.IsNull() {
		return false
	}
	if !eam.RequireSelfDestruct.IsNull() {
		return false
	}
	if !eam.StayLoggedInDefault.IsNull() {
		return false
	}
	if !eam.DisableSetupTour.IsNull() {
		return false
	}
	if !eam.AllowSecretsManager.IsNull() {
		return false
	}
	return true
}
