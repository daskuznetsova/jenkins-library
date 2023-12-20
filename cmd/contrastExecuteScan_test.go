package cmd

import (
	"github.com/SAP/jenkins-library/pkg/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type contrastExecuteScanMockUtils struct {
	*mock.ExecMockRunner
	*mock.FilesMock
}

func newContrastExecuteScanTestsUtils() contrastExecuteScanMockUtils {
	utils := contrastExecuteScanMockUtils{
		ExecMockRunner: &mock.ExecMockRunner{},
		FilesMock:      &mock.FilesMock{},
	}
	return utils
}

func TestRunContrastExecuteScan(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		// init
		config := contrastExecuteScanOptions{}

		utils := newContrastExecuteScanTestsUtils()
		utils.AddFile("file.txt", []byte("dummy content"))

		// test
		err := runContrastExecuteScan(&config, nil, utils)

		// assert
		assert.NoError(t, err)
	})

	t.Run("error path", func(t *testing.T) {
		t.Parallel()
		// init
		config := contrastExecuteScanOptions{}

		utils := newContrastExecuteScanTestsUtils()

		// test
		err := runContrastExecuteScan(&config, nil, utils)

		// assert
		assert.EqualError(t, err, "cannot run without important file")
	})
}
