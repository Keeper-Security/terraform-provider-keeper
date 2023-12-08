package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsPlatformDataSourceModel struct {
	RestrictCommanderAccess	types.Bool	`tfsdk:"restrict_commander_access"`
	RestrictChatMobileAccess	types.Bool	`tfsdk:"restrict_chat_mobile_access"`
	RestrictChatDesktopAccess	types.Bool	`tfsdk:"restrict_chat_desktop_access"`
	RestrictDesktopMacAccess	types.Bool	`tfsdk:"restrict_desktop_mac_access"`
	RestrictDesktopWinAccess	types.Bool	`tfsdk:"restrict_desktop_win_access"`
	RestrictMobileWindowsPhoneAccess	types.Bool	`tfsdk:"restrict_mobile_windows_phone_access"`
	RestrictMobileAndroidAccess	types.Bool	`tfsdk:"restrict_mobile_android_access"`
	RestrictMobileIosAccess	types.Bool	`tfsdk:"restrict_mobile_ios_access"`
	RestrictDesktopAccess	types.Bool	`tfsdk:"restrict_desktop_access"`
	RestrictMobileAccess	types.Bool	`tfsdk:"restrict_mobile_access"`
	RestrictExtensionsAccess	types.Bool	`tfsdk:"restrict_extensions_access"`
	RestrictWebVaultAccess	types.Bool	`tfsdk:"restrict_web_vault_access"`
}

func (model *EnforcementsPlatformDataSourceModel) FromKeeper(enforcements map[string]string) {
	model.RestrictCommanderAccess = types.BoolValue(enforcements["restrict_commander_access"] != "")
	model.RestrictChatMobileAccess = types.BoolValue(enforcements["restrict_chat_mobile_access"] != "")
	model.RestrictChatDesktopAccess = types.BoolValue(enforcements["restrict_chat_desktop_access"] != "")
	model.RestrictDesktopMacAccess = types.BoolValue(enforcements["restrict_desktop_mac_access"] != "")
	model.RestrictDesktopWinAccess = types.BoolValue(enforcements["restrict_desktop_win_access"] != "")
	model.RestrictMobileWindowsPhoneAccess = types.BoolValue(enforcements["restrict_mobile_windows_phone_access"] != "")
	model.RestrictMobileAndroidAccess = types.BoolValue(enforcements["restrict_mobile_android_access"] != "")
	model.RestrictMobileIosAccess = types.BoolValue(enforcements["restrict_mobile_ios_access"] != "")
	model.RestrictDesktopAccess = types.BoolValue(enforcements["restrict_desktop_access"] != "")
	model.RestrictMobileAccess = types.BoolValue(enforcements["restrict_mobile_access"] != "")
	model.RestrictExtensionsAccess = types.BoolValue(enforcements["restrict_extensions_access"] != "")
	model.RestrictWebVaultAccess = types.BoolValue(enforcements["restrict_web_vault_access"] != "")
}