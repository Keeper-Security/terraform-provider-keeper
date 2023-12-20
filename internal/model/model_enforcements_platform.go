package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsPlatformDataSourceModel struct {
	RestrictCommanderAccess          types.Bool `tfsdk:"restrict_commander_access"`
	RestrictChatMobileAccess         types.Bool `tfsdk:"restrict_chat_mobile_access"`
	RestrictChatDesktopAccess        types.Bool `tfsdk:"restrict_chat_desktop_access"`
	RestrictDesktopMacAccess         types.Bool `tfsdk:"restrict_desktop_mac_access"`
	RestrictDesktopWinAccess         types.Bool `tfsdk:"restrict_desktop_win_access"`
	RestrictMobileWindowsPhoneAccess types.Bool `tfsdk:"restrict_mobile_windows_phone_access"`
	RestrictMobileAndroidAccess      types.Bool `tfsdk:"restrict_mobile_android_access"`
	RestrictMobileIosAccess          types.Bool `tfsdk:"restrict_mobile_ios_access"`
	RestrictDesktopAccess            types.Bool `tfsdk:"restrict_desktop_access"`
	RestrictMobileAccess             types.Bool `tfsdk:"restrict_mobile_access"`
	RestrictExtensionsAccess         types.Bool `tfsdk:"restrict_extensions_access"`
	RestrictWebVaultAccess           types.Bool `tfsdk:"restrict_web_vault_access"`
}

func (epm *EnforcementsPlatformDataSourceModel) ToKeeper(enforcements map[string]string) {
	getBoolValue(epm.RestrictCommanderAccess, "restrict_commander_access", enforcements)
	getBoolValue(epm.RestrictChatMobileAccess, "restrict_chat_mobile_access", enforcements)
	getBoolValue(epm.RestrictChatDesktopAccess, "restrict_chat_desktop_access", enforcements)
	getBoolValue(epm.RestrictDesktopMacAccess, "restrict_desktop_mac_access", enforcements)
	getBoolValue(epm.RestrictDesktopWinAccess, "restrict_desktop_win_access", enforcements)
	getBoolValue(epm.RestrictMobileWindowsPhoneAccess, "restrict_mobile_windows_phone_access", enforcements)
	getBoolValue(epm.RestrictMobileAndroidAccess, "restrict_mobile_android_access", enforcements)
	getBoolValue(epm.RestrictMobileIosAccess, "restrict_mobile_ios_access", enforcements)
	getBoolValue(epm.RestrictDesktopAccess, "restrict_desktop_access", enforcements)
	getBoolValue(epm.RestrictMobileAccess, "restrict_mobile_access", enforcements)
	getBoolValue(epm.RestrictExtensionsAccess, "restrict_extensions_access", enforcements)
	getBoolValue(epm.RestrictWebVaultAccess, "restrict_web_vault_access", enforcements)
}

func (epm *EnforcementsPlatformDataSourceModel) FromKeeper(enforcements map[string]string) {
	setBoolValue(&epm.RestrictCommanderAccess, "restrict_commander_access", enforcements)
	setBoolValue(&epm.RestrictChatMobileAccess, "restrict_chat_mobile_access", enforcements)
	setBoolValue(&epm.RestrictChatDesktopAccess, "restrict_chat_desktop_access", enforcements)
	setBoolValue(&epm.RestrictDesktopMacAccess, "restrict_desktop_mac_access", enforcements)
	setBoolValue(&epm.RestrictDesktopWinAccess, "restrict_desktop_win_access", enforcements)
	setBoolValue(&epm.RestrictMobileWindowsPhoneAccess, "restrict_mobile_windows_phone_access", enforcements)
	setBoolValue(&epm.RestrictMobileAndroidAccess, "restrict_mobile_android_access", enforcements)
	setBoolValue(&epm.RestrictMobileIosAccess, "restrict_mobile_ios_access", enforcements)
	setBoolValue(&epm.RestrictDesktopAccess, "restrict_desktop_access", enforcements)
	setBoolValue(&epm.RestrictMobileAccess, "restrict_mobile_access", enforcements)
	setBoolValue(&epm.RestrictExtensionsAccess, "restrict_extensions_access", enforcements)
	setBoolValue(&epm.RestrictWebVaultAccess, "restrict_web_vault_access", enforcements)
}

func (epm *EnforcementsPlatformDataSourceModel) IsBlank() bool {
	if !epm.RestrictCommanderAccess.IsNull() {
		return false
	}
	if !epm.RestrictChatMobileAccess.IsNull() {
		return false
	}
	if !epm.RestrictChatDesktopAccess.IsNull() {
		return false
	}
	if !epm.RestrictDesktopMacAccess.IsNull() {
		return false
	}
	if !epm.RestrictDesktopWinAccess.IsNull() {
		return false
	}
	if !epm.RestrictMobileWindowsPhoneAccess.IsNull() {
		return false
	}
	if !epm.RestrictMobileAndroidAccess.IsNull() {
		return false
	}
	if !epm.RestrictMobileIosAccess.IsNull() {
		return false
	}
	if !epm.RestrictDesktopAccess.IsNull() {
		return false
	}
	if !epm.RestrictMobileAccess.IsNull() {
		return false
	}
	if !epm.RestrictExtensionsAccess.IsNull() {
		return false
	}
	if !epm.RestrictWebVaultAccess.IsNull() {
		return false
	}
	return true
}
