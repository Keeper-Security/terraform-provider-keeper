package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsKeeperFillDataSourceModel struct {
	KeeperFillAutoSuggest      types.String `tfsdk:"keeper_fill_auto_suggest"`
	RestrictHttpFillWarning    types.Bool   `tfsdk:"restrict_http_fill_warning"`
	RestrictPromptToDisable    types.Bool   `tfsdk:"restrict_prompt_to_disable"`
	KeeperFillMatchOnSubdomain types.String `tfsdk:"keeper_fill_match_on_subdomain"`
	KeeperFillAutoSubmit       types.String `tfsdk:"keeper_fill_auto_submit"`
	KeeperFillAutoFill         types.String `tfsdk:"keeper_fill_auto_fill"`
	KeeperFillHoverLocks       types.String `tfsdk:"keeper_fill_hover_locks"`
	MasterPasswordReentry      types.String `tfsdk:"master_password_reentry"`
	RestrictAutoFill           types.Bool   `tfsdk:"restrict_auto_fill"`
	RestrictPromptToChange     types.Bool   `tfsdk:"restrict_prompt_to_change"`
	RestrictPromptToSave       types.Bool   `tfsdk:"restrict_prompt_to_save"`
	RestrictAutoSubmit         types.Bool   `tfsdk:"restrict_auto_submit"`
	RestrictPromptToFill       types.Bool   `tfsdk:"restrict_prompt_to_fill"`
	RestrictPromptToLogin      types.Bool   `tfsdk:"restrict_prompt_to_login"`
	RestrictHoverLocks         types.Bool   `tfsdk:"restrict_hover_locks"`
	RestrictDomainCreate       types.String `tfsdk:"restrict_domain_create"`
	RestrictDomainAccess       types.String `tfsdk:"restrict_domain_access"`
}

func (ekfm *EnforcementsKeeperFillDataSourceModel) ToKeeper(enforcements map[string]string) {
	getBoolValue(ekfm.RestrictHttpFillWarning, "restrict_http_fill_warning", enforcements)
	getBoolValue(ekfm.RestrictPromptToDisable, "restrict_prompt_to_disable", enforcements)
	getBoolValue(ekfm.RestrictAutoFill, "restrict_auto_fill", enforcements)
	getBoolValue(ekfm.RestrictPromptToChange, "restrict_prompt_to_change", enforcements)
	getBoolValue(ekfm.RestrictPromptToSave, "restrict_prompt_to_save", enforcements)
	getBoolValue(ekfm.RestrictAutoSubmit, "restrict_auto_submit", enforcements)
	getBoolValue(ekfm.RestrictPromptToFill, "restrict_prompt_to_fill", enforcements)
	getBoolValue(ekfm.RestrictPromptToLogin, "restrict_prompt_to_login", enforcements)
	getBoolValue(ekfm.RestrictHoverLocks, "restrict_hover_locks", enforcements)
	getStringValue(ekfm.KeeperFillAutoSuggest, "keeper_fill_auto_suggest", enforcements)
	getStringValue(ekfm.KeeperFillMatchOnSubdomain, "keeper_fill_match_on_subdomain", enforcements)
	getStringValue(ekfm.KeeperFillAutoSubmit, "keeper_fill_auto_submit", enforcements)
	getStringValue(ekfm.KeeperFillAutoFill, "keeper_fill_auto_fill", enforcements)
	getStringValue(ekfm.KeeperFillHoverLocks, "keeper_fill_hover_locks", enforcements)
	getStringValue(ekfm.MasterPasswordReentry, "master_password_reentry", enforcements)
	getStringValue(ekfm.RestrictDomainCreate, "restrict_domain_create", enforcements)
	getStringValue(ekfm.RestrictDomainAccess, "restrict_domain_access", enforcements)
}

func (ekfm *EnforcementsKeeperFillDataSourceModel) FromKeeper(enforcements map[string]string) {
	setBoolValue(&ekfm.RestrictHttpFillWarning, "restrict_http_fill_warning", enforcements)
	setBoolValue(&ekfm.RestrictPromptToDisable, "restrict_prompt_to_disable", enforcements)
	setBoolValue(&ekfm.RestrictAutoFill, "restrict_auto_fill", enforcements)
	setBoolValue(&ekfm.RestrictPromptToChange, "restrict_prompt_to_change", enforcements)
	setBoolValue(&ekfm.RestrictPromptToSave, "restrict_prompt_to_save", enforcements)
	setBoolValue(&ekfm.RestrictAutoSubmit, "restrict_auto_submit", enforcements)
	setBoolValue(&ekfm.RestrictPromptToFill, "restrict_prompt_to_fill", enforcements)
	setBoolValue(&ekfm.RestrictPromptToLogin, "restrict_prompt_to_login", enforcements)
	setBoolValue(&ekfm.RestrictHoverLocks, "restrict_hover_locks", enforcements)
	setStringValue(&ekfm.KeeperFillAutoSuggest, "keeper_fill_auto_suggest", enforcements)
	setStringValue(&ekfm.KeeperFillMatchOnSubdomain, "keeper_fill_match_on_subdomain", enforcements)
	setStringValue(&ekfm.KeeperFillAutoSubmit, "keeper_fill_auto_submit", enforcements)
	setStringValue(&ekfm.KeeperFillAutoFill, "keeper_fill_auto_fill", enforcements)
	setStringValue(&ekfm.KeeperFillHoverLocks, "keeper_fill_hover_locks", enforcements)
	setStringValue(&ekfm.MasterPasswordReentry, "master_password_reentry", enforcements)
	setStringValue(&ekfm.RestrictDomainCreate, "restrict_domain_create", enforcements)
	setStringValue(&ekfm.RestrictDomainAccess, "restrict_domain_access", enforcements)
}

func (ekfm *EnforcementsKeeperFillDataSourceModel) IsBlank() bool {
	if !ekfm.KeeperFillAutoSuggest.IsNull() {
		return false
	}
	if !ekfm.RestrictHttpFillWarning.IsNull() {
		return false
	}
	if !ekfm.RestrictPromptToDisable.IsNull() {
		return false
	}
	if !ekfm.KeeperFillMatchOnSubdomain.IsNull() {
		return false
	}
	if !ekfm.KeeperFillAutoSubmit.IsNull() {
		return false
	}
	if !ekfm.KeeperFillAutoFill.IsNull() {
		return false
	}
	if !ekfm.KeeperFillHoverLocks.IsNull() {
		return false
	}
	if !ekfm.MasterPasswordReentry.IsNull() {
		return false
	}
	if !ekfm.RestrictAutoFill.IsNull() {
		return false
	}
	if !ekfm.RestrictPromptToChange.IsNull() {
		return false
	}
	if !ekfm.RestrictPromptToSave.IsNull() {
		return false
	}
	if !ekfm.RestrictAutoSubmit.IsNull() {
		return false
	}
	if !ekfm.RestrictPromptToFill.IsNull() {
		return false
	}
	if !ekfm.RestrictPromptToLogin.IsNull() {
		return false
	}
	if !ekfm.RestrictHoverLocks.IsNull() {
		return false
	}
	if !ekfm.RestrictDomainCreate.IsNull() {
		return false
	}
	if !ekfm.RestrictDomainAccess.IsNull() {
		return false
	}
	return true
}
