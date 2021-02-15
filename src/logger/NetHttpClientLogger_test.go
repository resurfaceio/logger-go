package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing NET HTTP Logger

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
