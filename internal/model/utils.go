package model

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
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
	level   zapcore.Level
	context context.Context
	encoder zapcore.Encoder
}

// NewTerraformCore returns a zap Core.
func NewTerraformCore(ctx context.Context, level zapcore.Level) zapcore.Core {
	var e = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	return nopCore{
		context: ctx,
		level:   level,
		encoder: e,
	}
}
func (n nopCore) Enabled(l zapcore.Level) bool  { return n.level <= l }
func (n nopCore) With([]zap.Field) zapcore.Core { return n }
func (n nopCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(e, n)
}
func (n nopCore) Write(entry zapcore.Entry, fields []zap.Field) (err error) {
	var b *buffer.Buffer
	if b, err = n.encoder.EncodeEntry(entry, fields); err == nil {
		switch entry.Level {
		case zapcore.DebugLevel:
			tflog.Debug(n.context, b.String())
		case zapcore.InfoLevel:
			tflog.Info(n.context, b.String())
		case zapcore.WarnLevel:
			tflog.Warn(n.context, b.String())
		case zapcore.ErrorLevel:
			tflog.Error(n.context, b.String())
		}
	}
	return
}
func (nopCore) Sync() error { return nil }
