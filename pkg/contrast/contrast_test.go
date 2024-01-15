package contrast

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type contrastClientMock struct{}

func (c *contrastClientMock) doRequest(url string, params map[string]string) (io.ReadCloser, error) {
	if url == "application" {
		appInfo := `{"id":"7cda8021-f371-42f0-b0e8-bd569afe1021","name":"owasp-benchmark","displayName":"","path":"/","language":"JAVA","importance":"MEDIUM","isArchived":false,"technologies":[],"tags":["DEMO-APPLICATION"],"metadata":{},"firstSeenTime":"2023-04-03T23:04:27Z","lastSeenTime":"2023-04-21T18:37:00Z"}`
		return io.NopCloser(strings.NewReader(appInfo)), nil
	}
	if url == "vulnerabilities.1" {
		appInfo := `{"success" : true, "messages" : [ "Organization Vulnerabilities loaded successfully" ], "traces" : [ {"app_version_tags" : [ ], "application" : {"primary" : false, "master" : false, "child" : false, "roles" : [ "ROLE_EDIT", "ROLE_RULES_ADMIN", "ROLE_ADMIN", "ROLE_ORG_AUDITOR", "ROLE_VIEW" ], "importance" : 2, "app_id" : "7cda8021-f371-42f0-b0e8-bd569afe1021", "name" : "owasp-benchmark", "parent_app_id" : null, "total_modules" : 1, "language" : "Java", "context_path" : "/", "last_seen" : 1682102220000, "license_level" : "Licensed", "importance_description" : "MEDIUM" }, "bugtracker_tickets" : [ ], "category" : "Cryptography", "category_label" : "Cryptography", "closed_time" : null, "confidence" : "Low", "confidence_label" : "Low", "default_severity" : "MEDIUM", "default_severity_label" : "Medium", "discovered" : 1680563040000, "evidence" : null, "first_time_seen" : 1680563040000, "hasParentApp" : false, "impact" : "Medium", "impact_label" : "Medium", "instance_uuid" : "E6W7-1BFS-6OHN-DSRP", "language" : "Java", "last_time_seen" : 1680869820000, "last_vuln_time_seen" : 1680869820000, "license" : "Licensed", "likelihood" : "Medium", "likelihood_label" : "Medium", "organization_name" : "SAP", "reported_to_bug_tracker" : false, "reported_to_bug_tracker_time" : null, "rule_name" : "crypto-bad-mac", "rule_title" : "Insecure Hash Algorithms", "severity" : "Medium", "severity_label" : "Medium", "status" : "Reported", "sub_status" : "", "sub_title" : "'SHA' hash algorithm used at KeyStoreUtil.java", "substatus_keycode" : null, "tags" : [ ], "title" : "'SHA' hash algorithm used at KeyStoreUtil.java", "total_traces_received" : 13, "uuid" : "E6W7-1BFS-6OHN-DSRP", "violations" : [ ], "visible" : true }, {"app_version_tags" : [ ], "application" : {"primary" : false, "master" : false, "child" : false, "roles" : [ "ROLE_EDIT", "ROLE_RULES_ADMIN", "ROLE_ADMIN", "ROLE_ORG_AUDITOR", "ROLE_VIEW" ], "importance" : 2, "app_id" : "7cda8021-f371-42f0-b0e8-bd569afe1021", "name" : "owasp-benchmark", "parent_app_id" : null, "total_modules" : 1, "language" : "Java", "context_path" : "/", "last_seen" : 1682102220000, "license_level" : "Licensed", "importance_description" : "MEDIUM" }, "bugtracker_tickets" : [ ], "category" : "Injection", "category_label" : "Injection", "closed_time" : null, "confidence" : "High", "confidence_label" : "High", "default_severity" : "HIGH", "default_severity_label" : "High", "discovered" : 1680563400000, "evidence" : null, "first_time_seen" : 1680563400000, "hasParentApp" : false, "impact" : "High", "impact_label" : "High", "instance_uuid" : "9I9R-AZEI-3OES-GLD0", "language" : "Java", "last_time_seen" : 1680789660000, "last_vuln_time_seen" : 1680789660000, "license" : "Licensed", "likelihood" : "Medium", "likelihood_label" : "Medium", "organization_name" : "SAP", "reported_to_bug_tracker" : false, "reported_to_bug_tracker_time" : null, "rule_name" : "cmd-injection", "rule_title" : "OS Command Injection", "severity" : "High", "severity_label" : "High", "status" : "Reported", "sub_status" : "", "sub_title" : "OS Command Injection from \"BenchmarkTest02411\" Parameter on \"/benchmark/cmdi-02/BenchmarkTest02411\" page", "substatus_keycode" : null, "tags" : [ ], "title" : "OS Command Injection from \"BenchmarkTest02411\" Parameter on \"/benchmark/cmdi-02/BenchmarkTest02411\" page", "total_traces_received" : 7, "uuid" : "9I9R-AZEI-3OES-GLD0", "violations" : [ ], "visible" : true }], "count" : 2, "licensedCount" : 2, "links" : [] }`
		return io.NopCloser(strings.NewReader(appInfo)), nil
	}
	if url == "vulnerabilities.2" {
		appInfo := `{"success" : true, "messages" : [ "Organization Vulnerabilities loaded successfully" ], "traces" : [ {"app_version_tags" : [ ], "application" : {"primary" : false, "master" : false, "child" : false, "roles" : [ "ROLE_EDIT", "ROLE_RULES_ADMIN", "ROLE_ADMIN", "ROLE_ORG_AUDITOR", "ROLE_VIEW" ], "importance" : 2, "app_id" : "7cda8021-f371-42f0-b0e8-bd569afe1021", "name" : "owasp-benchmark", "parent_app_id" : null, "total_modules" : 1, "language" : "Java", "context_path" : "/", "last_seen" : 1682102220000, "license_level" : "Licensed", "importance_description" : "MEDIUM" }, "bugtracker_tickets" : [ ], "category" : "Cryptography", "category_label" : "Cryptography", "closed_time" : null, "confidence" : "Low", "confidence_label" : "Low", "default_severity" : "MEDIUM", "default_severity_label" : "Medium", "discovered" : 1680563040000, "evidence" : null, "first_time_seen" : 1680563040000, "hasParentApp" : false, "impact" : "Medium", "impact_label" : "Medium", "instance_uuid" : "E6W7-1BFS-6OHN-DSRP", "language" : "Java", "last_time_seen" : 1680869820000, "last_vuln_time_seen" : 1680869820000, "license" : "Licensed", "likelihood" : "Medium", "likelihood_label" : "Medium", "organization_name" : "SAP", "reported_to_bug_tracker" : false, "reported_to_bug_tracker_time" : null, "rule_name" : "crypto-bad-mac", "rule_title" : "Insecure Hash Algorithms", "severity" : "Medium", "severity_label" : "Medium", "status" : "Reported", "sub_status" : "", "sub_title" : "'SHA' hash algorithm used at KeyStoreUtil.java", "substatus_keycode" : null, "tags" : [ ], "title" : "'SHA' hash algorithm used at KeyStoreUtil.java", "total_traces_received" : 13, "uuid" : "E6W7-1BFS-6OHN-DSRP", "violations" : [ ], "visible" : true }, {"app_version_tags" : [ ], "application" : {"primary" : false, "master" : false, "child" : false, "roles" : [ "ROLE_EDIT", "ROLE_RULES_ADMIN", "ROLE_ADMIN", "ROLE_ORG_AUDITOR", "ROLE_VIEW" ], "importance" : 2, "app_id" : "7cda8021-f371-42f0-b0e8-bd569afe1021", "name" : "owasp-benchmark", "parent_app_id" : null, "total_modules" : 1, "language" : "Java", "context_path" : "/", "last_seen" : 1682102220000, "license_level" : "Licensed", "importance_description" : "MEDIUM" }, "bugtracker_tickets" : [ ], "category" : "Cryptography", "category_label" : "Cryptography", "closed_time" : null, "confidence" : "Low", "confidence_label" : "Low", "default_severity" : "MEDIUM", "default_severity_label" : "Medium", "discovered" : 1680563040000, "evidence" : null, "first_time_seen" : 1680563040000, "hasParentApp" : false, "impact" : "Medium", "impact_label" : "Medium", "instance_uuid" : "TCH6-HW7N-5K2U-OMX2", "language" : "Java", "last_time_seen" : 1680869820000, "last_vuln_time_seen" : 1680869820000, "license" : "Licensed", "likelihood" : "Medium", "likelihood_label" : "Medium", "organization_name" : "SAP", "reported_to_bug_tracker" : false, "reported_to_bug_tracker_time" : null, "rule_name" : "crypto-bad-mac", "rule_title" : "Insecure Hash Algorithms", "severity" : "Medium", "severity_label" : "Medium", "status" : "Reported", "sub_status" : "", "sub_title" : "'MD5' hash algorithm used at GranteeManager", "substatus_keycode" : null, "tags" : [ ], "title" : "'MD5' hash algorithm used at GranteeManager", "total_traces_received" : 10, "uuid" : "TCH6-HW7N-5K2U-OMX2", "violations" : [ ], "visible" : true }, {"app_version_tags" : [ ], "application" : {"primary" : false, "master" : false, "child" : false, "roles" : [ "ROLE_EDIT", "ROLE_RULES_ADMIN", "ROLE_ADMIN", "ROLE_ORG_AUDITOR", "ROLE_VIEW" ], "importance" : 2, "app_id" : "7cda8021-f371-42f0-b0e8-bd569afe1021", "name" : "owasp-benchmark", "parent_app_id" : null, "total_modules" : 1, "language" : "Java", "context_path" : "/", "last_seen" : 1682102220000, "license_level" : "Licensed", "importance_description" : "MEDIUM" }, "bugtracker_tickets" : [ ], "category" : "Caching", "category_label" : "Caching", "closed_time" : null, "confidence" : "Low", "confidence_label" : "Low", "default_severity" : "NOTE", "default_severity_label" : "Note", "discovered" : 1680563040000, "evidence" : null, "first_time_seen" : 1680563040000, "hasParentApp" : false, "impact" : "Low", "impact_label" : "Low", "instance_uuid" : "T0KK-370G-8PFL-7FPC", "language" : "Java", "last_time_seen" : 1680869820000, "last_vuln_time_seen" : 1680869820000, "license" : "Licensed", "likelihood" : "Low", "likelihood_label" : "Low", "organization_name" : "SAP", "reported_to_bug_tracker" : false, "reported_to_bug_tracker_time" : null, "rule_name" : "cache-controls-missing", "rule_title" : "Anti-Caching Controls Missing", "severity" : "Note", "severity_label" : "Note", "status" : "Reported", "sub_status" : "", "sub_title" : "Anti-Caching Controls Missing detected", "substatus_keycode" : null, "tags" : [ ], "title" : "Anti-Caching Controls Missing detected", "total_traces_received" : 13, "uuid" : "T0KK-370G-8PFL-7FPC", "violations" : [ ], "visible" : true }, {"app_version_tags" : [ ], "application" : {"primary" : false, "master" : false, "child" : false, "roles" : [ "ROLE_EDIT", "ROLE_RULES_ADMIN", "ROLE_ADMIN", "ROLE_ORG_AUDITOR", "ROLE_VIEW" ], "importance" : 2, "app_id" : "7cda8021-f371-42f0-b0e8-bd569afe1021", "name" : "owasp-benchmark", "parent_app_id" : null, "total_modules" : 1, "language" : "Java", "context_path" : "/", "last_seen" : 1682102220000, "license_level" : "Licensed", "importance_description" : "MEDIUM" }, "bugtracker_tickets" : [ ], "category" : "Injection", "category_label" : "Injection", "closed_time" : null, "confidence" : "High", "confidence_label" : "High", "default_severity" : "HIGH", "default_severity_label" : "High", "discovered" : 1680563400000, "evidence" : null, "first_time_seen" : 1680563400000, "hasParentApp" : false, "impact" : "High", "impact_label" : "High", "instance_uuid" : "9I9R-AZEI-3OES-GLD0", "language" : "Java", "last_time_seen" : 1680789660000, "last_vuln_time_seen" : 1680789660000, "license" : "Licensed", "likelihood" : "Medium", "likelihood_label" : "Medium", "organization_name" : "SAP", "reported_to_bug_tracker" : false, "reported_to_bug_tracker_time" : null, "rule_name" : "cmd-injection", "rule_title" : "OS Command Injection", "severity" : "High", "severity_label" : "High", "status" : "Reported", "sub_status" : "", "sub_title" : "OS Command Injection from \"BenchmarkTest02411\" Parameter on \"/benchmark/cmdi-02/BenchmarkTest02411\" page", "substatus_keycode" : null, "tags" : [ ], "title" : "OS Command Injection from \"BenchmarkTest02411\" Parameter on \"/benchmark/cmdi-02/BenchmarkTest02411\" page", "total_traces_received" : 7, "uuid" : "9I9R-AZEI-3OES-GLD0", "violations" : [ ], "visible" : true }], "count" : 5, "licensedCount" : 5, "links" : [ {"rel" : "nextPage", "href" : "vulnerabilities.3", "hreflang" : null, "media" : null, "title" : null, "type" : null, "deprecation" : null, "method" : "GET" } ] }`
		return io.NopCloser(strings.NewReader(appInfo)), nil
	}
	if url == "vulnerabilities.3" {
		appInfo := `{"success" : true, "messages" : [ "Organization Vulnerabilities loaded successfully" ], "traces" : [ {"app_version_tags" : [ ], "application" : {"primary" : false, "master" : false, "child" : false, "roles" : [ "ROLE_EDIT", "ROLE_RULES_ADMIN", "ROLE_ADMIN", "ROLE_ORG_AUDITOR", "ROLE_VIEW" ], "importance" : 2, "app_id" : "7cda8021-f371-42f0-b0e8-bd569afe1021", "name" : "owasp-benchmark", "parent_app_id" : null, "total_modules" : 1, "language" : "Java", "context_path" : "/", "last_seen" : 1682102220000, "license_level" : "Licensed", "importance_description" : "MEDIUM" }, "bugtracker_tickets" : [ ], "category" : "Cryptography", "category_label" : "Cryptography", "closed_time" : null, "confidence" : "Low", "confidence_label" : "Low", "default_severity" : "MEDIUM", "default_severity_label" : "Medium", "discovered" : 1680563040000, "evidence" : null, "first_time_seen" : 1680563040000, "hasParentApp" : false, "impact" : "Medium", "impact_label" : "Medium", "instance_uuid" : "E6W7-1BFS-6OHN-DSRP", "language" : "Java", "last_time_seen" : 1680869820000, "last_vuln_time_seen" : 1680869820000, "license" : "Licensed", "likelihood" : "Medium", "likelihood_label" : "Medium", "organization_name" : "SAP", "reported_to_bug_tracker" : false, "reported_to_bug_tracker_time" : null, "rule_name" : "crypto-bad-mac", "rule_title" : "Insecure Hash Algorithms", "severity" : "Medium", "severity_label" : "Medium", "status" : "Reported", "sub_status" : "", "sub_title" : "'SHA' hash algorithm used at KeyStoreUtil.java", "substatus_keycode" : null, "tags" : [ ], "title" : "'SHA' hash algorithm used at KeyStoreUtil.java", "total_traces_received" : 13, "uuid" : "E6W7-1BFS-6OHN-DSRP", "violations" : [ ], "visible" : true }], "count" : 5, "licensedCount" : 5, "links" : [] }`
		return io.NopCloser(strings.NewReader(appInfo)), nil
	}
	return nil, fmt.Errorf("error")
}

type contrastClientErrorMock struct{}

func (c *contrastClientErrorMock) doRequest(url string, params map[string]string) (io.ReadCloser, error) {
	return nil, errors.New("error")
}

func TestGetVulnerabilitiesFromClient(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		contrastClient := &contrastClientMock{}
		findings, err := getVulnerabilitiesFromClient(contrastClient, "vulnerabilities.1", map[string]string{})
		assert.NoError(t, err)
		assert.NotEmpty(t, findings)
		assert.Equal(t, 1, len(findings))
		assert.Equal(t, 2, findings[0].Total)
		assert.Equal(t, 0, findings[0].Audited)
	})

	t.Run("Success with pagination results", func(t *testing.T) {
		contrastClient := &contrastClientMock{}
		findings, err := getVulnerabilitiesFromClient(contrastClient, "vulnerabilities.2", map[string]string{})
		assert.NoError(t, err)
		assert.NotEmpty(t, findings)
		assert.Equal(t, 1, len(findings))
		assert.Equal(t, 5, findings[0].Total)
		assert.Equal(t, 0, findings[0].Audited)
	})

	t.Run("Error", func(t *testing.T) {
		contrastClient := &contrastClientMock{}
		params := make(map[string]string)
		var url string
		_, err := getVulnerabilitiesFromClient(contrastClient, url, params)
		assert.Error(t, err)
	})
}
