package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsKeeperFillDataSourceModel struct {
	KeeperFillAutoSuggest	types.String	`tfsdk:"keeper_fill_auto_suggest"`
	RestrictHttpFillWarning	types.Bool	`tfsdk:"restrict_http_fill_warning"`
	RestrictPromptToDisable	types.Bool	`tfsdk:"restrict_prompt_to_disable"`
	KeeperFillMatchOnSubdomain	types.String	`tfsdk:"keeper_fill_match_on_subdomain"`
	KeeperFillAutoSubmit	types.String	`tfsdk:"keeper_fill_auto_submit"`
	KeeperFillAutoFill	types.String	`tfsdk:"keeper_fill_auto_fill"`
	KeeperFillHoverLocks	types.String	`tfsdk:"keeper_fill_hover_locks"`
	MasterPasswordReentry	types.String	`tfsdk:"master_password_reentry"`
	RestrictAutoFill	types.Bool	`tfsdk:"restrict_auto_fill"`
	RestrictPromptToChange	types.Bool	`tfsdk:"restrict_prompt_to_change"`
	RestrictPromptToSave	types.Bool	`tfsdk:"restrict_prompt_to_save"`
	RestrictAutoSubmit	types.Bool	`tfsdk:"restrict_auto_submit"`
	RestrictPromptToFill	types.Bool	`tfsdk:"restrict_prompt_to_fill"`
	RestrictPromptToLogin	types.Bool	`tfsdk:"restrict_prompt_to_login"`
	RestrictHoverLocks	types.Bool	`tfsdk:"restrict_hover_locks"`
	RestrictDomainCreate	types.String	`tfsdk:"restrict_domain_create"`
	RestrictDomainAccess	types.String	`tfsdk:"restrict_domain_access"`
}

func (model *EnforcementsKeeperFillDataSourceModel) FromKeeper(enforcements map[string]string) {
	model.RestrictHttpFillWarning = types.BoolValue(enforcements["restrict_http_fill_warning"] != "")
	model.RestrictPromptToDisable = types.BoolValue(enforcements["restrict_prompt_to_disable"] != "")
	model.RestrictAutoFill = types.BoolValue(enforcements["restrict_auto_fill"] != "")
	model.RestrictPromptToChange = types.BoolValue(enforcements["restrict_prompt_to_change"] != "")
	model.RestrictPromptToSave = types.BoolValue(enforcements["restrict_prompt_to_save"] != "")
	model.RestrictAutoSubmit = types.BoolValue(enforcements["restrict_auto_submit"] != "")
	model.RestrictPromptToFill = types.BoolValue(enforcements["restrict_prompt_to_fill"] != "")
	model.RestrictPromptToLogin = types.BoolValue(enforcements["restrict_prompt_to_login"] != "")
	model.RestrictHoverLocks = types.BoolValue(enforcements["restrict_hover_locks"] != "")
	model.KeeperFillAutoSuggest = types.StringValue(enforcements["keeper_fill_auto_suggest"])
	model.KeeperFillMatchOnSubdomain = types.StringValue(enforcements["keeper_fill_match_on_subdomain"])
	model.KeeperFillAutoSubmit = types.StringValue(enforcements["keeper_fill_auto_submit"])
	model.KeeperFillAutoFill = types.StringValue(enforcements["keeper_fill_auto_fill"])
	model.KeeperFillHoverLocks = types.StringValue(enforcements["keeper_fill_hover_locks"])
	model.MasterPasswordReentry = types.StringValue(enforcements["master_password_reentry"])
	model.RestrictDomainCreate = types.StringValue(enforcements["restrict_domain_create"])
	model.RestrictDomainAccess = types.StringValue(enforcements["restrict_domain_access"])
}