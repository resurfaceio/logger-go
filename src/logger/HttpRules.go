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

func newHttpRules(string rules) HttpRules{
	if(rules == null){
		rules = HttpRules.getDefaultRules()
	}

	//load rules from external files
	if(rules[0:6] == "file://"){
		string rfile = rules[0:6]
		//try to read rules from file
	}

	//force default rules if necessary
	strings.Replace(rules, "(?m)^\\s*include default\\s*$", /*Not sure exactly what its trying to replace it with*/)
	if(len(strings.Trimspace(rules)) == 0){
		rules = HttpRules.getDefaultRules
	}

	//expand rule includes
	strings.Replace(rules, "(?m)^\\s*include debug\\s*$", )
	strings.Replace(rules, "(?m)^\\s*include standard\\s*$", )
	strings.Replace(rules, "(?m)^\\s*include strict\\s*$", )

	// parse all rules
	// Not sure about this part

	// break out rules by verb
	// Not sure either but may need to add to the struct created at the top
}

func (HttpRules) getDefaultRules() string{
	return HttpRules.DEFAULT_RULES
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

	//https://www.geeksforgeeks.org/matching-using-regexp-in-golang/

}