package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"reflect"
	"terraform-provider-kepr/internal/model"
)

var (
	_ datasource.DataSource = &nodesDataSource{}
)

func newNodesDataSource() datasource.DataSource {
	return &nodesDataSource{}
}

type subNodesCriteria struct {
	NodeId        types.Int64 `tfsdk:"node_id"`
	IncludeParent types.Bool  `tfsdk:"include_parent"`
	Cascade       types.Bool  `tfsdk:"cascade"`
}

type nodesDataSourceModel struct {
	Filter          *model.FilterCriteria `tfsdk:"filter"`
	SubNodeCriteria *subNodesCriteria     `tfsdk:"subnodes"`
	Nodes           []*nodeModel          `tfsdk:"nodes"`
}

type nodesDataSource struct {
	nodes enterprise.IEnterpriseEntity[enterprise.INode, int64]
}

func (d *nodesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nodes"
}
func (d *nodesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"subnodes": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"node_id": schema.Int64Attribute{
						Required:    true,
						Description: "Base Node ID",
					},
					"include_parent": schema.BoolAttribute{
						Optional:    true,
						Description: "Include base node to result",
					},
					"cascade": schema.BoolAttribute{
						Optional:    true,
						Description: "Traverse node structure",
					},
				},
				Optional:    true,
				Description: "Load subnode tree",
			},
			"filter": schema.SingleNestedAttribute{
				Attributes:  model.FilterCriteriaAttributes,
				Optional:    true,
				Description: "Search By field filter",
			},
			"nodes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: nodeSchemaAttributes,
				},
			},
		},
	}
}

func (d *nodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var nq nodesDataSourceModel
	diags := req.Config.Get(ctx, &nq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = nq

	if nq.SubNodeCriteria != nil {
		var nodes []enterprise.INode
		var nodeId int64
		if !nq.SubNodeCriteria.NodeId.IsNull() && !nq.SubNodeCriteria.NodeId.IsUnknown() {
			nodeId = nq.SubNodeCriteria.NodeId.ValueInt64()
		}
		var isCascade bool
		if !nq.SubNodeCriteria.Cascade.IsNull() && !nq.SubNodeCriteria.Cascade.IsUnknown() {
			isCascade = nq.SubNodeCriteria.Cascade.ValueBool()
		}
		var includeParent bool
		if !nq.SubNodeCriteria.IncludeParent.IsNull() && !nq.SubNodeCriteria.IncludeParent.IsUnknown() {
			includeParent = nq.SubNodeCriteria.IncludeParent.ValueBool()
		}
		if nodeId > 0 {
			d.nodes.GetAllEntities(func(n enterprise.INode) bool {
				if n.NodeId() == nodeId {
					nodes = append(nodes, n)
					return false
				}
				return true
			})
			var pos = 0
			for pos < len(nodes) {
				nodeId = nodes[pos].NodeId()
				pos += 1
				tflog.Warn(ctx, fmt.Sprintf("pos=%d; len=%d", pos, len(nodes)))
				d.nodes.GetAllEntities(func(n enterprise.INode) bool {
					if n.ParentId() == nodeId {
						nodes = append(nodes, n)
					}
					return true
				})
				if !isCascade {
					break
				}
			}
			if !includeParent {
				nodes = nodes[1:]
			}
		}
		for _, n := range nodes {
			var nm = new(nodeModel)
			nm.fromKeeper(n)
			state.Nodes = append(state.Nodes, nm)
		}
	} else if nq.Filter != nil {
		var cb model.Matcher
		cb, diags = model.GetFieldMatcher(nq.Filter, reflect.TypeOf((*nodeModel)(nil)))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		d.nodes.GetAllEntities(func(n enterprise.INode) bool {
			var nm = new(nodeModel)
			nm.fromKeeper(n)
			if cb(nm) {
				state.Nodes = append(state.Nodes, nm)
			}
			return true
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *nodesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.nodes = ed.Nodes()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseData\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *nodesDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("subnodes"),
			path.MatchRoot("filter"),
		),
	}
}
