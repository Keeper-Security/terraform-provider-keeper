package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Enforcements2faDataSourceModel struct {
	RestrictTwoFactorChannelSecurityKeys	types.Bool	`tfsdk:"restrict_two_factor_channel_security_keys"`
	TwoFactorDurationDesktop	types.String	`tfsdk:"two_factor_duration_desktop"`
	TwoFactorDurationMobile	types.String	`tfsdk:"two_factor_duration_mobile"`
	TwoFactorDurationWeb	types.String	`tfsdk:"two_factor_duration_web"`
	RestrictTwoFactorChannelRsa	types.Bool	`tfsdk:"restrict_two_factor_channel_rsa"`
	RestrictTwoFactorChannelDuo	types.Bool	`tfsdk:"restrict_two_factor_channel_duo"`
	RestrictTwoFactorChannelDna	types.Bool	`tfsdk:"restrict_two_factor_channel_dna"`
	RestrictTwoFactorChannelGoogle	types.Bool	`tfsdk:"restrict_two_factor_channel_google"`
	RestrictTwoFactorChannelText	types.Bool	`tfsdk:"restrict_two_factor_channel_text"`
	RequireTwoFactor	types.Bool	`tfsdk:"require_two_factor"`
}

func (model *Enforcements2faDataSourceModel) FromKeeper(enforcements map[string]string) {
	model.RestrictTwoFactorChannelSecurityKeys = types.BoolValue(enforcements["restrict_two_factor_channel_security_keys"] != "")
	model.RestrictTwoFactorChannelRsa = types.BoolValue(enforcements["restrict_two_factor_channel_rsa"] != "")
	model.RestrictTwoFactorChannelDuo = types.BoolValue(enforcements["restrict_two_factor_channel_duo"] != "")
	model.RestrictTwoFactorChannelDna = types.BoolValue(enforcements["restrict_two_factor_channel_dna"] != "")
	model.RestrictTwoFactorChannelGoogle = types.BoolValue(enforcements["restrict_two_factor_channel_google"] != "")
	model.RestrictTwoFactorChannelText = types.BoolValue(enforcements["restrict_two_factor_channel_text"] != "")
	model.RequireTwoFactor = types.BoolValue(enforcements["require_two_factor"] != "")
	model.TwoFactorDurationDesktop = types.StringValue(enforcements["two_factor_duration_desktop"])
	model.TwoFactorDurationMobile = types.StringValue(enforcements["two_factor_duration_mobile"])
	model.TwoFactorDurationWeb = types.StringValue(enforcements["two_factor_duration_web"])
}