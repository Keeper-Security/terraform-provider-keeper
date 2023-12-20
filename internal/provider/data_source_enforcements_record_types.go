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
	_ datasource.DataSource = &enforcementsRecordTypesDataSource{}
)

type enforcementsRecordTypesDataSource struct {
}

func newEnforcementsRecordTypesDataSource() datasource.DataSource {
	return &enforcementsRecordTypesDataSource{}
}

func (d *enforcementsRecordTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_record_types"
}

func (d *enforcementsRecordTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsRecordTypesAttributes map[string]schema.Attribute
	enforcementsRecordTypesAttributes, diags = model.GenerateEnforcementDataSourceSchema(reflect.TypeOf((*model.EnforcementsRecordTypesDataSourceModel)(nil)))
	resp.Diagnostics.Append(diags...)

	resp.Schema = schema.Schema{
		Attributes: enforcementsRecordTypesAttributes,
	}
}

func (d *enforcementsRecordTypesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *enforcementsRecordTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsRecordTypesDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
