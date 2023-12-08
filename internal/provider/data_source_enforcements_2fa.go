
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcements2faDataSource{}
)

var enforcements2faAttributes = map[string]schema.Attribute{
	"restrict_two_factor_channel_security_keys": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict use of security keys (FIDO2 WebAuthn) for 2fa",
	},
	"two_factor_duration_desktop": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"2fa duration for desktop client app",
	},
	"two_factor_duration_mobile": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"2fa duration for mobile client app",
	},
	"two_factor_duration_web": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"2fa duration for web client app",
	},
	"restrict_two_factor_channel_rsa": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict use of RSA SecurID for 2fa",
	},
	"restrict_two_factor_channel_duo": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict use of DUO for 2fa",
	},
	"restrict_two_factor_channel_dna": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict use of KeeperDNA for 2fa",
	},
	"restrict_two_factor_channel_google": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict use of Google Authenticator for 2fa",
	},
	"restrict_two_factor_channel_text": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict use of SMS/text message for 2fa",
	},
	"require_two_factor": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Require 2fa for login",
	},
}

type enforcements2faDataSource struct {
}

func NewEnforcements2faDataSource() datasource.DataSource {
	return &enforcements2faDataSource{}
}

func (d *enforcements2faDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_2fa"
}

func (d *enforcements2faDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcements2faAttributes,
	}
}

func (d *enforcements2faDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *enforcements2faDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.Enforcements2faDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
