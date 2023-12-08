package model

type EnforcementsDataSourceModel struct {
	Account	*EnforcementsAccountDataSourceModel	`tfsdk:"account"`
	AllowIpList	*EnforcementsAllowIpListDataSourceModel	`tfsdk:"allow_ip_list"`
	Sharing	*EnforcementsSharingDataSourceModel	`tfsdk:"sharing"`
	KeeperFill	*EnforcementsKeeperFillDataSourceModel	`tfsdk:"keeper_fill"`
	Login	*EnforcementsLoginDataSourceModel	`tfsdk:"login"`
	Platform 	*EnforcementsPlatformDataSourceModel	`tfsdk:"platform"`
	RecordTypes *EnforcementsRecordTypesDataSourceModel	`tfsdk:"record_types"`
	TwoFactorAuthentication *Enforcements2faDataSourceModel	`tfsdk:"two_factor_authentication"`
	Vault 	*EnforcementsVaultDataSourceModel	`tfsdk:"vault"`
}

func (model *EnforcementsDataSourceModel) FromKeeper(enforcements map[string]string) {
	if enforcements == nil {
		return
	}

	model.Account = new(EnforcementsAccountDataSourceModel)
	model.Account.FromKeeper(enforcements)
	model.AllowIpList = new(EnforcementsAllowIpListDataSourceModel)
	model.AllowIpList.FromKeeper(enforcements)
	model.Sharing = new(EnforcementsSharingDataSourceModel)
	model.Sharing.FromKeeper(enforcements)
	model.KeeperFill = new(EnforcementsKeeperFillDataSourceModel)
	model.KeeperFill.FromKeeper(enforcements)
	model.Login = new(EnforcementsLoginDataSourceModel)
	model.Login.FromKeeper(enforcements)
	model.Platform = new(EnforcementsPlatformDataSourceModel)
	model.Platform.FromKeeper(enforcements)
	model.RecordTypes = new(EnforcementsRecordTypesDataSourceModel)
	model.RecordTypes.FromKeeper(enforcements)
	model.TwoFactorAuthentication = new(Enforcements2faDataSourceModel)
	model.TwoFactorAuthentication.FromKeeper(enforcements)
	model.Vault = new(EnforcementsVaultDataSourceModel)
	model.Vault.FromKeeper(enforcements)
}