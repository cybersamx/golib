package reflectutils

import "reflect"

// Indirect returns the value, after de-referencing as many times as necessary to reach
// the base type (or nil).
//
// Credit: https://github.com/golang/go/blob/master/src/html/template/content.go.
// Copyright 2011 The Go Authors. All rights reserved. Licensed under BSD.
func Indirect(a any) any {
	if a == nil {
		return nil
	}

	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}

	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	return v.Interface()
}
