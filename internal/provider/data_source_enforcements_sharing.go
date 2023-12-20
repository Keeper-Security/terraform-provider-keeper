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
	_ datasource.DataSource = &enforcementsSharingDataSource{}
)

type enforcementsSharingDataSource struct {
}

func newEnforcementsSharingDataSource() datasource.DataSource {
	return &enforcementsSharingDataSource{}
}

func (d *enforcementsSharingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_sharing"
}

func (d *enforcementsSharingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsSharingAttributes map[string]schema.Attribute
	enforcementsSharingAttributes, diags = model.GenerateEnforcementDataSourceSchema(reflect.TypeOf((*model.EnforcementsSharingDataSourceModel)(nil)))
	resp.Diagnostics.Append(diags...)
	resp.Schema = schema.Schema{
		Attributes: enforcementsSharingAttributes,
	}
}

func (d *enforcementsSharingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsSharingDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
