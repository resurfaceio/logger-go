package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/textproto"
	"regexp"
	"strings"
	"io/ioutil"
)

var httpRules *HttpRules

type HttpRules struct{
	debugRules		string
	standardRules	string
	strictRules		string
	defaultRules	string
	allowHttpUrl	[]HttpRule
	copySessionField	[]HttpRule
	remove			[]HttpRule
	removeIf 		[]HttpRule
	removeIfFound		[]HttpRule
	removeUnless		[]HttpRule
	removeUnlessFound 	[]HttpRule
	replace 			[]HttpeRule
	sample			[]HttpRule
	skipCompression		bool
	skipSubmission		bool
	size 			int
	stop			[]HttpRule
	stopIf 			[]HttpRule
	stopIfFound 	[]HttpRule
	stopUnless 		[]HttpRule
	stopUnlessFound 	[]HttpRule
	text 			string
}

func GetHttpRules() *HttpRules {

}

func NewHttpRules(string rules) *HttpRules {
	httpRules := GetHttpRules()
	if(rules == null){
		rules = httpRules.getDefaultRules()
	}

	//load rules from external files
	if(strings.HasPrefix(rules, "file://")){
		// obtain file name
		rfile := strings.TrimPrefix("file://")
		rfile := strings.TrimSpace(rfile)
		
		// read rules from file
		buffer, err := ioutil.ReadFile(rfile)
		if err != nil {
			// error handling
		}
		rules = string(buffer)
	}

	//force default rules if necessary
	regex := regexp.MustCompile("(?m)^\\s*include default\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.defaultRules)
	if(len(strings.TrimSpace(rules)) == 0){
		rules = HttpRules.getDefaultRules
	}

	//expand rule includes

	//include debug rules
	regex = regexp.MustCompile("(?m)^\\s*include debug\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.debugRules)
	// include standard rules
	regex = regexp.MustCompile("(?m)^\\s*include standard\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.standardRules)
	// include strict rules
	regex = regexp.MustCompile( "(?m)^\\s*include strict\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.strictRules)

	_text := rules

	// parse all rules
	// Not sure about this part
	var prs []HttpRule
	for _, rule := range regexp.MustCompile("\\r?\\n").Split(_text, -1) {
		parsed := parseRule(rule)
		if parsed != nil {
			prs = append(prs, parsed)
		}
	}

	_size := len(prs)

	// break out rules by verb
	_allowHttpUrl := len(ruleFilter(prs, "allow_http_url", ruleCompare)) > 0
	_copySessionField

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

/*
 ALL REGEX STRINGS NEED TO BE COVERTED TO STRINGS THAT GO CAN COMPILE

 THESE vvvv DON'T WORK vvvv
*/


const RegexAllowHttpUrl regexp.Regexp = regexp.MustCompile("/^\s*allow_http_url\s*(#.*)?$/")
const RegexBlankOrComment regexp.Regexp = regexp.MustCompile("/^\s*([#].*)*$/")
const RegexCopySessionField regexp.Regexp = regexp.MustCompile("/^\s*copy_session_field\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexRemove regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*remove\s*(#.*)?$/")
const RegexRemoveIf regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*remove_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexRemoveIfFound regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*remove_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexRemoveUnless regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*remove_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexRemoveUnlessFound regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*remove_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexReplace regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*replace[\s]+([~!%|\/].+[~!%|\/]),[\s]+([~!%|\/].*[~!%|\/])\s*(#.*)?$/")
const RegexSample regexp.Regexp = regexp.MustCompile("/^\s*sample\s+(\d+)\s*(#.*)?$/")
const RegexSkipCompression regexp.Regexp = regexp.MustCompile("/^\s*skip_compression\s*(#.*)?$/")
const RegexSkipSubmission regexp.Regexp = regexp.MustCompile("/^\s*skip_submission\s*(#.*)?$/")
const RegexStop regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*stop\s*(#.*)?$/")
const RegexStopIf regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*stop_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexStopIfFound regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*stop_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexStopUnless regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*stop_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")
const RegexStopUnlessFound regexp.Regexp = regexp.MustCompile("/^\s*([~!%|\/].+[~!%|\/])\s*stop_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$/")

func ruleCompare(ruleString string, s string) bool {
	return rule == s
}

func ruleFilter(parsedRules []HttpRule, ruleString string, cond func(string,HttpRule) bool) []HttpRule {
	result := []HttpRule{}
	for _, rule := range parsedRules {
		if cond(ruleString, rule) {
			result = append(result, rule)
		}
	}
	return result
}