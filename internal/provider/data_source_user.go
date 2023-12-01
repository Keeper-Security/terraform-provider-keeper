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
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

type userDataSourceModel struct {
	EnterpriseUserId       types.Int64       `tfsdk:"enterprise_user_id"`
	Username               types.String      `tfsdk:"username"`
	NodeId                 types.Int64       `tfsdk:"node_id"`
	Status                 types.String      `tfsdk:"status"`
	AccountShareExpiration types.Int64       `tfsdk:"account_share_expiration"`
	TfaEnabled             types.Bool        `tfsdk:"tfa_enabled"`
	FullName               types.String      `tfsdk:"full_name"`
	JobTitle               types.String      `tfsdk:"job_title"`
	IncludeTeams           types.Bool        `tfsdk:"include_teams"`
	Teams                  []*teamShortModel `tfsdk:"teams"`
	IncludeRoles           types.Bool        `tfsdk:"include_roles"`
	Roles                  []*roleShortModel `tfsdk:"roles"`
}

func (model *userDataSourceModel) fromKeeper(keeper enterprise.IUser) {
	model.EnterpriseUserId = types.Int64Value(keeper.EnterpriseUserId())
	model.Username = types.StringValue(keeper.Username())
	model.NodeId = types.Int64Value(keeper.NodeId())
	model.Status = types.StringValue(getUserStatus(keeper))
	model.TfaEnabled = types.BoolValue(keeper.TfaEnabled())
	if len(keeper.FullName()) > 0 {
		model.FullName = types.StringValue(keeper.FullName())
	} else {
		model.FullName = types.StringNull()
	}
	if len(keeper.JobTitle()) > 0 {
		model.JobTitle = types.StringValue(keeper.JobTitle())
	} else {
		model.JobTitle = types.StringNull()
	}
	if keeper.AccountShareExpiration() > 0 {
		model.AccountShareExpiration = types.Int64Value(keeper.AccountShareExpiration())
	} else {
		model.AccountShareExpiration = types.Int64Null()
	}
}

type userDataSource struct {
	users        enterprise.IEnterpriseEntity[enterprise.IUser, int64]
	teams        enterprise.IEnterpriseEntity[enterprise.ITeam, string]
	teamUsers    enterprise.IEnterpriseLink[enterprise.ITeamUser, string, int64]
	roles        enterprise.IEnterpriseEntity[enterprise.IRole, int64]
	roleUsers    enterprise.IEnterpriseLink[enterprise.IRoleUser, int64, int64]
	managedNodes enterprise.IEnterpriseLink[enterprise.IManagedNode, int64, int64]
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var filterAttributes = map[string]schema.Attribute{
		"enterprise_user_id": schema.Int64Attribute{
			Optional:    true,
			Description: "Enterprise User ID",
		},
		"username": schema.StringAttribute{
			Optional:    true,
			Description: "User email",
		},
		"include_teams": schema.BoolAttribute{
			Optional:    true,
			Description: "Include user teams",
		},
		"include_roles": schema.BoolAttribute{
			Optional:    true,
			Description: "Include user roles",
		},
	}
	var teamShortAttribute = map[string]schema.Attribute{
		"teams": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: teamShortSchemaAttributes,
			},
		},
	}
	var roleShortAttribute = map[string]schema.Attribute{
		"roles": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: roleShortSchemaAttributes,
			},
		},
	}

	resp.Schema = schema.Schema{
		Attributes: model.MergeMaps(filterAttributes, userSchemaAttributes, teamShortAttribute, roleShortAttribute),
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.users = ed.Users()
		d.teams = ed.Teams()
		d.teamUsers = ed.TeamUsers()
		d.roles = ed.Roles()
		d.roleUsers = ed.RoleUsers()
		d.managedNodes = ed.ManagedNodes()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseData\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var uq userDataSourceModel
	diags := req.Config.Get(ctx, &uq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var userMatcher func(enterprise.IUser) bool
	if !uq.Username.IsNull() && !uq.Username.IsUnknown() {
		var emailFilter = uq.Username.ValueString()
		userMatcher = func(user enterprise.IUser) bool {
			return strings.EqualFold(emailFilter, user.Username())
		}
	} else if !uq.EnterpriseUserId.IsNull() && !uq.EnterpriseUserId.IsUnknown() {
		var userIdFilter = uq.EnterpriseUserId.ValueInt64()
		userMatcher = func(user enterprise.IUser) bool {
			return userIdFilter == user.EnterpriseUserId()
		}
	}

	if userMatcher == nil {
		resp.Diagnostics.AddError(
			"Search criteria is not provided for \"user\" data source",
			fmt.Sprintf("Search criteria is not provided for \"%s\" data source", "kepr_user"),
		)
		return
	}

	var user enterprise.IUser
	d.users.GetAllEntities(func(u enterprise.IUser) bool {
		if userMatcher(u) {
			user = u
			return false
		}
		return true
	})

	if user == nil {
		resp.Diagnostics.AddError(
			"User not found",
			fmt.Sprintf("Cannot find a user according to the provided criteria"),
		)
		return
	}

	var state = uq
	var tm = &state
	tm.fromKeeper(user)
	if !state.IncludeTeams.IsNull() && state.IncludeTeams.ValueBool() {
		d.teamUsers.GetLinksByObject(user.EnterpriseUserId(), func(tu enterprise.ITeamUser) bool {
			var t = d.teams.GetEntity(tu.TeamUid())
			if t != nil {
				var tsm = new(teamShortModel)
				tsm.fromKeeper(t)
				state.Teams = append(state.Teams, tsm)
			}
			return true
		})
	}

	if !state.IncludeRoles.IsNull() && state.IncludeRoles.ValueBool() {
		d.roleUsers.GetLinksByObject(user.EnterpriseUserId(), func(ru enterprise.IRoleUser) bool {
			var r = d.roles.GetEntity(ru.RoleId())
			if r != nil {
				var rsm = new(roleShortModel)
				var isAdmin bool
				d.managedNodes.GetLinksBySubject(ru.RoleId(), func(_ enterprise.IManagedNode) bool {
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

func (d *userDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("enterprise_user_id"),
			path.MatchRoot("username"),
		),
	}
}
