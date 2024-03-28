package codeql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePattern(t *testing.T) {
	t.Parallel()

	t.Run("Include files, no rules", func(t *testing.T) {
		input := "+**/src/**/*"
		pattern, err := parsePattern(input)
		assert.NoError(t, err)
		assert.True(t, pattern.sign)
		assert.Equal(t, "**/src/**/*", pattern.filePattern)
		assert.Equal(t, "**", pattern.rulePattern)
	})

	t.Run("Exclude files, no rules", func(t *testing.T) {
		input := "-**/src/**/*"
		pattern, err := parsePattern(input)
		assert.NoError(t, err)
		assert.False(t, pattern.sign)
		assert.Equal(t, "**/src/**/*", pattern.filePattern)
		assert.Equal(t, "**", pattern.rulePattern)
	})

	t.Run("Include files with rule", func(t *testing.T) {
		input := "+**/src/**/*:security-rule"
		pattern, err := parsePattern(input)
		assert.NoError(t, err)
		assert.True(t, pattern.sign)
		assert.Equal(t, "**/src/**/*", pattern.filePattern)
		assert.Equal(t, "security-rule", pattern.rulePattern)
	})

	t.Run("Exclude files with rule", func(t *testing.T) {
		input := "-**/src/**/*:security-rule"
		pattern, err := parsePattern(input)
		assert.NoError(t, err)
		assert.False(t, pattern.sign)
		assert.Equal(t, "**/src/**/*", pattern.filePattern)
		assert.Equal(t, "security-rule", pattern.rulePattern)
	})
}
