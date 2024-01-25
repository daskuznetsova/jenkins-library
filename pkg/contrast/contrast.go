package contrast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/pkg/errors"
)

type Contrast interface {
	GetVulnerabilities() error
	GetAppInfo(appUIUrl, server string)
}

type ContrastInstance struct {
	url    string
	apiKey string
	auth   string
}

func NewContrastInstance(url, apiKey, auth string) ContrastInstance {
	return ContrastInstance{
		url:    url,
		apiKey: apiKey,
		auth:   auth,
	}
}

type Vulnerability struct {
	Category   string
	Id         string
	Severity   string
	Status     string
	Title      string
	RuleName   string
	Confidence string
	Impact     string
}

type Pageable struct {
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Paged      bool `json:"paged"`
	Unpaged    bool `json:"unpaged"`
	Offset     int  `json:"offset"`
}

type VulnsResponse struct {
	Pageable         Pageable `json:"pageable"`
	Size             int      `json:"size"`
	TotalElements    int      `json:"totalElements"`
	TotalPages       int      `json:"totalPages"`
	Empty            bool     `json:"empty"`
	First            bool     `json:"first"`
	Last             bool     `json:"last"`
	Number           int      `json:"number"`
	NumberOfElements int      `json:"numberOfElements"`
	Vulnerabilities  []Vuln   `json:"content"`
}

type Vuln struct {
	Severity string `json:"severity"`
	Status   string `json:"status"`
}

type VulnerabilitiesResponse struct {
	Success  bool                 `json:"success"`
	Messages []string             `json:"messages"`
	Traces   []VulnerabilityTrace `json:"traces"`
	Count    int                  `json:"count"`
	Links    []NextPageLink       `json:"links"`
}

type VulnerabilityTrace struct {
	Category   string `json:"category"`
	Confidence string `json:"confidence"`
	Impact     string `json:"impact"`
	RuleName   string `json:"rule_name"`
	Severity   string `json:"severity"`
	Status     string `json:"status"`
	Title      string `json:"title"`
	UUID       string `json:"uuid"`
}

type NextPageLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type ApplicationResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Path        string `json:"path"`
	Language    string `json:"language"`
	Importance  string `json:"importance"`
}

const (
	StatusFixed          = "FIXED"
	StatusNotAProblem    = "NOT_A_PROBLEM"
	StatusRemediated     = "REMEDIATED"
	StatusAutoRemediated = "AUTO_REMEDIATED"
	Critical             = "CRITICAL"
	High                 = "HIGH"
	Medium               = "MEDIUM"
	AuditAll             = "Audit All"
	Optional             = "Optional"
	pageSize             = 100
)

func (contrast *ContrastInstance) GetVulnerabilities() ([]ContrastFindings, error) {
	url := contrast.url + "/vulnerabilities"
	client := newContrastHTTPClient(contrast.apiKey, contrast.auth)

	return getVulnerabilitiesFromClient(client, url, 0)
}

func (contrast *ContrastInstance) GetAppInfo(appUIUrl, server string) (*ApplicationInfo, error) {
	client := newContrastHTTPClient(contrast.apiKey, contrast.auth)
	app, err := getApplicationFromClient(client, contrast.url)
	if err != nil {
		log.Entry().Errorf("failed to get application from client: %v", err)
		return nil, err
	}
	app.Url = appUIUrl
	app.Server = server
	return app, nil
}

func getApplicationFromClient(client contrastHTTPClient, url string) (*ApplicationInfo, error) {
	var appResponse ApplicationResponse

	response, err := client.doRequest(url, nil)
	if err != nil {
		return nil, err
	}
	defer response.Close()

	data, err := io.ReadAll(response)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &appResponse)
	if err != nil {
		return nil, err
	}

	return &ApplicationInfo{
		Id:   appResponse.Id,
		Name: appResponse.Name,
	}, nil
}

func getVulnerabilitiesFromClient(client contrastHTTPClient, url string, page int) ([]ContrastFindings, error) {
	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
		"size": fmt.Sprintf("%d", pageSize),
	}

	response, err := client.doRequest(url, params)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response)
	if err != nil {
		response.Close()
		return nil, err
	}

	var vulnsResponse VulnsResponse
	err = json.Unmarshal(data, &vulnsResponse)
	response.Close()
	if err != nil {
		return nil, err
	}
	if vulnsResponse.Empty {
		log.Entry().Debug("empty response")
		return nil, nil
	}
	auditAllFindings := ContrastFindings{
		ClassificationName: AuditAll,
		Total:              0,
		Audited:            0,
	}
	optionalFindings := ContrastFindings{
		ClassificationName: Optional,
		Total:              0,
		Audited:            0,
	}

	for _, vuln := range vulnsResponse.Vulnerabilities {
		if vuln.Severity == Critical || vuln.Severity == High || vuln.Severity == Medium {
			if vuln.Status == StatusFixed || vuln.Status == StatusNotAProblem ||
				vuln.Status == StatusRemediated || vuln.Status == StatusAutoRemediated {
				auditAllFindings.Audited += 1
			}
			auditAllFindings.Total += 1
		} else {
			if vuln.Status == StatusFixed || vuln.Status == StatusNotAProblem ||
				vuln.Status == StatusRemediated || vuln.Status == StatusAutoRemediated {
				optionalFindings.Audited += 1
			}
			optionalFindings.Total += 1
		}

	}
	if !vulnsResponse.Last {
		contrastFindings, err := getVulnerabilitiesFromClient(client, url, page+1)
		if err != nil {
			return nil, err
		}
		for i, fr := range contrastFindings {
			if fr.ClassificationName == AuditAll {
				contrastFindings[i].Total += auditAllFindings.Total
				contrastFindings[i].Audited += auditAllFindings.Audited
			}
			if fr.ClassificationName == Optional {
				contrastFindings[i].Total += optionalFindings.Total
				contrastFindings[i].Audited += optionalFindings.Audited
			}
		}
		return contrastFindings, nil
	}
	return []ContrastFindings{auditAllFindings, optionalFindings}, nil
}

type contrastHTTPClient interface {
	doRequest(url string, params map[string]string) (io.ReadCloser, error)
}

type contrastHTTPClientInstance struct {
	apiKey string
	auth   string
}

func newContrastHTTPClient(apiKey, auth string) *contrastHTTPClientInstance {
	return &contrastHTTPClientInstance{
		apiKey: apiKey,
		auth:   auth,
	}
}

func (c *contrastHTTPClientInstance) doRequest(url string, params map[string]string) (io.ReadCloser, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Add("API-Key", c.apiKey)
	req.Header.Add("Authorization", c.auth)

	q := req.URL.Query()
	for param, value := range params {
		q.Add(param, value)
	}
	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}
