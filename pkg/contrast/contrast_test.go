package contrast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type contrastHttpClientMock struct{}

func (c *contrastHttpClientMock) ExecuteRequest(url string, params map[string]string, dest interface{}) error {
	switch url {
	case appUrl:
		app, ok := dest.(*ApplicationResponse)
		if !ok {
			return fmt.Errorf("wrong destination type")
		}
		app.Id = "1"
		app.Name = "application"
	default:
		return fmt.Errorf("error")
	}
	return nil
}

const (
	appUrl   = "https://server.com/applications"
	vulnsUrl = "https://server.com/vulnerabilities"
)

func TestGetApplicationFromClient(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		contrastClient := &contrastHttpClientMock{}
		app, err := getApplicationFromClient(contrastClient, appUrl)
		assert.NoError(t, err)
		assert.NotEmpty(t, app)
		assert.Equal(t, "1", app.Id)
		assert.Equal(t, "application", app.Name)
		assert.Equal(t, "", app.Url)
		assert.Equal(t, "", app.Server)
	})

	t.Run("Fail", func(t *testing.T) {
		contrastClient := &contrastHttpClientMock{}
		_, err := getApplicationFromClient(contrastClient, "https://server.com/applications/fail")
		assert.Error(t, err)
	})
}
