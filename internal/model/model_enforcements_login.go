package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsLoginDataSourceModel struct {
	AllowAlternatePasswords               types.Bool  `tfsdk:"allow_alternate_passwords"`
	RestrictWindowsFingerprint            types.Bool  `tfsdk:"restrict_windows_fingerprint"`
	RestrictAndroidFingerprint            types.Bool  `tfsdk:"restrict_android_fingerprint"`
	RestrictMacFingerprint                types.Bool  `tfsdk:"restrict_mac_fingerprint"`
	RestrictIosFingerprint                types.Bool  `tfsdk:"restrict_ios_fingerprint"`
	MasterPasswordExpiredAsOf             types.Int64 `tfsdk:"master_password_expired_as_of"`
	MasterPasswordMaximumDaysBeforeChange types.Int64 `tfsdk:"master_password_maximum_days_before_change"`
	MasterPasswordRestrictDaysBeforeReuse types.Int64 `tfsdk:"master_password_restrict_days_before_reuse"`
	MasterPasswordMinimumDigits           types.Int64 `tfsdk:"master_password_minimum_digits"`
	MasterPasswordMinimumLower            types.Int64 `tfsdk:"master_password_minimum_lower"`
	MasterPasswordMinimumUpper            types.Int64 `tfsdk:"master_password_minimum_upper"`
	MasterPasswordMinimumSpecial          types.Int64 `tfsdk:"master_password_minimum_special"`
	MasterPasswordMinimumLength           types.Int64 `tfsdk:"master_password_minimum_length"`
}

func (elm *EnforcementsLoginDataSourceModel) ToKeeper(enforcements map[string]string) {
	getInt64Value(elm.MasterPasswordExpiredAsOf, "master_password_expired_as_of", enforcements)
	getInt64Value(elm.MasterPasswordMaximumDaysBeforeChange, "master_password_maximum_days_before_change", enforcements)
	getInt64Value(elm.MasterPasswordRestrictDaysBeforeReuse, "master_password_restrict_days_before_reuse", enforcements)
	getInt64Value(elm.MasterPasswordMinimumDigits, "master_password_minimum_digits", enforcements)
	getInt64Value(elm.MasterPasswordMinimumLower, "master_password_minimum_lower", enforcements)
	getInt64Value(elm.MasterPasswordMinimumUpper, "master_password_minimum_upper", enforcements)
	getInt64Value(elm.MasterPasswordMinimumSpecial, "master_password_minimum_special", enforcements)
	getInt64Value(elm.MasterPasswordMinimumLength, "master_password_minimum_length", enforcements)
	getBoolValue(elm.AllowAlternatePasswords, "allow_alternate_passwords", enforcements)
	getBoolValue(elm.AllowAlternatePasswords, "allow_alternate_passwords", enforcements)
	getBoolValue(elm.RestrictWindowsFingerprint, "restrict_windows_fingerprint", enforcements)
	getBoolValue(elm.RestrictAndroidFingerprint, "restrict_android_fingerprint", enforcements)
	getBoolValue(elm.RestrictMacFingerprint, "restrict_mac_fingerprint", enforcements)
	getBoolValue(elm.RestrictIosFingerprint, "restrict_ios_fingerprint", enforcements)
}

func (elm *EnforcementsLoginDataSourceModel) FromKeeper(enforcements map[string]string) {
	setInt64Value(&elm.MasterPasswordExpiredAsOf, "master_password_expired_as_of", enforcements)
	setInt64Value(&elm.MasterPasswordMaximumDaysBeforeChange, "master_password_maximum_days_before_change", enforcements)
	setInt64Value(&elm.MasterPasswordRestrictDaysBeforeReuse, "master_password_restrict_days_before_reuse", enforcements)
	setInt64Value(&elm.MasterPasswordMinimumDigits, "master_password_minimum_digits", enforcements)
	setInt64Value(&elm.MasterPasswordMinimumLower, "master_password_minimum_lower", enforcements)
	setInt64Value(&elm.MasterPasswordMinimumUpper, "master_password_minimum_upper", enforcements)
	setInt64Value(&elm.MasterPasswordMinimumSpecial, "master_password_minimum_special", enforcements)
	setInt64Value(&elm.MasterPasswordMinimumLength, "master_password_minimum_length", enforcements)
	setBoolValue(&elm.AllowAlternatePasswords, "allow_alternate_passwords", enforcements)
	setBoolValue(&elm.AllowAlternatePasswords, "allow_alternate_passwords", enforcements)
	setBoolValue(&elm.RestrictWindowsFingerprint, "restrict_windows_fingerprint", enforcements)
	setBoolValue(&elm.RestrictAndroidFingerprint, "restrict_android_fingerprint", enforcements)
	setBoolValue(&elm.RestrictMacFingerprint, "restrict_mac_fingerprint", enforcements)
	setBoolValue(&elm.RestrictIosFingerprint, "restrict_ios_fingerprint", enforcements)
}

func (elm *EnforcementsLoginDataSourceModel) IsBlank() bool {
	if !elm.AllowAlternatePasswords.IsNull() {
		return false
	}
	if !elm.RestrictWindowsFingerprint.IsNull() {
		return false
	}
	if !elm.RestrictAndroidFingerprint.IsNull() {
		return false
	}
	if !elm.RestrictMacFingerprint.IsNull() {
		return false
	}
	if !elm.RestrictIosFingerprint.IsNull() {
		return false
	}
	if !elm.MasterPasswordExpiredAsOf.IsNull() {
		return false
	}
	if !elm.MasterPasswordMaximumDaysBeforeChange.IsNull() {
		return false
	}
	if !elm.MasterPasswordRestrictDaysBeforeReuse.IsNull() {
		return false
	}
	if !elm.MasterPasswordMinimumDigits.IsNull() {
		return false
	}
	if !elm.MasterPasswordMinimumLower.IsNull() {
		return false
	}
	if !elm.MasterPasswordMinimumUpper.IsNull() {
		return false
	}
	if !elm.MasterPasswordMinimumSpecial.IsNull() {
		return false
	}
	if !elm.MasterPasswordMinimumLength.IsNull() {
		return false
	}
	return true
}
