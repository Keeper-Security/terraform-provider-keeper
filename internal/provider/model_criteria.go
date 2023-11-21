package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type nodeCriteria struct {
	NodeId  types.Int64 `tfsdk:"node_id"`
	Cascade types.Bool  `tfsdk:"cascade"`
}

var nodeCriteriaAttributes = map[string]schema.Attribute{
	"node_id": schema.Int64Attribute{
		Required:    true,
		Description: "Base Node ID",
	},
	"cascade": schema.BoolAttribute{
		Optional:    true,
		Description: "Include subnodes",
	},
}

type nodeMatcher = func(int64) bool

func getNodeMatcher(nc *nodeCriteria, nodes enterprise.IEnterpriseEntity[enterprise.INode, int64]) (mf nodeMatcher, diags diag.Diagnostics) {
	if nc != nil {
		var rootNodeId = nc.NodeId.ValueInt64()
		var n = nodes.GetEntity(rootNodeId)
		if n == nil {
			diags.AddError("Create a node matcher error.",
				fmt.Sprintf("Node ID \"%d\" not found", rootNodeId))
			return
		}
		var subNodes []int64
		var subnodeLookup = make(map[int64][]int64)
		nodes.GetAllEntities(func(node enterprise.INode) bool {
			subNodes = subnodeLookup[node.ParentId()]
			subNodes = append(subNodes, node.NodeId())
			subnodeLookup[node.ParentId()] = subNodes
			return true
		})

		if nc.Cascade.ValueBool() {
			var added = make(map[int64]bool)
			var nodeIds = append(([]int64)(nil), rootNodeId)
			added[rootNodeId] = true
			var pos = 0
			var ok bool
			for pos < len(nodeIds) {
				var nodeId = nodeIds[pos]
				pos++
				if subNodes, ok = subnodeLookup[nodeId]; ok {
					for _, nodeId = range subNodes {
						if _, ok = added[nodeId]; !ok {
							nodeIds = append(nodeIds, nodeId)
							added[nodeId] = true
						}
					}
				}
			}
			mf = func(nodeId int64) (ok bool) {
				_, ok = added[nodeId]
				return
			}
		} else {
			mf = func(nodeId int64) (ok bool) {
				ok = nodeId == rootNodeId
				return
			}
		}
	}
	return
}

type filterCriteria struct {
	Field types.String `tfsdk:"field"`
	Cmp   types.String `tfsdk:"cmp"`
	Value types.String `tfsdk:"value"`
}

var filterCriteriaAttributes = map[string]schema.Attribute{
	"field": schema.StringAttribute{
		Required:    true,
		Description: "Field Name",
	},
	"cmp": schema.StringAttribute{
		Optional:    true,
		Description: "Comparison operation",
		//		Description: "Comparison operation: ==, !=, >, <, <=, >=, matches, startswith, endswith",
	},
	"value": schema.StringAttribute{
		Required:    true,
		Description: "Comparison value",
	},
}

type IntOp = int32

const (
	IntOp_Invalid IntOp = iota
	IntOp_Equal
	IntOp_NotEqual
	IntOp_LessThan
	IntOp_GreaterThan
	IntOp_GreaterOrEqual
	IntOp_LessOrEqual
)

var (
	integerOperators = map[string]IntOp{
		"==": IntOp_Equal,
		"!=": IntOp_NotEqual,
		"<":  IntOp_LessThan,
		">":  IntOp_GreaterThan,
		"<=": IntOp_LessOrEqual,
		">=": IntOp_GreaterOrEqual,
	}
)

type StrOp = int32

const (
	StrOp_Invalid StrOp = iota
	StrOp_Equal
	StrOp_NotEqual
	StrOp_StartsWith
	StrOp_EndsWith
	StrOp_Matches
)

var (
	stringOperators = map[string]StrOp{
		"==":      StrOp_Equal,
		"!=":      StrOp_NotEqual,
		"starts":  StrOp_StartsWith,
		"ends":    StrOp_EndsWith,
		"matches": StrOp_Matches,
	}
)

type matcher = func(interface{}) bool

func getStrMatcher(field reflect.StructField, structType reflect.Type, strValue *string, op StrOp) func(interface{}) bool {
	return func(obj interface{}) bool {
		var vObj = reflect.ValueOf(obj)
		if vObj.Type().Kind() == reflect.Pointer {
			vObj = vObj.Elem()
		}
		if vObj.Type() != structType {
			return false
		}
		var val = vObj.FieldByIndex(field.Index)
		var sv *string
		var done = false
		if val.Kind() == reflect.Pointer {
			if val.IsNil() {
				done = true
			} else {
				val = val.Elem()
			}
		}
		if !done {
			switch val.Kind() {
			case reflect.String:
				sv = new(string)
				*sv = val.String()
				done = true
			case reflect.Struct:
				if val.CanInterface() {
					var i1 = val.Interface()
					var ok bool
					var av attr.Value
					if av, ok = i1.(attr.Value); ok {
						if !av.IsNull() {
							var i3 interface{ ValueString() string }
							if i3, ok = av.(interface{ ValueString() string }); ok {
								sv = new(string)
								*sv = i3.ValueString()
								done = true
							}
						} else {
							done = true
						}
					}
				}
			}
		}
		if !done {
			return false
		}

		if strValue == nil {
			switch op {
			case StrOp_Equal:
				return sv == nil
			case StrOp_NotEqual:
				return sv != nil
			}
		} else if sv != nil {
			switch op {
			case StrOp_Equal:
				return strings.ToLower(*sv) == strings.ToLower(*strValue)
			case StrOp_NotEqual:
				return strings.ToLower(*sv) != strings.ToLower(*strValue)
			case StrOp_StartsWith:
				return strings.HasPrefix(strings.ToLower(*sv), strings.ToLower(*strValue))
			case StrOp_EndsWith:
				return strings.HasSuffix(strings.ToLower(*sv), strings.ToLower(*strValue))
			case StrOp_Matches:
				var err error
				var matched bool
				if matched, err = path.Match(*strValue, *sv); err == nil {
					return matched
				}
			}
		}
		return false
	}
}

func getIntMatcher(field reflect.StructField, structType reflect.Type, intValue *int64, op IntOp) func(interface{}) bool {
	return func(obj interface{}) bool {
		var vObj = reflect.ValueOf(obj)
		if vObj.Type().Kind() == reflect.Pointer {
			vObj = vObj.Elem()
		}
		if vObj.Type().Kind() != reflect.Struct {
			return false
		}
		if structType != vObj.Type() {
			return false
		}
		var val = vObj.FieldByIndex(field.Index)
		var iv *int64
		var done = false
		if val.Type().Kind() == reflect.Pointer {
			if val.IsNil() {
				done = true
			} else {
				val = val.Elem()
			}
		}
		if !done {
			if val.CanInt() {
				iv = new(int64)
				*iv = val.Int()
				done = true
			} else if val.CanInterface() {
				var i1 = val.Interface()
				var ok bool
				var av attr.Value
				if av, ok = i1.(attr.Value); ok {
					if !av.IsNull() {
						var i3 interface{ ValueInt64() int64 }
						if i3, ok = av.(interface{ ValueInt64() int64 }); ok {
							iv = new(int64)
							*iv = i3.ValueInt64()
							done = true
						}
					} else {
						done = true
					}
				}
			}
		}
		if !done {
			return false
		}

		if intValue == nil {
			switch op {
			case IntOp_Equal:
				return iv == nil
			case IntOp_NotEqual:
				return iv != nil
			}
		} else if iv != nil {
			switch op {
			case IntOp_Equal:
				return *iv == *intValue
			case IntOp_NotEqual:
				return *iv != *intValue
			case IntOp_GreaterThan:
				return *iv > *intValue
			case IntOp_LessThan:
				return *iv < *intValue
			case IntOp_GreaterOrEqual:
				return *iv >= *intValue
			case IntOp_LessOrEqual:
				return *iv <= *intValue
			}
		}
		return false
	}
}

func getBoolMatcher(field reflect.StructField, structType reflect.Type, boolValue *bool) func(interface{}) bool {
	return func(obj interface{}) bool {
		var vObj = reflect.ValueOf(obj)
		if vObj.Type().Kind() == reflect.Pointer {
			vObj = vObj.Elem()
		}
		if vObj.Type().Kind() != reflect.Struct {
			return false
		}
		if structType != vObj.Type() {
			return false
		}
		var val = vObj.FieldByIndex(field.Index)
		var bv *bool
		var done = false
		if val.Type().Kind() == reflect.Pointer {
			if val.IsNil() {
				done = true
			} else {
				val = val.Elem()
			}
		}
		if !done {
			if val.Type().Kind() == reflect.Bool {
				bv = new(bool)
				*bv = val.Bool()
				done = true
			} else if val.CanInterface() {
				var i1 = val.Interface()
				var ok bool
				var av attr.Value
				if av, ok = i1.(attr.Value); ok {
					if !av.IsNull() {
						var i3 interface{ ValueBool() bool }
						if i3, ok = av.(interface{ ValueBool() bool }); ok {
							bv = new(bool)
							*bv = i3.ValueBool()
							done = true
						}
					}
				}
			}
		}
		if !done {
			return false
		}
		if bv == nil && boolValue == nil {
			return true
		}
		if bv == nil || boolValue == nil {
			return false
		}
		return *bv == *boolValue
	}
}

var int64Type = reflect.TypeOf((*types.Int64)(nil)).Elem()
var stringType = reflect.TypeOf((*types.String)(nil)).Elem()
var boolType = reflect.TypeOf((*types.Bool)(nil)).Elem()

type matcherType = int32

const (
	matcherType_Unsupported matcherType = iota
	matcherType_Int64
	matcherType_String
	matcherType_Bool
)

func getFieldMatcher(fc *filterCriteria, modelType reflect.Type) (cb matcher, diags diag.Diagnostics) {
	if fc == nil {
		return
	}
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		diags.AddError("Create a field matcher error: \"struct\" type expected.",
			fmt.Sprintf("Type \"%T\" should be a \"struct\"", modelType))
		return
	}
	var fieldName = fc.Field.ValueString()
	var fNum = modelType.NumField()
	var field reflect.StructField
	var ok = false
	for i := 0; i < fNum; i++ {
		var f = modelType.Field(i)
		var v = f.Tag.Get("tfsdk")
		if len(v) > 0 {
			var comps = strings.Split(v, ",")
			if comps[0] == fieldName {
				field = f
				ok = true
				break
			}
		}
	}
	if !ok {
		diags.AddError(fmt.Sprintf("Field \"%s\" not found on struct \"%s\"", fieldName, modelType.Name()),
			"Use only defined fields in filter criteria")
		return
	}

	var fieldType = field.Type
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}
	var mt = matcherType_Unsupported
	switch fieldType.Kind() {
	case reflect.Struct:
		switch field.Type {
		case int64Type:
			mt = matcherType_Int64
		case stringType:
			mt = matcherType_String
		case boolType:
			mt = matcherType_Bool
		}
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Int, reflect.Int64, reflect.Int32:
		mt = matcherType_Int64
	case reflect.String:
		mt = matcherType_String
	case reflect.Bool:
		mt = matcherType_Bool
	}
	if mt == matcherType_Unsupported {
		diags.AddError(fmt.Sprintf("Field \"%s\" has unsupported type on struct \"%s\"", fieldName, modelType.Name()),
			"Use only defined fields in filter criteria")
		return
	}
	var err error
	var op string
	if fc.Cmp.IsNull() {
		op = "=="
	} else {
		op = fc.Cmp.ValueString()
	}
	switch mt {
	case matcherType_Bool:
		var bv *bool
		var fieldValue = fc.Value.ValueString()
		switch fieldValue {
		case "true":
			bv = new(bool)
			*bv = true
		case "false":
			bv = new(bool)
			*bv = false
		case "null":
		default:
			diags.AddError(fmt.Sprintf("\"%s\" invalid value for boolean field: true, false, null", fieldValue),
				fmt.Sprintf("Field \"%s\" is defined as a boolean field.", fieldValue))
		}
		cb = getBoolMatcher(field, modelType, bv)

	case matcherType_String:
		var strOp = StrOp_Invalid
		if strOp, ok = stringOperators[op]; !ok {
			var ops []string
			for k := range stringOperators {
				ops = append(ops, fmt.Sprintf("\"%s\"", k))
			}
			diags.AddError(fmt.Sprintf("Unsupported string comparison operator: \"%s\"", op),
				fmt.Sprintf("Supported string comparison operators: %s", strings.Join(ops, ", ")))
			return
		}
		var sv *string
		var fieldValue = fc.Value.ValueString()
		if fieldValue != "null" {
			sv = new(string)
			*sv = fieldValue
		}
		if sv == nil {
			if strOp != IntOp_Equal && strOp != IntOp_NotEqual {
				diags.AddError(fmt.Sprintf("\"null\" comparison value can be used with \"==\" or \"!=\" operators"),
					"Field \"%s\" is defined as a string field.")
				return
			}
		}
		cb = getStrMatcher(field, modelType, sv, strOp)

	case matcherType_Int64:
		var intOp = IntOp_Invalid
		if intOp, ok = integerOperators[op]; !ok {
			var ops []string
			for k := range integerOperators {
				ops = append(ops, fmt.Sprintf("\"%s\"", k))
			}
			diags.AddError(fmt.Sprintf("Unsupported integer comparison operator: \"%s\"", op),
				fmt.Sprintf("Supported integer comparison operators: %s", strings.Join(ops, ", ")))
			return
		}
		var iv *int64
		var fieldValue = fc.Value.ValueString()
		if fieldValue != "null" {
			var ii int
			if ii, err = strconv.Atoi(fieldValue); err != nil {
				diags.AddError(fmt.Sprintf("\"%s\" field value should be an integer. Got \"%s\"", fieldName, fieldValue),
					"Field \"%s\" is defined as an integer field.")
				return
			} else {
				iv = new(int64)
				*iv = int64(ii)
			}
		}
		if iv == nil {
			if intOp != IntOp_Equal && intOp != IntOp_NotEqual {
				diags.AddError(fmt.Sprintf("\"null\" comparison value can be used with \"==\" or \"!=\" operators"),
					"Field \"%s\" is defined as an integer field.")
				return
			}
		}

		cb = getIntMatcher(field, modelType, iv, intOp)
	}

	return
}
