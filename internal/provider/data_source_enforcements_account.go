
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsAccountDataSource{}
)

var enforcementsAccountAttributes = map[string]schema.Attribute{
	"restrict_account_recovery": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict Account Recovery",
	},
	"resend_enterprise_invite_in_x_days": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Resend enterprise invite in X days",
	},
	"restrict_ip_autoapproval": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict IP auto-approval",
	},
	"disallow_v2_clients": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Disallow v2 clients",
	},
	"disable_onboarding": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Disable onboarding",
	},
	"restrict_personal_license": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict use of personal license",
	},
	"logout_timer_desktop": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Automatic-logout delay for desktop client",
	},
	"logout_timer_mobile": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Automatic-logout delay for mobile client",
	},
	"logout_timer_web": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Automatic-logout delay for web client",
	},
	"restrict_email_change": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict change of email address",
	},
	"send_invite_at_registration": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Send invite at registration",
	},
	"restrict_offline_access": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict offline access",
	},
	"automatic_backup_every_x_days": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"# of days between automatic-backups",
	},
	"require_account_recovery_approval": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Require account recovery approval",
	},
	"require_device_approval": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Require device approval",
	},
	"restrict_persistent_login": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict persistent login",
	},
	"max_session_login_time": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Max session login time",
	},
	"minimum_pbkdf2_iterations": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Minimum PBKDF2 iterations",
	},
	"require_security_key_pin": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Require security key PIN",
	},
	"restrict_import_shared_folders": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict shared-folder imports",
	},
	"allow_pam_discovery": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Allow PAM discovery",
	},
	"allow_pam_rotation": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Allow PAM rotation",
	},
	"maximum_record_size": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"Maximum record-size",
	},
	"require_self_destruct": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Require self-destruct",
	},
	"stay_logged_in_default": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Enable staying logged-in by default",
	},
}

type enforcementsAccountDataSource struct {
}

func NewEnforcementsAccountDataSource() datasource.DataSource {
	return &enforcementsAccountDataSource{}
}

func (d *enforcementsAccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_account"
}

func (d *enforcementsAccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcementsAccountAttributes,
	}
}

func (d *enforcementsAccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *enforcementsAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsAccountDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
