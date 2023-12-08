package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnforcementsRecordTypesDataSourceModel struct {
	RestrictRecordTypes	types.String	`tfsdk:"restrict_record_types"`
}

func (model *EnforcementsRecordTypesDataSourceModel) FromKeeper(enforcements map[string]string) {
	model.RestrictRecordTypes = types.StringValue(enforcements["restrict_record_types"])
}