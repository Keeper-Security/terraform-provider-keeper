package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"strings"
	"terraform-provider-kepr/internal/model"
)

var (
	_ datasource.DataSourceWithConfigure = &teamDataSource{}
)

type teamDataSourceModel struct {
	TeamUid        types.String      `tfsdk:"team_uid"`
	Name           types.String      `tfsdk:"name"`
	NodeId         types.Int64       `tfsdk:"node_id"`
	RestrictEdit   types.Bool        `tfsdk:"restrict_edit"`
	RestrictShare  types.Bool        `tfsdk:"restrict_share"`
	RestrictView   types.Bool        `tfsdk:"restrict_view"`
	IncludeMembers types.Bool        `tfsdk:"include_members"`
	Users          []*userShortModel `tfsdk:"users"`
	Roles          []*roleShortModel `tfsdk:"roles"`
}

func (model *teamDataSourceModel) fromKeeper(keeper enterprise.ITeam) {
	model.TeamUid = types.StringValue(keeper.TeamUid())
	model.Name = types.StringValue(keeper.Name())
	model.NodeId = types.Int64Value(keeper.NodeId())
	model.RestrictEdit = types.BoolValue(keeper.RestrictEdit())
	model.RestrictShare = types.BoolValue(keeper.RestrictShare())
	model.RestrictView = types.BoolValue(keeper.RestrictView())
}

type teamDataSource struct {
	teams        enterprise.IEnterpriseEntity[enterprise.ITeam, string]
	teamUsers    enterprise.IEnterpriseLink[enterprise.ITeamUser, string, int64]
	users        enterprise.IEnterpriseEntity[enterprise.IUser, int64]
	roleTeams    enterprise.IEnterpriseLink[enterprise.IRoleTeam, int64, string]
	roles        enterprise.IEnterpriseEntity[enterprise.IRole, int64]
	managedNodes enterprise.IEnterpriseLink[enterprise.IManagedNode, int64, int64]
}

func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

func (d *teamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the data source.
func (d *teamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var filterAttributes = map[string]schema.Attribute{
		"team_uid": schema.StringAttribute{
			Optional:    true,
			Description: "Team UID",
		},
		"name": schema.StringAttribute{
			Optional:    true,
			Description: "Team Name",
		},
		"include_members": schema.BoolAttribute{
			Optional:    true,
			Description: "Include team members",
		},
	}
	var usersAttribute = map[string]schema.Attribute{
		"users": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: userShortSchemaAttributes,
			},
		},
	}
	var rolesAttribute = map[string]schema.Attribute{
		"roles": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: roleShortSchemaAttributes,
			},
		},
	}
	resp.Schema = schema.Schema{
		Attributes: model.MergeMaps(filterAttributes, teamSchemaAttributes, usersAttribute, rolesAttribute),
	}
}

func (d *teamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.teams = ed.Teams()
		d.teamUsers = ed.TeamUsers()
		d.users = ed.Users()
		d.roleTeams = ed.RoleTeams()
		d.roles = ed.Roles()

	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseData\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tq teamDataSourceModel
	diags := req.Config.Get(ctx, &tq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var teamMatcher func(enterprise.ITeam) bool
	if !tq.Name.IsNull() && !tq.Name.IsUnknown() {
		teamMatcher = func(team enterprise.ITeam) bool {
			return strings.EqualFold(tq.Name.ValueString(), team.Name())
		}
	} else if !tq.TeamUid.IsNull() && !tq.TeamUid.IsUnknown() {
		teamMatcher = func(team enterprise.ITeam) bool {
			return strings.EqualFold(team.TeamUid(), tq.TeamUid.ValueString())
		}
	}

	if teamMatcher == nil {
		resp.Diagnostics.AddError(
			"Search criteria is not provided for \"team\" data source",
			fmt.Sprintf("Search criteria is not provided for \"team\" data source"),
		)
		return
	}

	var team enterprise.ITeam
	d.teams.GetAllEntities(func(t enterprise.ITeam) bool {
		if teamMatcher(t) {
			team = t
			return false
		}
		return true
	})

	if team == nil {
		resp.Diagnostics.AddError(
			"Team not found",
			fmt.Sprintf("Cannot find a team according to the provided criteria"),
		)
		return
	}

	var state = tq
	var tm = &state
	tm.fromKeeper(team)
	if !state.IncludeMembers.IsNull() && state.IncludeMembers.ValueBool() {
		d.teamUsers.GetLinksBySubject(team.TeamUid(), func(tu enterprise.ITeamUser) bool {
			var u = d.users.GetEntity(tu.EnterpriseUserId())
			if u != nil {
				var um = new(userShortModel)
				um.fromKeeper(u)
				state.Users = append(state.Users, um)
			}
			return true
		})
		d.roleTeams.GetLinksByObject(team.TeamUid(), func(rt enterprise.IRoleTeam) bool {
			var r = d.roles.GetEntity(rt.RoleId())
			if r != nil {
				var rsm = new(roleShortModel)
				var isAdmin bool
				d.managedNodes.GetLinksBySubject(rt.RoleId(), func(_ enterprise.IManagedNode) bool {
					isAdmin = true
					return false
				})
				rsm.fromKeeper(r, isAdmin)
				state.Roles = append(state.Roles, rsm)
			}
			return true
		})
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *teamDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("team_uid"),
			path.MatchRoot("name"),
		),
	}
}
