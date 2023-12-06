package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/api"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"strconv"
	"strings"
	"terraform-provider-kepr/internal/model"
)

func newManagedNodeResource() resource.Resource {
	return &managedNodeResource{}
}

var privilegesResourceAttributes = map[string]schema.Attribute{
	"manage_nodes": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Nodes",
	},
	"manage_users": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Users",
	},
	"manage_teams": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Teams",
	},
	"manage_roles": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Roles",
	},
	"manage_reports": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Reporting and Alerts",
	},
	"manage_sso": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Bridge/SSO",
	},
	"device_approval": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Perform Device Approvals",
	},
	"manage_record_types": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Record Types in Vault",
		MarkdownDescription: "This permission allows the admin rights to create, edit, or delete Record Types " +
			"which have pre-defined fields. Record Types appear during creation of records in the user's vault.",
	},
	"share_admin": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Share Admin",
	},
	"run_compliance_reports": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Run Compliance Reports",
	},
	"transfer_account": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Transfer Account",
	},
	"manage_companies": schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Manage Companies",
	},
}

type managedNodeResourceModel struct {
	RoleId                types.Int64                     `tfsdk:"role_id"`
	NodeId                types.Int64                     `tfsdk:"node_id"`
	CascadeNodeManagement types.Bool                      `tfsdk:"cascade_node_management"`
	Privileges            *model.PrivilegeDataSourceModel `tfsdk:"privileges"`
}

type managedNodeResource struct {
	management enterprise.IEnterpriseManagement
}

func (mnr *managedNodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_node"
}

func (mnr *managedNodeResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"role_id": schema.Int64Attribute{
				Required:    true,
				Description: "Role ID",
			},
			"managed_node_id": schema.Int64Attribute{
				Required:    true,
				Description: "Managed Node ID",
			},
			"cascade_node_management": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "The list of user IDs",
			},
			"privileges": schema.SingleNestedAttribute{
				Attributes:  privilegesResourceAttributes,
				Optional:    true,
				Computed:    true,
				Description: "Privileges",
			},
		},
	}
}

func (mnr *managedNodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	mgmt, ok := req.ProviderData.(enterprise.IEnterpriseManagement)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected IEnterpriseManagement, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	mnr.management = mgmt
}

func (mnr *managedNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var idParts = api.SliceSelect(strings.Split(req.ID, ","), func(x string) (id int64) {
		x = strings.TrimSpace(x)
		if i, er1 := strconv.Atoi(x); er1 == nil {
			id = int64(i)
		}
		return
	})

	if len(idParts) != 2 || idParts[0] == 0 || idParts[1] == 0 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: role_id,node_id. Integers. Got: %q", req.ID),
		)
		return
	}
	if r := mnr.management.EnterpriseData().Roles().GetEntity(idParts[0]); r == nil {
		resp.Diagnostics.AddError(
			"Import Managed Node: role not found",
			fmt.Sprintf("Role not found: RoleID=\"%d\"", idParts[0]),
		)
		return
	}
	if n := mnr.management.EnterpriseData().Nodes().GetEntity(idParts[1]); n == nil {
		resp.Diagnostics.AddError(
			"Import Managed Node: node not found",
			fmt.Sprintf("Node not found: NodeID=\"%d\"", idParts[1]),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("node_id"), idParts[1])...)
}

func (mnr *managedNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state managedNodeResourceModel
	var diags = req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var roleId = state.RoleId.ValueInt64()
	var nodeId = state.NodeId.ValueInt64()

	var mn = mnr.management.EnterpriseData().ManagedNodes().GetLink(roleId, nodeId)
	if mn == nil {
		resp.Diagnostics.AddError(
			"Managed Node not found",
			fmt.Sprintf("Managed Node: RoleID=\"%d\"; NodeId=\"%d\" does not exist", roleId, nodeId),
		)
		return
	}
	state.CascadeNodeManagement = types.BoolValue(mn.CascadeNodeManagement())
	var p = mnr.management.EnterpriseData().RolePrivileges().GetLink(roleId, nodeId)
	state.Privileges = new(model.PrivilegeDataSourceModel)
	if p != nil {
		state.Privileges.FromKeeper(p.Privileges())
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (mnr *managedNodeResource) applyManagedNode(model managedNodeResourceModel) (diags diag.Diagnostics) {
	var roleId = model.RoleId.ValueInt64()
	var nodeId = model.NodeId.ValueInt64()
	var cascade = true
	if !model.CascadeNodeManagement.IsNull() {
		cascade = model.CascadeNodeManagement.ValueBool()
	}

	var addManagedNodes []enterprise.IManagedNode
	var updateManagedNodes []enterprise.IManagedNode
	var rolePrivileges []enterprise.IRolePrivilege
	var mn = mnr.management.EnterpriseData().ManagedNodes().GetLink(roleId, nodeId)
	if mn == nil {
		if r := mnr.management.EnterpriseData().Roles().GetEntity(roleId); r == nil {
			diags.AddError("Create Managed Node: role not found",
				fmt.Sprintf("Role not found: RoleID=\"%d\"", roleId))
			return
		}
		if n := mnr.management.EnterpriseData().Roles().GetEntity(nodeId); n == nil {
			diags.AddError("Create Managed Node: node not found",
				fmt.Sprintf("Node not found: NodeID=\"%d\"", nodeId))
			return
		}
		var ok = true
		mnr.management.EnterpriseData().RoleUsers().GetLinksBySubject(roleId, func(x enterprise.IRoleUser) bool {
			if u := mnr.management.EnterpriseData().Users().GetEntity(x.EnterpriseUserId()); u != nil {
				if u.Status() != enterprise.UserStatus_Active {
					ok = false
				}
			}
			return ok
		})
		if !ok {
			diags.AddError("Create Managed Node: role cannot have pending users",
				fmt.Sprintf("Administrative Role cannot contain pending users: RoleID=\"%d\"", roleId))
			return
		}
		ok = true
		mnr.management.EnterpriseData().RoleTeams().GetLinksBySubject(roleId, func(x enterprise.IRoleTeam) bool {
			ok = false
			return ok
		})
		if !ok {
			diags.AddError("Create Managed Node: role cannot have teams",
				fmt.Sprintf("Administrative Role cannot contain teams: RoleID=\"%d\"", roleId))
			return
		}
		var imn = enterprise.NewManagedNode(roleId, nodeId)
		imn.SetCascadeNodeManagement(cascade)
		addManagedNodes = append(addManagedNodes, imn)
		if model.Privileges != nil {
			var rp = enterprise.NewRolePrivilege(roleId, nodeId)
			for _, e := range model.Privileges.ToKeeper() {
				rp.SetPrivilege(e)
			}
			rolePrivileges = append(rolePrivileges, rp)
		}
	} else {
		if cascade != mn.CascadeNodeManagement() {
			var imn = enterprise.NewManagedNode(roleId, nodeId)
			imn.SetCascadeNodeManagement(cascade)
			updateManagedNodes = append(updateManagedNodes, imn)
		}
	}
	var errs = mnr.management.ModifyManagedNodes(addManagedNodes, updateManagedNodes, nil)
	for _, er := range errs {
		diags.AddError(
			fmt.Sprintf("Create Managed Node error: RoleID=\"%d\"; NodeID=\"%d\"", roleId, nodeId),
			fmt.Sprintf("Error: %s", er),
		)
	}
	if diags.HasError() {
		return
	}

	if model.Privileges != nil {
		var ps = model.Privileges.ToKeeper()
		var modified = true
		if orp := mnr.management.EnterpriseData().RolePrivileges().GetLink(roleId, nodeId); orp != nil {
			var s = api.NewSet[string]()
			s.Union(ps)
			s.Difference(orp.Privileges())
			if len(s) == 0 {
				s.Union(orp.Privileges())
				s.Difference(ps)
				if len(s) == 0 {
					modified = false
				}
			}
		}
		if modified {
			var rp = enterprise.NewRolePrivilege(roleId, nodeId)
			for _, e := range ps {
				rp.SetPrivilege(e)
			}
			errs = mnr.management.ModifyRolePrivileges([]enterprise.IRolePrivilege{rp})
			for _, er := range errs {
				diags.AddError(
					fmt.Sprintf("Assign Managed Node privileges: RoleID=\"%d\"; NodeID=\"%d\"", roleId, nodeId),
					fmt.Sprintf("Error: %s", er),
				)
			}
		}
	}
	return
}

func (mnr *managedNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan managedNodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(mnr.applyManagedNode(plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (mnr *managedNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan managedNodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(mnr.applyManagedNode(plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (mnr *managedNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleMembershipResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var plan roleMembershipResourceModel
	plan.RoleId = state.RoleId
	//resp.Diagnostics.Append(rmr.applyMembership(plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	}
}
