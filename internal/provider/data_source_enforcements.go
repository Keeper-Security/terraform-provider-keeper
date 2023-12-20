package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-keeper/internal/model"
)

var (
	_ datasource.DataSource = &enforcementsDataSource{}
)

type enforcementsDataSource struct {
}

func (d *enforcementsDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_enforcements"
}

func (d *enforcementsDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	var diags diag.Diagnostics
	var enforcementsAttributes map[string]schema.Attribute
	enforcementsAttributes, diags = model.EnforcementsDataSourceAttributes()
	response.Diagnostics.Append(diags...)

	response.Schema = schema.Schema{
		Attributes:  enforcementsAttributes,
		Description: "Role Enforcements",
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

func newEnforcementsDataSource() datasource.DataSource {
	return &enforcementsDataSource{}
}
