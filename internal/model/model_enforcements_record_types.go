package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsRecordTypesDataSourceModel struct {
	RestrictRecordTypes types.String `tfsdk:"restrict_record_types"`
}

func (ertm *EnforcementsRecordTypesDataSourceModel) FromKeeper(enforcements map[string]string) {
	setStringValue(&ertm.RestrictRecordTypes, "restrict_record_types", enforcements)
}

func (ertm *EnforcementsRecordTypesDataSourceModel) ToKeeper(enforcements map[string]string) {
	getStringValue(ertm.RestrictRecordTypes, "restrict_record_types", enforcements)
}

func (ertm *EnforcementsRecordTypesDataSourceModel) IsBlank() bool {
	return ertm.RestrictRecordTypes.IsNull()
}
