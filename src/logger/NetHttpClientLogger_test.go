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
	assert.True(t, parsable(msg))
	assert.Contains(t, msg, "[\"request_method\",\"GET\"]")
	assert.Contains(t, msg, "[\"request_url\",\""+helper.mockURL+"\"]")
	assert.Contains(t, msg, "[\"response_body\",\""+helper.mockHTML+"\"]")
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
	filter := newLogger(queue, "includes standard")
	filter.init(nil)
	filter.doFilter(helper.mockRequest(), helper.mockResponse(), helper.mockJsonApp())
	assert.Equal(t, 1, len(queue))
	msg := queue[0]
	assert.True(t, parseable(msg))
	assert.True(t, parsable(msg))
	assert.Contains(t, msg, "[\"request_method\",\"GET\"]")
	assert.Contains(t, msg, "[\"response_body\",\""+helper.mockJSONescaped+"\"]")
	assert.Contains(t, msg, "[\"response_code\",\"200\"]")
	assert.Contains(t, msg, "[\"response_header:content-type\",\"application/json; charset=utf-8\"]")
	assert.NotContains(t, msg, "request_body")
	assert.NotContains(t, msg, "request_header")
	assert.NotContains(t, msg, "request_param")
}

func TestLogsJsonPost(t *testing.T) {
	queue := []string{}
	filter := newLogger(queue, "includes standard")
	helper := NewTestHelper()
	filter.init(nil)
	filter.doFilter(helper.mockRequestWithJson(), helper.mockResponse(), helper.mockJsonApp())
	assert.Equal(t, 1, len(queue))
	msg := queue[0]
	assert.True(t, parsable(msg))
	assert.Contains(t, msg, "[\"request_header:content-type\",\"Application/JSON\"]")
	assert.Contains(t, msg, "[\"request_method\",\"POST\"]")
	assert.Contains(t, msg, "[\"request_param:message\",\""+helper.mockJSONescaped+"\"]")
	assert.Contains(t, msg, "[\"request_url\",\""+helper.mockURL+"?"+helper.mockQueryString+"\"]")
	assert.Contains(t, msg, "[\"response_body\",\""+helper.mockJSONescaped+"\"]")
	assert.Contains(t, msg, "[\"response_code\",\"200\"]")
	assert.Contains(t, msg, "[\"response_header:content-type\",\"application/json; charset=utf-8\"]")
}

func TestJsonPostWithHeaders(t *testing.T) {
	helper := NewTestHelper()
	queue := []string{}
	filter := newLogger(queue, "includes standard")
	filter.init(nil)
	filter.doFilter(helper.mockRequest(), helper.mockResponse(), helper.mockJsonApp())
	assert.Equal(t, 1, len(queue))
	msg := queue[0]
	assert.True(t, parseable(msg))
	assert.Contains(t, msg, "[\"request_header:a\",\"1\"]")
	assert.Contains(t, msg, "[\"request_header:a\",\"2\"]")
	assert.Contains(t, msg, "[\"request_header:content-type\",\"Application/JSON\"]")
	assert.Contains(t, msg, "[\"request_method\",\"POST\"]")
	assert.Contains(t, msg, "[\"request_param:abc\",\"123\"]")
	assert.Contains(t, msg, "[\"request_param:abc\",\"234\"]")
	assert.Contains(t, msg, "[\"request_param:message\",\""+helper.mockJSONescaped+"\"]")
	assert.Contains(t, msg, "[\"request_url\",\""+helper.mockURL+"?"+helper.mockQueryString+"\"]")
	assert.Contains(t, msg, "[\"response_body\",\""+helper.mockHTML+"\"]")
	assert.Contains(t, msg, "[\"response_code\",\"200\"]")
	assert.Contains(t, msg, "[\"response_header:a\",\"Z\"]")
	assert.Contains(t, msg, "[\"response_header:content-type\",\"text/html\"]")
}

//need to figure out what these exceptions are doing
//as they can't be easily replicated in go
//so this is a bit of a fummy function
func TestSkipsException(t *testing.T) {
	queue := []string{}
	filter := newLogger(queue, "includes standard")
	filter.init(nil)
	//try {
	filter.doFilter(mockRequest(), mockResponse(), mockExceptionApp())
	//} catch (UnsupportedEncodingException uee) {
	assert.Equal(t, 0, len(queue))
	// } catch (Exception e) {
	// 	fail("Unexpected exception type");
	// }
}
