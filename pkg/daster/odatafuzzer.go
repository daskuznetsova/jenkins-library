package daster

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var ODataFuzzerType = "oDataFuzzer"

type ODataFuzzer struct {
	client     httpClient
	url        string
	verbose    bool
	maxRetries int
}

type NewODataFuzzerScanResponse struct {
	Url string `json:"url"`
}

type GetODataFuzzerResponse struct {
	Results     string      `json:"results"`
	RuntimeInfo interface{} `json:"runtimeInfo"`
	State       struct {
		Terminated bool `json:"terminated"`
	} `json:"state"`
}

func NewODataFuzzer(url string, verbose bool, maxRetires int) *ODataFuzzer {
	return &ODataFuzzer{
		url:        fmt.Sprintf("%s/%s", url, ODataFuzzerType),
		verbose:    verbose,
		maxRetries: maxRetires,
		client:     newHttpClient(),
	}
}

func (d *ODataFuzzer) TriggerScan(request map[string]interface{}) (string, error) {
	resp, err := callAPI(d.client, d.url, http.MethodPost, request, d.verbose, d.maxRetries)
	if err != nil {
		return "", err
	}

	var scan NewODataFuzzerScanResponse
	err = json.Unmarshal(resp, &scan)
	if err != nil {
		return "", err
	}
	return scan.Url, nil
}

func (d *ODataFuzzer) GetScan(scanId string) (*Scan, error) {
	resp, err := callAPI(d.client, fmt.Sprintf("%s/%s", d.url, scanId), http.MethodGet, nil, d.verbose, d.maxRetries)
	if err != nil {
		return nil, err
	}

	var scanResponse GetODataFuzzerResponse
	err = json.Unmarshal(resp, &scanResponse)
	if err != nil {
		return nil, err
	}

	return &Scan{
		Results: scanResponse.Results,
		State:   &ScanState{Terminated: scanResponse.State.Terminated},
		Summary: scanResponse.RuntimeInfo,
	}, nil
}

func (d *ODataFuzzer) DeleteScan(scanId string) error {
	_, err := callAPI(d.client, fmt.Sprintf("%s/%s", d.url, scanId), http.MethodDelete, nil, d.verbose, d.maxRetries)
	if err != nil {
		return err
	}
	return nil
}
