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
	_ datasource.DataSourceWithValidateConfig = &teamDataSource{}
)

type teamDataSourceModel struct {
	Name    types.String `tfsdk:"name"`
	TeamUid types.String `tfsdk:"team_uid"`
	Team    *teamModel   `tfsdk:"team"`
}

type teamDataSource struct {
	teams enterprise.IEnterpriseEntity[enterprise.Team]
}

func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

func (d *teamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the data source.
func (d *teamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Search By Team Name",
			},
			"team_uid": schema.StringAttribute{
				Optional:    true,
				Description: "Search by Team UID",
			},
			"team": schema.SingleNestedAttribute{
				Attributes:  teamSchemaAttributes,
				Computed:    true,
				Description: "A found team or nil",
			},
		},
	}
}

func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tq teamDataSourceModel
	diags := req.Config.Get(ctx, &tq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var teamMatcher func(*enterprise.Team) bool
	if !tq.Name.IsNull() && !tq.Name.IsUnknown() {
		teamMatcher = func(team *enterprise.Team) bool {
			return strings.EqualFold(tq.Name.ValueString(), team.Name)
		}
	} else if !tq.TeamUid.IsNull() && !tq.TeamUid.IsUnknown() {
		teamMatcher = func(team *enterprise.Team) bool {
			return strings.EqualFold(api.Base64UrlEncode(team.TeamUid), tq.TeamUid.ValueString())
		}
	}

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
