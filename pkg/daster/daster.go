package daster

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/SAP/jenkins-library/pkg/log"
)

type Daster struct {
	token    string
	url      string
	scanType string
	verbose  bool
}

func NewDaster(token, url, scanType string, verbose bool) *Daster {
	return &Daster{
		token:    token,
		url:      url,
		scanType: scanType,
		verbose:  verbose,
	}
}

type Scan struct {
	ScanId string
}

type ScanResponse struct {
	State *State
}

type State struct {
	Terminated *TerminatedState
}

type TerminatedState struct {
	ExitCode    int
	Reason      string
	ContainerId string
}

type ScanResult struct {
}

type ThresholdViolations struct {
	High   int
	Medium int
	Low    int
	Info   int
	All    int
}

func (d *Daster) TriggerScan(settings map[string]interface{}) (*Scan, error) {
	requestBody, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}
	resp, err := callApi(d.url+"/"+d.scanType, requestBody, http.MethodPost, d.verbose)
	if err != nil {
		return nil, err
	}

	var scan *Scan
	err = json.Unmarshal(resp, scan)
	return scan, err
}

func (d *Daster) GetScanResponse(scanId string) (*ScanResponse, error) {
	switch d.scanType {
	case "fioriDASTScan", "aemscan", "oDataFuzzer", "burpscan":
		resp, err := callApi(d.url+"/"+d.scanType+"/"+scanId, nil, http.MethodGet, d.verbose)
		if err != nil {
			return nil, err
		}
		var scanResponse *ScanResponse
		err = json.Unmarshal(resp, scanResponse)
		return scanResponse, err
	}
	return &ScanResponse{}, nil
}

/*
def result = [:]

	switch (this.config.scanType) {
	    case 'fioriDASTScan':
	        result.summary = scanResponse?.riskSummary
	        result. details = scanResponse?.riskReport
	        break
	    case  'aemscan':
	        result.details = scanResponse?.log
	        break
	}
	return result
*/
func (d *Daster) GetScanResult(scan *ScanResponse) (*ScanResult, error) {
	switch d.scanType {
	case "fioriDASTScan":

	}
	return &ScanResult{}, nil
}

func (d *Daster) DeleteScan(scanId string) error {
	return nil
}

func CheckThresholdViolations(violations *ThresholdViolations, scanResult *ScanResult) *ThresholdViolations {
	return nil
}

func callApi(url string, requestBody []byte, mode string, verbose bool) ([]byte, error) {
	var jsonStr = []byte("{}")
	if requestBody != nil {
		if verbose {
			log.Entry().Infof("request with body %s being sent.", requestBody)
		}
		jsonStr = requestBody
	}
	response, err := httpResource(url, mode, jsonStr)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func httpResource(url string, mode string, jsonStr []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(mode, url, strings.NewReader(string(jsonStr)))
	if err != nil {
		return nil, err
	}
	resp, err := performRequest(client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return readResponseBody(resp)
}

func performRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBytes, nil
}
