package codeql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendCustomFlags(t *testing.T) {
	t.Parallel()

	t.Run("All flags are valid", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "1",
			"--flag2": "2",
			"--flag3": "3",
		}
		expected := []string{"--flag1=1", "--flag2=2", "--flag3=3"}
		result, err := GetCustomFlags(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("No valid flags", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "1",
			"--flag2": "2",
			"--flag3": "3",
		}
		expected := []string{}
		result, err := GetCustomFlags(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Some flags are valid", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "1",
			"--flag2": "2",
			"--flag3": "3",
		}
		expected := []string{"--flag1=1", "--flag3=3"}
		result, err := GetCustomFlags(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Flags without values", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "",
			"--flag2": "",
			"--flag3": "",
		}
		expected := []string{"--flag1", "--flag2", "--flag3"}
		result, err := GetCustomFlags(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Empty input", func(t *testing.T) {
		input := map[string]string{}
		expected := []string{}
		result, err := GetCustomFlags(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestCheckIfFlagSetByUser(t *testing.T) {
	t.Parallel()

	customFlags := map[string]string{
		"--flag1": "--flag1=1",
		"-f2":     "-f2=2",
		"--flag3": "--flag3",
	}

	t.Run("Flag is not set by user", func(t *testing.T) {
		input := []string{"-f4"}
		assert.False(t, CheckIfFlagSetByUser(customFlags, input))
	})
	t.Run("Flag is set by user", func(t *testing.T) {
		input := []string{"-f2"}
		assert.True(t, CheckIfFlagSetByUser(customFlags, input))
	})
	t.Run("One of flags is set by user", func(t *testing.T) {
		input := []string{"--flag2", "-f2"}
		assert.True(t, CheckIfFlagSetByUser(customFlags, input))
	})
}

func TestGetFlags(t *testing.T) {
	t.Parallel()

	t.Run("Valid flags with values", func(t *testing.T) {
		inputStr := "--flag1=1 --flag2=2 --flag3=string"
		expected := map[string]bool{
			"--flag1=1":      true,
			"--flag2=2":      true,
			"--flag3=string": true,
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for _, f := range result {
			assert.True(t, expected[f])
		}
	})

	t.Run("Valid flags without values", func(t *testing.T) {
		inputStr := "--flag1 -flag2 -f3"
		expected := map[string]bool{
			"--flag1": true,
			"-flag2":  true,
			"-f3":     true,
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for _, f := range result {
			assert.True(t, expected[f])
		}
	})

	t.Run("Valid flags with spaces in value", func(t *testing.T) {
		inputStr := "--flag1='mvn install' --flag2=\"mvn clean install\" -f3='mvn clean install -DskipTests=true'"
		expected := map[string]bool{
			"--flag1=mvn install":                    true,
			"--flag2=mvn clean install":              true,
			"-f3=mvn clean install -DskipTests=true": true,
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for _, f := range result {
			assert.True(t, expected[f])
		}
	})
}

func TestRemoveDuplicateFlags(t *testing.T) {
	t.Parallel()

	longShortFlags := map[string]string{
		"--flag1": "-f1",
		"--flag2": "-f2",
		"--flag3": "-f3",
	}

	t.Run("No duplications", func(t *testing.T) {
		flags := map[string]string{
			"--flag1": "--flag1=1",
			"-f2":     "-f2=2",
			"--flag3": "--flag3",
		}
		expected := map[string]string{
			"--flag1": "--flag1=1",
			"-f2":     "-f2=2",
			"--flag3": "--flag3",
		}
		removeDuplicateFlags(flags, longShortFlags)
		assert.Equal(t, len(expected), len(flags))
		for k, v := range flags {
			assert.Equal(t, expected[k], v)
		}
	})

	t.Run("Duplications", func(t *testing.T) {
		flags := map[string]string{
			"--flag1": "--flag1=1",
			"-f1":     "-f1=2",
			"--flag2": "--flag2=1",
			"-f2":     "-f2=2",
			"--flag3": "--flag3",
			"-f3":     "-f3",
		}
		expected := map[string]string{
			"--flag1": "--flag1=1",
			"--flag2": "--flag2=1",
			"--flag3": "--flag3",
		}
		removeDuplicateFlags(flags, longShortFlags)
		assert.Equal(t, len(expected), len(flags))
		for k, v := range flags {
			assert.Equal(t, expected[k], v)
		}
	})
}
