package daster

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Daster struct {
	token string
	url   string
}

func NewDaster(token, url string) *Daster {
	return &Daster{
		token: token,
		url:   url,
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

/*
	def triggerScan() {
	        if (this.TRIGGER_ENDPOINTS.contains(this.config.scanType)) {
	            def body = transformConfiguration()
	            def parsedResponse = callApi("${this.config.scanType}", body)
	            return parsedResponse
	        }
	        return [:]
	    }
*/
func (d *Daster) TriggerScan() (*Scan, error) {
	return &Scan{}, nil
}

func (d *Daster) GetScanResponse(scanId string) (*ScanResponse, error) {
	return &ScanResponse{}, nil
}

func (d *Daster) GetScanResult(scan *ScanResponse) (*ScanResult, error) {
	return &ScanResult{}, nil
}

func (d *Daster) DeleteScan(scanId string) error {
	return nil
}

func CheckThresholdViolations(violations *ThresholdViolations, scanResult *ScanResult) *ThresholdViolations {
	return nil
}

/*
	private def transformConfiguration() {
	        def requestBody = [:].plus(config.settings)
	        return requestBody
	    }
*/
func transformConfiguration() {

}

/*
	private def callApi(endpoint, requestBody = null, mode = 'POST', contentType = 'APPLICATION_JSON', parseJsonResult = true){
	        def params = [
	            url                    : "${this.config.serviceUrl}${endpoint}",
	            httpMode               : mode,
	            acceptType             : 'APPLICATION_JSON',
	            contentType            : contentType,
	            quiet                  : !this.config.verbose,
	            consoleLogResponseBody : this.config.verbose,
	            validResponseCodes     : '100:499'
	        ]
	        if (requestBody) {
	            def requestBodyString = utils.jsonToString(requestBody)
	            if (this.config.verbose) this.script.echo "Request with body ${requestBodyString} being sent."
	            params.put('requestBody', requestBodyString)
	        }
	        def response = [status: 0]
	        def attempts = 0
	        while ((!response.status || RETRY_CODES.contains(response.status)) && attempts < this.config.maxRetries) {
	            response = httpResource(params)
	            attempts++
	        }
	        if (parseJsonResult)
	            return this.utils.parseJsonSerializable(response.content)
	        else
	            return response.content
	    }
*/
func callApi(url string, requestBody []byte, verbose bool) {
	var jsonStr = []byte("{}")
	if requestBody != nil {
		requestBodyString := utils(jsonStr)
		if verbose {
			fmt.Printf("Request with body %s being sent.\n", requestBodyString)
		}
		jsonStr = []byte(requestBodyString)
	}
	response, err := httpResource(url, mode, jsonStr)
	if err != nil {
		return nil, err
	}
	if parseJsonResult {
		parsedResponse, err := parser(string(response))
		if err != nil {
			return nil, err
		}
		return parsedResponse, nil
	} else {
		var responseContent map[string]interface{}
		json.Unmarshal(response, &responseContent)
		return responseContent, nil
	}
}

func httpResource(url string, mode string, jsonStr []byte) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(mode, url, strings.NewReader(string(jsonStr)))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBytes, _ := io.ReadAll(resp.Body)

	return responseBytes, nil
}
