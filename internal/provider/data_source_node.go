package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"strings"
)

var (
	_ datasource.DataSourceWithConfigValidators = &nodeDataSource{}
)

type nodeDataSourceModel struct {
	IsRoot types.Bool   `tfsdk:"is_root"`
	Name   types.String `tfsdk:"name"`
	NodeId types.Int64  `tfsdk:"node_id"`
	Node   *nodeModel   `tfsdk:"node"`
}

type nodeDataSource struct {
	nodes enterprise.IEnterpriseEntity[enterprise.Node]
}

func NewNodeDataSource() datasource.DataSource {
	return &nodeDataSource{}
}
func (d *nodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

// Schema defines the schema for the data source.
func (d *nodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"is_root": schema.BoolAttribute{
				Optional:    true,
				Description: "Return Root Node",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Search By Node Name",
			},
			"node_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Search by Node ID",
			},
			"node": schema.SingleNestedAttribute{
				Attributes:  nodeSchemaAttributes,
				Computed:    true,
				Description: "A found node or nil",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *nodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var nq nodeDataSourceModel
	diags := req.Config.Get(ctx, &nq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var m func(node *enterprise.Node) bool
	if !nq.IsRoot.IsNull() && nq.IsRoot.ValueBool() {
		m = func(node *enterprise.Node) bool {
			return node.ParentId == 0
		}
	} else if !nq.Name.IsNull() && !nq.Name.IsUnknown() {
		m = func(node *enterprise.Node) bool {
			return strings.EqualFold(nq.Name.ValueString(), node.Name)
		}
	} else if !nq.NodeId.IsNull() && !nq.NodeId.IsUnknown() {
		m = func(node *enterprise.Node) bool {
			return node.NodeId == nq.NodeId.ValueInt64()
		}
	}

	if m == nil {
		resp.Diagnostics.AddError(
			"Search criteria is not provided for \"node\" data source",
			fmt.Sprintf("Search criteria is not provided for \"node\" data source"),
		)
		return
	}

	var node *enterprise.Node
	for _, v := range d.nodes.GetData() {
		if m(v) {
			node = v
			break
		}
	}

	var state = nq
	if node != nil {
		state.Node = nodeModelFromKeeper(node)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *nodeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if loader, ok := req.ProviderData.(enterprise.IEnterpriseLoader); ok {
		d.nodes = loader.EnterpriseData().Nodes()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseLoader\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *nodeDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("is_root"),
			path.MatchRoot("node"),
			path.MatchRoot("name"),
		),
	}
}
