//go:build unit
// +build unit

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDasterExecuteScanCommand(t *testing.T) {
	t.Parallel()

	testCmd := DasterExecuteScanCommand()

	// only high level testing performed - details are tested in step generation procedure
	assert.Equal(t, "dasterExecuteScan", testCmd.Use, "command name incorrect")

}
