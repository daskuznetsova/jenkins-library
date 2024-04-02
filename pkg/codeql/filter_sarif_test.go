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

	t.Run("Components match, glob", func(t *testing.T) {
		filePath := "file.txt"
		pattern := "**"
		assert.True(t, matchComponent(pattern, filePath))
	})

	t.Run("Components match, glob", func(t *testing.T) {
		filePath := "file.txt"
		pattern := "file.txt"
		assert.True(t, matchComponent(pattern, filePath))
	})

	t.Run("Components don't match", func(t *testing.T) {
		filePath := "file"
		pattern := "file.txt"
		assert.False(t, matchComponent(pattern, filePath))
	})

	t.Run("Component with escape symbol", func(t *testing.T) {
		filePath := "/file\\ name.txt"
		pattern := "**"
		assert.True(t, matchComponent(pattern, filePath))
	})
}

func TestMatchComponents(t *testing.T) {
	t.Parallel()

	t.Run("Components match", func(t *testing.T) {
		file := []string{
			"path",
			"to",
			"src",
			"file.txt",
		}
		pattern := []string{"**", "src", "*"}
		assert.True(t, matchComponents(pattern, file))
	})
	t.Run("Components don't match", func(t *testing.T) {
		file := []string{
			"path",
			"to",
			"file.txt",
		}
		pattern := []string{"**", "src", "*"}
		assert.False(t, matchComponents(pattern, file))
	})
}

func TestMatch(t *testing.T) {
	t.Parallel()

	t.Run("Match", func(t *testing.T) {
		fileName := "path/to/src/file"
		pattern := "**/src/*"
		matches, err := match(pattern, fileName)
		assert.NoError(t, err)
		assert.True(t, matches)
	})

	t.Run("Match", func(t *testing.T) {
		fileName := "path/to/src/folder/some/files"
		pattern := "**/src/*"
		matches, err := match(pattern, fileName)
		assert.NoError(t, err)
		assert.True(t, matches)
	})

	t.Run("Don't match", func(t *testing.T) {
		fileName := "path/to/file"
		pattern := "**/src/*"
		matches, err := match(pattern, fileName)
		assert.NoError(t, err)
		assert.False(t, matches)
	})
}
