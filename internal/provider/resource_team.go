package provider

import (
	"bytes"
	"context"
	"fmt"
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

func NewTeamResource() resource.Resource {
	return &teamResource{}
}

type teamResource struct {
	enterprise enterprise.IEnterpriseLoader
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

	ent, ok := req.ProviderData.(enterprise.IEnterpriseLoader)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected IEnterpriseLoader, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.enterprise = ent
}

func (r *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}
	var teams = r.enterprise.EnterpriseData().Teams().GetData()
	if state.NodeId.IsNull() || state.NodeId.ValueInt64() == 0 {
		state.NodeId = types.Int64Value(r.enterprise.EnterpriseData().GetRootNode().NodeId)
	}
	var team *enterprise.Team
	if state.TeamUid.IsNull() {
		var teamName = state.Name.ValueString()
		var nodeId = state.NodeId.ValueInt64()
		for _, t := range teams {
			if t.NodeId == nodeId && strings.EqualFold(t.Name, teamName) {
				team = t
				break
			}
		}
	} else {
		var tUid = state.TeamUid.ValueString()
		if len(tUid) > 0 {
			var teamUid = api.Base64UrlDecode(tUid)
			for _, t := range teams {
				if bytes.Compare(t.TeamUid, teamUid) == 0 {
					team = t
					break
				}
			}
		}
	}

	if team != nil {
		state.fromKeeper(team)
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	} else {
		resp.State.RemoveResource(ctx)
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
		plan.NodeId = types.Int64Value(r.enterprise.EnterpriseData().GetRootNode().NodeId)
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

	var t = &enterprise.Team{
		TeamUid: api.Base64UrlDecode(plan.TeamUid.ValueString()),
	}
	plan.toKeeper(t)

	var er error
	if er1 := enterprise.PutTeams(r.enterprise, nil, []*enterprise.Team{t}, nil, func(_ []byte, err error) {
		er = err
	}); er1 != nil {
		er = er1
	}
	if er != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Update Team %s", plan.Name.ValueString()),
			"Error occurred while updating a team: "+er.Error(),
		)
		return
	}
	_ = r.enterprise.Load()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
		plan.NodeId = types.Int64Value(r.enterprise.EnterpriseData().GetRootNode().NodeId)
	}
	var nodeId = plan.NodeId.ValueInt64()
	var team *enterprise.Team
	for _, t := range r.enterprise.EnterpriseData().Teams().GetData() {
		if t.NodeId == nodeId && strings.EqualFold(t.Name, teamName) {
			team = t
			break
		}
	}
	if team == nil {
		team = new(enterprise.Team)
	} else {
		plan.TeamUid = types.StringValue(api.Base64UrlEncode(team.TeamUid))
	}
	plan.toKeeper(team)

	var forInsert []*enterprise.Team
	var forUpdate []*enterprise.Team
	if team.TeamUid == nil {
		forInsert = append(forInsert, team)
	} else {
		forUpdate = append(forUpdate, team)
	}
	var er error
	if er1 := enterprise.PutTeams(r.enterprise, forInsert, forUpdate, nil, func(_ []byte, err error) {
		er = err
	}); er1 != nil {
		er = er1
	}
	if er != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Team %s", teamName),
			fmt.Sprintf("Error: %s", er),
		)
		return
	}
	_ = r.enterprise.Load()
	plan.TeamUid = types.StringValue(api.Base64UrlEncode(team.TeamUid))
	plan.NodeId = types.Int64Value(team.NodeId)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}
	var forDelete = [][]byte{api.Base64UrlDecode(state.TeamUid.ValueString())}
	var er error
	if er1 := enterprise.PutTeams(r.enterprise, nil, nil, forDelete, func(_ []byte, err error) {
		er = err
	}); er1 != nil {
		er = er1
	}
	if er != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete Team %s", state.Name.ValueString()),
			"Error occurred while creating a team",
		)
		return
	}
	_ = r.enterprise.Load()
}
