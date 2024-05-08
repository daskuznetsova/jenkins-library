package daster

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHttpClient struct {
	Resp http.Response
	Err  error
}

func (m *MockHttpClient) sendHttpRequest(url, mode string, requestBody []byte) (*http.Response, error) {
	return &m.Resp, m.Err
}

func TestCallAPI(t *testing.T) {
	t.Parallel()
	t.Run("Test successful API call", func(t *testing.T) {
		mockClient := &MockHttpClient{
			Resp: http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			},
		}

		resp, err := callAPI(mockClient, "https://example.com", "POST", []byte{'0'}, false, 3)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
	})

	t.Run("Test API call with unsuccessful status code", func(t *testing.T) {
		mockClient := &MockHttpClient{
			Resp: http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"Internal Server Error"}`)),
			},
		}

		_, err := callAPI(mockClient, "https://example.com", "GET", nil, false, 3)

		assert.Error(t, err)
		assert.Equal(t, "API request failed with status code 500: Internal Server Error", err.Error())
	})

	t.Run("Test API call with error", func(t *testing.T) {
		mockClient := &MockHttpClient{
			Err: fmt.Errorf("request failed"),
		}

		_, err := callAPI(mockClient, "https://example.com", "DELETE", nil, false, 3)

		assert.Error(t, err)
		assert.Equal(t, "request failed", err.Error())
	})
}
