package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/api"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
)

func newRoleMembershipResource() resource.Resource {
	return &roleMembershipResource{}
}

type roleMembershipResourceModel struct {
	RoleId types.Int64 `tfsdk:"role_id"`
	Users  types.Set   `tfsdk:"users"`
	Teams  types.Set   `tfsdk:"teams"`
}

type roleMembershipResource struct {
	management enterprise.IEnterpriseManagement
}

func (rmr *roleMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_membership"
}

func (rmr *roleMembershipResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"role_id": schema.Int64Attribute{
				Required:    true,
				Description: "Role ID",
			},
			"users": schema.SetAttribute{
				Required:    true,
				ElementType: types.Int64Type,
				Description: "The list of user IDs",
			},
			"teams": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "The list of team UIDs",
			},
		},
	}
}

func (rmr *roleMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	rmr.management = mgmt
}

func (rmr *roleMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("role_id"), req, resp)
}

func (rmr *roleMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleMembershipResourceModel
	var diags = req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var roles = rmr.management.EnterpriseData().Roles()
	var roleId = state.RoleId.ValueInt64()
	var role = roles.GetEntity(roleId)
	if role == nil {
		resp.Diagnostics.AddError(
			"Role not found",
			fmt.Sprintf("Role ID \"%d\" does not exist", roleId),
		)
		return
	}

	state.RoleId = types.Int64Value(roleId)
	var members []attr.Value
	rmr.management.EnterpriseData().RoleUsers().GetLinksBySubject(roleId, func(ru enterprise.IRoleUser) bool {
		members = append(members, types.Int64Value(ru.EnterpriseUserId()))
		return true
	})
	state.Users, diags = types.SetValue(types.Int64Type, members)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	members = nil
	rmr.management.EnterpriseData().RoleTeams().GetLinksBySubject(roleId, func(rt enterprise.IRoleTeam) bool {
		members = append(members, types.StringValue(rt.TeamUid()))
		return true
	})
	state.Teams, diags = types.SetValue(types.StringType, members)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (rmr *roleMembershipResource) applyMembership(plan roleMembershipResourceModel) (diags diag.Diagnostics) {
	var roles = rmr.management.EnterpriseData().Roles()
	var roleId = plan.RoleId.ValueInt64()
	var role = roles.GetEntity(roleId)
	if role == nil {
		diags.AddError(
			"Role not found",
			fmt.Sprintf("Role ID \"%d\" does not exist", roleId),
		)
		return
	}

	var userId int64
	var origUsers = api.NewSet[int64]()
	rmr.management.EnterpriseData().RoleUsers().GetLinksBySubject(roleId, func(ru enterprise.IRoleUser) bool {
		origUsers.Add(ru.EnterpriseUserId())
		return true
	})

	var planUsers = api.NewSet[int64]()
	if !plan.Users.IsNull() {
		for _, k := range plan.Users.Elements() {
			if tv, ok := k.(types.Int64); ok {
				userId = tv.ValueInt64()
				planUsers.Add(userId)
			}
		}
	}

	var us = planUsers.Copy()
	us.Difference(origUsers.ToArray())
	var roleUserToAdd []enterprise.IRoleUser
	if len(us) > 0 {
		us.Enumerate(func(u int64) bool {
			roleUserToAdd = append(roleUserToAdd, enterprise.NewRoleUser(roleId, u))
			return true
		})
	}

	us = origUsers.Copy()
	us.Difference(planUsers.ToArray())
	var roleUserToDelete []enterprise.IRoleUser
	if len(us) > 0 {
		us.Enumerate(func(u int64) bool {
			roleUserToDelete = append(roleUserToDelete, enterprise.NewRoleUser(roleId, u))
			return true
		})
	}

	var errs = rmr.management.ModifyRoleUsers(roleUserToAdd, roleUserToDelete)
	for _, er := range errs {
		diags.AddError(
			fmt.Sprintf("Role ID \"%d\" Membership", roleId),
			fmt.Sprintf("Error: %s", er),
		)
	}

	var origTeams = api.NewSet[string]()
	rmr.management.EnterpriseData().RoleTeams().GetLinksBySubject(roleId, func(rt enterprise.IRoleTeam) bool {
		origTeams.Add(rt.TeamUid())
		return true
	})
	var planTeams = api.NewSet[string]()
	if !plan.Teams.IsNull() {
		for _, k := range plan.Teams.Elements() {
			if tv, ok := k.(types.String); ok {
				planTeams.Add(tv.ValueString())
			}
		}
	}
	var ts = planTeams.Copy()
	ts.Difference(origTeams.ToArray())
	var roleTeamToAdd []enterprise.IRoleTeam
	if len(ts) > 0 {
		ts.Enumerate(func(t string) bool {
			roleTeamToAdd = append(roleTeamToAdd, enterprise.NewRoleTeam(roleId, t))
			return true
		})
	}
	ts = origTeams.Copy()
	ts.Difference(planTeams.ToArray())
	var roleTeamToDelete []enterprise.IRoleTeam
	if len(ts) > 0 {
		ts.Enumerate(func(t string) bool {
			roleTeamToDelete = append(roleTeamToDelete, enterprise.NewRoleTeam(roleId, t))
			return true
		})
	}

	errs = rmr.management.ModifyRoleTeams(roleTeamToAdd, roleTeamToDelete)
	for _, er := range errs {
		diags.AddError(
			fmt.Sprintf("Role ID \"%d\" Membership", roleId),
			fmt.Sprintf("Error: %s", er),
		)
	}
	return
}

func (rmr *roleMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(rmr.applyMembership(plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (rmr *roleMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(rmr.applyMembership(plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (rmr *roleMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleMembershipResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var plan roleMembershipResourceModel
	plan.RoleId = state.RoleId
	resp.Diagnostics.Append(rmr.applyMembership(plan)...)
}
