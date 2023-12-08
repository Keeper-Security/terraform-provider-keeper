package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsAllowIpListDataSourceModel struct {
	TwoFactorByIp	types.String	`tfsdk:"two_factor_by_ip"`
	TipZoneRestrictAllowedIpRanges	types.String	`tfsdk:"tip_zone_restrict_allowed_ip_ranges"`
	RestrictVaultIpAddresses	types.String	`tfsdk:"restrict_vault_ip_addresses"`
	RestrictIpAddresses	types.String	`tfsdk:"restrict_ip_addresses"`
}

func (model *EnforcementsAllowIpListDataSourceModel) FromKeeper(enforcements map[string]string) {
	model.TwoFactorByIp = types.StringValue(enforcements["two_factor_by_ip"])
	model.TipZoneRestrictAllowedIpRanges = types.StringValue(enforcements["tip_zone_restrict_allowed_ip_ranges"])
	model.RestrictVaultIpAddresses = types.StringValue(enforcements["restrict_vault_ip_addresses"])
	model.RestrictIpAddresses = types.StringValue(enforcements["restrict_ip_addresses"])
}