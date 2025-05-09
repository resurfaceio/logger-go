// © 2016-2025 Graylog, Inc.

package logger

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var helloWorldJson = []byte{123, 32, 34, 104, 101, 108, 108, 111, 34, 58, 32, 34, 119, 111, 114, 108, 100, 34, 32, 125}

func TestLogRequestBody(t *testing.T) {
	req := MockGetRequestWithBody(helloWorldJson, "application/json")
	lr, err := NewLoggedReader(req.Body, 1024)

	assert.Nil(t, err)
	assert.Equal(t, helloWorldJson, lr.logged)

	b, err := io.ReadAll(lr.Reader)
	assert.Nil(t, err)
	assert.Equal(t, helloWorldJson, b)
}

func TestLogResponseBody(t *testing.T) {
	req := MockGetNoBodyRequest()
	resp := MockGetJSONResponse(&req)

	lr, err := NewLoggedReader(resp.Body, 1024)

	assert.Nil(t, err)
	assert.Equal(t, helperBodies["json"], lr.logged)

	b, err := io.ReadAll(lr.Reader)
	assert.Nil(t, err)
	assert.Equal(t, helperBodies["json"], b)
}

func TestLogEncodedResponse(t *testing.T) {
	req := MockGetNoBodyRequest()
	var resp http.Response

	// gzip
	resp = MockGetGzipResponse(&req)

	lr, err := NewLoggedReader(resp.Body, 1024)
	assert.Nil(t, err)
	assert.Equal(t, helperBodies["gzip"], lr.logged)

	b, err := io.ReadAll(lr.Reader)
	assert.Nil(t, err)
	assert.Equal(t, helperBodies["gzip"], b)

	// deflate
	resp = MockGetDeflateResponse(&req)

	lr, err = NewLoggedReader(resp.Body, 1024)
	assert.Nil(t, err)
	assert.Equal(t, helperBodies["deflate"], lr.logged)

	b, err = io.ReadAll(lr.Reader)
	assert.Nil(t, err)
	assert.Equal(t, helperBodies["deflate"], b)

	// brotli
	resp = MockGetBrotliResponse(&req)

	lr, err = NewLoggedReader(resp.Body, 1024)
	assert.Nil(t, err)
	assert.Equal(t, helperBodies["br"], lr.logged)

	b, err = io.ReadAll(lr.Reader)
	assert.Nil(t, err)
	assert.Equal(t, helperBodies["br"], b)
}
