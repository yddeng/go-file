package ternary

import (
	"go.luxshare-ict.com/pkg/time"
	xtime "time"
)

func String(v bool, t, f string) string {
	if v {
		return t
	}
	return f
}
func Int(v bool, t, f int) int {
	if v {
		return t
	}
	return f
}
func Int32(v bool, t, f int32) int32 {
	if v {
		return t
	}
	return f
}

func Int64(v bool, t, f int64) int64 {
	if v {
		return t
	}
	return f
}

func Time(v bool, t, f time.Time) time.Time {
	if v {
		return t
	}
	return f
}
func Timex(v bool, t, f xtime.Time) xtime.Time {
	if v {
		return t
	}
	return f
}

func Float(v bool, t, f float64) float64 {
	if v {
		return t
	}
	return f
}
func Float32(v bool, t, f float32) float32 {
	if v {
		return t
	}
	return f
}
