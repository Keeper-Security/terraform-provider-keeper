package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsAllowIpListDataSourceModel struct {
	TipZoneRestrictAllowedIpRanges types.List `tfsdk:"tip_zone_restrict_allowed_ip_ranges"`
	RestrictVaultIpAddresses       types.List `tfsdk:"restrict_vault_ip_addresses"`
	RestrictIpAddresses            types.List `tfsdk:"restrict_ip_addresses"`
}

func (eipm *EnforcementsAllowIpListDataSourceModel) ToKeeper(enforcements map[string]string) {
	getIpWhitelistValue(eipm.TipZoneRestrictAllowedIpRanges, "tip_zone_restrict_allowed_ip_ranges", enforcements)
	getIpWhitelistValue(eipm.RestrictVaultIpAddresses, "restrict_vault_ip_addresses", enforcements)
	getIpWhitelistValue(eipm.RestrictIpAddresses, "restrict_ip_addresses", enforcements)
}

func (eipm *EnforcementsAllowIpListDataSourceModel) FromKeeper(enforcements map[string]string) {
	setIpWhitelistValue(&eipm.TipZoneRestrictAllowedIpRanges, "tip_zone_restrict_allowed_ip_ranges", enforcements)
	setIpWhitelistValue(&eipm.RestrictVaultIpAddresses, "restrict_vault_ip_addresses", enforcements)
	setIpWhitelistValue(&eipm.RestrictIpAddresses, "restrict_ip_addresses", enforcements)
}

func (eipm *EnforcementsAllowIpListDataSourceModel) IsBlank() bool {
	if !eipm.TipZoneRestrictAllowedIpRanges.IsNull() {
		return false
	}
	if !eipm.RestrictVaultIpAddresses.IsNull() {
		return false
	}
	if !eipm.RestrictIpAddresses.IsNull() {
		return false
	}
	return true
}
