package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/SAP/jenkins-library/pkg/command"
	"github.com/SAP/jenkins-library/pkg/contrast"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperutils"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/pkg/errors"
)

type contrastExecuteScanUtils interface {
	command.ExecRunner
	piperutils.FileUtils
}

type contrastExecuteScanUtilsBundle struct {
	*command.Command
	*piperutils.Files
}

func newContrastExecuteScanUtils() contrastExecuteScanUtils {
	utils := contrastExecuteScanUtilsBundle{
		Command: &command.Command{},
		Files:   &piperutils.Files{},
	}
	utils.Stdout(log.Writer())
	utils.Stderr(log.Writer())
	return &utils
}

func contrastExecuteScan(config contrastExecuteScanOptions, telemetryData *telemetry.CustomData) {
	utils := newContrastExecuteScanUtils()

	reports, err := runContrastExecuteScan(&config, telemetryData, utils)
	piperutils.PersistReportsAndLinks("contrastExecuteScan", "./", utils, reports, nil)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}
}

func runContrastExecuteScan(config *contrastExecuteScanOptions, telemetryData *telemetry.CustomData, utils contrastExecuteScanUtils) (reports []piperutils.Path, err error) {
	auth, err := getAuth(config)
	if err != nil {
		return
	}
	appAPIUrl, appUIUrl, err := getApplicationUrls(config)
	if err != nil {
		return
	}

	contrastInstance := contrast.NewContrastInstance(appAPIUrl, config.UserAPIKey, auth)
	appInfo, err := contrastInstance.GetAppInfo(appUIUrl)
	if err != nil {
		return
	}

	findings, err := contrastInstance.GetVulnerabilities()
	if err != nil {
		return
	}
	log.Entry().Debugf("Contrast Findings:")
	for _, f := range findings {
		log.Entry().Debugf("Classification %s, total: %d, audited: %d", f.ClassificationName, f.Total, f.Audited)
	}

	contrastAudit := contrast.ContrastAudit{
		ToolName: "contrast",
		ApplicationUrl: fmt.Sprintf("%s/Contrast/static/ng/index.html#/%s/applications/%s",
			config.Server, config.OrganizationID, config.ApplicationID),
		ScanResults: findings,
	}
	paths, err := contrast.WriteJSONReport(contrastAudit, "./")
	if err != nil {
		return
	}
	reports = append(reports, paths...)

	if config.CheckForCompliance {
		for _, results := range findings {
			if results.ClassificationName == "Audit All" {
				unaudited := results.Total - results.Audited
				if unaudited > config.VulnerabilityThresholdTotal {
					msg := fmt.Sprintf("Your application %v in organization %v is not compliant. Total unaudited issues are %v which is greater than the VulnerabilityThresholdTotal count %v",
						config.ApplicationID, config.OrganizationID, unaudited, config.VulnerabilityThresholdTotal)
					return reports, errors.Errorf(msg)
				}
			}
		}
	}

	toolRecordFileName, err := contrast.CreateAndPersistToolRecord(utils, appInfo, "./")
	if err != nil {
		log.Entry().Warning("TR_CONTRAST: Failed to create toolrecord file ...", err)
	} else {
		reports = append(reports, piperutils.Path{Target: toolRecordFileName})
	}

	return
}

func getApplicationUrls(config *contrastExecuteScanOptions) (string, string, error) {
	if config.Server == "" {
		return "", "", errors.New("server is empty")
	}
	if config.OrganizationID == "" {
		return "", "", errors.New("organizationId is empty")
	}
	if config.ApplicationID == "" {
		return "", "", errors.New("applicationId is empty")
	}
	if !strings.HasPrefix(config.Server, "https://") {
		config.Server = "https://" + config.Server
	}

	return fmt.Sprintf("https://%s/api/v4/organizations/%s/applications/%s",
			config.Server, config.OrganizationID, config.ApplicationID),
		fmt.Sprintf("%s/Contrast/static/ng/index.html#/%s/applications/%s",
			config.Server, config.OrganizationID, config.ApplicationID), nil
}

func getAuth(config *contrastExecuteScanOptions) (string, error) {
	if config.UserAPIKey == "" {
		return "", errors.New("userApiKey is empty")
	}
	if config.Username == "" {
		return "", errors.New("username is empty")
	}
	if config.ServiceKey == "" {
		return "", errors.New("serviceKey is empty")
	}
	return base64.StdEncoding.EncodeToString([]byte(config.Username + ":" + config.ServiceKey)), nil
}
