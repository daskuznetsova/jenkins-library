package codeql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePattern(t *testing.T) {
	t.Parallel()

	t.Run("Include pattern", func(t *testing.T) {
		input := "+**/src/**/*"
		sign, filePattern, rulePattern, err := parsePattern(input)
		assert.NoError(t, err)
		assert.True(t, sign)
		assert.Equal(t, "**/src/**/*", filePattern)
		assert.Equal(t, "**", rulePattern)
	})
}
