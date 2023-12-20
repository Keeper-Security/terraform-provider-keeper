package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"reflect"
	"terraform-provider-keeper/internal/model"
)

var (
	_ datasource.DataSource = &rolesDataSource{}
)

type rolesDataSourceModel struct {
	FilterCriteria *model.FilterCriteria `tfsdk:"filter"`
	NodeCriteria   *model.NodeCriteria   `tfsdk:"nodes"`
	Roles          []*model.RoleModel    `tfsdk:"roles"`
}

type rolesDataSource struct {
	nodes        enterprise.IEnterpriseEntity[enterprise.INode, int64]
	roles        enterprise.IEnterpriseEntity[enterprise.IRole, int64]
	managedNodes enterprise.IEnterpriseLink[enterprise.IManagedNode, int64, int64]
}

func newRolesDataSource() datasource.DataSource {
	return &rolesDataSource{}
}

func (d *rolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}
func (d *rolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"roles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: model.RoleSchemaAttributes,
				},
			},
		},
	}
}

func (d *rolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.roles = ed.Roles()
		d.nodes = ed.Nodes()
		d.managedNodes = ed.ManagedNodes()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseData\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *rolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq rolesDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nm model.NodeMatcher
	nm, diags = model.GetNodeMatcher(rq.NodeCriteria, d.nodes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fm model.Matcher
	fm, diags = model.GetFieldMatcher(rq.FilterCriteria, reflect.TypeOf((*model.RoleModel)(nil)))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq

	state.Roles = make([]*model.RoleModel, 0)
	d.roles.GetAllEntities(func(r enterprise.IRole) bool {
		if nm != nil {
			if !nm(r.NodeId()) {
				return true
			}
		}
		var isAdmin bool
		d.managedNodes.GetLinksBySubject(r.RoleId(), func(mn enterprise.IManagedNode) bool {
			isAdmin = true
			return false
		})
		var role = new(model.RoleModel)
		role.FromKeeper(r, isAdmin)
		if fm != nil {
			if !fm(role) {
				return true
			}
		}
		state.Roles = append(state.Roles, role)
		return true
	})

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *rolesDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("nodes"),
			path.MatchRoot("filter"),
		),
	}
}
