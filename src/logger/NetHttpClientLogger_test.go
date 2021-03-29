package logger

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing NET HTTP Logger

func TestLogsGet(t *testing.T) {

	queue := make([]string, 1)
	options := Options{
		queue: queue,
	}
	netLogger := NewNetHttpClientLoggerOptions(options)
	helper := GetTestHelper()

	netLogger.Get("google.com")

	//netLogger.httpLogger.Queue()[0]
	//I believe this is correct however not sure if its functioning properly as nothing is in the queue.
	fmt.Println(netLogger.httpLogger.queue[0])

	fmt.Println(netLogger.httpLogger.Queue()[0])
	assert.True(t, parseable(netLogger.httpLogger.Queue()[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.Queue()[0], "[\"request_method\",\"GET\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.Queue()[0], "[\"request_url\",\""+helper.demoURL+"\"]"))
	//Dependding on what the get actually gets this could change the response body
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
		queue: queue,
	}
	netLogger := NewNetHttpClientLoggerOptions(options)
	helper := GetTestHelper()

	netLogger.Post(helper.demoURL, "text/html", bytes.NewBuffer([]byte(helper.mockJSON)))

	assert.True(t, parseable(queue[0]))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_header:content-type\",\"Application/JSON\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_method\",\"POST\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_param:message\",\""+helper.mockJSONescaped+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_url\",\""+helper.demoURL+"?"+helper.mockJSON+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_body\",\""+helper.mockJSONescaped+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_header:content-type\",\"application/json; charset=utf-8\"]"))
}

//Similar to Get Just doesn't return a body.
func TestLogsHead(t *testing.T) {
	queue := make([]string, 0)
	options := Options{
		queue: queue,
	}
	netLogger := NewNetHttpClientLoggerOptions(options)
	helper := GetTestHelper()

	netLogger.Head(helper.demoURL)

	assert.True(t, parseable(queue[0]))
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
	queue := make([]string, 0)
	options := Options{
		queue: queue,
	}
	netLogger := NewNetHttpClientLoggerOptions(options)
	helper := GetTestHelper()
	form := url.Values{}
	form.Add("username:", "resurfaceio")

	netLogger.PostForm(helper.demoURL, form)

	assert.True(t, parseable(queue[0]))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_header:content-type\",\"application/x-www-form-urlencoded\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_method\",\"POST\"]"))
	//Not sure where postform data is held in the param message or within the url. Will have to wait to see exactly
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_param:message\",\""+helper.mockFormData+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"request_url\",\""+helper.demoURL+"?"+helper.mockFormData+"\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(queue[0], "[\"response_header:content-type\",\"application/json; charset=utf-8\"]"))
}
