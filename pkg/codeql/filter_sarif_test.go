package codeql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePatterns(t *testing.T) {
	t.Parallel()

	t.Run("Empty input", func(t *testing.T) {
		input := ""
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 1, len(patterns))
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.True(t, patterns[0].sign)
		assert.Equal(t, "", patterns[0].filePattern)
	})

	t.Run("One pattern to exclude", func(t *testing.T) {
		input := "-**/src/**"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 1, len(patterns))
		assert.Equal(t, "**/src/**", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.False(t, patterns[0].sign)
	})

	t.Run("One pattern to include", func(t *testing.T) {
		input := "+**/src/**"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 1, len(patterns))
		assert.Equal(t, "**/src/**", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.True(t, patterns[0].sign)
	})

	t.Run("Several patterns to exclude", func(t *testing.T) {
		input := "-**/src/exclude1/* -**/src/exclude2/*"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "**/src/exclude1/*", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.False(t, patterns[0].sign)
		assert.Equal(t, "**/src/exclude2/*", patterns[1].filePattern)
		assert.Equal(t, "**", patterns[1].rulePattern)
		assert.False(t, patterns[1].sign)
	})

	t.Run("Several patterns to include", func(t *testing.T) {
		input := "+**/src/include1/* +**/src/include2/*"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "**/src/include1/*", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.True(t, patterns[0].sign)
		assert.Equal(t, "**/src/include2/*", patterns[1].filePattern)
		assert.Equal(t, "**", patterns[1].rulePattern)
		assert.True(t, patterns[1].sign)
	})

	t.Run("One pattern to exclude, one pattern to include", func(t *testing.T) {
		input := "-**/src/exclude/* +**/src/include/*"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "**/src/exclude/*", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.False(t, patterns[0].sign)
		assert.Equal(t, "**/src/include/*", patterns[1].filePattern)
		assert.Equal(t, "**", patterns[1].rulePattern)
		assert.True(t, patterns[1].sign)
	})

	t.Run("Several patterns to exclude and include", func(t *testing.T) {
		input := "-**/src/exclude1/* +**/src/include1/* -**/src/exclude2/* +**/src/include2/*"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 4, len(patterns))
		assert.Equal(t, "**/src/exclude1/*", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.False(t, patterns[0].sign)
		assert.Equal(t, "**/src/include1/*", patterns[1].filePattern)
		assert.Equal(t, "**", patterns[1].rulePattern)
		assert.True(t, patterns[1].sign)
		assert.Equal(t, "**/src/exclude2/*", patterns[2].filePattern)
		assert.Equal(t, "**", patterns[2].rulePattern)
		assert.False(t, patterns[2].sign)
		assert.Equal(t, "**/src/include2/*", patterns[3].filePattern)
		assert.Equal(t, "**", patterns[3].rulePattern)
		assert.True(t, patterns[3].sign)
	})

	t.Run("Pattern with spaces", func(t *testing.T) {
		input := "-**/src/exclude\\ 1/*"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 1, len(patterns))
		assert.Equal(t, "**/src/exclude\\ 1/*", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.False(t, patterns[0].sign)
	})

	t.Run("Patterns with spaces", func(t *testing.T) {
		input := "-**/src/exclude\\ 1/* -**/src/exclude\\ 2/*"
		patterns, err := ParsePatterns(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "**/src/exclude\\ 1/*", patterns[0].filePattern)
		assert.Equal(t, "**", patterns[0].rulePattern)
		assert.False(t, patterns[0].sign)
		assert.Equal(t, "**/src/exclude\\ 2/*", patterns[1].filePattern)
		assert.Equal(t, "**", patterns[1].rulePattern)
		assert.False(t, patterns[1].sign)
	})
}

func TestSplit(t *testing.T) {
	t.Run("Patterns with spaces", func(t *testing.T) {
		input := "-**/src/exclude\\ 1/* -**/src/exclude\\ 2/*"
		patterns := split(input)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "-**/src/exclude\\ 1/*", patterns[0])
		assert.Equal(t, "-**/src/exclude\\ 2/*", patterns[1])
	})
	t.Run("First pattern with space", func(t *testing.T) {
		input := "-**/src/exclude\\ 1/* -**/src/exclude/*"
		patterns := split(input)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "-**/src/exclude\\ 1/*", patterns[0])
		assert.Equal(t, "-**/src/exclude/*", patterns[1])
	})
	t.Run("Second pattern with space", func(t *testing.T) {
		input := "-**/src/exclude1/* -**/src/exclude\\ 2/*"
		patterns := split(input)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "-**/src/exclude1/*", patterns[0])
		assert.Equal(t, "-**/src/exclude\\ 2/*", patterns[1])
	})
	t.Run("Patterns without spaces", func(t *testing.T) {
		input := "-**/src/exclude1/* -**/src/exclude2/*"
		patterns := split(input)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "-**/src/exclude1/*", patterns[0])
		assert.Equal(t, "-**/src/exclude2/*", patterns[1])
	})
	t.Run("Second pattern with escape symbol", func(t *testing.T) {
		input := "-**/src/exclude1/* -**/src/exclude2/*\\"
		patterns := split(input)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 2, len(patterns))
		assert.Equal(t, "-**/src/exclude1/*", patterns[0])
		assert.Equal(t, "-**/src/exclude2/*\\", patterns[1])
	})
	t.Run("First pattern with escape symbol", func(t *testing.T) {
		input := "-**/src/exclude1/*\\ -**/src/exclude2/*"
		patterns := split(input)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 1, len(patterns))
		assert.Equal(t, "-**/src/exclude1/*\\ -**/src/exclude2/*", patterns[0])
	})
	t.Run("Pattern with several spaces", func(t *testing.T) {
		input := "-**/src\\ 1/exclude\\ 1/*"
		patterns := split(input)
		assert.NotEmpty(t, patterns)
		assert.Equal(t, 1, len(patterns))
		assert.Equal(t, "-**/src\\ 1/exclude\\ 1/*", patterns[0])
	})
}

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

	t.Run("File path matches pattern", func(t *testing.T) {
		fileName := "path/to/src/file"
		pattern := "**/src/*"
		matches, err := match(pattern, fileName)
		assert.NoError(t, err)
		assert.True(t, matches)
	})

	t.Run("'*' matches only within a single component", func(t *testing.T) {
		fileName := "path/to/src/folder/some/files"
		pattern := "**/src/*"
		matches, err := match(pattern, fileName)
		assert.NoError(t, err)
		assert.False(t, matches)
	})

	t.Run("'**' matches zero or more components in the complete file name", func(t *testing.T) {
		fileName := "path/to/src/folder/some/files"
		pattern := "**/src/**"
		matches, err := match(pattern, fileName)
		assert.NoError(t, err)
		assert.True(t, matches)
	})

	t.Run("Path doesn't match pattern", func(t *testing.T) {
		fileName := "path/to/file"
		pattern := "**/src/*"
		matches, err := match(pattern, fileName)
		assert.NoError(t, err)
		assert.False(t, matches)
	})
}
