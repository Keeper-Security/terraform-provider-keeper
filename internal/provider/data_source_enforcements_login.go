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
	_ datasource.DataSource = &enforcementsLoginDataSource{}
)

type enforcementsLoginDataSource struct {
}

func newEnforcementsLoginDataSource() datasource.DataSource {
	return &enforcementsLoginDataSource{}
}

func (d *enforcementsLoginDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_login"
}

func (d *enforcementsLoginDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsLoginAttributes map[string]schema.Attribute
	enforcementsLoginAttributes, diags = model.GenerateEnforcementDataSourceSchema(reflect.TypeOf((*model.EnforcementsLoginDataSourceModel)(nil)))
	resp.Diagnostics.Append(diags...)
	resp.Schema = schema.Schema{
		Attributes: enforcementsLoginAttributes,
	}
}

func (d *enforcementsLoginDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsLoginDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
