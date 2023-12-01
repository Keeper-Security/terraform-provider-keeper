package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/api"
	"strings"
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

func (model *PrivilegeDataSourceModel) FromKeeper(privileges []string) {
	var s = api.NewSet[string]()
	s.Union(api.SliceSelect(privileges, func(x string) string { return strings.ToLower(x) }))
	model.ManageNodes = types.BoolValue(s.Has("manage_nodes"))
	model.ManageUsers = types.BoolValue(s.Has("manage_user"))
	model.ManageRoles = types.BoolValue(s.Has("manage_roles"))
	model.ManageTeams = types.BoolValue(s.Has("manage_teams"))
	model.ManageReports = types.BoolValue(s.Has("run_reports"))
	model.ManageSso = types.BoolValue(s.Has("manage_bridge"))
	model.DeviceApproval = types.BoolValue(s.Has("approve_device"))
	model.ManageRecordTypes = types.BoolValue(s.Has("manage_record_types"))
	model.ShareAdmin = types.BoolValue(s.Has("sharing_administrator"))
	model.RunComplianceReports = types.BoolValue(s.Has("run_compliance_reports"))
	model.TransferAccount = types.BoolValue(s.Has("transfer_account"))
	model.ManageCompanies = types.BoolValue(s.Has("manage_companies"))
}

func (model *PrivilegeDataSourceModel) ToKeeper() (privileges []string) {
	if !model.ManageNodes.IsNull() && model.ManageNodes.ValueBool() {
		privileges = append(privileges, "manage_nodes")
	}
	if !model.ManageUsers.IsNull() && model.ManageUsers.ValueBool() {
		privileges = append(privileges, "manage_user")
	}
	if !model.ManageRoles.IsNull() && model.ManageRoles.ValueBool() {
		privileges = append(privileges, "manage_roles")
	}
	if !model.ManageTeams.IsNull() && model.ManageTeams.ValueBool() {
		privileges = append(privileges, "manage_teams")
	}
	if !model.ManageReports.IsNull() && model.ManageReports.ValueBool() {
		privileges = append(privileges, "run_reports")
	}
	if !model.ManageSso.IsNull() && model.ManageSso.ValueBool() {
		privileges = append(privileges, "manage_bridge")
	}
	if !model.DeviceApproval.IsNull() && model.DeviceApproval.ValueBool() {
		privileges = append(privileges, "approve_device")
	}
	if !model.ManageRecordTypes.IsNull() && model.ManageRecordTypes.ValueBool() {
		privileges = append(privileges, "manage_record_types")
	}
	if !model.ShareAdmin.IsNull() && model.ShareAdmin.ValueBool() {
		privileges = append(privileges, "sharing_administrator")
	}
	if !model.RunComplianceReports.IsNull() && model.RunComplianceReports.ValueBool() {
		privileges = append(privileges, "run_compliance_reports")
	}
	if !model.TransferAccount.IsNull() && model.TransferAccount.ValueBool() {
		privileges = append(privileges, "transfer_account")
	}
	if !model.ManageCompanies.IsNull() && model.ManageCompanies.ValueBool() {
		privileges = append(privileges, "manage_companies")
	}
	return
}
