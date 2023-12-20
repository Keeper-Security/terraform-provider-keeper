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

func newTeamMembershipResource() resource.Resource {
	return &teamMembershipResource{}
}

type teamMembershipResourceModel struct {
	TeamUid types.String `tfsdk:"team_uid"`
	Users   types.Set    `tfsdk:"users"`
}

type teamMembershipResource struct {
	management enterprise.IEnterpriseManagement
}

func (tmr *teamMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_membership"
}

func (tmr *teamMembershipResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"team_uid": schema.StringAttribute{
				Required:    true,
				Description: "Team UID",
			},
			"users": schema.SetAttribute{
				Required:    true,
				ElementType: types.Int64Type,
				Description: "The list of user IDs",
			},
		},
	}
}

func (tmr *teamMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	tmr.management = mgmt
}

func (tmr *teamMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("team_uid"), req, resp)
}

func (tmr *teamMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamMembershipResourceModel
	var diags = req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var teams = tmr.management.EnterpriseData().Teams()
	var teamUid = state.TeamUid.ValueString()
	var team = teams.GetEntity(teamUid)
	if team == nil {
		resp.Diagnostics.AddError(
			"Team not found",
			fmt.Sprintf("Team UID \"%s\" does not exist", teamUid),
		)
		return
	}

	var members []attr.Value
	tmr.management.EnterpriseData().TeamUsers().GetLinksBySubject(teamUid, func(tu enterprise.ITeamUser) bool {
		members = append(members, types.Int64Value(tu.EnterpriseUserId()))
		return true
	})
	state.TeamUid = types.StringValue(teamUid)
	state.Users, diags = types.SetValue(types.Int64Type, members)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (tmr *teamMembershipResource) applyMembership(plan teamMembershipResourceModel) (diags diag.Diagnostics) {
	var teams = tmr.management.EnterpriseData().Teams()
	var teamUid = plan.TeamUid.ValueString()
	var team = teams.GetEntity(teamUid)
	if team == nil {
		diags.AddError(
			"Team not found",
			fmt.Sprintf("Team UID \"%s\" does not exist", teamUid),
		)
		return
	}
	var origUsers = api.NewSet[int64]()
	tmr.management.EnterpriseData().TeamUsers().GetLinksBySubject(teamUid, func(tu enterprise.ITeamUser) bool {
		origUsers.Add(tu.EnterpriseUserId())
		return true
	})
	var planUsers = api.NewSet[int64]()
	if !plan.Users.IsNull() {
		for _, k := range plan.Users.Elements() {
			if tv, ok := k.(types.Int64); ok {
				planUsers.Add(tv.ValueInt64())
			}
		}
	}

	var toAdd []enterprise.ITeamUser
	var us = planUsers.Copy()
	us.Difference(origUsers.ToArray())
	if len(us) > 0 {
		us.Enumerate(func(u int64) bool {
			toAdd = append(toAdd, enterprise.NewTeamUser(teamUid, u))
			return true
		})
	}
	us = origUsers.Copy()
	us.Difference(planUsers.ToArray())
	var toDelete []enterprise.ITeamUser
	if len(us) > 0 {
		us.Enumerate(func(u int64) bool {
			toDelete = append(toDelete, enterprise.NewTeamUser(teamUid, u))
			return true
		})
	}

	var errs = tmr.management.ModifyTeamUsers(toAdd, toDelete)
	for _, er := range errs {
		diags.AddError(
			fmt.Sprintf("Team UID \"%s\" Membership", teamUid),
			fmt.Sprintf("Error: %s", er),
		)
	}
	return
}

func (tmr *teamMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan teamMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(tmr.applyMembership(plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (tmr *teamMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan teamMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(tmr.applyMembership(plan)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (tmr *teamMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamMembershipResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var plan teamMembershipResourceModel
	plan.TeamUid = state.TeamUid
	resp.Diagnostics.Append(tmr.applyMembership(plan)...)
}
