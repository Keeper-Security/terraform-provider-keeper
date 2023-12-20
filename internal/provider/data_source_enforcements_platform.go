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
	_ datasource.DataSource = &enforcementsPlatformDataSource{}
)

type enforcementsPlatformDataSource struct {
}

func newEnforcementsPlatformDataSource() datasource.DataSource {
	return &enforcementsPlatformDataSource{}
}

func (d *enforcementsPlatformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_platform"
}

func (d *enforcementsPlatformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsPlatformAttributes map[string]schema.Attribute
	enforcementsPlatformAttributes, diags = model.GenerateEnforcementDataSourceSchema(reflect.TypeOf((*model.EnforcementsPlatformDataSourceModel)(nil)))
	resp.Diagnostics.Append(diags...)

	resp.Schema = schema.Schema{
		Attributes: enforcementsPlatformAttributes,
	}
}

func (d *enforcementsPlatformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsPlatformDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
