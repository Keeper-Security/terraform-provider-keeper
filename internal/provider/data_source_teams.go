package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"reflect"
	"terraform-provider-kepr/internal/model"
)

var (
	_ datasource.DataSource = &teamsDataSource{}
)

func newTeamsDataSource() datasource.DataSource {
	return &teamsDataSource{}
}

type teamsDataSourceModel struct {
	FilterCriteria *model.FilterCriteria `tfsdk:"filter"`
	NodeCriteria   *model.NodeCriteria   `tfsdk:"nodes"`
	Teams          []*model.TeamModel    `tfsdk:"teams"`
}

type teamsDataSource struct {
	teams enterprise.IEnterpriseEntity[enterprise.ITeam, string]
	nodes enterprise.IEnterpriseEntity[enterprise.INode, int64]
}

func (d *teamsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

// Schema defines the schema for the data source.
func (d *teamsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filter": schema.SingleNestedAttribute{
				Attributes:  model.FilterCriteriaAttributes,
				Optional:    true,
				Description: "Search By field filter",
			},
			"nodes": schema.SingleNestedAttribute{
				Attributes:  model.NodeCriteriaAttributes,
				Optional:    true,
				Description: "Search By node filter",
			},
			"teams": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: model.TeamSchemaAttributes,
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

	var nm model.NodeMatcher
	nm, diags = model.GetNodeMatcher(tq.NodeCriteria, d.nodes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fm model.Matcher
	fm, diags = model.GetFieldMatcher(tq.FilterCriteria, reflect.TypeOf((*model.TeamModel)(nil)))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = tq
	d.teams.GetAllEntities(func(t enterprise.ITeam) bool {
		if nm != nil {
			if !nm(t.NodeId()) {
				return true
			}
		}
		var team = new(model.TeamModel)
		team.FromKeeper(t)
		if fm != nil {
			if !fm(team) {
				return true
			}
		}
		state.Teams = append(state.Teams, team)
		return true
	})

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *teamsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.teams = ed.Teams()
		d.nodes = ed.Nodes()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseData\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *teamsDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("nodes"),
			path.MatchRoot("filter"),
		),
	}
}
