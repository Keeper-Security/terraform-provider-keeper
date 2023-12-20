package model

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"time"
)

func GetUserStatus(u enterprise.IUser) (status string) {
	status = string(u.Status())
	if u.Status() == enterprise.UserStatus_Active {
		if u.Lock() != enterprise.UserLock_Unlocked {
			switch u.Lock() {
			case enterprise.UserLock_Locked:
				status = "locked"
			case enterprise.UserLock_Disabled:
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

type UserShortModel struct {
	EnterpriseUserId types.Int64  `tfsdk:"enterprise_user_id"`
	Username         types.String `tfsdk:"username"`
	Status           types.String `tfsdk:"status"`
}

func (usm *UserShortModel) FromKeeper(keeper enterprise.IUser) {
	usm.EnterpriseUserId = types.Int64Value(keeper.EnterpriseUserId())
	usm.Username = types.StringValue(keeper.Username())
	usm.Status = types.StringValue(GetUserStatus(keeper))
}

var UserShortSchemaAttributes = map[string]schema.Attribute{
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

type UserModel struct {
	EnterpriseUserId       types.Int64  `tfsdk:"enterprise_user_id"`
	Username               types.String `tfsdk:"username"`
	NodeId                 types.Int64  `tfsdk:"node_id"`
	Status                 types.String `tfsdk:"status"`
	AccountShareExpiration types.Int64  `tfsdk:"account_share_expiration"`
	TfaEnabled             types.Bool   `tfsdk:"tfa_enabled"`
	FullName               types.String `tfsdk:"full_name"`
	JobTitle               types.String `tfsdk:"job_title"`
}

func (model *UserModel) FromKeeper(keeper enterprise.IUser) {
	model.EnterpriseUserId = types.Int64Value(keeper.EnterpriseUserId())
	model.Username = types.StringValue(keeper.Username())
	model.NodeId = types.Int64Value(keeper.NodeId())
	model.Status = types.StringValue(GetUserStatus(keeper))
	model.TfaEnabled = types.BoolValue(keeper.TfaEnabled())
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

var UserSchemaAttributes = map[string]schema.Attribute{
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
		Description: "User Status: active | invited | locked | blocked | disabled",
	},
	"tfa_enabled": schema.BoolAttribute{
		Computed:    true,
		Description: "TFA Enabled flag",
	},
	"account_share_expiration": schema.Int64Attribute{
		Computed:    true,
		Description: "Account Share deadline: Timestamp in millis",
	},
	"full_name": schema.StringAttribute{
		Computed:    true,
		Description: "User Full Name",
	},
	"job_title": schema.StringAttribute{
		Computed:    true,
		Description: "User Job Title",
	},
}
