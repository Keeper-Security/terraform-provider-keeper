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
	"terraform-provider-kepr/internal/model"
)

var (
	_ datasource.DataSourceWithConfigValidators = &nodeDataSource{}
)

func newNodeDataSource() datasource.DataSource {
	return &nodeDataSource{}
}

type nodeDataSourceModel struct {
	NodeId             types.Int64            `tfsdk:"node_id"`
	Name               types.String           `tfsdk:"name"`
	ParentId           types.Int64            `tfsdk:"parent_id"`
	Bridge             *bridgeShortModel      `tfsdk:"bridge"`
	Scim               *scimShortModel        `tfsdk:"scim"`
	DuoEnabled         types.Bool             `tfsdk:"duo_enabled"`
	RsaEnabled         types.Bool             `tfsdk:"rsa_enabled"`
	RestrictVisibility types.Bool             `tfsdk:"restrict_visibility"`
	SsoOnPremise       *ssoProviderShortModel `tfsdk:"sso_provider_on_premise"`
	SsoInCloud         *ssoProviderShortModel `tfsdk:"sso_provider_in_cloud"`
	IsRoot             types.Bool             `tfsdk:"is_root"`
}

func (model *nodeDataSourceModel) fromKeeper(node enterprise.INode) {
	model.NodeId = types.Int64Value(node.NodeId())
	model.Name = types.StringValue(node.Name())
	model.DuoEnabled = types.BoolValue(node.DuoEnabled())
	model.RsaEnabled = types.BoolValue(node.RsaEnabled())
	if node.ParentId() > 0 {
		model.ParentId = types.Int64Value(node.ParentId())
	}
	if node.RestrictVisibility() {
		model.RestrictVisibility = types.BoolValue(true)
	}
}

func (model *nodeDataSourceModel) loadScimData(scim enterprise.IScim) {
	var sm = new(scimShortModel)
	sm.fromKeeper(scim)
	model.Scim = sm
}

func (model *nodeDataSourceModel) loadBridgeData(bridge enterprise.IBridge) {
	var bm = new(bridgeShortModel)
	bm.fromKeeper(bridge)
	model.Bridge = bm
}

func (model *nodeDataSourceModel) loadSsoServiceData(service enterprise.ISsoService) {
	var sm = new(ssoProviderShortModel)
	sm.fromKeeper(service)
	if service.IsCloud() {
		model.SsoOnPremise = sm
	} else {
		model.SsoInCloud = sm
	}
}

type nodeDataSource struct {
	rootNode     enterprise.INode
	nodes        enterprise.IEnterpriseEntity[enterprise.INode, int64]
	bridges      enterprise.IEnterpriseEntity[enterprise.IBridge, int64]
	scims        enterprise.IEnterpriseEntity[enterprise.IScim, int64]
	ssoProviders enterprise.IEnterpriseEntity[enterprise.ISsoService, int64]
}

func (d *nodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

// Schema defines the schema for the data source.
func (d *nodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var filterAttributes = map[string]schema.Attribute{
		"is_root": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Root Node",
		},
		"node_id": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Node ID",
		},
		"name": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Node Name",
		},
		"parent_id": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Parent Node ID",
		},
	}
	resp.Schema = schema.Schema{
		Attributes: model.MergeMaps(filterAttributes, nodeDetailedSchemaAttributes),
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

	var node enterprise.INode
	var m func(node enterprise.INode) bool

	if !nq.IsRoot.IsNull() && nq.IsRoot.ValueBool() {
		node = d.rootNode
	} else if !nq.Name.IsNull() && !nq.ParentId.IsNull() {
		var nodeName = nq.Name.ValueString()
		var parentId = nq.ParentId.ValueInt64()
		m = func(node enterprise.INode) bool {
			return node.ParentId() == parentId && strings.EqualFold(nodeName, node.Name())
		}
	} else if !nq.NodeId.IsNull() && !nq.NodeId.IsUnknown() {
		m = func(node enterprise.INode) bool {
			return node.NodeId() == nq.NodeId.ValueInt64()
		}
	}

	if node == nil && m == nil {
		resp.Diagnostics.AddError(
			"Search criteria is not provided for \"node\" data source",
			fmt.Sprintf("Search criteria is not provided for \"node\" data source"),
		)
		return
	}
	if node == nil && m != nil {
		d.nodes.GetAllEntities(func(n enterprise.INode) bool {
			if m(n) {
				node = n
				return false
			}
			return true
		})
	}

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

	var scimId = node.ScimId()
	if scimId > 0 {
		var scim = d.scims.GetEntity(scimId)
		if scim != nil {
			nm.loadScimData(scim)
		}
	}

	var bridgeId = node.BridgeId()
	if bridgeId > 0 {
		var bridge = d.bridges.GetEntity(bridgeId)
		if bridge != nil {
			nm.loadBridgeData(bridge)
		}
	}

	if len(node.SsoServiceProviderId()) > 0 {
		for _, x := range node.SsoServiceProviderId() {
			var service = d.ssoProviders.GetEntity(x)
			if service != nil {
				nm.loadSsoServiceData(service)
			}
		}
	}

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
		d.rootNode = ed.RootNode()
		d.scims = ed.Scims()
		d.ssoProviders = ed.SsoServices()
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
		datasourcevalidator.RequiredTogether(
			path.MatchRoot("name"),
			path.MatchRoot("parent_id"),
		),
	}
}
