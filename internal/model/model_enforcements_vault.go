package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsVaultDataSourceModel struct {
	RestrictBreachWatch                     types.Bool   `tfsdk:"restrict_breach_watch"`
	SendBreachWatchEvents                   types.Bool   `tfsdk:"send_breach_watch_events"`
	DaysBeforeDeletedRecordsAutoCleared     types.Int64  `tfsdk:"days_before_deleted_records_auto_cleared"`
	DaysBeforeDeletedRecordsClearedPerm     types.Int64  `tfsdk:"days_before_deleted_records_cleared_perm"`
	GeneratedSecurityQuestionComplexity     types.String `tfsdk:"generated_security_question_complexity"`
	GeneratedPasswordComplexity             types.String `tfsdk:"generated_password_complexity"`
	MaskPasswordsWhileEditing               types.Bool   `tfsdk:"mask_passwords_while_editing"`
	MaskNotes                               types.Bool   `tfsdk:"mask_notes"`
	MaskCustomFields                        types.Bool   `tfsdk:"mask_custom_fields"`
	RestrictCreateIdentityPaymentRecords    types.Bool   `tfsdk:"restrict_create_identity_payment_records"`
	RestrictCreateFolder                    types.Bool   `tfsdk:"restrict_create_folder"`
	RestrictCreateFolderToOnlySharedFolders types.Bool   `tfsdk:"restrict_create_folder_to_only_shared_folders"`
	RestrictCreateRecord                    types.Bool   `tfsdk:"restrict_create_record"`
	RestrictCreateRecordToSharedFolders     types.Bool   `tfsdk:"restrict_create_record_to_shared_folders"`
}

func (evm *EnforcementsVaultDataSourceModel) ToKeeper(enforcements map[string]string) {
	getInt64Value(evm.DaysBeforeDeletedRecordsAutoCleared, "days_before_deleted_records_auto_cleared", enforcements)
	getInt64Value(evm.DaysBeforeDeletedRecordsClearedPerm, "days_before_deleted_records_cleared_perm", enforcements)
	getBoolValue(evm.RestrictBreachWatch, "restrict_breach_watch", enforcements)
	getBoolValue(evm.SendBreachWatchEvents, "send_breach_watch_events", enforcements)
	getBoolValue(evm.MaskPasswordsWhileEditing, "mask_passwords_while_editing", enforcements)
	getBoolValue(evm.MaskNotes, "mask_notes", enforcements)
	getBoolValue(evm.MaskCustomFields, "mask_custom_fields", enforcements)
	getBoolValue(evm.RestrictCreateIdentityPaymentRecords, "restrict_create_identity_payment_records", enforcements)
	getBoolValue(evm.RestrictCreateFolder, "restrict_create_folder", enforcements)
	getBoolValue(evm.RestrictCreateFolderToOnlySharedFolders, "restrict_create_folder_to_only_shared_folders", enforcements)
	getBoolValue(evm.RestrictCreateRecord, "restrict_create_record", enforcements)
	getBoolValue(evm.RestrictCreateRecordToSharedFolders, "restrict_create_record_to_shared_folders", enforcements)
	getStringValue(evm.GeneratedSecurityQuestionComplexity, "generated_security_question_complexity", enforcements)
	getStringValue(evm.GeneratedPasswordComplexity, "generated_password_complexity", enforcements)
}

func (evm *EnforcementsVaultDataSourceModel) FromKeeper(enforcements map[string]string) {
	setInt64Value(&evm.DaysBeforeDeletedRecordsAutoCleared, "days_before_deleted_records_auto_cleared", enforcements)
	setInt64Value(&evm.DaysBeforeDeletedRecordsClearedPerm, "days_before_deleted_records_cleared_perm", enforcements)
	setBoolValue(&evm.RestrictBreachWatch, "restrict_breach_watch", enforcements)
	setBoolValue(&evm.SendBreachWatchEvents, "send_breach_watch_events", enforcements)
	setBoolValue(&evm.MaskPasswordsWhileEditing, "mask_passwords_while_editing", enforcements)
	setBoolValue(&evm.MaskNotes, "mask_notes", enforcements)
	setBoolValue(&evm.MaskCustomFields, "mask_custom_fields", enforcements)
	setBoolValue(&evm.RestrictCreateIdentityPaymentRecords, "restrict_create_identity_payment_records", enforcements)
	setBoolValue(&evm.RestrictCreateFolder, "restrict_create_folder", enforcements)
	setBoolValue(&evm.RestrictCreateFolderToOnlySharedFolders, "restrict_create_folder_to_only_shared_folders", enforcements)
	setBoolValue(&evm.RestrictCreateRecord, "restrict_create_record", enforcements)
	setBoolValue(&evm.RestrictCreateRecordToSharedFolders, "restrict_create_record_to_shared_folders", enforcements)
	setStringValue(&evm.GeneratedSecurityQuestionComplexity, "generated_security_question_complexity", enforcements)
	setStringValue(&evm.GeneratedPasswordComplexity, "generated_password_complexity", enforcements)
}

func (evm *EnforcementsVaultDataSourceModel) IsBlank() bool {
	if !evm.RestrictBreachWatch.IsNull() {
		return false
	}
	if !evm.SendBreachWatchEvents.IsNull() {
		return false
	}
	if !evm.DaysBeforeDeletedRecordsAutoCleared.IsNull() {
		return false
	}
	if !evm.DaysBeforeDeletedRecordsClearedPerm.IsNull() {
		return false
	}
	if !evm.GeneratedSecurityQuestionComplexity.IsNull() {
		return false
	}
	if !evm.GeneratedPasswordComplexity.IsNull() {
		return false
	}
	if !evm.MaskPasswordsWhileEditing.IsNull() {
		return false
	}
	if !evm.MaskNotes.IsNull() {
		return false
	}
	if !evm.MaskCustomFields.IsNull() {
		return false
	}
	if !evm.RestrictCreateIdentityPaymentRecords.IsNull() {
		return false
	}
	if !evm.RestrictCreateFolder.IsNull() {
		return false
	}
	if !evm.RestrictCreateFolderToOnlySharedFolders.IsNull() {
		return false
	}
	if !evm.RestrictCreateRecord.IsNull() {
		return false
	}
	if !evm.RestrictCreateRecordToSharedFolders.IsNull() {
		return false
	}
	return true
}
