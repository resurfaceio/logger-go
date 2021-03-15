package logger

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	//need to find a way to test these assertions, haven't used them before
	assert.Equal(t, 1, 0)
	fmt.Printf("hello world!")
	var b strings.Builder
	b.WriteString("hi")
	assert.Equal(t, b.String(), "")

	//using fmt's stringbuilder for this example
	//we need to figure out what we will be using for string concatenation
	//the java project uses specific functions like apend, escape, etc
	//ours might be different

	//https://blog.golang.org/slices
	//go strings.builder
}
