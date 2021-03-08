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
	options := Options{
		queue := queue,
	}
	netLogger := NewNetHttpClientLogger(options)
	helper := GetTestHelper()

	resp, err = netLogger.GET(helper.demoURL)

	assert.True(t, helper.parsable(resp))
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

	queue := make([]string, 0)
	options := Options{
		queue := queue,
	}
	netLogger := NewNetHttpClientLogger(options)
	helper := GetTestHelper()

	resp, err := netLogger.Post(helper.demoURL, "text/html" /*Need IO reader format for body*/)

	assert.True(t, parsable(resp))
	assert.Equal(t,true,strings.Contains(queue[0],"[\"request_header:content-type\",\"Application/JSON\"]"))
	assert.Equal(t,true,strings.Contains(queue[0],"[\"request_method\",\"POST\"]"))
	assert.Equal(t,true,strings.Contains(queue[0],"[\"request_param:message\",\""+helper.mockJSONescaped+"\"]"))
	assert.Equal(t,true,strings.Contains(queue[0],"[\"request_url\",\""+helper.mockURL+"?"+helper.mockQueryString+"\"]"))
	assert.Equal(t,true,strings.Contains(queue[0],"[\"response_body\",\""+helper.mockJSONescaped+"\"]"))
	assert.Equal(t,true,strings.Contains(queue[0],"[\"response_code\",\"200\"]"))
	assert.Equal(t,true,strings.Contains(queue[0],"[\"response_header:content-type\",\"application/json; charset=utf-8\"]"))
}

//Similar to Get Just doesn't return a body.
func TestLogsHead(t *testing.T) {
	queue := make([]string, 0)
	options := Options{
		queue := queue,
	}
	netLogger := NewNetHttpClientLogger(options)
	helper := GetTestHelper()

	resp, err := netLogger.Head(helper.demoURL)

	assert.True(t, helper.parsable(resp))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_method\",\"HEAD\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_url\",\""+helper.demoURL+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_header:a\",\"Z\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_header:content-type\",\"text/html\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"now\",\""))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"interval\",\""))
	assert.NotEqual(t, true, strings.Contains(queue[0], "request_body"))
	assert.NotEqual(t, true, strings.Contains(queue[0], "request_header"))
	assert.NotEqual(t, true, strings.Contains(queue[0], "request_param"))

}

func TestLogsPostForm(t *testing.T) {
	netLogger := NewNetHttpClientLogger(nil)
	helper := GetTestHelper()
	resp, err := netLogger.PostForm(helper.demoURL /*data url.values*/)
	fmt.Println(resp)
	fmt.Println(err)
}
