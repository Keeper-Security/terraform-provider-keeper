package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/api"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"strings"
)

func newTeamResource() resource.Resource {
	return &teamResource{}
}

type teamResource struct {
	management enterprise.IEnterpriseManagement
}

func (r *teamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *teamResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"team_uid": schema.StringAttribute{
				Computed:    true,
				Description: "Team UID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				}},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Team Name",
			},
			"node_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Team NodeID",
			},
			"restrict_edit": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Restrict Edit flag",
				Default:     booldefault.StaticBool(false),
			},
			"restrict_share": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Restrict Share flag",
				Default:     booldefault.StaticBool(false),
			},
			"restrict_view": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Restrict View flag",
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *teamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.management = mgmt
}

func (r *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var teams = r.management.EnterpriseData().Teams()
	if state.NodeId.IsNull() || state.NodeId.ValueInt64() == 0 {
		state.NodeId = types.Int64Value(r.management.EnterpriseData().RootNode().NodeId())
	}

	var team enterprise.ITeam
	if state.TeamUid.IsNull() {
		var teamName = state.Name.ValueString()
		var nodeId = state.NodeId.ValueInt64()

		teams.GetAllEntities(func(t enterprise.ITeam) bool {
			if t.NodeId() == nodeId && strings.EqualFold(t.Name(), teamName) {
				team = t
				return false
			}
			return true
		})
	} else {
		var teamUid = state.TeamUid.ValueString()
		if len(teamUid) > 0 {
			team = teams.GetEntity(teamUid)
		}
	}

	if team != nil {
		state.fromKeeper(team)
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r *teamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("team_uid"), req, resp)
}

func (r *teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan teamModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var teamName = plan.Name.ValueString()
	if len(teamName) == 0 {
		resp.Diagnostics.AddError(
			"Create Team: Team name cannot be empty",
			"Team name cannot be empty",
		)
		return
	}

	if plan.NodeId.IsNull() || plan.NodeId.ValueInt64() == 0 {
		plan.NodeId = types.Int64Value(r.management.EnterpriseData().RootNode().NodeId())
	}

	var team enterprise.ITeam
	var nodeId = plan.NodeId.ValueInt64()
	var teams = r.management.EnterpriseData().Teams()
	teams.GetAllEntities(func(t enterprise.ITeam) bool {
		if t.NodeId() == nodeId && strings.EqualFold(t.Name(), teamName) {
			team = t
			return false
		}
		return true
	})
	var te enterprise.ITeamEdit
	var toAdd []enterprise.ITeam
	var toUpdate []enterprise.ITeam
	if team == nil {
		te = enterprise.NewTeam(api.Base64UrlEncode(api.GenerateUid()))
		toAdd = append(toAdd, te)
	} else {
		te = enterprise.CloneTeam(team)
		toUpdate = append(toUpdate, te)
	}
	plan.toKeeper(te)

	var errs = r.management.ModifyTeams(toAdd, toUpdate, nil)
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Team %s", teamName),
			fmt.Sprintf("Error: %s", er),
		)
	}
	if !resp.Diagnostics.HasError() {
		plan.TeamUid = types.StringValue(te.TeamUid())
		plan.NodeId = types.Int64Value(te.NodeId())
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan teamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state teamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.NodeId.IsNull() || plan.NodeId.ValueInt64() == 0 {
		plan.NodeId = state.NodeId
	}
	if plan.NodeId.IsNull() || plan.NodeId.ValueInt64() == 0 {
		plan.NodeId = types.Int64Value(r.management.EnterpriseData().RootNode().NodeId())
	}
	if plan.TeamUid.IsNull() {
		plan.TeamUid = state.TeamUid
	}
	if plan.TeamUid.IsNull() {
		resp.Diagnostics.AddError(
			"Update Team error: cannot resolve team UID",
			"Could not resolve team UID",
		)
		return
	}

	var teamUid = plan.TeamUid.ValueString()
	var t = enterprise.NewTeam(teamUid)
	plan.toKeeper(t)

	errs := r.management.ModifyTeams(nil, []enterprise.ITeam{t}, nil)
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Update Team %s", plan.Name.ValueString()),
			"Error occurred while updating a team: "+er.Error(),
		)
	}
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}
	var forDelete = []string{state.TeamUid.ValueString()}
	var errs = r.management.ModifyTeams(nil, nil, forDelete)
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete Team \"%s\" error: %s", state.Name.ValueString(), er),
			"Error occurred while creating a team",
		)
	}
}
