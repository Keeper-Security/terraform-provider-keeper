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
	_ datasource.DataSource = &enforcementsKeeperFillDataSource{}
)

type enforcementsKeeperFillDataSource struct {
}

func newEnforcementsKeeperFillDataSource() datasource.DataSource {
	return &enforcementsKeeperFillDataSource{}
}

func (d *enforcementsKeeperFillDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_keeper_fill"
}

func (d *enforcementsKeeperFillDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsKeeperFillAttributes map[string]schema.Attribute
	enforcementsKeeperFillAttributes, diags = model.GenerateEnforcementDataSourceSchema(reflect.TypeOf((*model.EnforcementsKeeperFillDataSourceModel)(nil)))
	resp.Diagnostics.Append(diags...)
	resp.Schema = schema.Schema{
		Attributes: enforcementsKeeperFillAttributes,
	}
}

func (d *enforcementsKeeperFillDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsKeeperFillDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
