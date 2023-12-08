package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	//"github.com/hashicorp/terraform-plugin-framework/datasource"
	//"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	//"github.com/hashicorp/terraform-plugin-framework/path"
	//"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-kepr/internal/model"
)

var (
	_ datasource.DataSource = &enforcementsDataSource{}
)

var enforcementsAttributes =  map[string]schema.Attribute{
	"account": schema.SingleNestedAttribute{
		Attributes: enforcementsAccountAttributes,
		Computed: true,
		Optional: true,
		Description: "Account-related enforcements",
	},
	"allow_ip_list": schema.SingleNestedAttribute{
		Attributes: enforcementsAllowIpListAttributes,
		Computed: true,
		Optional: true,
		Description: "IP whitelist enforcements",
	},
	"sharing": schema.SingleNestedAttribute{
		Attributes: enforcementsSharingAttributes,
		Computed: true,
		Optional: true,
		Description: "Sharing enforcements",
	},
	"keeper_fill": schema.SingleNestedAttribute{
		Attributes: enforcementsKeeperFillAttributes,
		Computed: true,
		Optional: true,
		Description: "Keeper Fill enforcements",

	},
	"login": schema.SingleNestedAttribute{
		Attributes: enforcementsLoginAttributes,
		Computed: true,
		Optional: true,
		Description: "Login-related enforcements",

	},
	"platform": schema.SingleNestedAttribute{
		Attributes: enforcementsPlatformAttributes,
		Computed: true,
		Optional: true,
		Description: "Keeper platform enforcements",
	},
	"record_types": schema.SingleNestedAttribute{
		Attributes: enforcementsRecordTypesAttributes,
		Computed: true,
		Optional: true,
		Description: "Record-type enforcements",
	},
	"two_factor_authentication": schema.SingleNestedAttribute{
		Attributes: enforcements2faAttributes,
		Computed: true,
		Optional: true,
		Description: "2FA enforcments",
	},
	"vault": schema.SingleNestedAttribute{
		Attributes: enforcementsVaultAttributes,
		Computed: true,
		Optional: true,
		Description: "Vault-related enforcements",
	},
}

type enforcementsDataSource struct {
}

func (d *enforcementsDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_enforcements"
}

func (d *enforcementsDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: enforcementsAttributes,
	}
}

func (d *enforcementsDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var rq model.EnforcementsDataSourceModel

	diags := request.Config.Get(ctx, &rq)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (d *enforcementsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func NewEnforcementsDataSource() datasource.DataSource {
	return &enforcementsDataSource{}
}