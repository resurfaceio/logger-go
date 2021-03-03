package logger

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing NET HTTP Logger

func TestLogsGet(t *testing.T) {

	queue := make([]string, 0)
	netLogger := NewHttpLoggerQueue(queue)
	helper := GetTestHelper()
	//Don't think we will need this
	//assert.True(t, parsable(resp))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_method\",\"GET\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_url\",\""+helper.demoURL+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\""+helper.mockHTML+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_header:a\",\"Z\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_header:content-type\",\"text/html\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"now\",\""))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"interval\",\""))
	assert.NotEqual(t, true, strings.Contains(queue[0], "request_body"))
	assert.NotEqual(t, true, strings.Contains(queue[0], "request_header"))
	assert.NotEqual(t, true, strings.Contains(queue[0], "request_param"))
}

func TestLogsPost(t *testing.T) {
	netLogger := NewNetHttpClientLogger(nil)
	helper := GetTestHelper()
	resp, err := netLogger.Post(helper.demoURL, "text/html" /*Need IO reader format for body*/)
	fmt.Println(resp)
	fmt.Println(err)
	//Don't think we will need this
	//assert.True(t, parsable(resp))
	assert.Contains(t, resp, "[\"request_header:content-type\",\"Application/JSON\"]")
	assert.Contains(t, resp, "[\"request_method\",\"POST\"]")
	assert.Contains(t, resp, "[\"request_param:message\",\""+helper.mockJSONescaped+"\"]")
	assert.Contains(t, resp, "[\"request_url\",\""+helper.mockURL+"?"+helper.mockQueryString+"\"]")
	assert.Contains(t, resp, "[\"response_body\",\""+helper.mockJSONescaped+"\"]")
	assert.Contains(t, resp, "[\"response_code\",\"200\"]")
	assert.Contains(t, resp, "[\"response_header:content-type\",\"application/json; charset=utf-8\"]")
}

func TestLogsHead(t *testing.T) {
	netLogger := NewNetHttpClientLogger(nil)
	helper := GetTestHelper()
	resp, err := netLogger.Head(helper.demoURL)
	fmt.Println(resp)
	fmt.Println(err)
}

func TestLogsPostForm(t *testing.T) {
	netLogger := NewNetHttpClientLogger(nil)
	helper := GetTestHelper()
	resp, err := netLogger.PostForm(helper.demoURL /*data url.values*/)
	fmt.Println(resp)
	fmt.Println(err)
}
