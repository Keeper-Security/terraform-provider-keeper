package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/api"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"strings"
)

var (
	_ datasource.DataSourceWithValidateConfig = &teamMembershipDataSource{}
)

type teamUserModel struct {
	TeamUid          types.String `tfsdk:"team_uid"`
	TeamName         types.String `tfsdk:"team_name"`
	Username         types.String `tfsdk:"username"`
	EnterpriseUserId types.Int64  `tfsdk:"enterprise_user_id"`
}

var teamUserSchemaAttributes = map[string]schema.Attribute{
	"team_uid": schema.StringAttribute{
		Computed:    true,
		Description: "Team UID",
	},
	"enterprise_user_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Enterprise User ID",
	},
	"team_name": schema.StringAttribute{
		Computed:    true,
		Description: "Team Name",
	},
	"username": schema.StringAttribute{
		Computed:    true,
		Description: "User Account email",
	},
}

type teamMembershipDataSourceModel struct {
	Teams          types.Set        `tfsdk:"teams"`
	Users          types.Set        `tfsdk:"users"`
	TeamMembership []*teamUserModel `tfsdk:"team_users"`
}

type teamMembershipDataSource struct {
	teams enterprise.IEnterpriseLink[enterprise.TeamUser]
}

func NewTeamMembershipDDataSource() datasource.DataSource {
	return &teamMembershipDataSource{}
}

func (d *teamMembershipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_users"
}

func (d *teamMembershipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"teams": schema.SetAttribute{
				Optional:    true,
				Description: "List of teams",
				ElementType: types.StringType,
			},
			"users": schema.SetAttribute{
				Optional:    true,
				Description: "List of users",
				ElementType: types.StringType,
			},
			"team_users": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: teamUserSchemaAttributes,
				}, Computed: true,
				Description: "A found team or nil",
			},
		},
	}
}

func (d *teamMembershipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tmq teamMembershipDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tmq)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var teams map[string]bool
	var users map[string]bool
	var ok bool
	var sv string
	if !tmq.Teams.IsNull() {
		var ts = tmq.Teams.Elements()
		if len(ts) > 0 {
			teams = make(map[string]bool)
			for _, x := range ts {
				if
			}
		}
	} else if !tq.TeamUid.IsNull() && !tq.TeamUid.IsUnknown() {
		teamMatcher = func(team *enterprise.Team) bool {
			return strings.EqualFold(api.Base64UrlEncode(team.TeamUid), tq.TeamUid.ValueString())
		}
	}

	var teamMembershipMatcher func(user *enterprise.TeamUser) bool

	if teamMatcher == nil {
		resp.Diagnostics.AddError(
			"Search criteria is not provided for \"team\" data source",
			fmt.Sprintf("Search criteria is not provided for \"team\" data source"),
		)
		return
	}

	var team *enterprise.Team
	for _, v := range d.teams.GetData() {
		if teamMatcher(v) {
			team = v
			break
		}
	}

	var state = tq
	if team != nil {
		if state.Team == nil {
			state.Team = new(teamModel)
		}
		state.Team.fromKeeper(team)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *teamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if loader, ok := req.ProviderData.(enterprise.IEnterpriseLoader); ok {
		d.teams = loader.EnterpriseData().Teams()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseLoader\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *teamDataSource) ValidateConfig(ctx context.Context,
	req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var tq teamDataSourceModel
	diags := req.Config.Get(ctx, &tq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var criterias = 0
	if !tq.Name.IsNull() && !tq.Name.IsUnknown() {
		criterias++
	}
	if !tq.TeamUid.IsNull() && !tq.TeamUid.IsUnknown() {
		criterias++
	}
	if criterias == 1 {
		return
	}
	var attributes []string
	for k, v := range req.Config.Schema.GetAttributes() {
		if v.IsOptional() {
			attributes = append(attributes, k)
		}
	}

	if criterias == 0 {
		resp.Diagnostics.AddError(
			"\"team\" data source requires one of the following attributes: "+strings.Join(attributes, ", "),
			fmt.Sprintf("Invalid team request"),
		)
	} else if criterias > 1 {
		resp.Diagnostics.AddError(
			"\"team\" data source requires ONLY one of the following attributes: "+strings.Join(attributes, ", "),
			fmt.Sprintf("Invalid team request"),
		)
	}
}
