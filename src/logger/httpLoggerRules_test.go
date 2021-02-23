package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"strings"

	"net/http"
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
	request := http.Request{}

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)
	mockResponse := http.Response{Request: &request}

	// tests copy all of session field
	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "copy_session_field /.*/")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:butterfly\",\"poison\"]"), "queue did not contain expected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected values")
	// tests copy specifically session_id
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field /session_id/")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:butterfly\",\"poison\"]"), "queue contains unexpected value")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected values")
	// tests copy non matching term
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field /blah/")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:"), "queue contains unexpected value")
	// tests copy 2 specific values
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field /butterfly/\ncopy_session_field /session_id/")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:butterfly\",\"poison\"]"), "queue did not contain expected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected values")
}

// test uses copy session field and remove rules test

func TestUsesCopySessionFieldAndRemoveRules(t *testing.T) {
	// helper for function tests
	helper := GetTestHelper()
	// requests used for all tests in function
	request := http.Request{}

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)
	mockResponse := http.Response{Request: &request}
	//tests copy session field w/ remove
	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:.*! remove")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:"), "queue did contains an unexpected value")
	//
	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! remove")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:butterfly\","), "queue contains unexpected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected value")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:.*! remove_if !poi.*!")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:butterfly\","), "queue contains unexpected values")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"session_field:session_id\",\"asdf1234\"]"), "queue did not contain expected value")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:.*! remove_unless !sugar!")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"session_field:"), "queue contains unexpected value")
}

// test uses copy session field and stop rules

func TestUsesCopySessionFieldAndStopRules(t *testing.T) {
	// helper for function tests
	helper := GetTestHelper()
	// requests used for all tests in function
	request := http.Request{}

	var c1 *http.Cookie = &http.Cookie{Name: "butterfly", Value: "poison"}
	var c2 *http.Cookie = &http.Cookie{Name: "session_id", Value: "asdf1234"}
	request.AddCookie(c1)
	request.AddCookie(c2)
	mockResponse := http.Response{Request: &request}

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! stop")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! stop_if !poi.*!")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "copy_session_field !.*!\n!session_field:butterfly! stop_unless !sugar!")
	httpMessage.sendHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")
}

// test uses remove rules

func TestUsesRemoveRules(t *testing.T) {
	//helper for function
	helper := GetTestHelper()

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!.*! remove")
	httpMessage.Send(logger, MockRequestWithJson2(), MockResponseWithHtml(), helper.mockHTML, helper.mockJSON)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! remove")
	httpMessage.Send(logger, MockRequestWithJson2(), MockResponseWithHtml(), helper.mockHTML, helper.mockJSON)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body was not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")
}

// test uses remove if rules

// test uses remove if found rules

// test uses remove unless rules

// test uses remove unless found rules

// test uses replace rules

// test uses replace rules with complex expressions

// test uses sample rules

// test uses skip compression rules

// test uses skip submission rules

// test uses stop rules

// test uses stop if rules

// test uses stop if found rules

// test uses stop unless rules

// test uses stop unless found rules
