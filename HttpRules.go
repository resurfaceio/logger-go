// © 2016-2023 Graylog, Inc.

package logger

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// sync.Once for UsageLoggers
var onceHttpRules sync.Once

/*
package global containing default:

	debugRules
	standardRules
	strictRules
	defaultRules
*/
var httpRules *HttpRules

// struct for rules that are applied to http logging messages
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

// get package global httpRules containing default rules sets
func GetHttpRules() *HttpRules {
	onceHttpRules.Do(func() {
		_debugRules := "allow_http_url\ncopy_session_field /.*/\n"

		_standardRules := "/request_header:cookie|response_header:set-cookie/remove\n" +
			"/(request|response)_body|request_param/ replace /[a-zA-Z0-9.!#$%&’*+\\/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)/, /x@y.com/\n" +
			"/request_body|request_param|response_body/ replace /[0-9\\.\\-\\/]{9,}/, /xyxy/\n"

		_strictRules := "/request_url/ replace /([^\\?;]+).*/, /$1/\n" +
			"/request_body|response_body|request_param:.*|request_header:(user-agent).*|response_header:((content-length)|(content-type)).*/ remove\n"

		_defaultRules := _strictRules
		httpRules = &HttpRules{
			debugRules:    _debugRules,
			standardRules: _standardRules,
			strictRules:   _strictRules,
			defaultRules:  _defaultRules,
		}
	})
	return httpRules
}

// generate new set of rules based on given rule input string
func newHttpRules(rules string) (*HttpRules, error) {
	httpRules := GetHttpRules()
	if rules == "" {
		rules = httpRules.DefaultRules()
	}

	//load rules from external files
	if strings.HasPrefix(rules, "file://") {
		// obtain file name
		rfile := strings.TrimPrefix(rules, "file://")
		rfile = strings.TrimSpace(rfile)

		// read rules from file
		buffer, err := ioutil.ReadFile(rfile)
		if err != nil {
			return nil, fmt.Errorf("failed to load rules: %s", rfile)
		}
		rules = string(buffer)
	}

	//force default rules if necessary
	regex := regexp.MustCompile(`include default`)
	rules = regex.ReplaceAllString(rules, httpRules.defaultRules)
	if len(strings.TrimSpace(rules)) == 0 {
		rules = httpRules.defaultRules
	}

	//expand rule includes
	//include debug rules
	regex = regexp.MustCompile(`include debug`)
	rules = regex.ReplaceAllString(rules, httpRules.debugRules)
	// include standard rules
	regex = regexp.MustCompile(`include standard`)
	rules = regex.ReplaceAllString(rules, httpRules.standardRules)
	// include strict rules
	regex = regexp.MustCompile(`include strict`)
	rules = regex.ReplaceAllString(rules, httpRules.strictRules)

	_text := rules
	// fmt.Println("_text: " + _text)
	// parse all rules
	var prs []*HttpRule
	for _, rule := range regexp.MustCompile(`\r?\n`).Split(_text, -1) {
		// fmt.Println("for each rule: " + rules)
		parsed, err := parseRule(rule)
		if parsed != nil {
			prs = append(prs, parsed)
		} else if parsed == nil && err != nil {
			return nil, err
		}
	}

	_size := len(prs)

	_debugRules := "allow_http_url\ncopy_session_field /.*/\n"

	_standardRules := "/request_header:cookie|response_header:set-cookie/remove\n" +
		"/(request|response)_body|request_param/ replace /[a-zA-Z0-9.!#$%&’*+\\/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)/, /x@y.com/\n" +
		"/request_body|request_param|response_body/ replace /[0-9\\.\\-\\/]{9,}/, /xyxy/\n"

	_strictRules := "/request_url/ replace /([^\\?;]+).*/, /$1/\n" +
		"/request_body|response_body|request_param:.*|request_header:(user-agent).*|response_header:((content-length)|(content-type)).*/ remove\n"

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
		return nil, fmt.Errorf("multiple sample rules")
	}

	return &HttpRules{
		debugRules:        _debugRules,
		standardRules:     _standardRules,
		strictRules:       _strictRules,
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
	}, nil // error is nil
}

func (rules *HttpRules) DefaultRules() string {
	return rules.defaultRules
}

// *HttpRules.SetDefaultRules(r string) sets the default rules of the logger to rule(s) r
func (rules *HttpRules) SetDefaultRules(r string) {
	regex := regexp.MustCompile(`(?m)^\s*include default\s*$`)
	rules.defaultRules = regex.ReplaceAllString(r, "")
}

func (rules *HttpRules) DebugRules() string {
	return rules.debugRules
}

func (rules *HttpRules) StandardRules() string {
	return rules.standardRules
}

func (rules *HttpRules) StrictRules() string {
	return rules.strictRules
}

func (rules *HttpRules) AllowHttpUrl() bool {
	return rules.allowHttpUrl
}

func (rules *HttpRules) CopySessionField() []*HttpRule {
	return rules.copySessionField
}

func (rules *HttpRules) Remove() []*HttpRule {
	return rules.remove
}

func (rules *HttpRules) RemoveIf() []*HttpRule {
	return rules.removeIf
}

func (rules *HttpRules) RemoveIfFound() []*HttpRule {
	return rules.removeIfFound
}

func (rules *HttpRules) RemoveUnless() []*HttpRule {
	return rules.removeUnless
}

func (rules *HttpRules) RemoveUnlessFound() []*HttpRule {
	return rules.removeUnlessFound
}

func (rules *HttpRules) Replace() []*HttpRule {
	return rules.replace
}

func (rules *HttpRules) Sample() []*HttpRule {
	return rules.sample
}

func (rules *HttpRules) SkipCompression() bool {
	return rules.skipCompression
}

func (rules *HttpRules) SkipSubmissio() bool {
	return rules.skipSubmission
}

func (rules *HttpRules) Size() int {
	return rules.size
}

func (rules *HttpRules) Stop() []*HttpRule {
	return rules.stop
}

func (rules *HttpRules) StopIf() []*HttpRule {
	return rules.stopIf
}

func (rules *HttpRules) StopIfFound() []*HttpRule {
	return rules.stopIfFound
}

func (rules *HttpRules) StopUnless() []*HttpRule {
	return rules.stopUnless
}

func (rules *HttpRules) StopUnlessFound() []*HttpRule {
	return rules.stopUnlessFound
}

func (rules *HttpRules) Text() string {
	return rules.text
}

// parse rule from single line
func parseRule(r string) (*HttpRule, error) {
	if r == "" || regexBlankOrComment.MatchString(r) {
		return nil, nil
	}
	if regexAllowHttpUrl.MatchString(r) {
		return NewHttpRule("allow_http_url", nil, nil, nil), nil
	}
	m := regexCopySessionField.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("copy_session_field", nil, parsedRegex, nil), nil
	}
	m = regexRemove.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("remove", parsedRegex, nil, nil), nil
	}
	m = regexRemoveIf.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex1, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegex2, err := parseRegex(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("remove_if", parsedRegex1, parsedRegex2, nil), nil
	}
	m = regexRemoveIfFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegexFind, err := parseRegexFind(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("remove_if_found", parsedRegex, parsedRegexFind, nil), nil
	}
	m = regexRemoveUnless.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex1, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegex2, err := parseRegex(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("remove_unless", parsedRegex1, parsedRegex2, nil), nil
	}
	m = regexRemoveUnlessFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegexFind, err := parseRegexFind(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("remove_unless_found", parsedRegex, parsedRegexFind, nil), nil
	}
	m = regexReplace.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegexFind, err := parseRegexFind(r, m[0][2])
		if err != nil {
			return nil, err
		}
		parsedString, err := parseString(r, m[0][3])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("replace", parsedRegex, parsedRegexFind, parsedString), nil
	}
	m = regexSample.FindAllStringSubmatch(r, -1)
	if m != nil {
		m1, err := strconv.Atoi(m[0][1])
		if err != nil {
			return nil, fmt.Errorf("error parsing sample rule: %s", r)
		}
		if m1 < 1 || m1 > 99 {
			return nil, fmt.Errorf("invalid sample percent: %d", m1)
		}
		return NewHttpRule("sample", nil, m1, nil), nil
	}
	m = regexSkipCompression.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("skip_compression", nil, nil, nil), nil
	}
	m = regexSkipSubmission.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("skip_submission", nil, nil, nil), nil
	}
	m = regexStop.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("stop", parsedRegex, nil, nil), nil
	}
	m = regexStopIf.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex1, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegex2, err := parseRegex(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("stop_if", parsedRegex1, parsedRegex2, nil), nil
	}
	m = regexStopIfFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegexFind, err := parseRegexFind(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("stop_if_found", parsedRegex, parsedRegexFind, nil), nil
	}
	m = regexStopUnless.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex1, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegex2, err := parseRegex(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("stop_unless", parsedRegex1, parsedRegex2, nil), nil
	}
	m = regexStopUnlessFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		parsedRegex, err := parseRegex(r, m[0][1])
		if err != nil {
			return nil, err
		}
		parsedRegexFind, err := parseRegexFind(r, m[0][2])
		if err != nil {
			return nil, err
		}
		return NewHttpRule("stop_unless_found", parsedRegex, parsedRegexFind, nil), nil
	}
	return nil, fmt.Errorf("invalid rule: %s", r)
}

// Parses regex for matching.
func parseRegex(r string, regex string) (*regexp.Regexp, error) {
	s, err := parseString(r, regex)
	if err != nil {
		return nil, err
	}
	if s == "*" || s == "+" || s == "?" {
		return nil, fmt.Errorf("invalid regex (%s) in rule: %s", regex, r)
	}
	if !strings.HasPrefix(s, "^") {
		s = "^" + s
	}
	if !strings.HasSuffix(s, "$") {
		s = s + "$"
	}

	regexp, err := regexp.Compile(s)
	if err != nil {
		return nil, fmt.Errorf("invalid regex (%s) in rule: %s", regex, r)
	}
	return regexp, nil
}

// Parses regex for finding.
func parseRegexFind(r string, regex string) (*regexp.Regexp, error) {
	parsedString, err := parseString(r, regex)
	if err != nil {
		return nil, err
	}
	regexp, err := regexp.Compile(parsedString)
	if err != nil {
		return nil, fmt.Errorf("invalid regex (%s) in rule: %s", regex, r)
	}
	return regexp, nil
}

// Parses delimited string expression
func parseString(r string, expr string) (string, error) {
	separators := []string{"~", "!", "%", "|", "/"}
	for _, sep := range separators {
		regex := regexp.MustCompile(fmt.Sprintf("^[%s](.*)[%s]$", sep, sep))
		m := regex.FindAllStringSubmatch(expr, -1)
		if m != nil {
			m1 := m[0][1]
			regex = regexp.MustCompile(fmt.Sprintf("^[%s].*|.*[^\\\\][%s].*", sep, sep))
			if regex.MatchString(m1) {
				// return error '"Unescaped separator (%s) in rule: %s", sep, r'
				return "", fmt.Errorf("unescaped separator (%s) in rule: %s", sep, r)
			}
			return strings.Replace(m1, "\\"+sep, sep, -1), nil
		}
	}
	return "", fmt.Errorf("invalid expression (%s) in rule: %s", expr, r)
}

// Apply current rules to message details.
func (rules *HttpRules) apply(details [][]string) [][]string {
	// stop rules come first
	for _, r := range rules.stop {
		for _, d := range details {
			if r.scope.FindAllStringSubmatch(d[0], -1) != nil {
				return nil
			}
		}
	}
	for _, r := range rules.stopIfFound {
		for _, d := range details {
			regex := r.param1.(*regexp.Regexp)
			if r.scope.FindAllStringSubmatch(d[0], -1) != nil && regex.FindAllStringSubmatch(d[1], -1) != nil {
				return nil
			}
		}
	}
	for _, r := range rules.stopIf {
		for _, d := range details {
			regex := r.param1.(*regexp.Regexp)
			if r.scope.FindAllStringSubmatch(d[0], -1) != nil && regex.FindAllStringSubmatch(d[1], -1) != nil {
				return nil
			}
		}
	}
	passed := 0
	for _, r := range rules.stopUnlessFound {
		for _, d := range details {
			regex := r.param1.(*regexp.Regexp)
			if r.scope.FindAllStringSubmatch(d[0], -1) != nil && regex.FindAllStringSubmatch(d[1], -1) != nil {
				passed++
			}
		}
	}
	if passed != len(rules.stopUnlessFound) {
		return nil
	}
	passed = 0
	for _, r := range rules.stopUnless {
		for _, d := range details {
			regex := r.param1.(*regexp.Regexp)
			if r.scope.FindAllStringSubmatch(d[0], -1) != nil && regex.FindAllStringSubmatch(d[1], -1) != nil {
				passed++
			}
		}
	}
	if passed != len(rules.stopUnless) {
		return nil
	}

	// do sampling if configured
	if len(rules.sample) == 1 && rand.Intn(100) >= rules.sample[0].param1.(int) {
		return nil
	}

	// winnow sensitive details based on remove rules if configured
	for _, r := range rules.remove {
		details = removeDetailIf(details, [][]interface{}{{true, r.scope}})
	}
	for _, r := range rules.removeUnlessFound {
		details = removeDetailIf(details, [][]interface{}{{true, r.scope}, {false, r.param1.(*regexp.Regexp)}})
	}
	for _, r := range rules.removeIfFound {
		details = removeDetailIf(details, [][]interface{}{{true, r.scope}, {true, r.param1.(*regexp.Regexp)}})
	}
	for _, r := range rules.removeUnless {
		details = removeDetailIf(details, [][]interface{}{{true, r.scope}, {false, r.param1.(*regexp.Regexp)}})
	}
	for _, r := range rules.removeIf {
		details = removeDetailIf(details, [][]interface{}{{true, r.scope}, {true, r.param1.(*regexp.Regexp)}})
	}
	if len(details) == 0 {
		return nil
	}

	// mask sensitive details based on replace rules if configured
	for _, r := range rules.replace {
		for _, d := range details {
			if r.scope.FindAllStringSubmatch(d[0], -1) != nil {
				d[1] = r.param1.(*regexp.Regexp).ReplaceAllString(d[1], r.param2.(string))
			}
		}
	}

	// remove any details with empty values
	i := 0
	for _, d := range details {
		if d[1] != "" {
			details[i] = d
			i++
		}
	}
	details = details[:i]
	if len(details) == 0 {
		return nil
	}

	return details
}

/*
The following unexported Regexps should be treated as constants
and remain unchanged throughout package usage
*/
var regexAllowHttpUrl *regexp.Regexp = regexp.MustCompile(`^\s*allow_http_url\s*(#.*)?$`)
var regexBlankOrComment *regexp.Regexp = regexp.MustCompile(`^\s*([#].*)*$`)
var regexCopySessionField *regexp.Regexp = regexp.MustCompile(`^\s*copy_session_field\s+([~!%|\/].+[~!%|\/])\s*(#.*)?`)
var regexRemove *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove\s*(#.*)?$`)
var regexRemoveIf *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexRemoveIfFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexRemoveUnless *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexRemoveUnlessFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexReplace *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*replace[\s]+([~!%|\/].+[~!%|\/]),[\s]+([~!%|\/].*[~!%|\/])\s*(#.*)?$`)
var regexSample *regexp.Regexp = regexp.MustCompile(`^\s*sample\s+(\d+)\s*(#.*)?$`)
var regexSkipCompression *regexp.Regexp = regexp.MustCompile(`^\s*skip_compression\s*(#.*)?$`)
var regexSkipSubmission *regexp.Regexp = regexp.MustCompile(`^\s*skip_submission\s*(#.*)?$`)
var regexStop *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop\s*(#.*)?$`)
var regexStopIf *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexStopIfFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexStopUnless *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexStopUnlessFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)

/*
used in ruleFilter method to compare given rule string
and a parse rule's verb
*/
func ruleCompare(ruleString string, s string) bool {
	return ruleString == s
}

// filter a slice of HttpRules comparing a given rule string with an HttpRule's verb
func ruleFilter(parsedRules []*HttpRule, ruleString string, cond func(string, string) bool) []*HttpRule {
	result := []*HttpRule{}
	for _, rule := range parsedRules {
		if cond(ruleString, rule.verb) {
			result = append(result, rule)
		}
	}
	return result
}

// https://stackoverflow.com/questions/20545743/delete-entries-from-a-slice-while-iterating-over-it-in-go/20551116
// remove detail form details slice if all regexp.Regexp are matched from the given slice of *regexp.Regexp
func removeDetailIf(details [][]string, condRegex [][]interface{}) [][]string {
	i := 0
	for _, d := range details {
		allMatched := true
		for i, exp := range condRegex {
			if (exp[1].(*regexp.Regexp).FindAllStringSubmatch(d[i], -1) == nil) == exp[0].(bool) {
				allMatched = false
			}
		}
		/*
			for current details item, no
			expression was matched so
			include it in the details
			slice
		*/
		if !allMatched {
			details[i] = d
			i++
		}
	}
	details = details[:i]
	return details
}
