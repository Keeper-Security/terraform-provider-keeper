
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsVaultDataSource{}
)

var enforcementsVaultAttributes = map[string]schema.Attribute{
	"allow_secrets_manager": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Allow Keeper Secret Manager access",
	},
	"restrict_breach_watch": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict BreachWatch",
	},
	"send_breach_watch_events": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Send BreachWatch events",
	},
	"disable_setup_tour": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Disable setup-tour",
	},
	"days_before_deleted_records_auto_cleared": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"# of days before deleted records are automatically cleared",
	},
	"days_before_deleted_records_cleared_perm": schema.Int64Attribute{
		Optional:	true,
		Computed:	true,
		Description:	"# of days before deleted records are automatically cleared permanently",
	},
	"generated_security_question_complexity": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Generated security question complexity",
	},
	"generated_password_complexity": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Generated password complexity",
	},
	"mask_passwords_while_editing": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Mask passwords while editing",
	},
	"mask_notes": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Mask notes",
	},
	"mask_custom_fields": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Mask custom fields",
	},
	"restrict_create_identity_payment_records": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict creation of identity payment records",
	},
	"restrict_create_folder": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict creation of folders",
	},
}

type enforcementsVaultDataSource struct {
}

func NewEnforcementsVaultDataSource() datasource.DataSource {
	return &enforcementsVaultDataSource{}
}

func (d *enforcementsVaultDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_vault"
}

func (d *enforcementsVaultDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcementsVaultAttributes,
	}
}

func (d *enforcementsVaultDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *enforcementsVaultDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsVaultDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
