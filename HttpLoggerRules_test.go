// © 2016-2023 Graylog, Inc.

package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"strings"

	"net/http"

	"io/ioutil"

	"bytes"
)

// test override default rules

func TestOverrideDefaultRules(t *testing.T) {
	rules := GetHttpRules()
	assert.Equal(t, rules.strictRules, rules.defaultRules, "HTTP default rules are not strict rules")

	options := Options{
		Url: "https://mysite.com",
	}
	logger, _ := NewHttpLogger(options)
	assert.Equal(t, httpRules.StrictRules(), logger.rules.text, "logger rules are not set to default rules")
	options = Options{
		Url:   "https://mysite.com",
		Rules: "# 123",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "# 123", logger.rules.text, "logger default rules not overriden")

	httpRules.SetDefaultRules("")
	options = Options{
		Url: "https://mysite.com",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "", logger.rules.text, "logger default rules were not applied")
	options = Options{
		Url:   "https://mysite.com",
		Rules: "   ",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "", logger.rules.text, "logger default rules not overriden or blank space not ignored")
	options = Options{
		Url:   "https://mysite.com",
		Rules: " sample 42",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, " sample 42", logger.rules.text, "logger default rules not overriden")

	httpRules.SetDefaultRules("skip_compression")
	options = Options{
		Url: "https://mysite.com",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "skip_compression", logger.rules.text, "logger default rules not applied")
	options = Options{
		Url:   "https://mysite.com",
		Rules: "include default\nskip_submission\n",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "skip_compression\nskip_submission\n", logger.rules.text, ":logger default rules not overriden")

	httpRules.SetDefaultRules("sample 42\n")
	options = Options{
		Url: "https://mysite.com",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "sample 42\n", logger.rules.text, "logger default rules not applied")
	options = Options{
		Url:   "https://mysite.com",
		Rules: "   ",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "sample 42\n", logger.rules.text, "white space not ignored")
	options = Options{
		Url:   "https://mysite.com",
		Rules: "include default\nskip_submission\n",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, "sample 42\n\nskip_submission\n", logger.rules.text, "logger rules not applied correctly")

	httpRules.SetDefaultRules("include debug")
	options = Options{
		Url:   "https://mysite.com",
		Rules: httpRules.strictRules,
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, httpRules.strictRules, logger.rules.text, "logger default rules not overriden")

	httpRules.SetDefaultRules(httpRules.strictRules)
}

// test uses allow http url rules

func TestUsesAllowHttpUrlRules(t *testing.T) {
	// requires url, rules, and enableable to be in logger struct !!!
	options := Options{
		Url: "http://mysite.com",
	}
	logger, _ := NewHttpLogger(options)
	assert.Equal(t, false, logger.enableable, "Logger enableable flag should be set to false")

	options = Options{
		Url:   "http://mysite.com",
		Rules: "",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, false, logger.enableable, "Logger enableable flag should be set to false")

	options = Options{
		Url: "https://mysite.com",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, true, logger.enableable, "Logger enableable flag should be set to true")

	options = Options{
		Url:   "https://mysite.com",
		Rules: "allow_http_url",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, true, logger.enableable, "Logger enableable flag should be set to true")

	options = Options{
		Url:   "https://mysite.com",
		Rules: "allow_http_url\nallow_http_url",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, true, logger.enableable, "Logger enableable flag should be set to true")
}

// test uses copy session field rules test
func TestUsesCopySessionFieldRules(t *testing.T) {
	// helper for function tests
	helper := newTestHelper()
	// requests used for all tests in function
	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML))
	mockResponse.Request = request

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)

	// tests copy all of session field
	_queue := make([]string, 0)
	options := Options{
		Rules: "copy_session_field /.*/",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:butterfly\",\"poison\"]"), "_queue did not contain expected values")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "_queue did not contain expected values")
	// tests copy specifically session_id
	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field /session_id/",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:butterfly\",\"poison\"]"), "_queue contains unexpected value")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "_queue did not contain expected values")
	// tests copy non matching term
	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field /blah/",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:"), "_queue contains unexpected value")
	// tests copy 2 specific values
	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field /butterfly/\ncopy_session_field /session_id/",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:butterfly\",\"poison\"]"), "_queue did not contain expected values")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "_queue did not contain expected values")
}

// test uses copy session field and remove rules test

func TestUsesCopySessionFieldAndRemoveRules(t *testing.T) {
	// helper for function tests
	helper := newTestHelper()
	// requests used for all tests in function
	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML))
	mockResponse.Request = request

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)

	//tests copy session field w/ remove
	_queue := make([]string, 0)
	options := Options{
		Rules: "copy_session_field !.*!\n!session_field:.*! remove",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:"), "_queue did contains an unexpected value")
	//
	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field !.*!\n!session_field:butterfly! remove",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:butterfly\","), "_queue contains unexpected values")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "_queue did not contain expected value")

	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field !.*!\n!session_field:.*! remove_if !poi.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:butterfly\","), "_queue contains unexpected values")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "_queue did not contain expected value")

	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field !.*!\n!session_field:.*! remove_unless !sugar!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"session_field:"), "_queue contains unexpected value")
}

// test uses copy session field and stop rules

func TestUsesCopySessionFieldAndStopRules(t *testing.T) {
	// helper for function tests
	helper := newTestHelper()
	// requests used for all tests in function
	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML))
	mockResponse.Request = request

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)

	_queue := make([]string, 0)
	options := Options{
		Rules: "copy_session_field !.*!\n!session_field:butterfly! stop",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field !.*!\n!session_field:butterfly! stop_if !poi.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	_queue = make([]string, 0)
	options = Options{
		Rules: "copy_session_field !.*!\n!session_field:butterfly! stop_unless !sugar!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")
}

// test uses remove rules

func TestUsesRemoveRules(t *testing.T) {
	//helper for function
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!.*! remove",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! remove",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! remove",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body|response_body! remove",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_header:.*! remove",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_header:"), "request_header not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_header:abc! remove\n!response_body! remove",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_header:"), "request_header not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_header:abc\","), "request_header:abc not removed")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")
}

// test uses remove if rules

func TestUsesRemoveIfRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! remove_if !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "length of _queue is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!.*! remove_if !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! remove_if !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! remove_if !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_if !.*World.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_if !.*blahblahblah.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! remove_if !.*!\n!response_body! remove_if !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")
}

// test uses remove if found rules
func TestUsesRemoveIfFoundRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! remove_if_found !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!.*! remove_if_found !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! remove_if_found !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! remove_if_found !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_if_found !World!", //mock response should contain "World" and should therefore be removed
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_if_found !.*World.*!", //mock response should contain "World" and should therefore be removed
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_if_found !blahblahblah!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")
}

// test uses remove unless rules

func TestUsesRemoveUnlessRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! remove_unless !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!.*! remove_unless !.*blahblahblah.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! remove_unless !.*blahblahblah.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! remove_unless !.*blahblahblah.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_unless !.*World.*!", //mock response should contain "World" and therefore should not be removed
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_unless !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! remove_unless !.*!\n!request_body! remove_unless !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")
}

// test uses remove unless found rules

func TestUsesRemoveUnlessFoundRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! remove_unless_found !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!.*! remove_unless_found !blahblahblah!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! remove_unless_found !blahblahblah!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! remove_unless_found !blahblahblah!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not removed")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_unless_found !World!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_unless_found !.*World.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body|request_body! remove_unless_found !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "response_body not found")
}

// test uses replace rules

func TestUsesReplaceRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_body! replace !blahblahblah!, !ZZZZZ!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "World"), "_queue was altered unexpectedly") //default mock response should contain "World"
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "ZZZZZ"), "_queue was altered unexpectedly")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !World!, !Mundo!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>Hello Mundo!</html>\"],"), "_queue was not altered")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body|response_body! replace !^.*!, !ZZZZZ!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\",\"ZZZZZ\"],"), "_queue was not altered")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"ZZZZZ\"],"), "_queue was not altered")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! replace !^.*!, !QQ!\n!response_body! replace !^.*!, !SS!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"request_body\",\"QQ\"],"), "_queue was not altered")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"SS\"],"), "_queue was not altered")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !World!, !!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>Hello !</html>\"],"), "_queue was not altered")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !.*!, !!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, false, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\","), "_queue was not altered")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML3)) //change html used from helper to mockHtml3
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !World!, !Z!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>1 Z 2 Z Red Z Blue Z!</html>\"],"), "_queue was not altered")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML4)) //change html used from helper to mockHtml4
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !World!, !Z!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>1 Z\\n2 Z\\nRed Z \\nBlue Z!\\n</html>\"],"), "_queue was not altered")
}

// test uses replace rules with complex expressions
func TestUsesReplaceRulesWithComplexExpressions(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()

	mockHtml := strings.Replace(helper.mockHTML, "World", "rob@resurface.io", 1)
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(mockHtml))
	_queue := make([]string, 0)
	options := Options{
		Rules: "/response_body/ replace /[a-zA-Z0-9.!#$%&’*+\\/=?^_'{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)/, /x@y.com/",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>Hello x@y.com!</html>\"],"), "email not replaced in _queue")

	mockResponse = helper.MockResponseWithHtml()
	mockHtml = strings.Replace(helper.mockHTML, "World", "123-45-1343", 1)
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(mockHtml))
	_queue = make([]string, 0)
	options = Options{
		Rules: "/response_body/ replace /[0-9\\.\\-\\/]{9,}/, /xyxy/",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>Hello xyxy!</html>\"],"), "custom string not replaced in _queue")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML))
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !World!, !<b>$0</b>!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>Hello <b>World</b>!</html>\"],"), "custom string not replaced in _queue")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !(World)!, !<b>$1</b>!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>Hello <b>World</b>!</html>\"],"), "custom string not replaced in _queue")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML5))
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! replace !<input([^>]*)>([^<]*)</input>!, !<input$1></input>!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
	assert.Equal(t, true, strings.Contains(logger.baseLogger.queue[0], "[\"response_body\",\"<html>\\n<input type=\\\"hidden\\\"></input>\\n<input class='foo' type=\\\"hidden\\\"></input>\\n</html>\"],"), "custom string not replaced in _queue")
}

// test uses sample rules
func TestUsesSampleRules(t *testing.T) {
	helper := newTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML))
	mockResponse.Request = request

	_queue := make([]string, 0)

	options := Options{
		Rules: "sample 10\nsample 99",
		Queue: _queue,
	}
	_, err := NewHttpLogger(options)
	if err != nil {
		assert.Equal(t, "multiple sample rules", err.Error(), "multiple sample rule error not correct") //This is what I came up with as Go has no Try & Catch functionality for errors
	}

	options = Options{
		Rules: "sample 10",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	for i := 1; i <= 100; i++ {
		SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	}
	assert.GreaterOrEqual(t, len(logger.baseLogger.queue), 2, "sample amount is less than specified 10")
	assert.LessOrEqual(t, len(logger.baseLogger.queue), 20, "sample amount is greater than specified 10")
}

// test uses skip compression rules
func TestUsesSkipCompressionRules(t *testing.T) {
	options := Options{
		Url: "http://mysite.com",
	}
	logger, _ := NewHttpLogger(options)
	assert.Equal(t, false, logger.skipCompression, "Logger skipCompression flag should be set to false")

	options = Options{
		Url:   "http://mysite.com",
		Rules: "",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, false, logger.skipCompression, "Logger skipCompression flag should be set to false")

	options = Options{
		Url:   "http://mysite.com",
		Rules: "skip_compression",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, true, logger.skipCompression, "Logger skipCompression flag should be set to true")
}

// test uses skip submission rules
func TestUsesSkipSubmission(t *testing.T) {
	options := Options{
		Url: "http://mysite.com",
	}
	logger, _ := NewHttpLogger(options)
	assert.Equal(t, false, logger.skipSubmission, "Logger skipSubmission flag should be set to false")

	options = Options{
		Url:   "http://mysite.com",
		Rules: "",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, false, logger.skipSubmission, "Logger skipSubmission flag should be set to false")

	options = Options{
		Url:   "http://mysite.com",
		Rules: "skip_submission",
	}
	logger, _ = NewHttpLogger(options)
	assert.Equal(t, true, logger.skipSubmission, "Logger skipSubmission flag should be set to true")
}

// test uses stop rules

func TestUsesStopRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! stop",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	_queue = make([]string, 0)
	options = Options{
		Rules: "!.*! stop",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Body = nil
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! stop",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	fmt.Println()
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHTML))
	mockResponse.Request.Body = nil
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Body = nil
	mockResponse.Request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	_queue = make([]string, 0)
	options = Options{
		Rules: "!request_body! stop\n!response_body! stop",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")
}

// test uses stop if rules

func TestUsesStopIfRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! stop_if !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_if !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_if !.*World.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_if !.*blahblahblah.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
}

// test uses stop if found rules

func TestUsesStopIfFoundRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! stop_if_found !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	mockResponse.Request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJSON))
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_if_found !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_if_found !World!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_if_found !.*World.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_if_found !blahblahblah!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")
}

// test uses stop unless rules
func TestUsesStopUnlessRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! stop_unless !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_unless !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_unless !.*World.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_unless !.*blahblahblah.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")
}

// test uses stop unless found rules
func TestUsesStopUnlessFoundRules(t *testing.T) {
	helper := newTestHelper()

	mockResponse := helper.MockResponseWithHtml()
	_queue := make([]string, 0)
	options := Options{
		Rules: "!response_header:blahblahblah! stop_unless_found !.*!",
		Queue: _queue,
	}
	logger, _ := NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_unless_found !.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_unless_found !World!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_unless_found !.*World.*!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 1, len(logger.baseLogger.queue), "_queue length is not 1")

	mockResponse = helper.MockResponseWithHtml()
	_queue = make([]string, 0)
	options = Options{
		Rules: "!response_body! stop_unless_found !blahblahblah!",
		Queue: _queue,
	}
	logger, _ = NewHttpLogger(options)
	SendHttpMessage(logger, mockResponse, mockResponse.Request, 0, 0, nil)
	assert.Equal(t, 0, len(logger.baseLogger.queue), "_queue is not empty")
}
