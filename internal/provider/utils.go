package provider

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func mergeMaps[K comparable, V any](maps ...map[K]V) (m map[K]V) {
	m = make(map[K]V)
	for _, mm := range maps {
		for k, v := range mm {
			m[k] = v
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
