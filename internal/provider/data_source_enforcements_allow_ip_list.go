
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsAllowIpListDataSource{}
)

var enforcementsAllowIpListAttributes = map[string]schema.Attribute{
	"two_factor_by_ip": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"2fa by IP",
	},
	"tip_zone_restrict_allowed_ip_ranges": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict allowed IP ranges for tip zone",
	},
	"restrict_vault_ip_addresses": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict vault-access from IP addresses",
	},
	"restrict_ip_addresses": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict IP addresses",
	},
}

type enforcementsAllowIpListDataSource struct {
}

func NewEnforcementsAllowIpListDataSource() datasource.DataSource {
	return &enforcementsAllowIpListDataSource{}
}

func (d *enforcementsAllowIpListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_allow_ip_list"
}

func (d *enforcementsAllowIpListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcementsAllowIpListAttributes,
	}
}

func (d *enforcementsAllowIpListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
