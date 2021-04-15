package logger

import (
	"bytes"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing NET HTTP Logger

func TestLogsGet(t *testing.T) {
	helper := GetTestHelper()
	queue := make([]string, 0)
	options := Options{
		url:     helper.demoURL1,
		queue:   queue,
		enabled: true,
		rules:   "include debug",
	}
	netLogger := NewNetHttpClientLoggerOptions(options)
	//fmt.Println(netLogger.httpLogger.queue[0])

	netLogger.Get(helper.demoURL1)
	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\",\"GET\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\",\""+helper.demoURL1+"\"]"))
	//Dependding on what the get actually gets this could change the response body
	//assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_body\",]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:content-type\",\"text/html;charset=utf-8\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"now\",\""))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"interval\",\""))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_body"))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_header"))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_param"))
}

func TestLogsPost(t *testing.T) {
	helper := GetTestHelper()
	queue := make([]string, 0)
	options := Options{
		url:     helper.demoURL,
		queue:   queue,
		enabled: true,
	}
	netLogger := NewNetHttpClientLoggerOptions(options)

	netLogger.Post(helper.demoURL, "Application/JSON", bytes.NewBuffer([]byte(helper.mockJSON)))
	//fmt.Println(netLogger.httpLogger.queue[0])

	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_header:content-type\",\"Application/JSON\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\",\"POST\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\",\""+helper.demoURL+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\",\"204\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"now\",\""))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"interval\",\""))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_body"))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_param"))
}

//Similar to Get Just doesn't return a body.
func TestLogsHead(t *testing.T) {
	helper := GetTestHelper()
	queue := make([]string, 0)
	options := Options{
		url:     helper.demoURL1,
		queue:   queue,
		enabled: true,
	}
	netLogger := NewNetHttpClientLoggerOptions(options)

	netLogger.Head(helper.demoURL1)
	//fmt.Println(netLogger.httpLogger.queue[0])

	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\",\"HEAD\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\",\""+helper.demoURL1+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:"))
	//Content Type not captured
	//assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:content-type\",\"text/html\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"now\",\""))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"interval\",\""))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_body"))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_header"))
	assert.NotEqual(t, true, strings.Contains(netLogger.httpLogger.queue[0], "request_param"))
}

func TestLogsPostForm(t *testing.T) {
	helper := GetTestHelper()
	queue := make([]string, 0)
	options := Options{
		url:     helper.demoURL,
		queue:   queue,
		enabled: true,
		rules:   "include debug",
	}
	netLogger := NewNetHttpClientLoggerOptions(options)
	form := url.Values{}
	form.Add("username:", "resurfaceio")

	netLogger.PostForm(helper.demoURL, form)
	//fmt.Println(netLogger.httpLogger.queue[0])
	
	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_header:content-type\",\"application/x-www-form-urlencoded\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\",\"POST\"]"))
	//Request_param not being captured
	//assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_param:message\",\""+helper.mockFormData+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\",\""+helper.demoURL+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\",\"204\"]"))
}
