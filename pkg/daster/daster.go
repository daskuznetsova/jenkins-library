package daster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SAP/jenkins-library/pkg/log"
)

var RetryCodes = []int{100, 101, 102, 103, 404, 408, 425,
	/* not really common but a DASTer specific issue*/ 500, 503, 504}

type Daster interface {
	TriggerScan(request map[string]interface{}) (string, error)
	GetScan(scanId string) (*Scan, error)
	DeleteScan(scanId string) error
}

type Scan struct {
	Results string
	State   *ScanState
	Summary interface{}
}

type ScanState struct {
	Terminated bool
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func callAPI(url, mode string, requestBody interface{}, verbose bool, maxRetries int) ([]byte, error) {
	var requestBodyString []byte
	var err error
	if requestBody != nil {
		requestBodyString, err = json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		if verbose {
			log.Entry().Infof("request with body %s being sent.", string(requestBodyString))
		}
	}

	response := &http.Response{StatusCode: 0}
	attempts := 0

	for (response.StatusCode == 0 || IsInRetryCodes(response.StatusCode)) && attempts < maxRetries {
		response, err = SendHTTPRequest(url, mode, requestBodyString)
		if err != nil {
			return nil, err
		}
		attempts += 1
		time.Sleep(1 * time.Second)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		errResponse := ErrorResponse{}
		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("API request failed with status code %d: %s", response.StatusCode, errResponse.Message)
	}

	return body, nil
}

func SendHTTPRequest(url, mode string, requestBody []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(mode, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return client.Do(req)
}

func IsInRetryCodes(statusCode int) bool {
	for _, code := range RetryCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}
