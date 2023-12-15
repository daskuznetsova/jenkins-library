package contrast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SAP/jenkins-library/pkg/log"
)

type Contrast interface {
	GetVulnerabilities(applicationId string) error
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

type ContrastResponse struct {
	Content       []Vulnerability `json:"content"`
	Pageable      Pageable        `json:"pageable"`
	Last          bool            `json:"last"`
	TotalPages    int             `json:"totalPages"`
	TotalElements int             `json:"totalElements"`
	Empty         bool            `json:"empty"`
}

func (contrast *ContrastInstance) GetVulnerabilities(applicationId string) (*ContrastFindings, error) {

	pageSize := 100
	pageNumber := 0
	audited := 0
	totalAlerts := 0

	var vulnerabilities []Vulnerability
	for {
		client := http.Client{}
		req, err := http.NewRequest("GET", contrast.url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Add("API-Key", contrast.apiKey)
		req.Header.Add("Authorization", contrast.auth)
		q := req.URL.Query()
		q.Add("page", fmt.Sprintf("%d", pageNumber))
		q.Add("size", fmt.Sprintf("%d", pageSize))
		req.URL.RawQuery = q.Encode()

		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if response == nil {
			log.Entry().Warn("response is empty")
			break
		}
		defer response.Body.Close()

		bodyText, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		resp := ContrastResponse{}
		err = json.Unmarshal(bodyText, &resp)
		if err != nil {
			return nil, err
		}

		log.Entry().Infof("page %d from %d", resp.Pageable.PageNumber+1, resp.TotalPages)

		for _, vuln := range resp.Content {
			if vuln.ApplicationId == applicationId {
				vulnerabilities = append(vulnerabilities, vuln)
				if vuln.Status == "Fixed" || vuln.Status == "Not a problem" {
					audited += 1
				}
				totalAlerts += 1
			}
		}
		if resp.Last {
			break
		}
		pageNumber++
	}

	auditAll := &ContrastFindings{
		ClassificationName: "Audit All",
		Total:              totalAlerts,
		Audited:            audited,
	}

	return auditAll, nil
}
