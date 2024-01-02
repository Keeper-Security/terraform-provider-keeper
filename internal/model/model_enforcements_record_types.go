package model

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/api"
	"github.com/keeper-security/keeper-sdk-golang/vault"
	"sort"
	"strings"
)

type EnforcementsRecordTypesDataSourceModel struct {
	RestrictRecordTypes types.Set `tfsdk:"restrict_record_types"`
}

type RecordTypeEnforcement struct {
	Std []int64 `json:"std"`
	Ent []int64 `json:"ent"`
}

func (ertm *EnforcementsRecordTypesDataSourceModel) VerifyNames(keeperRecordTypes []vault.IRecordType) (diags diag.Diagnostics) {
	if !ertm.RestrictRecordTypes.IsNull() {
		var ok bool
		var sv types.String
		var names = make(map[string]string)
		for _, v := range keeperRecordTypes {
			if v.Scope() == vault.RecordTypeScope_Standard || v.Scope() == vault.RecordTypeScope_Enterprise {
				names[strings.ToLower(v.Name())] = v.Name()
			}
		}
		for _, e := range ertm.RestrictRecordTypes.Elements() {
			if sv, ok = e.(types.String); ok {
				var name = sv.ValueString()
				var rName string
				if rName, ok = names[strings.ToLower(name)]; ok {
					if rName != name {
						diags.AddError(fmt.Sprintf("Invalid record type: %s", name),
							fmt.Sprintf("Record type name are case-sensitive. Valid name is \"%s\"", rName))
					}
				} else {
					diags.AddError(fmt.Sprintf("Invalid record type: %s", name),
						fmt.Sprintf("Record type \"%s\" does not exist", name))
				}
			} else {
				diags.AddError(fmt.Sprintf("Invalid record type: %s", e.String()),
					"Record type name should be a string")
			}
		}
	}
	return
}

func (ertm *EnforcementsRecordTypesDataSourceModel) FromKeeper(enforcements map[string]string, keeperRecordTypes []vault.IRecordType) {
	var ok bool
	var err error
	var strValue string
	if strValue, ok = enforcements["restrict_record_types"]; ok {
		var rte = new(RecordTypeEnforcement)
		if err = json.Unmarshal([]byte(strValue), rte); err == nil {
			var recordTypes = api.NewSet[string]()
			var rtl map[int64]string
			if len(rte.Std) > 0 {
				rtl = make(map[int64]string)
				for _, v := range keeperRecordTypes {
					if v.Scope() == vault.RecordTypeScope_Standard {
						rtl[v.Id()] = v.Name()
					}
				}
				for _, rtId := range rte.Std {
					if strValue, ok = rtl[rtId]; ok {
						recordTypes.Add(strValue)
					}
				}
			}
			if len(rte.Ent) > 0 {
				rtl = make(map[int64]string)
				for _, v := range keeperRecordTypes {
					if v.Scope() == vault.RecordTypeScope_Enterprise {
						rtl[v.Id()] = v.Name()
					}
				}
				for _, rtId := range rte.Ent {
					if strValue, ok = rtl[rtId]; ok {
						recordTypes.Add(strValue)
					}
				}
			}
			if len(recordTypes) > 0 {
				ertm.RestrictRecordTypes, _ = types.SetValue(types.StringType, api.SliceSelect(recordTypes.ToArray(), func(x string) attr.Value {
					return types.StringValue(x)
				}))
				return
			}
		}
	}
	ertm.RestrictRecordTypes = types.SetNull(types.StringType)
}

func (ertm *EnforcementsRecordTypesDataSourceModel) ToKeeper(enforcements map[string]string, keeperRecordTypes []vault.IRecordType) {
	if !ertm.RestrictRecordTypes.IsNull() {
		var rts = make(map[string]int64)
		var rte = make(map[string]int64)
		for _, v := range keeperRecordTypes {
			if v.Scope() == vault.RecordTypeScope_Standard {
				rts[strings.ToLower(v.Name())] = v.Id()
			}
			if v.Scope() == vault.RecordTypeScope_Enterprise {
				rte[strings.ToLower(v.Name())] = v.Id()
			}
		}
		var enf = new(RecordTypeEnforcement)
		enf.Std = make([]int64, 0)
		enf.Ent = make([]int64, 0)
		var ok bool
		var irt int64
		var sv types.String
		for _, e := range ertm.RestrictRecordTypes.Elements() {
			if sv, ok = e.(types.String); ok {
				var rt = strings.ToLower(sv.ValueString())
				if irt, ok = rts[rt]; ok {
					enf.Std = append(enf.Std, irt)
				} else if irt, ok = rte[rt]; ok {
					enf.Ent = append(enf.Ent, irt)
				}
			}
		}
		if len(enf.Std) > 0 || len(enf.Ent) > 0 {
			if len(enf.Std) > 1 {
				sort.Slice(enf.Std, func(i, j int) bool {
					return i < j
				})
			}
			if len(enf.Ent) > 1 {
				sort.Slice(enf.Ent, func(i, j int) bool {
					return i < j
				})
			}
			if data, err := json.Marshal(enf); err == nil {
				enforcements["restrict_record_types"] = string(data)
			}
		}
	}
}

func (ertm *EnforcementsRecordTypesDataSourceModel) IsBlank() bool {
	return ertm.RestrictRecordTypes.IsNull() || len(ertm.RestrictRecordTypes.Elements()) == 0
}
