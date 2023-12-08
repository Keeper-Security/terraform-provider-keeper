package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

type EnforcementsLoginDataSourceModel struct {
	AllowAlternatePasswords	types.Bool	`tfsdk:"allow_alternate_passwords"`
	RestrictWindowsFingerprint	types.Bool	`tfsdk:"restrict_windows_fingerprint"`
	RestrictAndroidFingerprint	types.Bool	`tfsdk:"restrict_android_fingerprint"`
	RestrictMacFingerprint	types.Bool	`tfsdk:"restrict_mac_fingerprint"`
	RestrictIosFingerprint	types.Bool	`tfsdk:"restrict_ios_fingerprint"`
	MasterPasswordExpiredAsOf	types.Int64	`tfsdk:"master_password_expired_as_of"`
	MasterPasswordMaximumDaysBeforeChange	types.Int64	`tfsdk:"master_password_maximum_days_before_change"`
	MasterPasswordRestrictDaysBeforeReuse	types.Int64	`tfsdk:"master_password_restrict_days_before_reuse"`
	MasterPasswordMinimumDigits	types.Int64	`tfsdk:"master_password_minimum_digits"`
	MasterPasswordMinimumLower	types.Int64	`tfsdk:"master_password_minimum_lower"`
	MasterPasswordMinimumUpper	types.Int64	`tfsdk:"master_password_minimum_upper"`
	MasterPasswordMinimumSpecial	types.Int64	`tfsdk:"master_password_minimum_special"`
	MasterPasswordMinimumLength	types.Int64	`tfsdk:"master_password_minimum_length"`
}

func (model *EnforcementsLoginDataSourceModel) FromKeeper(enforcements map[string]string) {
	var intValue int64
	var err error
	intValue, err = strconv.ParseInt(enforcements["master_password_expired_as_of"], 0, 64)
	if err != nil {
		model.MasterPasswordExpiredAsOf = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["master_password_maximum_days_before_change"], 0, 64)
	if err != nil {
		model.MasterPasswordMaximumDaysBeforeChange = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["master_password_restrict_days_before_reuse"], 0, 64)
	if err != nil {
		model.MasterPasswordRestrictDaysBeforeReuse = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["master_password_minimum_digits"], 0, 64)
	if err != nil {
		model.MasterPasswordMinimumDigits = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["master_password_minimum_lower"], 0, 64)
	if err != nil {
		model.MasterPasswordMinimumLower = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["master_password_minimum_upper"], 0, 64)
	if err != nil {
		model.MasterPasswordMinimumUpper = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["master_password_minimum_special"], 0, 64)
	if err != nil {
		model.MasterPasswordMinimumSpecial = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["master_password_minimum_length"], 0, 64)
	if err != nil {
		model.MasterPasswordMinimumLength = types.Int64Value(intValue)
	}
	model.AllowAlternatePasswords = types.BoolValue(enforcements["allow_alternate_passwords"] != "")
	model.RestrictWindowsFingerprint = types.BoolValue(enforcements["restrict_windows_fingerprint"] != "")
	model.RestrictAndroidFingerprint = types.BoolValue(enforcements["restrict_android_fingerprint"] != "")
	model.RestrictMacFingerprint = types.BoolValue(enforcements["restrict_mac_fingerprint"] != "")
	model.RestrictIosFingerprint = types.BoolValue(enforcements["restrict_ios_fingerprint"] != "")
}