package contrast

import (
	"encoding/json"
	"io"

	piperhttp "github.com/SAP/jenkins-library/pkg/http"
	"github.com/SAP/jenkins-library/pkg/log"
)

type Contrast interface {
	GetVulnerabilities() error
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

func (contrast *ContrastInstance) GetVulnerabilities() error {
	client := piperhttp.Client{}
	header := make(map[string][]string)
	header["API-Key"] = []string{contrast.apiKey}
	header["Authorization"] = []string{contrast.auth}

	response, err := client.SendRequest("GET", contrast.url, nil, header, nil)
	if err != nil {
		return err
	}
	if response == nil {
		log.Entry().Warn("response is empty")
		return nil
	}
	defer response.Body.Close()

	bodyText, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	vulns := []Vulnerability{}
	err = json.Unmarshal(bodyText, &vulns)
	if err != nil {
		return err
	}

	for _, v := range vulns {
		log.Entry().Info(v)
	}

	return nil
}
