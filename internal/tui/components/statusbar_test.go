package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusBar_ContainsMessage(t *testing.T) {
	result := StatusBar("Press ? for help", 80)
	assert.Contains(t, result, "Press ? for help")
}

func TestErrorBar_ContainsError(t *testing.T) {
	result := ErrorBar("connection failed", 80)
	assert.Contains(t, result, "Error:")
	assert.Contains(t, result, "connection failed")
}
