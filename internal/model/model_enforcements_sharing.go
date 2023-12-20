package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsSharingDataSourceModel struct {
	RestrictSharingOutsideOfIsolatedNodes types.Bool   `tfsdk:"restrict_sharing_outside_of_isolated_nodes"`
	RestrictLinkSharing                   types.Bool   `tfsdk:"restrict_link_sharing"`
	RestrictSharingRecordWithAttachments  types.Bool   `tfsdk:"restrict_sharing_record_with_attachments"`
	RequireAccountShare                   types.String `tfsdk:"require_account_share"`
	RestrictFileUpload                    types.Bool   `tfsdk:"restrict_file_upload"`
	RestrictExport                        types.Bool   `tfsdk:"restrict_export"`
	RestrictImport                        types.Bool   `tfsdk:"restrict_import"`
	RestrictSharingEnterpriseIncoming     types.Bool   `tfsdk:"restrict_sharing_enterprise_incoming"`
	RestrictSharingEnterpriseOutgoing     types.Bool   `tfsdk:"restrict_sharing_enterprise_outgoing"`
	RestrictSharingAllIncoming            types.Bool   `tfsdk:"restrict_sharing_all_incoming"`
	RestrictSharingAllOutgoing            types.Bool   `tfsdk:"restrict_sharing_all_outgoing"`
	RestrictSharingRecordToSharedFolders  types.Bool   `tfsdk:"restrict_sharing_record_to_shared_folders"`
	RestrictCreateSharedFolder            types.Bool   `tfsdk:"restrict_create_shared_folder"`
}

func (esm *EnforcementsSharingDataSourceModel) ToKeeper(enforcements map[string]string) {
	getBoolValue(esm.RestrictSharingOutsideOfIsolatedNodes, "restrict_sharing_outside_of_isolated_nodes", enforcements)
	getBoolValue(esm.RestrictLinkSharing, "restrict_link_sharing", enforcements)
	getBoolValue(esm.RestrictSharingRecordWithAttachments, "restrict_sharing_record_with_attachments", enforcements)
	getBoolValue(esm.RestrictSharingAllIncoming, "restrict_sharing_all_incoming", enforcements)
	getBoolValue(esm.RestrictFileUpload, "restrict_file_upload", enforcements)
	getBoolValue(esm.RestrictExport, "restrict_export", enforcements)
	getBoolValue(esm.RestrictImport, "restrict_import", enforcements)
	getBoolValue(esm.RestrictSharingEnterpriseIncoming, "restrict_sharing_enterprise_incoming", enforcements)
	getBoolValue(esm.RestrictSharingEnterpriseOutgoing, "restrict_sharing_enterprise_outgoing", enforcements)
	getBoolValue(esm.RestrictSharingAllOutgoing, "restrict_sharing_all_outgoing", enforcements)
	getBoolValue(esm.RestrictSharingRecordToSharedFolders, "restrict_sharing_record_to_shared_folders", enforcements)
	getBoolValue(esm.RestrictCreateSharedFolder, "restrict_create_shared_folder", enforcements)
	getStringValue(esm.RequireAccountShare, "require_account_share", enforcements)
}

func (esm *EnforcementsSharingDataSourceModel) FromKeeper(enforcements map[string]string) {
	setBoolValue(&esm.RestrictSharingOutsideOfIsolatedNodes, "restrict_sharing_outside_of_isolated_nodes", enforcements)
	setBoolValue(&esm.RestrictLinkSharing, "restrict_link_sharing", enforcements)
	setBoolValue(&esm.RestrictSharingRecordWithAttachments, "restrict_sharing_record_with_attachments", enforcements)
	setBoolValue(&esm.RestrictSharingAllIncoming, "restrict_sharing_all_incoming", enforcements)
	setBoolValue(&esm.RestrictFileUpload, "restrict_file_upload", enforcements)
	setBoolValue(&esm.RestrictExport, "restrict_export", enforcements)
	setBoolValue(&esm.RestrictImport, "restrict_import", enforcements)
	setBoolValue(&esm.RestrictSharingEnterpriseIncoming, "restrict_sharing_enterprise_incoming", enforcements)
	setBoolValue(&esm.RestrictSharingEnterpriseOutgoing, "restrict_sharing_enterprise_outgoing", enforcements)
	setBoolValue(&esm.RestrictSharingAllOutgoing, "restrict_sharing_all_outgoing", enforcements)
	setBoolValue(&esm.RestrictSharingRecordToSharedFolders, "restrict_sharing_record_to_shared_folders", enforcements)
	setBoolValue(&esm.RestrictCreateSharedFolder, "restrict_create_shared_folder", enforcements)
	setStringValue(&esm.RequireAccountShare, "require_account_share", enforcements)
}

func (esm *EnforcementsSharingDataSourceModel) IsBlank() bool {
	if !esm.RestrictSharingOutsideOfIsolatedNodes.IsNull() {
		return false
	}
	if !esm.RestrictLinkSharing.IsNull() {
		return false
	}
	if !esm.RestrictSharingRecordWithAttachments.IsNull() {
		return false
	}
	if !esm.RestrictSharingAllIncoming.IsNull() {
		return false
	}
	if !esm.RequireAccountShare.IsNull() {
		return false
	}
	if !esm.RestrictFileUpload.IsNull() {
		return false
	}
	if !esm.RestrictExport.IsNull() {
		return false
	}
	if !esm.RestrictImport.IsNull() {
		return false
	}
	if !esm.RestrictSharingEnterpriseIncoming.IsNull() {
		return false
	}
	if !esm.RestrictSharingEnterpriseOutgoing.IsNull() {
		return false
	}
	if !esm.RestrictSharingAllOutgoing.IsNull() {
		return false
	}
	if !esm.RestrictCreateSharedFolder.IsNull() {
		return false
	}

	return true
}
