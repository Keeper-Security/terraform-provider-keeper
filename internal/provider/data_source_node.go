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
	NodeId               types.Int64   `tfsdk:"node_id"`
	Name                 types.String  `tfsdk:"name"`
	ParentId             types.Int64   `tfsdk:"parent_id"`
	BridgeId             types.Int64   `tfsdk:"bridge_id"`
	ScimId               types.Int64   `tfsdk:"scim_id"`
	DuoEnabled           types.Bool    `tfsdk:"duo_enabled"`
	RsaEnabled           types.Bool    `tfsdk:"rsa_enabled"`
	RestrictVisibility   types.Bool    `tfsdk:"restrict_visibility"`
	SsoServiceProviderId []types.Int64 `tfsdk:"sso_provider_ids"`
	IsRoot               types.Bool    `tfsdk:"is_root"`
}

func (model *nodeDataSourceModel) fromKeeper(node enterprise.INode) {
	model.NodeId = types.Int64Value(node.NodeId())
	model.Name = types.StringValue(node.Name())
	model.DuoEnabled = types.BoolValue(node.DuoEnabled())
	model.RsaEnabled = types.BoolValue(node.RsaEnabled())
	if node.ParentId() > 0 {
		model.ParentId = types.Int64Value(node.ParentId())
	}
	if node.BridgeId() > 0 {
		model.BridgeId = types.Int64Value(node.BridgeId())
	}
	if node.ScimId() > 0 {
		model.ScimId = types.Int64Value(node.ScimId())
	}
	if node.RestrictVisibility() {
		model.RestrictVisibility = types.BoolValue(true)
	}
	if len(node.SsoServiceProviderId()) > 0 {
		for _, x := range node.SsoServiceProviderId() {
			model.SsoServiceProviderId = append(model.SsoServiceProviderId, types.Int64Value(x))
		}
	}
}

type nodeDataSource struct {
	nodes   enterprise.IEnterpriseEntity[enterprise.INode, int64]
	bridges enterprise.IEnterpriseEntity[enterprise.IBridge, int64]
}

func NewNodeDataSource() datasource.DataSource {
	return &nodeDataSource{}
}
func (d *nodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

// Schema defines the schema for the data source.
func (d *nodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var filterAttributes = map[string]schema.Attribute{
		"is_root": schema.BoolAttribute{
			Optional:    true,
			Description: "Root Node",
		},
		"node_id": schema.Int64Attribute{
			Optional:    true,
			Description: "Node ID",
		},
		"name": schema.StringAttribute{
			Optional:    true,
			Description: "Node Name",
		},
	}
	resp.Schema = schema.Schema{
		Attributes: mergeMaps(filterAttributes, nodeSchemaAttributes),
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

	var m func(node enterprise.INode) bool
	if !nq.IsRoot.IsNull() && nq.IsRoot.ValueBool() {
		m = func(node enterprise.INode) bool {
			return node.ParentId() == 0
		}
	} else if !nq.Name.IsNull() && !nq.Name.IsUnknown() {
		m = func(node enterprise.INode) bool {
			return strings.EqualFold(nq.Name.ValueString(), node.Name())
		}
	} else if !nq.NodeId.IsNull() && !nq.NodeId.IsUnknown() {
		m = func(node enterprise.INode) bool {
			return node.NodeId() == nq.NodeId.ValueInt64()
		}
	}

	if m == nil {
		resp.Diagnostics.AddError(
			"Search criteria is not provided for \"node\" data source",
			fmt.Sprintf("Search criteria is not provided for \"node\" data source"),
		)
		return
	}

	var node enterprise.INode
	d.nodes.GetAllEntities(func(n enterprise.INode) bool {
		if m(n) {
			node = n
			return false
		}
		return true
	})

	if node == nil {
		resp.Diagnostics.AddError(
			"Node not found",
			fmt.Sprintf("Cannot find a node according to the provided criteria"),
		)
		return
	}

	var state = nq
	var nm = &state
	nm.fromKeeper(node)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *nodeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if ed, ok := req.ProviderData.(enterprise.IEnterpriseData); ok {
		d.nodes = ed.Nodes()
		d.bridges = ed.Bridges()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected \"IEnterpriseData\", got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (d *nodeDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("is_root"),
			path.MatchRoot("node_id"),
			path.MatchRoot("name"),
		),
	}
}
