package model

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
)

type RoleModel struct {
	RoleId         types.Int64  `tfsdk:"role_id"`
	Name           types.String `tfsdk:"name"`
	NodeId         types.Int64  `tfsdk:"node_id"`
	VisibleBelow   types.Bool   `tfsdk:"visible_below"`
	NewUserInherit types.Bool   `tfsdk:"new_user_inherit"`
	IsAdmin        types.Bool   `tfsdk:"is_admin"`
}

func (rm *RoleModel) FromKeeper(role enterprise.IRole, isAdmin bool) {
	rm.RoleId = types.Int64Value(role.RoleId())
	rm.Name = types.StringValue(role.Name())
	rm.NodeId = types.Int64Value(role.NodeId())
	rm.VisibleBelow = types.BoolValue(role.VisibleBelow())
	rm.NewUserInherit = types.BoolValue(role.NewUserInherit())
	rm.IsAdmin = types.BoolValue(isAdmin)
}

var RoleSchemaAttributes = map[string]schema.Attribute{
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

type ManagedNodeModel struct {
	NodeId                types.Int64              `tfsdk:"node_id"`
	Name                  types.String             `tfsdk:"name"`
	CascadeNodeManagement types.Bool               `tfsdk:"cascade_node_management"`
	Privileges            PrivilegeDataSourceModel `tfsdk:"privileges"`
}

var ManagedNodeDataSourceSchemaAttributes = map[string]schema.Attribute{
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
	"privileges": schema.SingleNestedAttribute{
		Attributes:  PrivilegesDataSourceAttributes,
		Optional:    true,
		Description: "Privileges",
	},
}

type RoleShortModel struct {
	RoleId  types.Int64  `tfsdk:"role_id"`
	Name    types.String `tfsdk:"name"`
	NodeId  types.Int64  `tfsdk:"node_id"`
	IsAdmin types.Bool   `tfsdk:"is_admin"`
}

func (rsm *RoleShortModel) FromKeeper(role enterprise.IRole, isAdmin bool) {
	rsm.RoleId = types.Int64Value(role.RoleId())
	rsm.Name = types.StringValue(role.Name())
	rsm.NodeId = types.Int64Value(role.NodeId())
	rsm.IsAdmin = types.BoolValue(isAdmin)
}

var RoleShortSchemaAttributes = map[string]schema.Attribute{
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
