package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
)

type nodeModel struct {
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

func (model *nodeModel) fromKeeper(node enterprise.INode) {
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

var nodeSchemaAttributes = map[string]schema.Attribute{
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

var nodeDetailedSchemaAttributes = map[string]schema.Attribute{
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
		Attributes: bridgeShortSchemaAttributes,
		Computed:    true,
		Description: "Bridge",
	},
	"scim": schema.SingleNestedAttribute{
		Attributes: scimShortSchemaAttributes,
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
		Attributes:		ssoProviderShortSchemaAttributes,
		Computed:		true,
		Description: 	"On-premise SSO Service Provider",
	},
	"sso_provider_in_cloud": schema.SingleNestedAttribute{
		Attributes:		ssoProviderShortSchemaAttributes,
		Computed:		true,
		Description: 	"Cloud-based SSO Service Provider",
	},
}
