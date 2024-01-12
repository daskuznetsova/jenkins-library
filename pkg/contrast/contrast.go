package contrast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Contrast interface {
	GetVulnerabilities(applicationId string) error
	GetApplication(applicationId string) (ApplicationInfo, error)
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
	StatusFixed       = "Fixed"
	StatusNotAProblem = "Not a problem"
)

func (contrast *ContrastInstance) GetVulnerabilities(organizationId, applicationId string) ([]ContrastFindings, error) {
	pageSize := 100
	pageNumber := 0

	url := fmt.Sprintf("https://cs003.contrastsecurity.com/Contrast/api/ng/%s/orgtraces/filter", organizationId)
	client := newContrastHTTPClient(contrast.apiKey, contrast.auth)

	params := map[string]string{
		"expand":  "application",
		"modules": applicationId,
		"offset":  fmt.Sprintf("%d", pageSize*pageNumber),
		"limit":   fmt.Sprintf("%d", pageSize),
	}

	audited, total, _, err := getVulnerabilitiesFromClient(client, url, params)
	if err != nil {
		return nil, err
	}

	auditAll := ContrastFindings{
		ClassificationName: "Audit All",
		Total:              total,
		Audited:            audited,
	}

	return []ContrastFindings{auditAll}, nil
}

func (contrast *ContrastInstance) GetApplication(server, organization, applicationId string) (*ApplicationInfo, error) {
	var appResponse ApplicationResponse

	url := fmt.Sprintf("%s/applications/%s", contrast.url, applicationId)

	client := newContrastHTTPClient(contrast.apiKey, contrast.auth)
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
		ServerUrl:      server,
		OrganizationId: organization,
		Id:             applicationId,
		Name:           appResponse.Name,
		DisplayName:    appResponse.DisplayName,
		Path:           appResponse.Path,
		ApplicationUrl: fmt.Sprintf("%s/Contrast/static/ng/index.html#/%s/applications/%s",
			server, organization, applicationId),
	}, nil
}

func getVulnerabilitiesFromClient(client *contrastHTTPClientInstance, url string, params map[string]string) (int, int, []*Vulnerability, error) {
	var auditedAll, totalAll int
	var vulnerabilities []*Vulnerability

	response, err := client.doRequest(url, params)
	if err != nil {
		return auditedAll, totalAll, vulnerabilities, err
	}

	data, err := io.ReadAll(response)
	if err != nil {
		response.Close()
		return auditedAll, totalAll, vulnerabilities, err
	}

	var vulnsResponse VulnerabilitiesResponse
	err = json.Unmarshal(data, &vulnsResponse)
	response.Close()
	if err != nil {
		return auditedAll, totalAll, vulnerabilities, err
	}

	for _, vuln := range vulnsResponse.Traces {
		vulnerabilities = append(vulnerabilities, &Vulnerability{
			Category:   vuln.Category,
			Confidence: vuln.Confidence,
			Id:         vuln.UUID,
			Impact:     vuln.Impact,
			Severity:   vuln.Severity,
			Status:     vuln.Status,
			Title:      vuln.Title,
			RuleName:   vuln.RuleName,
		})
		if vuln.Status == StatusFixed || vuln.Status == StatusNotAProblem {
			auditedAll += 1
		}
		totalAll += 1
	}
	for _, link := range vulnsResponse.Links {
		if link.Rel == "nextPage" {
			audited, total, vulns, err := getVulnerabilitiesFromClient(client, link.Href, nil)
			if err != nil {
				return auditedAll, totalAll, vulnerabilities, err
			}
			auditedAll += audited
			totalAll += total
			vulnerabilities = append(vulnerabilities, vulns...)
		}
	}

	return auditedAll, totalAll, vulnerabilities, err
}

type contrastHTTPClient interface {
	doRequest(url, apiKey, auth string, params map[string]string) (io.ReadCloser, error)
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
