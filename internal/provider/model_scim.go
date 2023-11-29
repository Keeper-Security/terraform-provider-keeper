package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
)

type scimShortModel struct {
	ScimId 		types.Int64		`tfsdk:"scim_id"`
	Status		types.String	`tfsdk:"status"`
	RolePrefix	types.String	`tfsdk:"role_prefix"`
}

var scimShortSchemaAttributes = map[string]schema.Attribute{
	"scim_id": schema.Int64Attribute{
		Computed: 		true,
		Description: 	"SCIM ID",
	},
	"status": schema.StringAttribute{
		Computed: 		true,
		Description: 	"SCIM Status",
	},
	"role_prefix": schema.StringAttribute{
		Computed: 		true,
		Description: 	"SCIM Role Prefix",
	},
}

func (model *scimShortModel) fromKeeper(keeper enterprise.IScim) {
	model.ScimId = types.Int64Value(keeper.ScimId())
	model.Status = types.StringValue(keeper.Status())
	model.RolePrefix = types.StringValue(keeper.RolePrefix())
}