// © 2016-2022 Resurface Labs Inc.

package logger

import (
	"regexp"
)

type HttpRule struct {
	verb   string
	scope  *regexp.Regexp
	param1 interface{}
	param2 interface{}
}

func NewHttpRule(_verb string, _scope *regexp.Regexp,
	_param1 interface{}, _param2 interface{}) *HttpRule {
	return &HttpRule{
		verb:   _verb,
		scope:  _scope,
		param1: _param1,
		param2: _param2,
	}
}
