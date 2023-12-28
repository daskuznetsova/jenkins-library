package contrast

import (
	"testing"

	"github.com/SAP/jenkins-library/pkg/mock"
	"github.com/stretchr/testify/assert"
)

type contrastExecuteScanMockUtils struct {
	*mock.ExecMockRunner
	*mock.FilesMock
}

func newContrastExecuteScanTestsUtils() contrastExecuteScanMockUtils {
	return contrastExecuteScanMockUtils{
		ExecMockRunner: &mock.ExecMockRunner{},
		FilesMock:      &mock.FilesMock{},
	}
}

func TestCreateToolRecordContrast(t *testing.T) {
	modulePath := "./"

	t.Run("Valid toolrun file", func(t *testing.T) {
		appInfo := &ApplicationInfo{
			ServerUrl:      "https://contrastsecurity.com",
			OrganizationId: "organization-id",
			Id:             "application-id",
			Name:           "app name",
			DisplayName:    "application name",
			Path:           "/",
		}
		toolRecord, err := createToolRecordContrast(newContrastExecuteScanTestsUtils(), appInfo, modulePath)
		assert.NoError(t, err)
		assert.Equal(t, toolRecord.ToolName, "contrast")
		assert.Equal(t, toolRecord.ToolInstance, appInfo.ServerUrl)
		assert.Equal(t, toolRecord.DisplayName, appInfo.DisplayName)
		assert.Equal(t, toolRecord.DisplayURL, appInfo.ApplicationUrl)
	})

	t.Run("Empty server", func(t *testing.T) {
		appInfo := &ApplicationInfo{
			OrganizationId: "organization-id",
			Id:             "application-id",
			Name:           "app name",
			DisplayName:    "application name",
			Path:           "/",
			ApplicationUrl: "",
		}
		_, err := createToolRecordContrast(newContrastExecuteScanTestsUtils(), appInfo, modulePath)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Contrast server is not set")
	})

	t.Run("Empty organization id", func(t *testing.T) {
		appInfo := &ApplicationInfo{
			ServerUrl:      "https://contrastsecurity.com",
			Id:             "application-id",
			Name:           "app name",
			DisplayName:    "application name",
			Path:           "/",
			ApplicationUrl: "",
		}
		_, err := createToolRecordContrast(newContrastExecuteScanTestsUtils(), appInfo, modulePath)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Organization Id is not set")
	})

	t.Run("Empty application id", func(t *testing.T) {
		appInfo := &ApplicationInfo{
			ServerUrl:      "https://contrastsecurity.com",
			OrganizationId: "organization-id",
			Name:           "app name",
			DisplayName:    "application name",
			Path:           "/",
			ApplicationUrl: "",
		}
		_, err := createToolRecordContrast(newContrastExecuteScanTestsUtils(), appInfo, modulePath)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Application Id is not set")
	})
}
