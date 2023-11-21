package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
)

type teamModel struct {
	TeamUid       types.String `tfsdk:"team_uid"`
	Name          types.String `tfsdk:"name"`
	NodeId        types.Int64  `tfsdk:"node_id"`
	RestrictEdit  types.Bool   `tfsdk:"restrict_edit"`
	RestrictShare types.Bool   `tfsdk:"restrict_share"`
	RestrictView  types.Bool   `tfsdk:"restrict_view"`
}

func (model *teamModel) toKeeper(keeper enterprise.ITeamEdit) {
	keeper.SetName(model.Name.ValueString())
	if model.NodeId.IsNull() {
		keeper.SetNodeId(0)
	} else {
		keeper.SetNodeId(model.NodeId.ValueInt64())
	}
	keeper.SetRestrictEdit(!model.RestrictEdit.IsNull() && model.RestrictEdit.ValueBool())
	keeper.SetRestrictShare(!model.RestrictShare.IsNull() && model.RestrictShare.ValueBool())
	keeper.SetRestrictView(!model.RestrictView.IsNull() && model.RestrictView.ValueBool())
}

func (model *teamModel) fromKeeper(keeper enterprise.ITeam) {
	model.TeamUid = types.StringValue(keeper.TeamUid())
	model.Name = types.StringValue(keeper.Name())
	model.NodeId = types.Int64Value(keeper.NodeId())
	model.RestrictEdit = types.BoolValue(keeper.RestrictEdit())
	model.RestrictShare = types.BoolValue(keeper.RestrictShare())
	model.RestrictView = types.BoolValue(keeper.RestrictView())
}

var teamSchemaAttributes = map[string]schema.Attribute{
	"team_uid": schema.StringAttribute{
		Computed:    true,
		Description: "Team UID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Team Name",
	},
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Team NodeID",
	},
	"restrict_edit": schema.BoolAttribute{
		Computed:    true,
		Description: "Restrict Edit flag",
	},
	"restrict_share": schema.BoolAttribute{
		Computed:    true,
		Description: "Restrict Share flag",
	},
	"restrict_view": schema.BoolAttribute{
		Computed:    true,
		Description: "Restrict View flag",
	},
}

type teamShortModel struct {
	TeamUid types.String `tfsdk:"team_uid"`
	Name    types.String `tfsdk:"name"`
}

func (model *teamShortModel) fromKeeper(keeper enterprise.ITeam) {
	model.TeamUid = types.StringValue(keeper.TeamUid())
	model.Name = types.StringValue(keeper.Name())
}

var teamShortSchemaAttributes = map[string]schema.Attribute{
	"team_uid": schema.StringAttribute{
		Computed:    true,
		Description: "Team UID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Team Name",
	},
}
