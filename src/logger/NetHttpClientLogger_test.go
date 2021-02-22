package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing NET HTTP Logger

func TestLogsGet(t *testing.T) {

	netLogger := NewNetHttpClientLogger()
	helper := GetTestHelper()
	resp, err := netLogger.Get(helper.demoURL)
	fmt.Println(resp)
	fmt.Println(err)
	//Don't think we will need this
	//assert.True(t, parsable(resp))
	assert.Contains(t, resp, "[\"request_method\",\"GET\"]")
	assert.Contains(t, resp, "[\"request_url\",\""+helper.demoURL+"\"]")
	assert.Contains(t, resp, "[\"response_body\",\""+helper.mockHTML+"\"]")
	assert.Contains(t, resp, "[\"response_code\",\"200\"]")
	assert.Contains(t, resp, "[\"response_header:a\",\"Z\"]")
	assert.Contains(t, resp, "[\"response_header:content-type\",\"text/html\"]")
	assert.Contains(t, resp, "[\"now\",\"")
	assert.Contains(t, resp, "[\"interval\",\"")
	assert.NotContains(t, resp, "request_body")
	assert.NotContains(t, resp, "request_header")
	assert.NotContains(t, resp, "request_param")
}

func TestLogsPost(t *testing.T) {
	netLogger := NewNetHttpClientLogger()
	helper := GetTestHelper()
	resp, err := netLogger.Post(helper.demoURL)
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
	netLogger := NewNetHttpClientLogger()
	helper := GetTestHelper()
	resp, err := netLogger.Head(helper.demoURL)
	fmt.Println(resp)
	fmt.Println(err)
}

func TestLogsPostForm(t *testing.T) {
	netLogger := NewNetHttpClientLogger()
	helper := GetTestHelper()
	resp, err := netLogger.PostForm(helper.demoURL)
	fmt.Println(resp)
	fmt.Println(err)
}
