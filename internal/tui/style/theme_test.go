package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoveryColor(t *testing.T) {
	assert.Equal(t, ColorGreen, RecoveryColor(78.0))
	assert.Equal(t, ColorGreen, RecoveryColor(67.0))
	assert.Equal(t, ColorYellow, RecoveryColor(66.9))
	assert.Equal(t, ColorYellow, RecoveryColor(45.0))
	assert.Equal(t, ColorYellow, RecoveryColor(34.0))
	assert.Equal(t, ColorRed, RecoveryColor(33.9))
	assert.Equal(t, ColorRed, RecoveryColor(25.0))
	assert.Equal(t, ColorRed, RecoveryColor(0.0))
}
