package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-keeper/internal/model"
)

var (
	_ datasource.DataSource = &privilegeDataSource{}
)

func newPrivilegeDataSource() datasource.DataSource {
	return &privilegeDataSource{}
}

type privilegeDataSource struct {
}

func (d *privilegeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privileges"
}

// Schema defines the schema for the data source.
func (d *privilegeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: model.PrivilegesDataSourceAttributes,
	}
}

func (d *privilegeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *privilegeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.PrivilegeDataSourceModel
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
