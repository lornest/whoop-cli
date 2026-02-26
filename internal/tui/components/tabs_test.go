package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTabs_ContainsAllNames(t *testing.T) {
	result := RenderTabs(0, 80)
	for _, name := range TabNames {
		assert.Contains(t, result, name)
	}
}

func TestRenderTabs_HighlightsActive(t *testing.T) {
	for i := range TabNames {
		result := RenderTabs(i, 80)
		assert.Contains(t, result, TabNames[i])
	}
}
