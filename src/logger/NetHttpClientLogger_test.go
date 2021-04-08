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
	helper := GetTestHelper()
	queue := make([]string, 0) // argument of 1 vs 0
	options := Options{
		url:     helper.demoURL,
		queue:   queue,
		enabled: true,
	}
	netLogger := NewNetHttpClientLoggerOptions(options)

	netLogger.Get(helper.demoURL1)
	//queue = append(queue, "[\"request_method\",\"GET\"]") // this is for testing that the queue can hold strings
	//Populate the queue by building these requests and responses
	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	fmt.Println(netLogger.httpLogger.queue[0])
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\", \"GET\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\", \""+helper.demoURL1+"\"]"))
	//Dependding on what the get actually gets this could change the response body
	// assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_body\", \""+helper.mockHTML+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\", \"200\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:a\", \"Z\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:content-type\", \"text/html\"]"))
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

	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_header:content-type\", \"Application/JSON\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\", \"POST\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\", \""+helper.demoURL+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\", \"204\"]"))
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

	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\",\"HEAD\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\",\""+helper.demoURL+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:a\",\"Z\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:content-type\",\"text/html\"]"))
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
	}
	netLogger := NewNetHttpClientLoggerOptions(options)
	form := url.Values{}
	form.Add("username:", "resurfaceio")

	netLogger.PostForm(helper.demoURL, form)
	fmt.Println(netLogger.httpLogger.queue[0])
	assert.True(t, parseable(netLogger.httpLogger.queue[0]))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_header:content-type\",\"application/x-www-form-urlencoded\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_method\",\"POST\"]"))
	//Not sure where postform data is held in the param message or within the url. Will have to wait to see exactly
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_param:message\",\""+helper.mockFormData+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"request_url\",\""+helper.demoURL+"\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_code\",\"200\"]"))
	assert.Equal(t, true, strings.Contains(netLogger.httpLogger.queue[0], "[\"response_header:content-type\",\"application/json; charset=utf-8\"]"))
}
