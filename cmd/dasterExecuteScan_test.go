package cmd

import (
	"github.com/SAP/jenkins-library/pkg/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dasterExecuteScanMockUtils struct {
	*mock.ExecMockRunner
	*mock.FilesMock
}

func newDasterExecuteScanTestsUtils() dasterExecuteScanMockUtils {
	utils := dasterExecuteScanMockUtils{
		ExecMockRunner: &mock.ExecMockRunner{},
		FilesMock:      &mock.FilesMock{},
	}
	return utils
}

func TestRunDasterExecuteScan(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		// init
		config := dasterExecuteScanOptions{}

		utils := newDasterExecuteScanTestsUtils()
		utils.AddFile("file.txt", []byte("dummy content"))

		// test
		err := runDasterExecuteScan(&config, nil, utils)

		// assert
		assert.NoError(t, err)
	})

	t.Run("error path", func(t *testing.T) {
		t.Parallel()
		// init
		config := dasterExecuteScanOptions{}

		utils := newDasterExecuteScanTestsUtils()

		// test
		err := runDasterExecuteScan(&config, nil, utils)

		// assert
		assert.EqualError(t, err, "cannot run without important file")
	})
}
