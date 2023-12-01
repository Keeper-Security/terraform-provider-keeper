package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
	"reflect"
	"strconv"
	"testing"
)

type strModel struct {
	F1 string       `tfsdk:"f1"`
	F2 *string      `tfsdk:"f2"`
	F3 types.String `tfsdk:"f3"`
}

type intModel struct {
	F1 int         `tfsdk:"f1"`
	F2 *int64      `tfsdk:"f2"`
	F3 types.Int64 `tfsdk:"f3"`
}

type boolModel struct {
	F1 bool       `tfsdk:"f1"`
	F2 *bool      `tfsdk:"f2"`
	F3 types.Bool `tfsdk:"f3"`
}

func TestIntNullMatcher(t *testing.T) {
	var wm = &intModel{}

	var nt = reflect.TypeOf(wm)
	var fields = []string{"f3"}
	var fc = &FilterCriteria{
		Field: types.StringValue(""),
		Value: types.StringValue("null"),
	}
	for _, f := range fields {
		fc.Field = types.StringValue(f)
		cb, diags := GetFieldMatcher(fc, nt)
		assert.Assert(t, len(diags) == 0)
		assert.Assert(t, cb != nil)
		res := cb(wm)
		assert.Assert(t, res)
	}

}

func TestBoolMatcher(t *testing.T) {
	var r = true
	var wm = &boolModel{
		F1: r,
		F2: &r,
		F3: types.BoolValue(r),
	}

	var nt = reflect.TypeOf(wm)
	var fields = []string{"f1", "f2", "f3"}
	var fc = &FilterCriteria{
		Field: types.StringValue(""),
		Value: types.StringValue("true"),
	}
	for _, f := range fields {
		fc.Field = types.StringValue(f)
		cb, diags := GetFieldMatcher(fc, nt)
		assert.Assert(t, len(diags) == 0)
		assert.Assert(t, cb != nil)
		res := cb(wm)
		assert.Assert(t, res)
	}
	r = false
	wm.F1 = r
	wm.F2 = &r
	wm.F3 = types.BoolValue(r)
	for _, f := range fields {
		fc.Field = types.StringValue(f)
		cb, diags := GetFieldMatcher(fc, nt)
		assert.Assert(t, len(diags) == 0)
		assert.Assert(t, cb != nil)
		res := cb(wm)
		assert.Assert(t, !res)
	}
}

func TestStringMatcher(t *testing.T) {
	var r = "valid"
	var wm = &strModel{
		F1: r,
		F2: &r,
		F3: types.StringValue(r),
	}

	var nt = reflect.TypeOf(wm)
	var fields = []string{"f1", "f2", "f3"}
	var fc = &FilterCriteria{
		Field: types.StringValue(""),
		Cmp:   types.StringValue("=="),
		Value: types.StringValue(r),
	}
	for _, f := range fields {
		fc.Field = types.StringValue(f)
		cb, diags := GetFieldMatcher(fc, nt)
		assert.Assert(t, len(diags) == 0)
		assert.Assert(t, cb != nil)
		res := cb(wm)
		assert.Assert(t, res)
	}
	r = "invalid"
	wm.F1 = r
	wm.F2 = &r
	wm.F3 = types.StringValue(r)
	for _, f := range fields {
		fc.Field = types.StringValue(f)
		cb, diags := GetFieldMatcher(fc, nt)
		assert.Assert(t, len(diags) == 0)
		assert.Assert(t, cb != nil)
		res := cb(wm)
		assert.Assert(t, !res)
	}
}

func TestIntegerMatcher(t *testing.T) {
	var r int64 = 32432423423423
	var wm = &intModel{
		F1: int(r),
		F2: &r,
		F3: types.Int64Value(r),
	}

	var nt = reflect.TypeOf(wm)
	var fields = []string{"f1", "f2", "f3"}
	var fc = &FilterCriteria{
		Field: types.StringValue(""),
		Cmp:   types.StringValue("=="),
		Value: types.StringValue(strconv.Itoa(int(r))),
	}
	for _, f := range fields {
		fc.Field = types.StringValue(f)
		cb, diags := GetFieldMatcher(fc, nt)
		assert.Assert(t, len(diags) == 0)
		assert.Assert(t, cb != nil)
		res := cb(wm)
		assert.Assert(t, res)
	}
	r = 222222
	wm.F1 = int(r)
	wm.F2 = &r
	wm.F3 = types.Int64Value(r)
	for _, f := range fields {
		fc.Field = types.StringValue(f)
		cb, diags := GetFieldMatcher(fc, nt)
		assert.Assert(t, len(diags) == 0)
		assert.Assert(t, cb != nil)
		res := cb(wm)
		assert.Assert(t, !res)
	}
}
