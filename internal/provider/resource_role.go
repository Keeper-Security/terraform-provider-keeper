package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"strconv"
	"strings"
)

func newRoleResource() resource.Resource {
	return &roleResource{}
}

type roleResourceModel struct {
	RoleId         types.Int64  `tfsdk:"role_id"`
	Name           types.String `tfsdk:"name"`
	NodeId         types.Int64  `tfsdk:"node_id"`
	VisibleBelow   types.Bool   `tfsdk:"visible_below"`
	NewUserInherit types.Bool   `tfsdk:"new_user_inherit"`
}

func (r *roleResourceModel) fromKeeper(role enterprise.IRole) {
	r.NodeId = types.Int64Value(role.RoleId())
	r.Name = types.StringValue(role.Name())
	r.NodeId = types.Int64Value(role.NodeId())
	r.VisibleBelow = types.BoolValue(role.VisibleBelow())
	r.NewUserInherit = types.BoolValue(role.NewUserInherit())
}

func (r *roleResourceModel) toKeeper(role enterprise.IRoleEdit) {
	role.SetName(r.Name.ValueString())
	if !r.NodeId.IsNull() {
		role.SetNodeId(r.NodeId.ValueInt64())
	}
	if !r.VisibleBelow.IsNull() {
		role.SetVisibleBelow(r.VisibleBelow.ValueBool())
	} else {
		role.SetVisibleBelow(true)
	}
	if !r.NewUserInherit.IsNull() {
		role.SetNewUserInherit(r.NewUserInherit.ValueBool())
	} else {
		role.SetNewUserInherit(false)
	}
}

type roleResource struct {
	management enterprise.IEnterpriseManagement
}

func (nr *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (nr *roleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"role_id": schema.Int64Attribute{
				Computed:      true,
				Description:   "Role ID",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Role Name",
			},
			"node_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Node ID",
			},
			"visible_below": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Visible Below",
			},
			"new_user_inherit": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "New User Inherit",
			},
		},
	}
}

func (nr *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	mgmt, ok := req.ProviderData.(enterprise.IEnterpriseManagement)
	if ok {
		nr.management = mgmt
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected IEnterpriseManagement, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (nr *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var roleId int
	var err error

	if roleId, err = strconv.Atoi(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid role ID value.",
			fmt.Sprintf("Interger value expected. Got \"%s\"", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_id"), int64(roleId))...)
}

func (nr *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var roles = nr.management.EnterpriseData().Roles()
	var role enterprise.IRole
	if state.RoleId.IsNull() {
		var nodeId int64 = 0
		if !state.NodeId.IsNull() {
			nodeId = state.NodeId.ValueInt64()
		}
		var roleName = state.Name.ValueString()
		roles.GetAllEntities(func(r enterprise.IRole) bool {
			if strings.EqualFold(r.Name(), roleName) {
				if nodeId > 0 {
					if nodeId == r.NodeId() {
						role = r
						return false
					}
				} else {
					if role == nil {
						role = r
					} else {
						role = nil
						resp.Diagnostics.AddError(fmt.Sprintf("Role name \"%s\" is not unique", roleName),
							"Please specify role's node ID or use Role ID to reference a role resource")
						return false
					}
				}
			}
			return true
		})
	} else {
		var roleId = state.RoleId.ValueInt64()
		if roleId > 0 {
			role = roles.GetEntity(roleId)
		}
	}

	if !resp.Diagnostics.HasError() {
		if role != nil {
			state.fromKeeper(role)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		} else {
			resp.State.RemoveResource(ctx)
		}
	}
}

func (nr *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roleName = plan.Name.ValueString()
	if len(roleName) == 0 {
		resp.Diagnostics.AddError(
			"Create Role: Role name cannot be empty",
			"Role name cannot be empty",
		)
		return
	}
	if plan.NodeId.IsNull() || plan.NodeId.ValueInt64() == 0 {
		plan.NodeId = types.Int64Value(nr.management.EnterpriseData().RootNode().NodeId())
	}
	var nodeId = plan.NodeId.ValueInt64()
	var role enterprise.IRole
	nr.management.EnterpriseData().Roles().GetAllEntities(func(r enterprise.IRole) bool {
		if nodeId == r.NodeId() && strings.EqualFold(r.Name(), roleName) {
			role = r
			return false
		}
		return true
	})

	var re enterprise.IRoleEdit
	var toAdd []enterprise.IRole
	var toUpdate []enterprise.IRole
	if role == nil {
		var roleId int64
		var err error
		if roleId, err = nr.management.GetEnterpriseId(); err != nil {
			resp.Diagnostics.AddError(
				"Create Role: Error getting enterprise ID",
				fmt.Sprintf("Error getting enterprise ID: %s", err.Error()),
			)
			return
		}
		re = enterprise.NewRole(roleId)
		toAdd = append(toAdd, re)
	} else {
		re = enterprise.CloneRole(role)
		toUpdate = append(toUpdate, re)
	}
	plan.toKeeper(re)

	var errs = nr.management.ModifyRoles(toAdd, toUpdate, nil)
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Role %s", roleName),
			fmt.Sprintf("Error: %s", er),
		)
	}
	if !resp.Diagnostics.HasError() {
		plan.RoleId = types.Int64Value(re.RoleId())
		plan.NodeId = types.Int64Value(re.NodeId())
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (nr *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.NodeId.IsNull() || plan.NodeId.ValueInt64() == 0 {
		plan.NodeId = state.NodeId
	}
	if plan.NodeId.IsNull() || plan.NodeId.ValueInt64() == 0 {
		plan.NodeId = types.Int64Value(nr.management.EnterpriseData().RootNode().NodeId())
	}
	if plan.RoleId.IsNull() {
		plan.RoleId = state.RoleId
	}
	if plan.NodeId.IsNull() {
		resp.Diagnostics.AddError(
			"Update Role error: cannot resolve role ID",
			"Could not resolve role ID",
		)
		return
	}

	var roleId = plan.RoleId.ValueInt64()
	var role = nr.management.EnterpriseData().Roles().GetEntity(roleId)
	if role != nil {
		var r = enterprise.CloneRole(role)
		plan.toKeeper(r)

		errs := nr.management.ModifyRoles(nil, []enterprise.IRole{r}, nil)
		for _, er := range errs {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Update Role %s", plan.Name.ValueString()),
				"Error occurred while updating a role: "+er.Error(),
			)
		}
	} else {
		resp.Diagnostics.AddError(fmt.Sprintf("Update Role error: role ID \"%d\" not found", roleId),
			"Role not found")
	}
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (nr *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var errs = nr.management.ModifyRoles(nil, nil, []int64{state.RoleId.ValueInt64()})
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete Role \"%s\" error: %s", state.Name.ValueString(), er),
			"Error occurred while deleting a role",
		)
	}
}
