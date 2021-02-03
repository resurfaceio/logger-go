package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestBaseLogger(t *testing.T) {
	logger := newLogger()
	assert.NotNil(t, logger)

}
