package model

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
)

type PrivilegeDataSourceModel struct {
	ManageNodes          types.Bool `tfsdk:"manage_nodes"`
	ManageUsers          types.Bool `tfsdk:"manage_users"`
	ManageRoles          types.Bool `tfsdk:"manage_roles"`
	ManageTeams          types.Bool `tfsdk:"manage_teams"`
	ManageReports        types.Bool `tfsdk:"manage_reports"`
	ManageSso            types.Bool `tfsdk:"manage_sso"`
	DeviceApproval       types.Bool `tfsdk:"device_approval"`
	ManageRecordTypes    types.Bool `tfsdk:"manage_record_types"`
	ShareAdmin           types.Bool `tfsdk:"share_admin"`
	RunComplianceReports types.Bool `tfsdk:"run_compliance_reports"`
	TransferAccount      types.Bool `tfsdk:"transfer_account"`
	ManageCompanies      types.Bool `tfsdk:"manage_companies"`
}

func (pds *PrivilegeDataSourceModel) FromKeeper(rp enterprise.IRolePrivilege) {
	pds.ManageNodes = types.BoolValue(rp.ManageNodes())
	pds.ManageUsers = types.BoolValue(rp.ManageUsers())
	pds.ManageRoles = types.BoolValue(rp.ManageRoles())
	pds.ManageTeams = types.BoolValue(rp.ManageTeams())
	pds.ManageReports = types.BoolValue(rp.RunReports())
	pds.ManageSso = types.BoolValue(rp.ManageBridge())
	pds.DeviceApproval = types.BoolValue(rp.ApproveDevices())
	pds.ManageRecordTypes = types.BoolValue(rp.ManageRecordTypes())
	pds.ShareAdmin = types.BoolValue(rp.SharingAdministrator())
	pds.RunComplianceReports = types.BoolValue(rp.RunComplianceReport())
	pds.TransferAccount = types.BoolValue(rp.TransferAccount())
	pds.ManageCompanies = types.BoolValue(rp.ManageCompanies())
}

func (pds *PrivilegeDataSourceModel) ToKeeper(result enterprise.IRolePrivilegeEdit) {
	if !pds.ManageNodes.IsNull() {
		result.SetManageNodes(pds.ManageNodes.ValueBool())
	}
	if !pds.ManageUsers.IsNull() {
		result.SetManageUsers(pds.ManageUsers.ValueBool())
	}
	if !pds.ManageRoles.IsNull() {
		result.SetManageRoles(pds.ManageRoles.ValueBool())
	}
	if !pds.ManageTeams.IsNull() {
		result.SetManageTeams(pds.ManageTeams.ValueBool())
	}
	if !pds.ManageReports.IsNull() {
		result.SetRunReports(pds.ManageReports.ValueBool())
	}
	if !pds.ManageSso.IsNull() {
		result.SetManageBridge(pds.ManageSso.ValueBool())
	}
	if !pds.DeviceApproval.IsNull() {
		result.SetApproveDevices(pds.DeviceApproval.ValueBool())
	}
	if !pds.ManageRecordTypes.IsNull() {
		result.SetManageRecordTypes(pds.ManageRecordTypes.ValueBool())
	}
	if !pds.ShareAdmin.IsNull() {
		result.SetSharingAdministrator(pds.ShareAdmin.ValueBool())
	}
	if !pds.RunComplianceReports.IsNull() {
		result.SetRunComplianceReport(pds.RunComplianceReports.ValueBool())
	}
	if !pds.TransferAccount.IsNull() {
		result.SetTransferAccount(pds.TransferAccount.ValueBool())
	}
	if !pds.ManageCompanies.IsNull() {
		result.SetManageCompanies(pds.ManageCompanies.ValueBool())
	}
}

var PrivilegesDataSourceAttributes = map[string]schema.Attribute{
	"manage_nodes": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Nodes",
	},
	"manage_users": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Users",
	},
	"manage_teams": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Teams",
	},
	"manage_roles": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Roles",
	},
	"manage_reports": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Reporting and Alerts",
	},
	"manage_sso": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Bridge/SSO",
	},
	"device_approval": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Perform Device Approvals",
	},
	"manage_record_types": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Record Types in Vault",
		MarkdownDescription: "This permission allows the admin rights to create, edit, or delete Record Types " +
			"which have pre-defined fields. Record Types appear during creation of records in the user's vault.",
	},
	"share_admin": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Share Admin",
	},
	"run_compliance_reports": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Run Compliance Reports",
	},
	"transfer_account": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Transfer Account",
	},
	"manage_companies": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Companies",
	},
}
