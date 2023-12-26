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
	ApplicationId string `json:"applicationId"`
	Category      string `json:"category"`
	Id            string `json:"id"`
	Severity      string `json:"severity"`
	Status        string `json:"status"`
	Substatus     string `json:"substatus"`
	Title         string `json:"title"`
	RuleName      string `json:"ruleName"`
}

type Pageable struct {
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Paged      bool `json:"paged"`
	Unpaged    bool `json:"unpaged"`
	Offset     int  `json:"offset"`
}

type VulnerabilitiesResponse struct {
	Content       []Vulnerability `json:"content"`
	Pageable      Pageable        `json:"pageable"`
	Last          bool            `json:"last"`
	TotalPages    int             `json:"totalPages"`
	TotalElements int             `json:"totalElements"`
	Empty         bool            `json:"empty"`
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

func (contrast *ContrastInstance) GetVulnerabilities(applicationId string) ([]ContrastFindings, error) {
	pageSize := 100
	pageNumber := 0
	audited := 0
	totalAlerts := 0

	var vulnerabilities []Vulnerability
	for {
		params := map[string]string{
			"page": fmt.Sprintf("%d", pageNumber),
			"size": fmt.Sprintf("%d", pageSize),
		}
		response, err := doRequest(contrast.url+"/vulnerabilities", contrast.apiKey, contrast.auth, params)
		if err != nil {
			return nil, err
		}
		defer response.Close()

		data, err := io.ReadAll(response)
		if err != nil {
			return nil, err
		}

		var vulnsResponse VulnerabilitiesResponse
		err = json.Unmarshal(data, &vulnsResponse)
		if err != nil {
			return nil, err
		}

		for _, vuln := range vulnsResponse.Content {
			if vuln.ApplicationId == applicationId {
				vulnerabilities = append(vulnerabilities, vuln)
				if vuln.Status == StatusFixed || vuln.Status == StatusNotAProblem {
					audited += 1
				}
				totalAlerts += 1
			}
		}
		if vulnsResponse.Last {
			break
		}
		pageNumber++
	}

	auditAll := ContrastFindings{
		ClassificationName: "Audit All",
		Total:              totalAlerts,
		Audited:            audited,
	}

	return []ContrastFindings{auditAll}, nil
}

func (contrast *ContrastInstance) GetApplication(server, organization, applicationId string) (*ApplicationInfo, error) {
	var appResponse ApplicationResponse

	url := fmt.Sprintf("%s/applications/%s", contrast.url, applicationId)

	response, err := doRequest(url, contrast.apiKey, contrast.auth, nil)
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

func doRequest(url, apiKey, auth string, params map[string]string) (io.ReadCloser, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Add("API-Key", apiKey)
	req.Header.Add("Authorization", auth)

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
