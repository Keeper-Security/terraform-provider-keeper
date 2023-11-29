package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/api"
	"strings"
)

var (
	_ datasource.DataSource = &privilegeDataSource{}
)

type privilegeDataSourceModel struct {
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

func (model *privilegeDataSourceModel) fromKeeper(privileges []string) {
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

type privilegeDataSource struct {
}

func NewPrivilegeDataSource() datasource.DataSource {
	return &privilegeDataSource{}
}
func (d *privilegeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privileges"
}

// Schema defines the schema for the data source.
func (d *privilegeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: privilegesAttributes,
	}
}

func (d *privilegeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *privilegeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq privilegeDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rq.ManageNodes.IsNull() {
		rq.ManageNodes = types.BoolValue(false)
	}
	if rq.ManageUsers.IsNull() {
		rq.ManageUsers = types.BoolValue(false)
	}
	if rq.ManageRoles.IsNull() {
		rq.ManageRoles = types.BoolValue(false)
	}
	if rq.ManageTeams.IsNull() {
		rq.ManageTeams = types.BoolValue(false)
	}
	if rq.ManageReports.IsNull() {
		rq.ManageReports = types.BoolValue(false)
	}
	if rq.ManageSso.IsNull() {
		rq.ManageSso = types.BoolValue(false)
	}
	if rq.DeviceApproval.IsNull() {
		rq.DeviceApproval = types.BoolValue(false)
	}
	if rq.ManageRecordTypes.IsNull() {
		rq.ManageRecordTypes = types.BoolValue(false)
	}
	if rq.ShareAdmin.IsNull() {
		rq.ShareAdmin = types.BoolValue(false)
	} else {
		if rq.ShareAdmin.ValueBool() {
			rq.ManageUsers = types.BoolValue(true)
		}
	}
	if rq.RunComplianceReports.IsNull() {
		rq.RunComplianceReports = types.BoolValue(false)
	}
	if rq.TransferAccount.IsNull() {
		rq.TransferAccount = types.BoolValue(false)
	}
	if rq.ManageCompanies.IsNull() {
		rq.ManageCompanies = types.BoolValue(false)
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
