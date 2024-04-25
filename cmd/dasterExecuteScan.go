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
	token, err := fetchOAuthToken(config)
	if err != nil {
		return err
	}

	dasterInstance := daster.NewDaster(token, config.ServiceURL)
	scan, err := dasterInstance.TriggerScan()
	if err != nil {
		log.Entry().WithError(err).Error("failed to trigger scan")
		return err
	}
	if scan.ScanId == "" {
		return nil
	}

	if config.Synchronous && config.ScanType != "burpscan" {
		var scanResponse *daster.ScanResponse
		for {
			scanResponse, err = dasterInstance.GetScanResponse(scan.ScanId)
			if err != nil {
				return err
			}
			if scanResponse.State.Terminated != nil {
				break
			}
			time.Sleep(15 * time.Second)
		}
		scanResult, err := dasterInstance.GetScanResult(scanResponse)
		if err != nil {
			return err
		}
		violations := daster.CheckThresholdViolations(&daster.ThresholdViolations{}, scanResult)
		if violations != nil {
			//error "[${STEP_NAME}][ERROR] Threshold(s) ${thresholdViolations} violated by findings '${scanResult.summary}'"
		} else if scanResponse.State.Terminated.ExitCode != 0 {
			//error "[${STEP_NAME}][ERROR] Scan failed with code '${scanResponse?.state?.terminated?.exitCode}', reason '${scanResponse?.state?.terminated?.reason}' on container '${scanResponse?.state?.terminated?.containerID}'"
		} else {
			log.Entry().Infof("Result of scan is %v", scanResponse)
		}
		err = dasterInstance.DeleteScan(scan.ScanId)
		if err != nil {
			log.Entry().WithError(err).Warn("failed to delete scan")
		}
	}
	if config.ScanType == "burpscan" {

	}
	return nil
}

/*
def runScan(parameters, utils, config, body) {
    withCredentials([string(
        credentialsId: config.settings.dasterTokenCredentialsId,
        variable: 'token',
    )]) {
        def extendedConfig = [:].plus(config)
        extendedConfig.settings.remove('dasterTokenCredentialsId')
        extendedConfig.settings.dasterToken = token

        def daster = parameters.dasterStub ?: new Daster(this, utils, extendedConfig)
        def scan = daster.triggerScan()
        echo "[${STEP_NAME}][INFO] Triggered scan of type ${config.scanType}${scan.message ? ' and received message: \'' + scan.message  + '\'' : ''}: ${scan.url ?: scan.scanId + ' and waiting for it to complete'}"

	if (scan?.scanId) {
            if(config.synchronous && config.scanType != 'burpscan') {
                try {
                    def scanResponse = [:]
                    while (scanResponse?.state?.terminated == null) {
                        scanResponse = daster.getScanResponse(scan?.scanId)
                        sleep(15)
                    }
                    def scanResult = daster.getScanResult(scanResponse)
                    def thresholdViolations = checkThresholdViolations(config, scanResult)
                    if (thresholdViolations) {
                        error "[${STEP_NAME}][ERROR] Threshold(s) ${thresholdViolations} violated by findings '${scanResult.summary}'"
                    } else if (scanResponse?.state?.terminated?.exitCode) {
                        error "[${STEP_NAME}][ERROR] Scan failed with code '${scanResponse?.state?.terminated?.exitCode}', reason '${scanResponse?.state?.terminated?.reason}' on container '${scanResponse?.state?.terminated?.containerID}'"
                    } else {
                        echo "Result of scan is ${scanResponse}"
                    }
                } finally {
                    if (config.deleteScan)
                        daster.deleteScan(scan?.scanId)
                }
            } else if (config.scanType == 'burpscan') {
                try {
                    withEnv(["BURP_PROXY=${scan.proxyURL}".toString()]) {
                        body()
                    }
                    def scanResponse = daster.getScanResponse(scan.scanId)
                    def scanResult = daster.getScanResult(scanResponse)
                    def thresholdViolations = checkThresholdViolations(config, scanResult)
                    if (thresholdViolations) {
                        error "[${STEP_NAME}][ERROR] Threshold(s) ${thresholdViolations} violated by findings '${scanResult.summary}'"
                    } else if (scanResponse?.state?.terminated?.exitCode) {
                        error "[${STEP_NAME}][ERROR] Scan failed with code '${scanResponse?.state?.terminated?.exitCode}', reason '${scanResponse?.state?.terminated?.reason}' on container '${scanResponse?.state?.terminated?.containerID}'"
                    } else {
                        echo "Result of scan is ${scanResponse}"
                    }
                } finally {
                    daster.stopBurpScan(scan.scanId)
                }
            }
        }
    }
}
*/

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
