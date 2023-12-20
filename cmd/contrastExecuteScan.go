package cmd

import (
	"encoding/base64"
	"fmt"

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

func runContrastExecuteScan(config *contrastExecuteScanOptions, telemetryData *telemetry.CustomData, utils contrastExecuteScanUtils) ([]piperutils.Path, error) {
	var reports []piperutils.Path

	contrastInstance := contrast.NewContrastInstance(getUrl(config), config.UserAPIKey, getAuth(config))
	findings, err := contrastInstance.GetVulnerabilities(config.ApplicationID)
	if err != nil {
		return reports, err
	}

	if config.CheckForCompliance {
		unaudited := findings.Total - findings.Audited
		if unaudited > config.VulnerabilityThresholdTotal {
			msg := fmt.Sprintf("Your application %v in organization %v is not compliant. Total unaudited issues are %v which is greater than the VulnerabilityThresholdTotal count %v",
				config.ApplicationID, config.OrganizationID, unaudited, config.VulnerabilityThresholdTotal)
			return reports, errors.Errorf(msg)
		}
	}

	//contrastAudit := contrast.ContrastAudit{
	//	ToolName:       "contrast",
	//	ApplicationURL: fmt.Sprintf("%s/Contrast/static/"),
	//}

	return reports, nil
}

func getUrl(config *contrastExecuteScanOptions) string {
	return fmt.Sprintf("https://%s/api/ng/organizations/%s/applications/%s/vulnerabilities",
		config.Server, config.OrganizationID, config.ApplicationID)
	//return fmt.Sprintf("https://%s/api/v4/organizations/%s/vulnerabilities",
	//	config.Server, config.OrganizationID)
}

func getAuth(config *contrastExecuteScanOptions) string {
	return base64.StdEncoding.EncodeToString([]byte(config.Username + ":" + config.ServiceKey))
}
