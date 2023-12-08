
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsKeeperFillDataSource{}
)

var enforcementsKeeperFillAttributes = map[string]schema.Attribute{
	"keeper_fill_auto_suggest": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Keeper auto-fill suggestion",
	},
	"restrict_http_fill_warning": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict HTTP fill warning",
	},
	"restrict_prompt_to_disable": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict prompt to disable Keeper Fill",
	},
	"keeper_fill_match_on_subdomain": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Keeper Fill subdomains to match on",
	},
	"keeper_fill_auto_submit": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Keeper Fill auto-submit",
	},
	"keeper_fill_auto_fill": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Keeper Fill auto-fill",
	},
	"keeper_fill_hover_locks": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Keeper Fill hover locks",
	},
	"master_password_reentry": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Master password re-entry",
	},
	"restrict_auto_fill": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict auto-fill",
	},
	"restrict_prompt_to_change": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict prompt to change",
	},
	"restrict_prompt_to_save": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict prmpt to save",
	},
	"restrict_auto_submit": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict auto-submit",
	},
	"restrict_prompt_to_fill": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict prompt to fill",
	},
	"restrict_prompt_to_login": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict prompt to login",
	},
	"restrict_hover_locks": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict hover-locks",
	},
	"restrict_domain_create": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict creation of new records for domain(s)",
	},
	"restrict_domain_access": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict access to domain(s)",
	},
}

type enforcementsKeeperFillDataSource struct {
}

func NewEnforcementsKeeperFillDataSource() datasource.DataSource {
	return &enforcementsKeeperFillDataSource{}
}

func (d *enforcementsKeeperFillDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_keeper_fill"
}

func (d *enforcementsKeeperFillDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcementsKeeperFillAttributes,
	}
}

func (d *enforcementsKeeperFillDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
