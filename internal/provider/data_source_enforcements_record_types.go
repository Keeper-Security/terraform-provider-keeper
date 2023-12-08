
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsRecordTypesDataSource{}
)

var enforcementsRecordTypesAttributes = map[string]schema.Attribute{
	"restrict_record_types": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict record-types",
	},
}

type enforcementsRecordTypesDataSource struct {
}

func NewEnforcementsRecordTypesDataSource() datasource.DataSource {
	return &enforcementsRecordTypesDataSource{}
}

func (d *enforcementsRecordTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_record_types"
}

func (d *enforcementsRecordTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
