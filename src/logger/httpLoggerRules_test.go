package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"strings"

	"net/http"

	"io/ioutil"

	"bytes"
)

// test override default rules

func TestOverrideDefaultRules(t *testing.T) {
	assert.Equal(t, httpRules.StrictRules, httpRules.DefaultRules, "HTTP default rules are not strict rules")

	logger := NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, httpRules.StrictRules, logger.rules.Text, "logger rules are not set to default rules")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("# 123")
	assert.Equal(t, "# 123", logger.rules.Text, "logger default rules not overriden")

	httpRules.SetDefaultRules("")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, "", logger.rules.Text, "logger default rules were not applied")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("   ")
	assert.Equal(t, "", logger.rules.Text, "logger default rules not overriden or blank space not ignored")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules(" sample 42")
	assert.Equal(t, " sample 42", logger.rules.Text, "logger default rules not overriden")

	httpRules.SetDefaultRules("skip_compression")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, "skip_compression", logger.rules.Text, "logger default rules not applied")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("include default\nskip_submission\n")
	assert.Equal(t, "include default\nskip_submission\n", logger.rules.Text, ":logger default rules not overriden")

	httpRules.SetDefaultRules("sample 42\n")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, "sample 42\n", logger.rules.Text, "logger default rules not applied")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("   ")
	assert.Equal(t, "sample 42\n", logger.rules.Text, "white space not ignored")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("include default\nskip_submission\n")
	assert.Equal(t, "sample 42\n\nskip_submission\n", logger.rules.Text, "logger rules not applied correctly")

	httpRules.SetDefaultRules("inlude debug")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules(httpRules.StrictRules)
	assert.Equal(t, httpRules.StrictRules, logger.rules.text, "logger default rules not overriden")

	httpRules.SetDefaultRules(httpRules.StrictRules)
}

// test uses allow http url rules

func TestUsesAllowHttpUrlRules(t *testing.T) {
	// requires url, rules, and enableable to be in logger struct !!!
	logger := NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	assert.Equal(t, false, logger.Enableable(), "Logger enableable flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("")
	assert.Equal(t, false, logger.Enableable(), "Logger enableable flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, true, logger.Enableable(), "Logger enableable flag should be set to true")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("allow_http_url")
	assert.Equal(t, true, logger.Enableable(), "Logger enableable flag should be set to true")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("allow_http_url\nallow_http_url")
	assert.Equal(t, true, logger.Enableable(), "Logger enableable flag should be set to true")
}

// test uses copy session field rules test

func TestUsesCopySessionFieldRules(t *testing.T) {
	// helper for function tests
	helper := GetTestHelper()
	// requests used for all tests in function
	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)

	// tests copy all of session field
	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "copy_session_field /.*/")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:butterfly\",\"poison\"]"), "queue did not contain expected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected values")
	// tests copy specifically session_id
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field /session_id/")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:butterfly\",\"poison\"]"), "queue contains unexpected value")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected values")
	// tests copy non matching term
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field /blah/")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:"), "queue contains unexpected value")
	// tests copy 2 specific values
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field /butterfly/\ncopy_session_field /session_id/")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:butterfly\",\"poison\"]"), "queue did not contain expected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected values")
}

// test uses copy session field and remove rules test

func TestUsesCopySessionFieldAndRemoveRules(t *testing.T) {
	// helper for function tests
	helper := GetTestHelper()
	// requests used for all tests in function
	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)

	//tests copy session field w/ remove
	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:.*! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:"), "queue did contains an unexpected value")
	//
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:butterfly\","), "queue contains unexpected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected value")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:.*! remove_if !poi.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:butterfly\","), "queue contains unexpected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected value")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:.*! remove_unless !sugar!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:"), "queue contains unexpected value")
}

// test uses copy session field and stop rules

func TestUsesCopySessionFieldAndStopRules(t *testing.T) {
	// helper for function tests
	helper := GetTestHelper()
	// requests used for all tests in function
	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! stop")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! stop_if !poi.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! stop_unless !sugar!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")
}

// test uses remove rules

func TestUsesRemoveRules(t *testing.T) {
	//helper for function
	helper := GetTestHelper()

	request := helper.MockRequstWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!.*! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueules(queue, "!request_body! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "qeue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! remove")
	httpMessage.sendNetHttpReuestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(qeue, "!request_body|response_body! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "qeue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_header:.*! remove")
	httpMessage.sendNetHttpRequestRsponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_header:"), "request_header not removed")
	assert.Equal(t, true, strings.Contans(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_header:abc! remove\n!response_body! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_header:"), "request_header not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_header:abc\","), "request_header:abc not removed")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")
}

// test uses remove if rules

// test uses remove if found ruls

// test uses remove unless rules

// test uses remove unlessfound rules

// test uses replace rules

// test uses replace rule with complex expressions

// test uses sample rules

// test uses skip compression rule

// test uses skip submision rules

// test uses stop rules

// test uses stop if rules

// test uses stop if found ruls

// test uses stop unless rules

// test uses stop unless found rules
