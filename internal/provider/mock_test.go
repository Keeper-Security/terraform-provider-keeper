package provider

import "time"

type mockContext struct{}

func (_ mockContext) Deadline() (deadline time.Time, ok bool) { return }
func (_ mockContext) Done() (ch <-chan struct{})              { return }
func (_ mockContext) Err() error                              { return nil }
func (_ mockContext) Value(key any) any                       { return nil }
