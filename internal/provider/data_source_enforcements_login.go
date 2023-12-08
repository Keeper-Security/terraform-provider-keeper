
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsLoginDataSource{}
)

var enforcementsLoginAttributes = map[string]schema.Attribute{
	"allow_alternate_passwords": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Allow alternate passwords",
	},
	"restrict_windows_fingerprint": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict Windows fingerprint login",
	},
	"restrict_android_fingerprint": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict Android fingerprint login",
	},
	"restrict_mac_fingerprint": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict MacOS fingerprint login",
	},
	"restrict_ios_fingerprint": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict iOS fingerprint login",
	},
	"master_password_expired_as_of": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Master password expiration",
	},
	"master_password_maximum_days_before_change": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Maximum days before master password change",
	},
	"master_password_restrict_days_before_reuse": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"# of days before master password can be re-used",
	},
	"master_password_minimum_digits": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Minimum # of digits required for master password",
	},
	"master_password_minimum_lower": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	" Minimum # of lower-case characters required for master password",
	},
	"master_password_minimum_upper": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Minimum # of upper-case characters required for master password",
	},
	"master_password_minimum_special": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Minimum # of special characters required for master password",
	},
	"master_password_minimum_length": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Minimum length required for master password",
	},
}

type enforcementsLoginDataSource struct {
}

func NewEnforcementsLoginDataSource() datasource.DataSource {
	return &enforcementsLoginDataSource{}
}

func (d *enforcementsLoginDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_login"
}

func (d *enforcementsLoginDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcementsLoginAttributes,
	}
}

func (d *enforcementsLoginDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
