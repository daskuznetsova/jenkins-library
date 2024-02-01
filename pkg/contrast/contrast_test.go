package contrast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type contrastHttpClientMock struct {
	page *int
}

func (c *contrastHttpClientMock) ExecuteRequest(url string, params map[string]string, dest interface{}) error {
	switch url {
	case appUrl:
		app, ok := dest.(*ApplicationResponse)
		if !ok {
			return fmt.Errorf("wrong destination type")
		}
		app.Id = "1"
		app.Name = "application"
	case vulnsUrl:
		vulns, ok := dest.(*VulnerabilitiesResponse)
		if !ok {
			return fmt.Errorf("wrong destination type")
		}
		vulns.Size = 6
		vulns.TotalElements = 6
		vulns.TotalPages = 1
		vulns.Empty = false
		vulns.First = true
		vulns.Last = true
		vulns.Vulnerabilities = []Vulnerability{
			{Severity: "HIGH", Status: "FIXED"},
			{Severity: "MEDIUM", Status: "REMEDIATED"},
			{Severity: "HIGH", Status: "REPORTED"},
			{Severity: "MEDIUM", Status: "REPORTED"},
			{Severity: "HIGH", Status: "CONFIRMED"},
			{Severity: "NOTE", Status: "SUSPICIOUS"},
		}
	case vulnsUrlPaginated:
		vulns, ok := dest.(*VulnerabilitiesResponse)
		if !ok {
			return fmt.Errorf("wrong destination type")
		}
		vulns.Size = 100
		vulns.TotalElements = 300
		vulns.TotalPages = 3
		vulns.Empty = false
		vulns.Last = false
		if *c.page == 3 {
			vulns.Last = true
			return nil
		}
		for i := 0; i < 20; i++ {
			vulns.Vulnerabilities = append(vulns.Vulnerabilities, Vulnerability{Severity: "HIGH", Status: "FIXED"})
			vulns.Vulnerabilities = append(vulns.Vulnerabilities, Vulnerability{Severity: "NOTE", Status: "FIXED"})
			vulns.Vulnerabilities = append(vulns.Vulnerabilities, Vulnerability{Severity: "MEDIUM", Status: "REPORTED"})
			vulns.Vulnerabilities = append(vulns.Vulnerabilities, Vulnerability{Severity: "LOW", Status: "REPORTED"})
			vulns.Vulnerabilities = append(vulns.Vulnerabilities, Vulnerability{Severity: "CRITICAL", Status: "NOT_A_PROBLEM"})
		}
		*c.page++
	default:
		return fmt.Errorf("error")
	}
	return nil
}

const (
	appUrl            = "https://server.com/applications"
	errorUrl          = "https://server.com/error"
	vulnsUrl          = "https://server.com/vulnerabilities"
	vulnsUrlPaginated = "https://server.com/vulnerabilities/pagination"
)

func TestGetApplicationFromClient(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		contrastClient := &contrastHttpClientMock{}
		app, err := getApplicationFromClient(contrastClient, appUrl)
		assert.NoError(t, err)
		assert.NotEmpty(t, app)
		assert.Equal(t, "1", app.Id)
		assert.Equal(t, "application", app.Name)
		assert.Equal(t, "", app.Url)
		assert.Equal(t, "", app.Server)
	})

	t.Run("Error", func(t *testing.T) {
		contrastClient := &contrastHttpClientMock{}
		_, err := getApplicationFromClient(contrastClient, errorUrl)
		assert.Error(t, err)
	})
}

func TestGetVulnerabilitiesFromClient(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		contrastClient := &contrastHttpClientMock{}
		findings, err := getVulnerabilitiesFromClient(contrastClient, vulnsUrl, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, findings)
		assert.Equal(t, 2, len(findings))
		for _, f := range findings {
			assert.True(t, f.ClassificationName == AuditAll || f.ClassificationName == Optional)
			if f.ClassificationName == AuditAll {
				assert.Equal(t, 5, f.Total)
				assert.Equal(t, 2, f.Audited)
			}
			if f.ClassificationName == Optional {
				assert.Equal(t, 1, f.Total)
				assert.Equal(t, 0, f.Audited)
			}
		}
	})

	t.Run("Success with pagination results", func(t *testing.T) {
		page := 0
		contrastClient := &contrastHttpClientMock{page: &page}
		findings, err := getVulnerabilitiesFromClient(contrastClient, vulnsUrlPaginated, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, findings)
		assert.Equal(t, 2, len(findings))
		for _, f := range findings {
			assert.True(t, f.ClassificationName == AuditAll || f.ClassificationName == Optional)
			if f.ClassificationName == AuditAll {
				assert.Equal(t, 180, f.Total)
				assert.Equal(t, 120, f.Audited)
			}
			if f.ClassificationName == Optional {
				assert.Equal(t, 120, f.Total)
				assert.Equal(t, 60, f.Audited)
			}
		}
	})

	t.Run("Error", func(t *testing.T) {
		contrastClient := &contrastHttpClientMock{}
		_, err := getVulnerabilitiesFromClient(contrastClient, errorUrl, 0)
		assert.Error(t, err)
	})
}

func TestGetFindings(t *testing.T) {
	t.Parallel()
	t.Run("Audit All", func(t *testing.T) {
		vulns := []Vulnerability{
			{
				Severity: "CRITICAL",
				Status:   "REPORTED",
			},
		}
		auditAll, optional := getFindings(vulns)
		assert.Equal(t, 1, auditAll.Total)
		assert.Equal(t, 0, auditAll.Audited)
		assert.Equal(t, 0, optional.Total)
		assert.Equal(t, 0, optional.Audited)
	})
}

func TestIsVulnerabilityResolved(t *testing.T) {
	t.Parallel()
	t.Run("Vulnerability is resolved", func(t *testing.T) {
		assert.True(t, isVulnerabilityResolved("FIXED"))
		assert.True(t, isVulnerabilityResolved("REMEDIATED"))
		assert.True(t, isVulnerabilityResolved("NOT_A_PROBLEM"))
		assert.True(t, isVulnerabilityResolved("AUTO_REMEDIATED"))
	})
	t.Run("Vulnerability isn't resolved", func(t *testing.T) {
		assert.False(t, isVulnerabilityResolved("REPORTED"))
		assert.False(t, isVulnerabilityResolved("SUSPICIOUS"))
		assert.False(t, isVulnerabilityResolved("CONFIRMED"))
	})
}

func TestAccumulateFindings(t *testing.T) {

}
