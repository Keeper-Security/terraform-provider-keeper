package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
)

type roleModel struct {
	RoleId         types.Int64  `tfsdk:"role_id"`
	Name           types.String `tfsdk:"name"`
	NodeId         types.Int64  `tfsdk:"node_id"`
	VisibleBelow   types.Bool   `tfsdk:"visible_below"`
	NewUserInherit types.Bool   `tfsdk:"new_user_inherit"`
}

func (model *roleModel) fromKeeper(role enterprise.IRole) {
	model.RoleId = types.Int64Value(role.RoleId())
	model.Name = types.StringValue(role.Name())
	model.NodeId = types.Int64Value(role.NodeId())
	model.VisibleBelow = types.BoolValue(role.VisibleBelow())
	model.NewUserInherit = types.BoolValue(role.NewUserInherit())
}

var roleSchemaAttributes = map[string]schema.Attribute{
	"role_id": schema.StringAttribute{
		Computed:    true,
		Description: "Role ID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Role Name",
	},
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Role Node ID",
	},
	"visible_below": schema.BoolAttribute{
		Computed:    true,
		Description: "Visible Below",
	},
	"new_user_inherit": schema.BoolAttribute{
		Computed:    true,
		Description: "New User Inherit",
	},
}

type managedNodeModel struct {
	NodeId                types.Int64    `tfsdk:"node_id"`
	CascadeNodeManagement types.Bool     `tfsdk:"cascade_node_management"`
	Privileges            []types.String `tfsdk:"privileges"`
}

var managedNodeSchemaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Managed Node ID",
	},
	"cascade_node_management": schema.BoolAttribute{
		Computed:    true,
		Description: "Cascade Node Management",
	},
	"privileges": schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Description: "Privileges",
	},
}
