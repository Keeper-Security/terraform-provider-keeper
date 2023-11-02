package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"reflect"
)

var (
	_ datasource.DataSource = &teamsDataSource{}
)

type teamsDataSourceModel struct {
	Filter *filterCriteria `tfsdk:"filter"`
	Teams  []*teamModel    `tfsdk:"teams"`
}

type teamsDataSource struct {
	teams enterprise.IEnterpriseEntity[enterprise.Team]
}

func NewTeamsDataSource() datasource.DataSource {
	return &teamsDataSource{}
}

func (d *teamsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

// Schema defines the schema for the data source.
func (d *teamsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filter": schema.SingleNestedAttribute{
				Attributes:  filterCriteriaAttributes,
				Required:    true,
				Description: "Search By field filter",
			},
			"teams": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: teamSchemaAttributes,
				},
			},
		},
	}
}

func (d *teamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tq teamsDataSourceModel
	diags := req.Config.Get(ctx, &tq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tq.Filter == nil {
		resp.Diagnostics.AddError(
			"Search criteria \"filter\" is not provided for \"teams\" data source",
			fmt.Sprintf("Search criteria is not provided for \"teams\" data source"),
		)
		return
	}

	var cb matcher
	cb, diags = getFieldMatcher(tq.Filter, reflect.TypeOf((*teamModel)(nil)))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = tq
	for _, v := range d.teams.GetData() {
		var team = new(teamModel)
		team.fromKeeper(v)
		if cb(team) {
			state.Teams = append(state.Teams, team)
		}
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *teamsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

/*
func (d *teamsDataSource) ValidateConfig(ctx context.Context,
	req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var tq teamsDataSourceModel
	diags := req.Config.Get(ctx, &tq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var criterias = 0
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
*/
