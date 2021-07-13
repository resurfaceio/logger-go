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

func (rule *HttpRule) Verb() string {
	return rule.verb
}

func (rule *HttpRule) Scope() *regexp.Regexp {
	return rule.scope
}

func (rule *HttpRule) Param1() interface{} {
	return rule.param1
}

func (rule *HttpRule) Param2() interface{} {
	return rule.param2
}
