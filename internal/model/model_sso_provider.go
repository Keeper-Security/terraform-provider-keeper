package model

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
)

type SsoProviderShortModel struct {
	ProviderId types.Int64  `tfsdk:"provider_id"`
	Name       types.String `tfsdk:"name"`
	IsActive   types.Bool   `tfsdk:"is_active"`
	Url        types.String `tfsdk:"url"`
}

var SsoProviderShortSchemaAttributes = map[string]schema.Attribute{
	"provider_id": schema.Int64Attribute{
		Computed:    true,
		Description: "SSO service provider ID",
	},
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "SSO service provider name",
	},
	"is_active": schema.BoolAttribute{
		Computed:    true,
		Description: "Is SSO service active?",
	},
	"url": schema.StringAttribute{
		Computed:    true,
		Description: "SSO service provider URL",
	},
}

func (sso *SsoProviderShortModel) FromKeeper(keeper enterprise.ISsoService) {
	sso.ProviderId = types.Int64Value(keeper.SsoServiceProviderId())
	sso.Name = types.StringValue(keeper.Name())
	sso.IsActive = types.BoolValue(keeper.Active())
	sso.Url = types.StringValue(keeper.SpUrl())
}
