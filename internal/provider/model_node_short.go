package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
)

type nodeShortModel struct {
	NodeId types.Int64  `tfsdk:"node_id"`
	Name   types.String `tfsdk:"name"`
}

func (model *nodeShortModel) fromKeeper(role enterprise.INode) {
	model.NodeId = types.Int64Value(role.NodeId())
	model.Name = types.StringValue(role.Name())
}

var nodeShortSchemaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Role Node ID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Role Name",
	},
}
