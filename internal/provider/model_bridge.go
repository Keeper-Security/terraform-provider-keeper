package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
)

type bridgeShortModel struct {
	BridgeId	types.Int64		`tfsdk:"bridge_id"`
	Status		types.String	`tfsdk:"status"`
}

var bridgeShortSchemaAttributes = map[string]schema.Attribute{
	"bridge_id": schema.Int64Attribute{
		Computed: true,
		Description: "Bridge ID",
	},
	"status": schema.StringAttribute{
		Computed: true,
		Description: "Bridge Status",
	},
}

func (model *bridgeShortModel) fromKeeper(bridge enterprise.IBridge) {
	model.BridgeId = types.Int64Value(bridge.BridgeId())
	model.Status = types.StringValue(bridge.Status())
}