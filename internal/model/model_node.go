package model

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
)

type NodeModel struct {
	NodeId               types.Int64   `tfsdk:"node_id"`
	Name                 types.String  `tfsdk:"name"`
	ParentId             types.Int64   `tfsdk:"parent_id"`
	BridgeId             types.Int64   `tfsdk:"bridge_id"`
	ScimId               types.Int64   `tfsdk:"scim_id"`
	DuoEnabled           types.Bool    `tfsdk:"duo_enabled"`
	RsaEnabled           types.Bool    `tfsdk:"rsa_enabled"`
	RestrictVisibility   types.Bool    `tfsdk:"restrict_visibility"`
	SsoServiceProviderId []types.Int64 `tfsdk:"sso_provider_ids"`
}

func (nm *NodeModel) FromKeeper(node enterprise.INode) {
	nm.NodeId = types.Int64Value(node.NodeId())
	nm.Name = types.StringValue(node.Name())
	nm.DuoEnabled = types.BoolValue(node.DuoEnabled())
	nm.RsaEnabled = types.BoolValue(node.RsaEnabled())
	if node.ParentId() > 0 {
		nm.ParentId = types.Int64Value(node.ParentId())
	}
	if node.BridgeId() > 0 {
		nm.BridgeId = types.Int64Value(node.BridgeId())
	}
	if node.ScimId() > 0 {
		nm.ScimId = types.Int64Value(node.ScimId())
	}
	if node.RestrictVisibility() {
		nm.RestrictVisibility = types.BoolValue(true)
	}
	if len(node.SsoServiceProviderId()) > 0 {
		for _, x := range node.SsoServiceProviderId() {
			nm.SsoServiceProviderId = append(nm.SsoServiceProviderId, types.Int64Value(x))
		}
	}
}

var NodeSchemaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Node ID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Node Name",
	},
	"parent_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Parent Node ID",
	},
	"bridge_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Bridge ID",
	},
	"scim_id": schema.Int64Attribute{
		Computed:    true,
		Description: "SCIM ID",
	},
	"duo_enabled": schema.BoolAttribute{
		Computed:    true,
		Description: "DUO Enabled",
	},
	"rsa_enabled": schema.BoolAttribute{
		Computed:    true,
		Description: "RSA Configured",
	},
	"restrict_visibility": schema.BoolAttribute{
		Computed:    true,
		Description: "Restrict Node Visibility",
	},
	"sso_provider_ids": schema.ListAttribute{
		Computed:    true,
		ElementType: types.Int64Type,
	},
}

var NodeDetailedSchemaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Node ID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Node Name",
	},
	"parent_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Parent Node ID",
	},
	"bridge": schema.SingleNestedAttribute{
		Attributes:  bridgeShortSchemaAttributes,
		Computed:    true,
		Description: "Bridge",
	},
	"scim": schema.SingleNestedAttribute{
		Attributes:  ScimShortSchemaAttributes,
		Computed:    true,
		Description: "SCIM",
	},
	"duo_enabled": schema.BoolAttribute{
		Computed:    true,
		Description: "DUO Enabled",
	},
	"rsa_enabled": schema.BoolAttribute{
		Computed:    true,
		Description: "RSA Configured",
	},
	"restrict_visibility": schema.BoolAttribute{
		Computed:    true,
		Description: "Restrict Node Visibility",
	},
	"sso_provider_on_premise": schema.SingleNestedAttribute{
		Attributes:  SsoProviderShortSchemaAttributes,
		Computed:    true,
		Description: "On-premise SSO Service Provider",
	},
	"sso_provider_in_cloud": schema.SingleNestedAttribute{
		Attributes:  SsoProviderShortSchemaAttributes,
		Computed:    true,
		Description: "Cloud-based SSO Service Provider",
	},
}

type NodeShortModel struct {
	NodeId types.Int64  `tfsdk:"node_id"`
	Name   types.String `tfsdk:"name"`
}

func (nsm *NodeShortModel) FromKeeper(role enterprise.INode) {
	nsm.NodeId = types.Int64Value(role.NodeId())
	nsm.Name = types.StringValue(role.Name())
}

var NodeShortSchemaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Role Node ID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Role Name",
	},
}
