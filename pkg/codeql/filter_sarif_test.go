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

func TestMatchComponent(t *testing.T) {
	t.Parallel()

	t.Run("Path matches pattern", func(t *testing.T) {
		filePath := "path/to/src/file.txt"
		pattern := "**/src/*"
		assert.True(t, matchComponent(pattern, filePath))
	})

	t.Run("Path matches exact pattern", func(t *testing.T) {
		filePath := "file.txt"
		pattern := "file.txt"
		assert.True(t, matchComponent(pattern, filePath))
	})

	t.Run("Path with escape symbols matches pattern", func(t *testing.T) {
		filePath := "/file\\ name.txt"
		pattern := "**"
		assert.True(t, matchComponent(pattern, filePath))
	})

	t.Run("Path doesn't match pattern", func(t *testing.T) {
		filePath := "path/to/file.txt"
		pattern := "**/src/*"
		assert.False(t, matchComponent(pattern, filePath))
	})

	t.Run("Path doesn't match exact pattern", func(t *testing.T) {
		filePath := "path/to/file.txt"
		pattern := "file"
		assert.False(t, matchComponent(pattern, filePath))
	})
}
