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
		flags := map[string]bool{"--flag1": true, "--flag2": true, "--flag3": true}
		expected := []string{"--flag1=1", "--flag2=2", "--flag3=3"}
		result, err := AppendCustomFlags(input, flags)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("No valid flags", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "1",
			"--flag2": "2",
			"--flag3": "3",
		}
		flags := map[string]bool{}
		expected := []string{}
		result, err := AppendCustomFlags(input, flags)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Some flags are valid", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "1",
			"--flag2": "2",
			"--flag3": "3",
		}
		flags := map[string]bool{"--flag1": true, "--flag3": true}
		expected := []string{"--flag1=1", "--flag3=3"}
		result, err := AppendCustomFlags(input, flags)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Flags without values", func(t *testing.T) {
		input := map[string]string{
			"--flag1": "",
			"--flag2": "",
			"--flag3": "",
		}
		flags := map[string]bool{"--flag1": true, "--flag2": true, "--flag3": true}
		expected := []string{"--flag1", "--flag2", "--flag3"}
		result, err := AppendCustomFlags(input, flags)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Empty input", func(t *testing.T) {
		input := map[string]string{}
		flags := map[string]bool{"--flag1": true, "--flag2": true, "--flag3": true}
		expected := []string{}
		result, err := AppendCustomFlags(input, flags)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
