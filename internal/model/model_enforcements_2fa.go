package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Enforcements2faDataSourceModel struct {
	RestrictTwoFactorChannelSecurityKeys types.Bool   `tfsdk:"restrict_two_factor_channel_security_keys"`
	TwoFactorDurationDesktop             types.String `tfsdk:"two_factor_duration_desktop"`
	TwoFactorDurationMobile              types.String `tfsdk:"two_factor_duration_mobile"`
	TwoFactorDurationWeb                 types.String `tfsdk:"two_factor_duration_web"`
	RestrictTwoFactorChannelRsa          types.Bool   `tfsdk:"restrict_two_factor_channel_rsa"`
	RestrictTwoFactorChannelDuo          types.Bool   `tfsdk:"restrict_two_factor_channel_duo"`
	RestrictTwoFactorChannelDna          types.Bool   `tfsdk:"restrict_two_factor_channel_dna"`
	RestrictTwoFactorChannelGoogle       types.Bool   `tfsdk:"restrict_two_factor_channel_google"`
	RestrictTwoFactorChannelText         types.Bool   `tfsdk:"restrict_two_factor_channel_text"`
	RequireTwoFactor                     types.Bool   `tfsdk:"require_two_factor"`
}

func (e2fa *Enforcements2faDataSourceModel) ToKeeper(enforcements map[string]string) {
	getBoolValue(e2fa.RestrictTwoFactorChannelSecurityKeys, "restrict_two_factor_channel_security_keys", enforcements)
	getBoolValue(e2fa.RestrictTwoFactorChannelRsa, "restrict_two_factor_channel_rsa", enforcements)
	getBoolValue(e2fa.RestrictTwoFactorChannelDuo, "restrict_two_factor_channel_duo", enforcements)
	getBoolValue(e2fa.RestrictTwoFactorChannelDna, "restrict_two_factor_channel_dna", enforcements)
	getBoolValue(e2fa.RestrictTwoFactorChannelGoogle, "restrict_two_factor_channel_google", enforcements)
	getBoolValue(e2fa.RestrictTwoFactorChannelText, "restrict_two_factor_channel_text", enforcements)
	getBoolValue(e2fa.RequireTwoFactor, "require_two_factor", enforcements)
	getStringValue(e2fa.TwoFactorDurationDesktop, "two_factor_duration_desktop", enforcements)
	getStringValue(e2fa.TwoFactorDurationMobile, "two_factor_duration_mobile", enforcements)
	getStringValue(e2fa.TwoFactorDurationWeb, "two_factor_duration_web", enforcements)
}

func (e2fa *Enforcements2faDataSourceModel) FromKeeper(enforcements map[string]string) {
	setBoolValue(&e2fa.RestrictTwoFactorChannelSecurityKeys, "restrict_two_factor_channel_security_keys", enforcements)
	setBoolValue(&e2fa.RestrictTwoFactorChannelRsa, "restrict_two_factor_channel_rsa", enforcements)
	setBoolValue(&e2fa.RestrictTwoFactorChannelDuo, "restrict_two_factor_channel_duo", enforcements)
	setBoolValue(&e2fa.RestrictTwoFactorChannelDna, "restrict_two_factor_channel_dna", enforcements)
	setBoolValue(&e2fa.RestrictTwoFactorChannelGoogle, "restrict_two_factor_channel_google", enforcements)
	setBoolValue(&e2fa.RestrictTwoFactorChannelText, "restrict_two_factor_channel_text", enforcements)
	setBoolValue(&e2fa.RequireTwoFactor, "require_two_factor", enforcements)
	setStringValue(&e2fa.TwoFactorDurationDesktop, "two_factor_duration_desktop", enforcements)
	setStringValue(&e2fa.TwoFactorDurationMobile, "two_factor_duration_mobile", enforcements)
	setStringValue(&e2fa.TwoFactorDurationWeb, "two_factor_duration_web", enforcements)
}

func (e2fa *Enforcements2faDataSourceModel) IsBlank() bool {
	if !e2fa.RestrictTwoFactorChannelSecurityKeys.IsNull() {
		return false
	}
	if !e2fa.RestrictTwoFactorChannelRsa.IsNull() {
		return false
	}
	if !e2fa.RestrictTwoFactorChannelDuo.IsNull() {
		return false
	}
	if !e2fa.RestrictTwoFactorChannelDna.IsNull() {
		return false
	}
	if !e2fa.RestrictTwoFactorChannelGoogle.IsNull() {
		return false
	}
	if !e2fa.RestrictTwoFactorChannelText.IsNull() {
		return false
	}
	if !e2fa.RequireTwoFactor.IsNull() {
		return false
	}
	if !e2fa.TwoFactorDurationDesktop.IsNull() {
		return false
	}
	if !e2fa.TwoFactorDurationMobile.IsNull() {
		return false
	}
	if !e2fa.TwoFactorDurationWeb.IsNull() {
		return false
	}

	return true
}
