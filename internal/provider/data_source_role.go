package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"github.com/keeper-security/keeper-sdk-golang/vault"
	"strings"
	"terraform-provider-keeper/internal/model"
)

var (
	_ datasource.DataSourceWithConfigValidators = &roleDataSource{}
)

type roleDataSourceModel struct {
	RoleId         types.Int64                        `tfsdk:"role_id"`
	Name           types.String                       `tfsdk:"name"`
	Node           *model.NodeShortModel              `tfsdk:"node"`
	VisibleBelow   types.Bool                         `tfsdk:"visible_below"`
	NewUserInherit types.Bool                         `tfsdk:"new_user_inherit"`
	IsAdmin        types.Bool                         `tfsdk:"is_admin"`
	ManagedNodes   []*model.ManagedNodeModel          `tfsdk:"managed_nodes"`
	IncludeMembers types.Bool                         `tfsdk:"include_members"`
	Users          []*model.UserShortModel            `tfsdk:"users"`
	Teams          []*model.TeamShortModel            `tfsdk:"teams"`
	Enforcements   *model.EnforcementsDataSourceModel `tfsdk:"enforcements"`
}

func (rm *roleDataSourceModel) fromKeeper(role enterprise.IRole, isAdmin bool, node enterprise.INode) {
	rm.RoleId = types.Int64Value(role.RoleId())
	rm.Name = types.StringValue(role.Name())
	rm.VisibleBelow = types.BoolValue(role.VisibleBelow())
	rm.NewUserInherit = types.BoolValue(role.NewUserInherit())
	rm.IsAdmin = types.BoolValue(isAdmin)
	if node != nil {
		rm.Node = new(model.NodeShortModel)
		rm.Node.FromKeeper(node)
	}
}

type roleDataSource struct {
	roles            enterprise.IEnterpriseEntity[enterprise.IRole, int64]
	roleUsers        enterprise.IEnterpriseLink[enterprise.IRoleUser, int64, int64]
	users            enterprise.IEnterpriseEntity[enterprise.IUser, int64]
	nodes            enterprise.IEnterpriseEntity[enterprise.INode, int64]
	managedNodes     enterprise.IEnterpriseLink[enterprise.IManagedNode, int64, int64]
	rolePrivileges   enterprise.IEnterpriseLink[enterprise.IRolePrivilege, int64, int64]
	roleEnforcements enterprise.IEnterpriseLink[enterprise.IRoleEnforcement, int64, string]
	roleTeams        enterprise.IEnterpriseLink[enterprise.IRoleTeam, int64, string]
	teams            enterprise.IEnterpriseEntity[enterprise.ITeam, string]
	recordTypes      enterprise.IEnterpriseEntity[vault.IRecordType, string]
}

func newRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}
func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the data source.
func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsAttributes map[string]schema.Attribute
	enforcementsAttributes, diags = model.EnforcementsDataSourceAttributes()
	resp.Diagnostics.Append(diags...)

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"role_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Role ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Role Name",
			},
			"node": schema.SingleNestedAttribute{
				Attributes:  model.NodeShortSchemaAttributes,
				Computed:    true,
				Description: "Role Node",
			},
			"visible_below": schema.BoolAttribute{
				Computed:    true,
				Description: "Visible Below",
			},
			"new_user_inherit": schema.BoolAttribute{
				Computed:    true,
				Description: "New User Inherit",
			},
			"is_admin": schema.BoolAttribute{
				Computed:    true,
				Description: "Is Administrative Role",
			},
			"include_members": schema.BoolAttribute{
				Optional:    true,
				Description: "Include role members",
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: model.UserShortSchemaAttributes,
				},
			},
			"teams": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: model.TeamShortSchemaAttributes,
				},
			},
			"managed_nodes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: model.ManagedNodeDataSourceSchemaAttributes,
				},
			},
			"enforcements": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Role enforcemnts",
				Attributes:  enforcementsAttributes,
			},
		},
	}
}

func (d *roleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.roles = ed.Roles()
		d.roleUsers = ed.RoleUsers()
		d.users = ed.Users()
		d.managedNodes = ed.ManagedNodes()
		d.rolePrivileges = ed.RolePrivileges()
		d.roleEnforcements = ed.RoleEnforcements()
		d.nodes = ed.Nodes()
		d.roleTeams = ed.RoleTeams()
		d.teams = ed.Teams()
		d.recordTypes = ed.RecordTypes()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseData\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq roleDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var m func(node enterprise.IRole) bool
	if !rq.Name.IsNull() && !rq.Name.IsUnknown() {
		m = func(role enterprise.IRole) bool {
			return strings.EqualFold(rq.Name.ValueString(), role.Name())
		}
	} else if !rq.RoleId.IsNull() && !rq.RoleId.IsUnknown() {
		m = func(role enterprise.IRole) bool {
			return role.RoleId() == rq.RoleId.ValueInt64()
		}
	}

	if m == nil {
		resp.Diagnostics.AddError(
			"Search criteria is not provided for \"role\" data source",
			"Search criteria is not provided for \"role\" data source",
		)
		return
	}

	var role enterprise.IRole
	d.roles.GetAllEntities(func(r enterprise.IRole) bool {
		if m(r) {
			role = r
			return false
		}
		return true
	})

	if role == nil {
		resp.Diagnostics.AddError("Role not found", "Cannot find a role according to the provided criteria")
		return
	}

	var state = rq
	var isAdmin = false
	d.managedNodes.GetLinksBySubject(role.RoleId(), func(x enterprise.IManagedNode) bool {
		isAdmin = true
		return false
	})
	var node = d.nodes.GetEntity(role.NodeId())
	state.fromKeeper(role, isAdmin, node)
	d.managedNodes.GetLinksBySubject(role.RoleId(), func(mn enterprise.IManagedNode) bool {
		var mnm = new(model.ManagedNodeModel)
		mnm.NodeId = types.Int64Value(mn.ManagedNodeId())
		var node = d.nodes.GetEntity(mn.ManagedNodeId())
		if node != nil {
			mnm.Name = types.StringValue(node.Name())
		}
		mnm.CascadeNodeManagement = types.BoolValue(mn.CascadeNodeManagement())
		var privileges = d.rolePrivileges.GetLink(mn.RoleId(), mn.ManagedNodeId())
		if privileges != nil {
			mnm.Privileges.FromKeeper(privileges)
		}
		state.ManagedNodes = append(state.ManagedNodes, mnm)
		return true
	})
	if !rq.IncludeMembers.IsNull() && rq.IncludeMembers.ValueBool() {
		d.roleUsers.GetLinksBySubject(role.RoleId(), func(ru enterprise.IRoleUser) bool {
			var user = d.users.GetEntity(ru.EnterpriseUserId())
			if user != nil {
				var usm = new(model.UserShortModel)
				usm.FromKeeper(user)
				state.Users = append(state.Users, usm)
			}
			return true
		})
		d.roleTeams.GetLinksBySubject(role.RoleId(), func(rt enterprise.IRoleTeam) bool {
			var team = d.teams.GetEntity(rt.TeamUid())
			if team != nil {
				var tsm = new(model.TeamShortModel)
				tsm.FromKeeper(team)
				state.Teams = append(state.Teams, tsm)
			}
			return true
		})
	}
	enforcements := make(map[string]string)
	d.roleEnforcements.GetLinksBySubject(role.RoleId(), func(re enterprise.IRoleEnforcement) bool {
		if re != nil {
			enforcements[re.EnforcementType()] = re.Value()
		}
		return true
	})
	state.Enforcements = new(model.EnforcementsDataSourceModel)
	var rts []vault.IRecordType
	if d.recordTypes != nil {
		d.recordTypes.GetAllEntities(func(x vault.IRecordType) bool {
			rts = append(rts, x)
			return true
		})
	}
	state.Enforcements.FromKeeper(enforcements, rts)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *roleDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("role_id"),
			path.MatchRoot("name"),
		),
	}
}
