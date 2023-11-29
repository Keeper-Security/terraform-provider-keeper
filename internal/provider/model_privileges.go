package provider

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

var privilegesAttributes = map[string]schema.Attribute{
	"manage_nodes": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Nodes",
	},
	"manage_users": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Users",
	},
	"manage_teams": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Teams",
	},
	"manage_roles": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Roles",
	},
	"manage_reports": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Reporting and Alerts",
	},
	"manage_sso": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Bridge/SSO",
	},
	"device_approval": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Perform Device Approvals",
	},
	"manage_record_types": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Record Types in Vault",
		MarkdownDescription: "This permission allows the admin rights to create, edit, or delete Record Types " +
			"which have pre-defined fields. Record Types appear during creation of records in the user's vault.",
	},
	"share_admin": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Share Admin",
	},
	"run_compliance_reports": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Run Compliance Reports",
	},
	"transfer_account": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Transfer Account",
	},
	"manage_companies": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Manage Companies",
	},
}
