package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"strings"
)

func newNodeResource() resource.Resource {
	return &nodeResource{}
}

type nodeResourceModel struct {
	NodeId             types.Int64  `tfsdk:"node_id"`
	Name               types.String `tfsdk:"name"`
	ParentId           types.Int64  `tfsdk:"parent_id"`
	RestrictVisibility types.Bool   `tfsdk:"restrict_visibility"`
}

func (model *nodeResourceModel) fromKeeper(node enterprise.INode) {
	model.NodeId = types.Int64Value(node.NodeId())
	model.Name = types.StringValue(node.Name())
	if node.ParentId() > 0 {
		model.ParentId = types.Int64Value(node.ParentId())
	}
	if node.RestrictVisibility() {
		model.RestrictVisibility = types.BoolValue(true)
	}
}

func (model *nodeResourceModel) toKeeper(node enterprise.INodeEdit) {
	node.SetName(model.Name.ValueString())
	node.SetParentId(model.ParentId.ValueInt64())
	if !model.RestrictVisibility.IsNull() {
		node.SetRestrictVisibility(model.RestrictVisibility.ValueBool())
	}
}

var nodeResourceSchemaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Computed:      true,
		Description:   "Node ID",
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	},
	"name": schema.StringAttribute{
		Required:    true,
		Description: "Node Name",
	},
	"parent_id": schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: "Parent Node ID",
	},
	"restrict_visibility": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Restrict Node Visibility",
	},
}

type nodeResource struct {
	management enterprise.IEnterpriseManagement
}

func (nr *nodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (nr *nodeResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: mergeMaps(nodeResourceSchemaAttributes),
	}
}

func (nr *nodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	mgmt, ok := req.ProviderData.(enterprise.IEnterpriseManagement)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected IEnterpriseManagement, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	nr.management = mgmt
}

func (nr *nodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("node_id"), req, resp)
}

func (nr *nodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state nodeResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	if state.ParentId.IsNull() || state.ParentId.ValueInt64() == 0 {
		state.ParentId = types.Int64Value(nr.management.EnterpriseData().GetRootNode().NodeId())
	}

	var nodes = nr.management.EnterpriseData().Nodes()
	var node enterprise.INode
	if state.NodeId.IsNull() {
		var nodeName = state.Name.ValueString()
		var parentId = state.ParentId.ValueInt64()
		nodes.GetAllEntities(func(t enterprise.INode) bool {
			if t.NodeId() == parentId && strings.EqualFold(t.Name(), nodeName) {
				node = t
				return false
			}
			return true
		})
	} else {
		var nodeId = state.NodeId.ValueInt64()
		if nodeId > 0 {
			node = nodes.GetEntity(nodeId)
		}
	}

	if node != nil {
		state.fromKeeper(node)
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (nr *nodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan nodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state nodeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ParentId.IsNull() || plan.ParentId.ValueInt64() == 0 {
		plan.ParentId = state.NodeId
	}
	if plan.ParentId.IsNull() || plan.ParentId.ValueInt64() == 0 {
		plan.ParentId = types.Int64Value(nr.management.EnterpriseData().GetRootNode().NodeId())
	}
	if plan.NodeId.IsNull() {
		plan.NodeId = state.NodeId
	}
	if plan.NodeId.IsNull() {
		resp.Diagnostics.AddError(
			"Update Node error: cannot resolve node ID",
			"Could not resolve node ID",
		)
		return
	}

	var nodeId = plan.NodeId.ValueInt64()
	var n = enterprise.NewNode(nodeId)
	plan.toKeeper(n)

	errs := nr.management.ModifyNodes(nil, []enterprise.INode{n}, nil)
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Update Node %s", plan.Name.ValueString()),
			"Error occurred while updating a node: "+er.Error(),
		)
	}
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (nr *nodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan nodeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nodeName = plan.Name.ValueString()
	if len(nodeName) == 0 {
		resp.Diagnostics.AddError(
			"Create Node: Node name cannot be empty",
			"Node name cannot be empty",
		)
		return
	}
	if plan.ParentId.IsNull() || plan.ParentId.ValueInt64() == 0 {
		plan.ParentId = types.Int64Value(nr.management.EnterpriseData().GetRootNode().NodeId())
	}
	var parentId = plan.ParentId.ValueInt64()
	var node enterprise.INode
	nr.management.EnterpriseData().Nodes().GetAllEntities(func(n enterprise.INode) bool {
		if parentId == n.ParentId() && strings.EqualFold(n.Name(), nodeName) {
			node = n
			return false
		}
		return true
	})

	var ne enterprise.INodeEdit
	var toAdd []enterprise.INode
	var toUpdate []enterprise.INode
	if node == nil {
		var nodeId int64
		var err error
		if nodeId, err = nr.management.GetEnterpriseId(); err != nil {
			resp.Diagnostics.AddError(
				"Create Node: Error getting enterprise ID",
				fmt.Sprintf("Error getting enterprise ID: %s", err.Error()),
			)
			return
		}
		ne = enterprise.NewNode(nodeId)
		toAdd = append(toAdd, ne)
	} else {
		ne = enterprise.CloneNode(node)
		toUpdate = append(toUpdate, ne)
	}
	plan.toKeeper(ne)

	var errs = nr.management.ModifyNodes(toAdd, toUpdate, nil)
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Create Node %s", nodeName),
			fmt.Sprintf("Error: %s", er),
		)
	}
	if !resp.Diagnostics.HasError() {
		plan.NodeId = types.Int64Value(ne.NodeId())
		plan.ParentId = types.Int64Value(ne.ParentId())
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}
}

func (nr *nodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state nodeResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var errs = nr.management.ModifyNodes(nil, nil, []int64{state.NodeId.ValueInt64()})
	for _, er := range errs {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Delete Node \"%s\" error: %s", state.Name.ValueString(), er),
			"Error occurred while creating a node",
		)
	}
}
