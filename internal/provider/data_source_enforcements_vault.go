package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"reflect"
	"terraform-provider-keeper/internal/model"
)

var (
	_ datasource.DataSource = &enforcementsVaultDataSource{}
)

type enforcementsVaultDataSource struct {
}

func newEnforcementsVaultDataSource() datasource.DataSource {
	return &enforcementsVaultDataSource{}
}

func (d *enforcementsVaultDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_vault"
}

func (d *enforcementsVaultDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsVaultAttributes map[string]schema.Attribute
	enforcementsVaultAttributes, diags = model.GenerateEnforcementDataSourceSchema(reflect.TypeOf((*model.EnforcementsVaultDataSourceModel)(nil)))
	resp.Diagnostics.Append(diags...)
	resp.Schema = schema.Schema{
		Attributes: enforcementsVaultAttributes,
	}
}

func (d *enforcementsVaultDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsVaultDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
