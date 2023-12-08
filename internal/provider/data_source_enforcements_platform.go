
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsPlatformDataSource{}
)

var enforcementsPlatformAttributes = map[string]schema.Attribute{
	"restrict_commander_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Commander",
	},
	"restrict_chat_mobile_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Chat for mobile",
	},
	"restrict_chat_desktop_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Chat for desktop",
	},
	"restrict_desktop_mac_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Desktop for MacOS",
	},
	"restrict_desktop_win_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Desktop for Windows",
	},
	"restrict_mobile_windows_phone_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Vault for Windows mobile",
	},
	"restrict_mobile_android_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Vault for Android",
	},
	"restrict_mobile_ios_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Vault for iOS",
	},
	"restrict_desktop_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Vault for desktop",
	},
	"restrict_mobile_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Vault for mobile",
	},
	"restrict_extensions_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper browser extensions",
	},
	"restrict_web_vault_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to Keeper Vault for web",
	},
}

type enforcementsPlatformDataSource struct {
}

func NewEnforcementsPlatformDataSource() datasource.DataSource {
	return &enforcementsPlatformDataSource{}
}

func (d *enforcementsPlatformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_platform"
}

func (d *enforcementsPlatformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcementsPlatformAttributes,
	}
}

func (d *enforcementsPlatformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
