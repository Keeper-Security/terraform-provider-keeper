package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
)

type ssoProviderShortModel struct {
	ProviderId	types.Int64		`tfsdk:"provider_id"`
	Name		types.String	`tfsdk:"name"`
	IsActive	types.Bool		`tfsdk:"is_active"`
	Url			types.String	`tfsdk:"url"`
}

var ssoProviderShortSchemaAttributes = map[string]schema.Attribute{
	"provider_id": schema.Int64Attribute{
		Computed: true,
		Description: "SSO service provider ID",
	},
	"name": schema.StringAttribute{
		Computed: true,
		Description: "SSO service provider name",
	},
	"is_active": schema.BoolAttribute{
		Computed: true,
		Description: "Is SSO service active?",
	},
	"url": schema.StringAttribute{
		Computed: true,
		Description: "SSO service provider URL",
	},
}

func (model *ssoProviderShortModel) fromKeeper(keeper enterprise.ISsoService) {
	model.ProviderId = types.Int64Value(keeper.SsoServiceProviderId())
	model.Name = types.StringValue(keeper.Name())
	model.IsActive = types.BoolValue(keeper.Active())
	model.Url = types.StringValue(keeper.SpUrl())
}