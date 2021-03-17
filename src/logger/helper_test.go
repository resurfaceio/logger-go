package logger

import (
	"testing"
)

func test_good_json(t *testing.T){
	assert.True(t, parseable("[\n]"))
	assert.True(t, parseable("[\n\t\n"))
	assert.True(t, parseable('["A"]')) // check double quotes
	assert.True(t, parseable('["A","B"]')) // check double quotes
}

func test_invalid_json(t *testing.T){
	assert.False(t, parseable(nil))
	assert.False(t, parseable(""))
	assert.False(t, parseable(" "))
	assert.False(t, parseable("\n\n\n\n"))
	assert.False(t, parseable("1234"))
	assert.False(t, parseable("archer"))
	assert.False(t, parseable('"sterling archer"'))
	assert.False(t, parseable(",,"))
	assert.False(t, parseable("[]"))
	assert.False(t, parseable("[,,]"))
	assert.False(t, parseable('["]'))
	assert.False(t, parseable("[:,]"))
	assert.False(t, parseable(","))
	assert.False(t, parseable("exact words"))
	assert.False(t, parseable("his exact words"))
	assert.False(t, parseable('"exact words'))
	assert.False(t, parseable('his exact words"'))
	assert.False(t, parseable('"hello":"world" }'))
	assert.False(t, parseable('{ "hello":"world"'))
	assert.False(t, parseable('{ "hello world"}'))
}

func TestHelper(t *testing.T) {
	testHelper := GetTestHelper()
	if testHelper.demoURL == "" {
		t.Error("Helper DEMO_URL is empty")
	}
}
