package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing NET HTTP Logger

func TestLogsHtml(t *testing.T) {

	queue := []string{}
	filter := newLogger(queue, "includes standard")
	helper := NewTestHelper()
	filter.init(nil)
	filter.doFilter(helper.mockRequest(), helper.mockResponse(), helper.mockJsonApp())
	assert.Equal(t, 1, len(queue))
	msg := queue[0]
	assert.True(t, parseable(msg))
	assert.Contains(t, msg, "[\"request_method\",\"GET\"]")
	assert.Contains(t, msg, "[\"request_url\",\"" + helper.mockURL + "\"]")
	assert.Contains(t, msg, "[\"response_body\",\"" + helper.mockHTML + "\"]")
	assert.Contains(t, msg, "[\"response_code\",\"200\"]")
	assert.Contains(t, msg, "[\"response_header:a\",\"Z\"]")
	assert.Contains(t, msg, "[\"response_header:content-type\",\"text/html\"]")
	assert.Contains(t, msg, "[\"now\",\"")
	assert.Contains(t, msg, "[\"interval\",\"")
	assert.NotContains(t, msg, "request_body")
	assert.NotContains(t, msg, "request_header")
	assert.NotContains(t, msg, "request_param")

}

func TestlogJson(t *testing.T) {
	helper := NewTestHelper()
	queue := []string{}
	//do filter
	//this block will do the actual logging
	//but we don't know how this works yet
	//make filter, initialize, do filter
	assert.Equal(t, 1, len(queue))
	msg := queue[0]
	//must implement parsable function
	//assert.True(t,parsable(msg))
	assert.Contains(t, msg, "[\"request_method\",\"GET\"]")
	assert.Contains(t, msg, "[\"response_body\",\""+helper.mockJSONescaped+"\"]")
	assert.Contains(t, msg, "[\"response_code\",\"200\"]")
	assert.Contains(t, msg, "[\"response_header:content-type\",\"application/json; charset=utf-8\"]")
	assert.NotContains(t, msg, "request_body")
	assert.NotContains(t, msg, "request_header")
	assert.NotContains(t, msg, "request_param")
}
