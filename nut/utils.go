package nut

import (
	"reflect"
	"runtime"
)

// GetFunctionName get function name
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
