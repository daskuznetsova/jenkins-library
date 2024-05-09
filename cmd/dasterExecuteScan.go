package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/SAP/jenkins-library/pkg/command"
	"github.com/SAP/jenkins-library/pkg/daster"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperutils"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

type dasterExecuteScanUtils interface {
	command.ExecRunner

	piperutils.FileUtils
}

type dasterExecuteScanUtilsBundle struct {
	*command.Command
	*piperutils.Files
}

func newDasterExecuteScanUtils() dasterExecuteScanUtils {
	utils := dasterExecuteScanUtilsBundle{
		Command: &command.Command{},
		Files:   &piperutils.Files{},
	}
	// Reroute command output to logging framework
	utils.Stdout(log.Writer())
	utils.Stderr(log.Writer())
	return &utils
}

func dasterExecuteScan(config dasterExecuteScanOptions, telemetryData *telemetry.CustomData) {
	utils := newDasterExecuteScanUtils()

	err := runDasterExecuteScan(&config, telemetryData, utils)
	if err != nil {
		log.Entry().WithError(err).Fatal("daster execution failed")
	}
}

func runDasterExecuteScan(config *dasterExecuteScanOptions, telemetryData *telemetry.CustomData, utils dasterExecuteScanUtils) error {
	var dasterInstance daster.Daster
	switch config.ScanType {
	case daster.FioriDASTScanType:
		dasterInstance = daster.NewFioriDASTScan(config.ServiceURL, config.Verbose, config.MaxRetries)
	case daster.ODataFuzzerType:
		dasterInstance = daster.NewODataFuzzer(config.ServiceURL, config.Verbose, config.MaxRetries)
	default:
		log.Entry().Errorf("scan type %s is currently unavailable", config.ScanType)
	}
	if config.Settings == nil {
		config.Settings = map[string]interface{}{}
	}

	if config.OAuthServiceURL != "" && config.ClientID != "" && config.ClientSecret != "" {
		token, err := fetchOAuthToken(config)
		if err != nil {
			return err
		}
		config.Settings["parameterRules"] = token
	}

	if config.TargetURL != "" {
		config.Settings["targetURL"] = config.TargetURL
	}
	if config.DasterToken != "" {
		config.Settings["dasterToken"] = config.DasterToken
	}
	if config.UserCredentials != "" {
		config.Settings["userCredentials"] = config.UserCredentials
	}
	scanId, err := dasterInstance.TriggerScan(config.Settings)
	if err != nil {
		log.Entry().WithError(err).Error("failed to trigger scan")
		return err
	}
	if scanId == "" {
		return nil
	}

	if !config.Synchronous {
		return nil
	}

	var scan *daster.Scan
	for {
		scan, err = dasterInstance.GetScan(scanId)
		if err != nil {
			log.Entry().WithError(err).Error("failed to get scan")
			return err
		}
		if scan.State.Terminated {
			break
		}
		time.Sleep(15 * time.Second)
	}

	// TODO: CheckThresholdViolations

	if config.DeleteScan {
		err = dasterInstance.DeleteScan(scanId)
		if err != nil {
			log.Entry().WithError(err).Warn("failed to delete scan")
		}
	}

	return nil
}

func fetchOAuthToken(config *dasterExecuteScanOptions) (string, error) {
	resp, err := http.DefaultClient.PostForm(config.OAuthServiceURL,
		url.Values{
			"grant_type":    []string{config.OAuthGrantType},
			"scope":         []string{config.OAuthSource},
			"client_id":     []string{config.ClientID},
			"client_secret": []string{config.ClientSecret},
		})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("unable to fetch access token")
	}
	return accessToken, nil
}
