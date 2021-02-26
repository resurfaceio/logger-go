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

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!.*! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body|response_body! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_header:.*! remove")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_header:"), "request_header not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

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

func TestUsesRemoveIfRules(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!response_header:blahblahblah! remove_if !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "length of queue is not 1")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!.*! remove_if !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! remove_if !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! remove_if !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_if !.*World.*!") //mock response should contain "World" and therefore be removed
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_if !.*blahblahblah.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! remove_if !.*!\n!response_body! remove_if !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")
}

// test uses remove if found rules

func TestUsesRemoveIfFoundRules(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!response_header:blahblahblah! remove_if_found !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!.*! remove_if_found !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! remove_if_found !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! remove_if_found !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_if_found !World!") //mock response should contain "World" and should therefore be removed
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_if_found !.*World.*!") //mock response should contain "World" and should therefore be removed
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_if_found !blahblahblah!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")
}

// test uses remove unless rules

func TestUsesRemoveUnlessRules(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!response_header:blahblahblah! remove_unless !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!.*! remove_unless !.*blahblahblah.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! remove_unless !.*blahblahblah.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! remove_unless !.*blahblahblah.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_unless !.*World.*!") //mock response should contain "World" and therefore should not be removed
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_unless !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! remove_unless !.*!\n!request_body! remove_unless !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")
}

// test uses remove unless found rules

func TestUsesRemoveUnlessFoundRules(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!response_header:blahblahblah! remove_unless_found !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!.*! remove_unless_found !blahblahblah!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 0, len(queue), "queue is not empty")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! remove_unless_found !blahblahblah!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! remove_unless_found !blahblahblah!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "response_body not removed")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_unless_found !World!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_unless_found !.*World.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"request_body\","), "request_body not removed")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body|request_body! remove_unless_found !.*!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\","), "request_body not found")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\","), "response_body not found")
}

// test uses replace rules

func TestUsesReplaceRules(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!response_body! replace !blahblahblah!, !ZZZZZ!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "World"), "queue was altered unexpectedly") //default mock response should contain "World"
	assert.Equal(t, false, strings.Contains(queue[0], "ZZZZZ"), "queue was altered unexpectedly")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !World!, !Mundo!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>Hello Mundo!</html>\"],"), "queue was not altered")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body|response_body! replace !^.*!, !ZZZZZ!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\",\"ZZZZZ\"],"), "queue was not altered")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"ZZZZZ\"],"), "queue was not altered")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!request_body! replace !^.*!, !QQ!\n!response_body! replace !^.*!, !SS!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_body\",\"QQ\"],"), "queue was not altered")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"SS\"],"), "queue was not altered")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !World!, !!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>Hello !</html>\"],"), "queue was not altered")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !.*!, !!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, false, strings.Contains(queue[0], "[\"response_body\","), "queue was not altered")

	request = helper.MockRequestWithJson2()
	mockResponse = helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml3)) //change html used from helper to mockHtml3
	mockResponse.Request = request

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !World!, !Z!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>1 Z 2 Z Red Z Blue Z!</html>\"],"), "queue was not altered")

	request = helper.MockRequestWithJson2()
	mockResponse = helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml4)) //change html used from helper to mockHtml4
	mockResponse.Request = request

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !World!, !Z!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>1 Z\\n2 Z\\nRed Z \\nBlue Z!\\n</html>\"],"), "queue was not altered")
}

// test uses replace rules with complex expressions

func TestUsesReplaceRulesWithComplexExpressions(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	mockHtml := strings.Replace(helper.mockJSON, "World", "rob@resurface.io", 1)

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(mockHtml))
	mockResponse.Request = request

	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "/response_body/ replace /[a-zA-Z0-9.!#$%&’*+\\/=?^_\'{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)/, /x@y.com/")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>Hello x@y.com!</html>\"],"), "email not replaced in queue")

	request = helper.MockRequestWithJson2()
	mockResponse = helper.MockResponseWithHtml()
	
	mockHtml = strings.Replace(helper.mockJSON, "World", "123-45-1343", 1)

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(mockHtml))
	mockResponse.Request = request

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "/response_body/ replace /[0-9\\.\\-\\/]{9,}/, /xyxy/")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>Hello xyxy!</html>\"],"), "custom string not replaced in queue")

	request = helper.MockRequestWithJson2()
	mockResponse = helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !World!, !<b>$&</b>!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>Hello <b>World</b>!</html>\"],"), "custom string not replaced in queue")

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !(World)!, !<b>$1</b>!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>Hello <b>World</b>!</html>\"],"), "custom string not replaced in queue")	

	request = helper.MockRequestWithJson2()
	mockResponse = helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml5))
	mockResponse.Request = request

	queue = make([]string, 0)
	logger = NewHttpLoggerQueueRules(queue, "!response_body! replace !<input([^>]*)>([^<]*)</input>!, !<input$1></input>!")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\"<html>\\n<input type=\\\"hidden\\\"></input>\\n<input class='foo' type=\\\"hidden\\\"></input>\\n</html>\"],"), "custom string not replaced in queue")
}

// test uses sample rules

func TestUsesSampleRules(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson))
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request
	
	queue := make([]string, 0)

	logger, err := NewHttpLoggerQueueRules(queue, "sample 10\nsample 99")
	if err != nil {
		assert.Equal(t, "Multiple sample rules", err.message, "multiple sample rule error not correct") //This is what I came up with as Go has no Try & Catch functionality for errors
	}
	
	logger = NewHttpLoggerQueueRules(queue, "sample 10")
	for (i := 1; i <= 100; i++) {
		httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	}
	assert.Greater(t, len(queue), 2, "sample amount is less than specified 10")
	assert.Less(t, len(queue), 20, "sample amount is greater than specified 10")
}

// test uses skip compression rules

func TestUsesSkipCompressionRules(t *testing.T) {
	logger := NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	assert.Equal(t, false, logger.SkipCompression(), "Logger skipCompression flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("")
	assert.Equal(t, false, logger.SkipCompression(), "Logger skipCompression flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("skip_compression")
	assert.Equal(t, true, logger.SkipCompression(), "Logger skipCompression flag should be set to true")
}

// test uses skip submission rules

func TestUsesSkipSubmission(t *testing.T) {
	logger := NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	assert.Equal(t, false, logger.SkipSubmission(), "Logger skipSubmission flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("")
	assert.Equal(t, false, logger.SkipSubmission(), "Logger skipSubmission flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("skip_submission")
	assert.Equal(t, true, logger.SkipSubmission(), "Logger skipSubmission flag should be set to true")
}

// test uses stop rules

func TestUsesStopRules(t *testing.T) {
	helper := GetTestHelper()

	request := helper.MockRequestWithJson2()
	mockResponse := helper.MockResponseWithHtml()

	request.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockJson)) //not sure how to do null here
	mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(helper.mockHtml))
	mockResponse.Request = request
	
	queue := make([]string, 0)
	logger := NewHttpLoggerQueueRules(queue, "!response_header:blahblahblah! stop")
	httpMessage.sendNetHttpRequestResponseMessage(logger, mockResponse, 0, 0)
	assert.Equal(t, 1, len(queue), "queue length is not 1")
}

// test uses stop if rules

// test uses stop if found rules

// test uses stop unless rules

// test uses stop unless found rules
