package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpRule struct{
	string verb
	pattern scope
	//not sure how to an object
}

type HttpRules struct{
	DEBUG_RULES		string
	STANDARD_RULES	string
	STRICT_RULES	string
	DEFAULT_RULES	string
}

func newHttpRules() HttpRules{
	return 
}

func (HttpRules) getDefaultRules() string{
	rreturn HttpRules.DEFAULT_RULES;
}

func (HttpRules) setDefaultRules(string r){
	HttpRules.DEFAULT_RULES := strings.Replace(r,"(?m)^\\s*include default\\s*$","",-1)
}

func (HttpRules) getDebugRules() string{
	return HttpRules.DEBUG_RULES
}

func (HttpRules) getStandardRules() string{
	return HttpRules.STANDARD_RULES
}

func (HttpRules) getStrictRules() string{
	return HttpRules.STRICT_RULES
}

func parseRule(string r) HttpRule{
	
}