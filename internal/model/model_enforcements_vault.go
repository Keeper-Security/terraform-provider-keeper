package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

type EnforcementsVaultDataSourceModel struct {
	AllowSecretsManager	types.Bool	`tfsdk:"allow_secrets_manager"`
	RestrictBreachWatch	types.Bool	`tfsdk:"restrict_breach_watch"`
	SendBreachWatchEvents	types.Bool	`tfsdk:"send_breach_watch_events"`
	DisableSetupTour	types.Bool	`tfsdk:"disable_setup_tour"`
	DaysBeforeDeletedRecordsAutoCleared	types.Int64	`tfsdk:"days_before_deleted_records_auto_cleared"`
	DaysBeforeDeletedRecordsClearedPerm	types.Int64	`tfsdk:"days_before_deleted_records_cleared_perm"`
	GeneratedSecurityQuestionComplexity	types.String	`tfsdk:"generated_security_question_complexity"`
	GeneratedPasswordComplexity	types.String	`tfsdk:"generated_password_complexity"`
	MaskPasswordsWhileEditing	types.Bool	`tfsdk:"mask_passwords_while_editing"`
	MaskNotes	types.Bool	`tfsdk:"mask_notes"`
	MaskCustomFields	types.Bool	`tfsdk:"mask_custom_fields"`
	RestrictCreateIdentityPaymentRecords	types.Bool	`tfsdk:"restrict_create_identity_payment_records"`
	RestrictCreateFolder	types.Bool	`tfsdk:"restrict_create_folder"`
}

func (model *EnforcementsVaultDataSourceModel) FromKeeper(enforcements map[string]string) {
	var intValue int64
	var err error
	intValue, err = strconv.ParseInt(enforcements["days_before_deleted_records_auto_cleared"], 0, 64)
	if err != nil {
		model.DaysBeforeDeletedRecordsAutoCleared = types.Int64Value(intValue)
	}
	intValue, err = strconv.ParseInt(enforcements["days_before_deleted_records_cleared_perm"], 0, 64)
	if err != nil {
		model.DaysBeforeDeletedRecordsClearedPerm = types.Int64Value(intValue)
	}
	model.AllowSecretsManager = types.BoolValue(enforcements["allow_secrets_manager"] != "")
	model.RestrictBreachWatch = types.BoolValue(enforcements["restrict_breach_watch"] != "")
	model.SendBreachWatchEvents = types.BoolValue(enforcements["send_breach_watch_events"] != "")
	model.DisableSetupTour = types.BoolValue(enforcements["disable_setup_tour"] != "")
	model.MaskPasswordsWhileEditing = types.BoolValue(enforcements["mask_passwords_while_editing"] != "")
	model.MaskNotes = types.BoolValue(enforcements["mask_notes"] != "")
	model.MaskCustomFields = types.BoolValue(enforcements["mask_custom_fields"] != "")
	model.RestrictCreateIdentityPaymentRecords = types.BoolValue(enforcements["restrict_create_identity_payment_records"] != "")
	model.RestrictCreateFolder = types.BoolValue(enforcements["restrict_create_folder"] != "")
	model.GeneratedSecurityQuestionComplexity = types.StringValue(enforcements["generated_security_question_complexity"])
	model.GeneratedPasswordComplexity = types.StringValue(enforcements["generated_password_complexity"])
}