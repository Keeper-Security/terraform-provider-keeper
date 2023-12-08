
package provider

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    //"github.com/hashicorp/terraform-plugin-framework/types"
    "terraform-provider-kepr/internal/model"
)
    

var (
	_ datasource.DataSource = &enforcementsSharingDataSource{}
)

var enforcementsSharingAttributes = map[string]schema.Attribute{
	"restrict_sharing_outside_of_isolated_nodes": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict sharing outside of isolated nodes",
	},
	"restrict_link_sharing": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict link-sharing",
	},
	"restrict_sharing_record_with_attachments": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict sharing records with attachments",
	},
	"restirct_sharing_record_and_folder": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict sharing records and folders",
	},
	"restrict_sharing_incoming_all": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict all incoming-shares0",
	},
	"require_account_share": schema.StringAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Require account-share",
	},
	"restrict_file_upload": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict file uploads",
	},
	"restrict_import": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict record imports",
	},
	"restrict_export": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict record exports",
	},
	"restrict_sharing_enterprise": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict sharing within the enterprise",
	},
	"restrict_sharing_all": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict all sharing",
	},
	"restrict_create_record_to_shared_folders": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict record creation within shared-folders",
	},
	"restrict_create_record": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict creation of records",
	},
	"restrict_create_folder_to_only_shared_folders": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict creation of folders to within shared-folders only",
	},
	"restrict_direct_sharing": schema.BoolAttribute{
		Optional:	true,
		Computed:	true,
		Description:	"Restrict direct-shares",
	},
}

type enforcementsSharingDataSource struct {
}

func NewEnforcementsSharingDataSource() datasource.DataSource {
	return &enforcementsSharingDataSource{}
}

func (d *enforcementsSharingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enforcements_sharing"
}

func (d *enforcementsSharingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: enforcementsSharingAttributes,
	}
}

func (d *enforcementsSharingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *enforcementsSharingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var rq model.EnforcementsSharingDataSourceModel
	diags := req.Config.Get(ctx, &rq)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state = rq
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
