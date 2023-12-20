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
	_ datasource.DataSource = &enforcementsAllowIpListDataSource{}
)

type enforcementsAllowIpListDataSource struct {
}

func newEnforcementsAllowIpListDataSource() datasource.DataSource {
	return &enforcementsAllowIpListDataSource{}
}

func (d *enforcementsAllowIpListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_allow_ip_list"
}

func (d *enforcementsAllowIpListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsAllowIpListAttributes map[string]schema.Attribute
	enforcementsAllowIpListAttributes, diags = model.GenerateEnforcementDataSourceSchema(reflect.TypeOf((*model.EnforcementsAllowIpListDataSourceModel)(nil)))
	resp.Diagnostics.Append(diags...)
	resp.Schema = schema.Schema{
		Attributes: enforcementsAllowIpListAttributes,
	}
}

func (d *enforcementsAllowIpListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsAllowIpListDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
