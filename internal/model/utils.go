package model

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type TestInstanceStateFunc func(state *terraform.InstanceState) error

func TestCheckImportStateAttr(key string, value string) TestInstanceStateFunc {
	return func(is *terraform.InstanceState) error {
		var ok bool
		var v string
		v, ok = is.Attributes[key]

		if !ok {
			if value == "0" && (strings.HasSuffix(key, ".#") || strings.HasSuffix(key, ".%")) {
				return nil
			}
			return fmt.Errorf("attribute '%s' not found", key)
		}

		if v != value {
			return fmt.Errorf("attribute '%s' expected %#v, got %#v", key, value, v)
		}

		return nil
	}
}

func ComposeAggregateImportStateCheckFunc(fn ...TestInstanceStateFunc) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) (err error) {
		for _, f := range fn {
			if err = f(states[0]); err != nil {
				break
			}
		}
		return
	}
}

func MergeMaps[K comparable, V any](maps ...map[K]V) (m map[K]V) {
	m = make(map[K]V)
	var ok bool
	for _, mm := range maps {
		for k, v := range mm {
			if _, ok = m[k]; !ok {
				m[k] = v
			}
		}
	}
	return
}

type nopCore struct {
	level zapcore.Level
}

// NewNopCore returns a no-op Core.
func NewNopCore() zapcore.Core {
	return nopCore{
		level: zapcore.InfoLevel,
	}
}
func (n nopCore) Enabled(l zapcore.Level) bool  { return n.level <= l }
func (n nopCore) With([]zap.Field) zapcore.Core { return n }
func (n nopCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(e, n)
}
func (nopCore) Write(zapcore.Entry, []zap.Field) error { return nil }
func (nopCore) Sync() error                            { return nil }
