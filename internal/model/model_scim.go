package model

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
)

type ScimShortModel struct {
	ScimId     types.Int64  `tfsdk:"scim_id"`
	Status     types.String `tfsdk:"status"`
	RolePrefix types.String `tfsdk:"role_prefix"`
}

var ScimShortSchemaAttributes = map[string]schema.Attribute{
	"scim_id": schema.Int64Attribute{
		Computed:    true,
		Description: "SCIM ID",
	},
	"status": schema.StringAttribute{
		Computed:    true,
		Description: "SCIM Status",
	},
	"role_prefix": schema.StringAttribute{
		Computed:    true,
		Description: "SCIM Role Prefix",
	},
}

func (scim *ScimShortModel) FromKeeper(keeper enterprise.IScim) {
	scim.ScimId = types.Int64Value(keeper.ScimId())
	scim.Status = types.StringValue(keeper.Status())
	scim.RolePrefix = types.StringValue(keeper.RolePrefix())
}
