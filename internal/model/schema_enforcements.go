package model

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"reflect"
	"strings"
)

var EnforcementLookup map[string]enterprise.IEnforcement
var availableDurations []string

func init() {
	EnforcementLookup = make(map[string]enterprise.IEnforcement)
	enterprise.AvailableRoleEnforcements(func(enforcement enterprise.IEnforcement) bool {
		EnforcementLookup[strings.ToLower(enforcement.Name())] = enforcement
		return true
	})
	for k := range tfaDurations {
		availableDurations = append(availableDurations, k)
	}
}

func extractTfSdkFields(modelType reflect.Type, cb func(string)) (err error) {
	if modelType == nil {
		return
	}
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		err = fmt.Errorf("type \"%T\" should be a \"struct\"", modelType)
		return
	}
	var fNum = modelType.NumField()
	for i := 0; i < fNum; i++ {
		var f = modelType.Field(i)
		var v = f.Tag.Get("tfsdk")
		if len(v) > 0 {
			var comps = strings.Split(v, ",")
			cb(comps[0])
		}
	}
	return
}

func GenerateEnforcementResourceSchema(modelType reflect.Type) (rs map[string]rsschema.Attribute, diags diag.Diagnostics) {
	rs = make(map[string]rsschema.Attribute)
	var err = extractTfSdkFields(modelType, func(name string) {
		if enf, ok := EnforcementLookup[strings.ToLower(name)]; ok {
			var description = enforcementDescriptions[name]
			switch strings.ToLower(enf.ValueType()) {
			case "boolean":
				rs[name] = rsschema.BoolAttribute{
					Optional:    true,
					Description: description,
				}
			case "long":
				rs[name] = rsschema.Int64Attribute{
					Optional:    true,
					Description: description,
				}
			case "ternary_edn", "ternary_den":
				rs[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
					Validators:  []validator.String{stringvalidator.OneOf("enforce", "disable", "e", "d")},
				}
			case "ip_whitelist":
				rs[name] = rsschema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: description,
				}
			case "record_types":
				rs[name] = rsschema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: description,
				}
			case "two_factor_duration":
				rs[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
					Validators:  []validator.String{stringvalidator.OneOf(availableDurations...)},
				}
			case "account_share":
				rs[name] = rsschema.Int64Attribute{
					Optional:    true,
					Description: description,
				}
			case "string":
				rs[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
				}
			default:
				rs[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
				}
			}
		} else {
			diags.AddError(fmt.Sprintf("Resource schema: enforcement \"%s\" description not found", name),
				fmt.Sprintf("Model \"%s\" defines unsupported enforcement \"%s\"", modelType.String(), name))
		}
	})
	if err != nil {
		diags.AddError("Resource schema: extract fields error", err.Error())
	}
	return
}

func EnforcementResourceAttributes() (rs map[string]rsschema.Attribute, diags diag.Diagnostics) {
	rs = make(map[string]rsschema.Attribute)
	var attributes map[string]rsschema.Attribute
	var d diag.Diagnostics

	var enforcementType = reflect.TypeOf((*EnforcementsDataSourceModel)(nil))
	enforcementType = enforcementType.Elem()

	var fNum = enforcementType.NumField()
	for i := 0; i < fNum; i++ {
		var f = enforcementType.Field(i)
		var v = f.Tag.Get("tfsdk")
		if len(v) > 0 {
			var comps = strings.Split(v, ",")
			attributes, d = GenerateEnforcementResourceSchema(f.Type)
			diags.Append(d...)
			var name = comps[0]
			rs[name] = rsschema.SingleNestedAttribute{
				Attributes:  attributes,
				Optional:    true,
				Description: enforcementDescriptions[name],
			}
		}
	}

	return
}

type noFalseValueValidator struct{}

func (bv noFalseValueValidator) Description(ctx context.Context) string {
	return bv.MarkdownDescription(ctx)
}

func (bv noFalseValueValidator) MarkdownDescription(_ context.Context) string {
	return "Ensure boolean data source enforcement does not have \"false\" value"
}
func (bv noFalseValueValidator) ValidateBool(_ context.Context, rq validator.BoolRequest, rs *validator.BoolResponse) {
	if !rq.ConfigValue.IsNull() && !rq.ConfigValue.ValueBool() {
		rs.Diagnostics.AddError("Invalid boolean enforcement value: cannot be false",
			fmt.Sprintf("Boolean enforcement \"%s\" cannot have \"false\" value. Delete it (comment out) to unset", rq.Path.String()))
	}
}

func GenerateEnforcementDataSourceSchema(modelType reflect.Type) (ds map[string]dsschema.Attribute, diags diag.Diagnostics) {
	ds = make(map[string]dsschema.Attribute)
	var err = extractTfSdkFields(modelType, func(name string) {
		if enf, ok := EnforcementLookup[strings.ToLower(name)]; ok {
			var description = enforcementDescriptions[name]
			switch strings.ToLower(enf.ValueType()) {
			case "boolean":
				ds[name] = rsschema.BoolAttribute{
					Optional:    true,
					Description: description,
					Validators:  []validator.Bool{new(noFalseValueValidator)},
				}
			case "long":
				ds[name] = rsschema.Int64Attribute{
					Optional:    true,
					Description: description,
				}
			case "ternary_edn", "ternary_den":
				ds[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
					Validators:  []validator.String{stringvalidator.OneOf("enforce", "disable", "e", "d")},
				}
			case "ip_whitelist":
				ds[name] = rsschema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: description,
				}
			case "record_types":
				ds[name] = rsschema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: description,
				}
			case "two_factor_duration":
				ds[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
					Validators:  []validator.String{stringvalidator.OneOf(availableDurations...)},
				}
			case "account_share":
				ds[name] = rsschema.Int64Attribute{
					Optional:    true,
					Description: description,
				}
			case "string":
				ds[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
				}
			default:
				ds[name] = rsschema.StringAttribute{
					Optional:    true,
					Description: description,
				}
			}
		} else {
			diags.AddError(fmt.Sprintf("Data source schema: enforcement \"%s\" not found", name),
				fmt.Sprintf("Model \"%s\" defines unsupported enforcement \"%s\"", modelType.Name(), name))
		}
	})
	if err != nil {
		diags.AddError(" Data source schema: extract fields error", err.Error())
	}
	return
}

func EnforcementsDataSourceAttributes() (rs map[string]dsschema.Attribute, diags diag.Diagnostics) {
	rs = make(map[string]dsschema.Attribute)
	var attributes map[string]dsschema.Attribute
	var d diag.Diagnostics

	var enforcementType = reflect.TypeOf((*EnforcementsDataSourceModel)(nil))
	enforcementType = enforcementType.Elem()

	var fNum = enforcementType.NumField()
	for i := 0; i < fNum; i++ {
		var f = enforcementType.Field(i)
		var v = f.Tag.Get("tfsdk")
		if len(v) > 0 {
			var comps = strings.Split(v, ",")
			attributes, d = GenerateEnforcementDataSourceSchema(f.Type)
			diags.Append(d...)
			var name = comps[0]
			rs[name] = dsschema.SingleNestedAttribute{
				Attributes:  attributes,
				Optional:    true,
				Computed:    true,
				Description: enforcementDescriptions[name],
			}
		}
	}

	return
}
