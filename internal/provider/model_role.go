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
	IsAdmin        types.Bool   `tfsdk:"is_admin"`
}

func (model *roleModel) fromKeeper(role enterprise.IRole, isAdmin bool) {
	model.RoleId = types.Int64Value(role.RoleId())
	model.Name = types.StringValue(role.Name())
	model.NodeId = types.Int64Value(role.NodeId())
	model.VisibleBelow = types.BoolValue(role.VisibleBelow())
	model.NewUserInherit = types.BoolValue(role.NewUserInherit())
	model.IsAdmin = types.BoolValue(isAdmin)
}

var roleSchemaAttributes = map[string]schema.Attribute{
	"role_id": schema.Int64Attribute{
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
	"is_admin": schema.BoolAttribute{
		Computed:    true,
		Description: "Is Administrative Role",
	},
}

type managedNodeModel struct {
	NodeId                types.Int64    `tfsdk:"node_id"`
	Name                  types.String   `tfsdk:"name"`
	CascadeNodeManagement types.Bool     `tfsdk:"cascade_node_management"`
	Privileges            []types.String `tfsdk:"privileges"`
}

var managedNodeSchemaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Managed Node ID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Managed Node Name",
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

type roleShortModel struct {
	RoleId  types.Int64  `tfsdk:"role_id"`
	Name    types.String `tfsdk:"name"`
	NodeId  types.Int64  `tfsdk:"node_id"`
	IsAdmin types.Bool   `tfsdk:"is_admin"`
}

func (model *roleShortModel) fromKeeper(role enterprise.IRole, isAdmin bool) {
	model.RoleId = types.Int64Value(role.RoleId())
	model.Name = types.StringValue(role.Name())
	model.NodeId = types.Int64Value(role.NodeId())
	model.IsAdmin = types.BoolValue(isAdmin)
}

var roleShortSchemaAttributes = map[string]schema.Attribute{
	"role_id": schema.Int64Attribute{
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
	"is_admin": schema.BoolAttribute{
		Computed:    true,
		Description: "Is Administrative Role",
	},
}
