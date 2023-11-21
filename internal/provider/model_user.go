package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"time"
)

func getUserStatus(u enterprise.IUser) (status string) {
	status = u.Status()
	if u.Status() == "active" {
		if u.Lock() > 0 {
			switch u.Lock() {
			case 1:
				status = "locked"
			case 2:
				status = "disabled"
			}
		} else if u.AccountShareExpiration() > 0 {
			if time.Now().UnixMilli() > u.AccountShareExpiration() {
				status = "blocked"
			} else {
				status = "pending"
			}
		}
	}
	return
}

type userShortModel struct {
	EnterpriseUserId types.Int64  `tfsdk:"enterprise_user_id"`
	Username         types.String `tfsdk:"username"`
	Status           types.String `tfsdk:"status"`
}

func (model *userShortModel) fromKeeper(keeper enterprise.IUser) {
	model.EnterpriseUserId = types.Int64Value(keeper.EnterpriseUserId())
	model.Username = types.StringValue(keeper.Username())
	model.Status = types.StringValue(getUserStatus(keeper))
}

var userShortSchemaAttributes = map[string]schema.Attribute{
	"enterprise_user_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Enterprise User ID",
	},
	"username": schema.StringAttribute{
		Computed:    true,
		Description: "Email",
	},
	"status": schema.StringAttribute{
		Computed:    true,
		Description: "User Status: active | invited | locked | blocked | pending",
	},
}

type userModel struct {
	EnterpriseUserId       types.Int64  `tfsdk:"enterprise_user_id"`
	Username               types.String `tfsdk:"username"`
	NodeId                 types.Int64  `tfsdk:"node_id"`
	Status                 types.String `tfsdk:"status"`
	Lock                   types.Int64  `tfsdk:"lock"` // TODO merge lock and status
	AccountShareExpiration types.Int64  `tfsdk:"account_share_expiration"`
	TfaEnabled             types.Bool   `tfsdk:"tfa_enabled"`
	//UserId                 types.Int64  `tfsdk:"user_id"`
	FullName types.String `tfsdk:"full_name"`
	JobTitle types.String `tfsdk:"job_title"`
}

func (model *userModel) fromKeeper(keeper enterprise.IUser) {
	model.EnterpriseUserId = types.Int64Value(keeper.EnterpriseUserId())
	model.Username = types.StringValue(keeper.Username())
	model.NodeId = types.Int64Value(keeper.NodeId())
	model.Status = types.StringValue(keeper.Status())
	model.Lock = types.Int64Value(int64(keeper.Lock()))
	model.TfaEnabled = types.BoolValue(keeper.TfaEnabled())
	//model.UserId = types.Int64Value(int64(keeper.UserId))
	if len(keeper.FullName()) > 0 {
		model.FullName = types.StringValue(keeper.FullName())
	} else {
		model.FullName = types.StringNull()
	}
	if len(keeper.JobTitle()) > 0 {
		model.JobTitle = types.StringValue(keeper.JobTitle())
	} else {
		model.JobTitle = types.StringNull()
	}
	if keeper.AccountShareExpiration() > 0 {
		model.AccountShareExpiration = types.Int64Value(keeper.AccountShareExpiration())
	} else {
		model.AccountShareExpiration = types.Int64Null()
	}
}

var userSchemaAttributes = map[string]schema.Attribute{
	"enterprise_user_id": schema.Int64Attribute{
		Computed:    true,
		Description: "Enterprise User ID",
	},
	"username": schema.StringAttribute{
		Computed:    true,
		Description: "Email",
	},
	"node_id": schema.Int64Attribute{
		Computed:    true,
		Description: "User Node ID",
	},
	"status": schema.StringAttribute{
		Computed:    true,
		Description: "User Status: active | invited",
	},
	"lock": schema.Int64Attribute{
		Computed:    true,
		Description: "User Lock Status: 0 - unlocked, 1 - locked, 2 - disabled",
	},
	"tfa_enabled": schema.BoolAttribute{
		Computed:    true,
		Description: "TFA Enabled flag",
	},
	"account_share_expiration": schema.Int64Attribute{
		Computed:    true,
		Description: "Account Share deadline: Timestamp in millis",
	},
	//"user_id": schema.Int64Attribute{
	//	Computed:    true,
	//	Description: "Keeper User ID",
	//},
	"full_name": schema.StringAttribute{
		Computed:    true,
		Description: "User Full Name",
	},
	"job_title": schema.StringAttribute{
		Computed:    true,
		Description: "User Job Title",
	},
}
