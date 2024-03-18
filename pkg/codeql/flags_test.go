package codeql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFlags(t *testing.T) {
	t.Parallel()

	t.Run("All flags are valid", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "1",
			"--flag2": "2",
			"--flag3": "3",
		}
		expected := []string{"--flag1=1", "--flag2=2", "--flag3=3"}
		result, err := AppendCustomFlags(input)
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
		result, err := AppendCustomFlags(input)
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
		result, err := AppendCustomFlags(input)
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
		result, err := AppendCustomFlags(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Empty input", func(t *testing.T) {
		input := map[string]string{}
		expected := []string{}
		result, err := AppendCustomFlags(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestParseCustomFlags(t *testing.T) {
	t.Parallel()

	t.Run("Valid flags with values", func(t *testing.T) {
		inputStr := "--flag1=1 --flag2=2 --flag3=string"
		expected := map[string]string{
			"--flag1": "--flag1=1",
			"--flag2": "--flag2=2",
			"--flag3": "--flag3=string",
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for k, v := range result {
			assert.Equal(t, expected[k], v)
		}
	})

	t.Run(".", func(t *testing.T) {
		inputStr := "--no-db-cluster -l=java --threads=1 --command='mvn clean package -Dmaven.test.skip=true'"
		expected := map[string]string{
			"--no-db-cluster": "--no-db-cluster",
			"-l":              "-l=java",
			"--threads":       "--threads=1",
			"--command":       "--command='mvn clean package -Dmaven.test.skip=true'",
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for k, v := range result {
			assert.Equal(t, expected[k], v)
		}
	})

	t.Run("Valid flags without values", func(t *testing.T) {
		inputStr := "--flag1 -flag2 -f3"
		expected := map[string]string{
			"--flag1": "--flag1",
			"-flag2":  "-flag2",
			"-f3":     "-f3",
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for k, v := range result {
			assert.Equal(t, expected[k], v)
		}
	})

	t.Run("Duplications with short flags", func(t *testing.T) {
		inputStr := "--language=java -l=python -s=. --ram=2000"
		expected := map[string]string{
			"--language": "--language=java",
			"-s":         "-s=.",
			"--ram":      "--ram=2000",
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for k, v := range result {
			assert.Equal(t, expected[k], v)
		}
	})

	t.Run("Valid flags with spaces in value", func(t *testing.T) {
		inputStr := "--flag1='mvn install' --flag2='mvn clean install'"
		expected := map[string]string{
			"--flag1": "--flag1='mvn install'",
			"--flag2": "--flag2='mvn clean install'",
		}
		result := ParseCustomFlags(inputStr)
		assert.Equal(t, len(expected), len(result))
		for k, v := range result {
			assert.Equal(t, expected[k], v)
		}
	})
}
