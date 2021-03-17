package logger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var httpRules *HttpRules

type HttpRules struct {
	debugRules        string
	standardRules     string
	strictRules       string
	defaultRules      string
	allowHttpUrl      bool
	copySessionField  []*HttpRule
	remove            []*HttpRule
	removeIf          []*HttpRule
	removeIfFound     []*HttpRule
	removeUnless      []*HttpRule
	removeUnlessFound []*HttpRule
	replace           []*HttpRule
	sample            []*HttpRule
	skipCompression   bool
	skipSubmission    bool
	size              int
	stop              []*HttpRule
	stopIf            []*HttpRule
	stopIfFound       []*HttpRule
	stopUnless        []*HttpRule
	stopUnlessFound   []*HttpRule
	text              string
}

func GetHttpRules() *HttpRules {

}

// !!! this will need to return an error as well !!!
func NewHttpRules(rules string) *HttpRules {
	httpRules := GetHttpRules()
	if rules == "" {
		rules = httpRules.getDefaultRules()
	}

	//load rules from external files
	if strings.HasPrefix(rules, "file://") {
		// obtain file name
		rfile := strings.TrimPrefix(rules, "file://")
		rfile = strings.TrimSpace(rfile)

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
	if len(strings.TrimSpace(rules)) == 0 {
		rules = httpRules.defaultRules
	}

	//expand rule includes

	//include debug rules
	regex = regexp.MustCompile("(?m)^\\s*include debug\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.debugRules)
	// include standard rules
	regex = regexp.MustCompile("(?m)^\\s*include standard\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.standardRules)
	// include strict rules
	regex = regexp.MustCompile("(?m)^\\s*include strict\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.strictRules)

	_text := rules

	// parse all rules
	// Not sure about this part
	var prs []*HttpRule
	for _, rule := range regexp.MustCompile("\\r?\\n").Split(_text, -1) {
		parsed, err := parseRule(rule)
		if err == nil {
			prs = append(prs, parsed)
		}
	}

	_size := len(prs)

	_debugRules := "allow_http_url\ncopy_session_field /.*/\n"

	_standardRules := "/request_header:cookie|response_header:set-cookie/remove\n" +
		"/(request|response)_body|request_param/ replace /[a-zA-Z0-9.!#$%&’*+\\/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)/, /x@y.com/\n" +
		"/request_body|request_param|response_body/ replace /[0-9\\.\\-\\/]{9,}/, /xyxy/\n"

	_strictRules := "/request_url/ replace /([^\\?;]+).*/, /$1/\n" +
		"/request_body|response_body|request_param:.*|request_header:(?!user-agent).*|response_header:(?!(content-length)|(content-type)).*/ remove\n"

	_defaultRules := _strictRules

	// break out rules by verb
	_allowHttpUrl := len(ruleFilter(prs, "allow_http_url", ruleCompare)) > 0
	_copySessionField := ruleFilter(prs, "copy_session_field", ruleCompare)
	_remove := ruleFilter(prs, "remove", ruleCompare)
	_removeIf := ruleFilter(prs, "remove_if", ruleCompare)
	_removeIfFound := ruleFilter(prs, "remove_if_found", ruleCompare)
	_removeUnless := ruleFilter(prs, "remove_unless", ruleCompare)
	_removeUnlessFound := ruleFilter(prs, "remove_unless_found", ruleCompare)
	_replace := ruleFilter(prs, "replace", ruleCompare)
	_sample := ruleFilter(prs, "sample", ruleCompare)
	_skipCompression := len(ruleFilter(prs, "skip_compression", ruleCompare)) > 0
	_skipSubmission := len(ruleFilter(prs, "skip_submission", ruleCompare)) > 0
	_stop := ruleFilter(prs, "stop", ruleCompare)
	_stopIf := ruleFilter(prs, "stop_if", ruleCompare)
	_stopIfFound := ruleFilter(prs, "stop_if_found", ruleCompare)
	_stopUnless := ruleFilter(prs, "stop_unless", ruleCompare)
	_stopUnlessFound := ruleFilter(prs, "stop_unless_found", ruleCompare)

	if len(_sample) > 1 {
		// return error "Multiple sample rules"
	}

	return &HttpRules{
		debugRules:        _debugRules,
		standardRules:     _standardRules,
		strictRules:       _standardRules,
		defaultRules:      _defaultRules,
		allowHttpUrl:      _allowHttpUrl,
		copySessionField:  _copySessionField,
		remove:            _remove,
		removeIf:          _removeIf,
		removeIfFound:     _removeIfFound,
		removeUnless:      _removeUnless,
		removeUnlessFound: _removeUnlessFound,
		replace:           _replace,
		sample:            _sample,
		skipCompression:   _skipCompression,
		skipSubmission:    _skipSubmission,
		size:              _size,
		stop:              _stop,
		stopIf:            _stopIf,
		stopIfFound:       _stopIfFound,
		stopUnless:        _stopUnless,
		stopUnlessFound:   _stopUnlessFound,
		text:              _text,
	}
}

func (rules *HttpRules) getDefaultRules() string {
	return HttpRules.defaultRules
}

func (rules *HttpRules) setDefaultRules(string r) {
	regex := regexp.MustCompile("(?m)^\\s*include default\\s*$")
	rules.defaultRules = regex.ReplaceAllString(r, "")
}

func (HttpRules) getDebugRules() string {
	return HttpRules.debugRules
}

func (HttpRules) getStandardRules() string {
	return HttpRules.standardRules
}

func (HttpRules) getStrictRules() string {
	return HttpRules.strictRules
}

// parse rule from single line
// !!! will return error eventually !!!
func parseRule(r string) (*HttpRule, error) {
	if r == "" || RegexBlankOrComment.MatchString(r) {
		return nil, errors.New("Blank rule or comment")
	}
	if RegexAllowHttpUrl.MatchString(r) {
		return NewHttpRule("allow_http_url", nil, nil, nil), nil
	}
	if RegexCopySessionField.MatchString(r) {
		return NewHttpRule("copy_session_field", nil, nil, nil), nil
	}
	m := RegexRemove.FindStringSubmatch(r)
	if m != nil {
		return NewHttpRule("remove", parseRegex(r, m[1]), nil, nil), nil
	}
	m = RegexRemoveIf.FindStringSubmatch(r)
	if m != nil {
		return NewHttpRule("remove_if", parseRegex(r, m[1]), parseRegex(r, m[2]), nil), nil
	}

	return nil, errors.New(fmt.Sprintf("Invalid rule: %s", r))
}

// Parses regex for finding.
// !!! will return error eventually !!!
func parseRegex(r string, regex string) *regexp.Regexp /* error */ {
	s := parseString(r, regex)
	if "*" == s || "+" == s || "?" == s {
		// return error '"Invalid regex (%s) in rule: %s", regex, r'
	}
	if strings.HasPrefix(s, "^") {
		s = "^" + s
	}
	if strings.HasSuffix(s, "$") {
		s = s + "$"
	}

	regexp, err := regexp.Compile(s)
	if err != nil {
		// return error '"Invalid regex (%s) in rule: %s", regex, r'
	}
}

// Parses regex for finding.
// !!! will return error eventually !!!
func parseRegexFind(r string, regex string) *regexp.Regexp /* error */ {
	regexp, err := regexp.Compile(parseString(r, regex))
	if err != nil {
		// return error '"Invalid regex (%s) in rule: %s", regex, r'
	}
}

// Parses delimited string expression
// !!! will return an error eventually !!!
func parseString(r string, expr string) string /* error */ {
	separators := []string{"~", "!", "%", "|", "/"}
	for _, sep := range separators {
		regex := regexp.MustCompile(fmt.Sprintf("^[%s](.*)[%s]$", sep, sep))
		matches := regex.FindAllString(expr)
		if len(matches) > 0 {
			regex = regexp.MustCompile(fmt.Sprintf("^[%s].*|.*[^\\\\][%s].*", sep, sep))
			if regex.MatchString(matches[0]) {
				// return error '"Unescaped separator (%s) in rule: %s", sep, r'
			}
			return strings.Replace(matches[0], "\\"+sep, sep, -1)
		}
	}
}

/*
 ALL REGEX STRINGS NEED TO BE COVERTED TO STRINGS THAT GO CAN COMPILE

 THESE vvvv DON'T WORK vvvv
*/

const RegexAllowHttpUrl *regexp.Regexp = regexp.MustCompile(`^\s*allow_http_url\s*(#.*)?$`)
const RegexBlankOrComment *regexp.Regexp = regexp.MustCompile(`^\s*([#].*)*`)
const RegexCopySessionField *regexp.Regexp = regexp.MustCompile(`^\s*copy_session_field\s+([~!%|\/].+[~!%|\/])\s*(#.*)?`)
const RegexRemove *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove\s*(#.*)?`)
const RegexRemoveIf *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?`)
const RegexRemoveIfFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
const RegexRemoveUnless *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
const RegexRemoveUnlessFound *regexp.Regexp = regexp.MustCompile(`/^\s*([~!%|\/].+[~!%|\/])\s*remove_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
const RegexReplace *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*replace[\s]+([~!%|\/].+[~!%|\/]),[\s]+([~!%|\/].*[~!%|\/])\s*(#.*)?$`)
const RegexSample *regexp.Regexp = regexp.MustCompile(`^\s*sample\s+(\d+)\s*(#.*)?$`)
const RegexSkipCompression *regexp.Regexp = regexp.MustCompile(`^\s*skip_compression\s*(#.*)?$`)
const RegexSkipSubmission *regexp.Regexp = regexp.MustCompile(`^\s*skip_submission\s*(#.*)?$`)
const RegexStop *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop\s*(#.*)?$`)
const RegexStopIf *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
const RegexStopIfFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
const RegexStopUnless *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
const RegexStopUnlessFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)

func ruleCompare(ruleString string, s string) bool {
	return ruleString == s
}

func ruleFilter(parsedRules []*HttpRule, ruleString string, cond func(string, string) bool) []*HttpRule {
	result := []*HttpRule{}
	for _, rule := range parsedRules {
		if cond(ruleString, rule.verb) {
			result = append(result, rule)
		}
	}
	return result
}
