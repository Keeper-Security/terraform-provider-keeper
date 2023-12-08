package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsSharingDataSourceModel struct {
	RestrictSharingOutsideOfIsolatedNodes	types.Bool	`tfsdk:"restrict_sharing_outside_of_isolated_nodes"`
	RestrictLinkSharing	types.Bool	`tfsdk:"restrict_link_sharing"`
	RestrictSharingRecordWithAttachments	types.Bool	`tfsdk:"restrict_sharing_record_with_attachments"`
	RestirctSharingRecordAndFolder	types.Bool	`tfsdk:"restirct_sharing_record_and_folder"`
	RestrictSharingIncomingAll	types.Bool	`tfsdk:"restrict_sharing_incoming_all"`
	RequireAccountShare	types.String	`tfsdk:"require_account_share"`
	RestrictFileUpload	types.Bool	`tfsdk:"restrict_file_upload"`
	RestrictImport	types.Bool	`tfsdk:"restrict_import"`
	RestrictExport	types.Bool	`tfsdk:"restrict_export"`
	RestrictSharingEnterprise	types.Bool	`tfsdk:"restrict_sharing_enterprise"`
	RestrictSharingAll	types.Bool	`tfsdk:"restrict_sharing_all"`
	RestrictCreateRecordToSharedFolders	types.Bool	`tfsdk:"restrict_create_record_to_shared_folders"`
	RestrictCreateRecord	types.Bool	`tfsdk:"restrict_create_record"`
	RestrictCreateFolderToOnlySharedFolders	types.Bool	`tfsdk:"restrict_create_folder_to_only_shared_folders"`
	RestrictDirectSharing	types.Bool	`tfsdk:"restrict_direct_sharing"`
}

func (model *EnforcementsSharingDataSourceModel) FromKeeper(enforcements map[string]string) {
	model.RestrictSharingOutsideOfIsolatedNodes = types.BoolValue(enforcements["restrict_sharing_outside_of_isolated_nodes"] != "")
	model.RestrictLinkSharing = types.BoolValue(enforcements["restrict_link_sharing"] != "")
	model.RestrictSharingRecordWithAttachments = types.BoolValue(enforcements["restrict_sharing_record_with_attachments"] != "")
	model.RestirctSharingRecordAndFolder = types.BoolValue(enforcements["restirct_sharing_record_and_folder"] != "")
	model.RestrictSharingIncomingAll = types.BoolValue(enforcements["restrict_sharing_incoming_all"] != "")
	model.RestrictFileUpload = types.BoolValue(enforcements["restrict_file_upload"] != "")
	model.RestrictImport = types.BoolValue(enforcements["restrict_import"] != "")
	model.RestrictExport = types.BoolValue(enforcements["restrict_export"] != "")
	model.RestrictSharingEnterprise = types.BoolValue(enforcements["restrict_sharing_enterprise"] != "")
	model.RestrictSharingAll = types.BoolValue(enforcements["restrict_sharing_all"] != "")
	model.RestrictCreateRecordToSharedFolders = types.BoolValue(enforcements["restrict_create_record_to_shared_folders"] != "")
	model.RestrictCreateRecord = types.BoolValue(enforcements["restrict_create_record"] != "")
	model.RestrictCreateFolderToOnlySharedFolders = types.BoolValue(enforcements["restrict_create_folder_to_only_shared_folders"] != "")
	model.RestrictDirectSharing = types.BoolValue(enforcements["restrict_direct_sharing"] != "")
	model.RequireAccountShare = types.StringValue(enforcements["require_account_share"])
}