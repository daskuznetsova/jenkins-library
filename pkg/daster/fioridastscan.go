package daster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SAP/jenkins-library/pkg/log"
)

var FioriDASTScanType = "fioriDASTScan"

type FioriDASTScan struct {
	client     httpClient
	url        string
	verbose    bool
	maxRetries int
}

type NewFioriDASTScanResponse struct {
	ScanId string `json:"scanId"`
}

type GetFioriDASTScanResponse struct {
	Results     string      `json:"results"`
	RiskSummary interface{} `json:"riskSummary"`
	State       struct {
		Terminated bool `json:"terminated"`
	} `json:"state"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
}

func NewFioriDASTScan(url string, verbose bool, maxRetires int) *FioriDASTScan {
	return &FioriDASTScan{
		verbose:    verbose,
		maxRetries: maxRetires,
		url:        fmt.Sprintf("%s/%s", strings.TrimSuffix(url, "/"), FioriDASTScanType),
		client:     newHttpClient(),
	}
}

func (d *FioriDASTScan) TriggerScan(request map[string]interface{}) (string, error) {
	log.Entry().Debugf("TriggerScan start")
	resp, err := callAPI(d.client, d.url, http.MethodPost, request, d.verbose, d.maxRetries)
	if err != nil {
		return "", err
	}

	var scan NewFioriDASTScanResponse
	err = json.Unmarshal(resp, &scan)
	if err != nil {
		return "", err
	}
	return scan.ScanId, nil
}

func (d *FioriDASTScan) GetScan(scanId string) (*Scan, error) {
	resp, err := callAPI(d.client, fmt.Sprintf("%s/%s", d.url, scanId), http.MethodGet, nil, d.verbose, d.maxRetries)
	if err != nil {
		return nil, err
	}

	var scanResponse GetFioriDASTScanResponse
	err = json.Unmarshal(resp, &scanResponse)
	if err != nil {
		return nil, err
	}

	return &Scan{
		Results: scanResponse.Results,
		State:   &ScanState{Terminated: scanResponse.State.Terminated},
		Summary: scanResponse.RiskSummary,
	}, nil
}

func (d *FioriDASTScan) DeleteScan(scanId string) error {
	_, err := callAPI(d.client, fmt.Sprintf("%s/%s", d.url, scanId), http.MethodDelete, nil, d.verbose, d.maxRetries)
	if err != nil {
		return err
	}
	return nil
}
