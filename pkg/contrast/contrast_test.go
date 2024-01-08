package contrast

import (
	"fmt"
	"io"
	"strings"
)

var (
	applicationRequest     = "application"
	vulnerabilitiesRequest = "vulnerabilities"
)

type contrastClientMock struct {
}

func (c *contrastClientMock) doRequest(url string, params map[string]string) (io.ReadCloser, error) {
	if url == applicationRequest {
		appInfo := `{"id":"7cda8021-f371-42f0-b0e8-bd569afe1021","name":"owasp-benchmark","displayName":"","path":"/","language":"JAVA","importance":"MEDIUM","isArchived":false,"technologies":[],"tags":["DEMO-APPLICATION"],"metadata":{},"firstSeenTime":"2023-04-03T23:04:27Z","lastSeenTime":"2023-04-21T18:37:00Z"}`
		return io.NopCloser(strings.NewReader(appInfo)), nil
	}
	if url == vulnerabilitiesRequest {
		appInfo := `{"id":"7cda8021-f371-42f0-b0e8-bd569afe1021","name":"owasp-benchmark","displayName":"","path":"/","language":"JAVA","importance":"MEDIUM","isArchived":false,"technologies":[],"tags":["DEMO-APPLICATION"],"metadata":{},"firstSeenTime":"2023-04-03T23:04:27Z","lastSeenTime":"2023-04-21T18:37:00Z"}`
		return io.NopCloser(strings.NewReader(appInfo)), nil
	}
	return nil, fmt.Errorf("error")
}
