// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type contrastExecuteScanOptions struct {
	UserAPIKey                  string `json:"userApiKey,omitempty"`
	ServiceKey                  string `json:"serviceKey,omitempty"`
	Username                    string `json:"username,omitempty"`
	Server                      string `json:"server,omitempty"`
	OrganizationID              string `json:"organizationId,omitempty"`
	ApplicationID               string `json:"applicationId,omitempty"`
	VulnerabilityThresholdTotal int    `json:"vulnerabilityThresholdTotal,omitempty"`
	CheckForCompliance          bool   `json:"checkForCompliance,omitempty"`
}

// ContrastExecuteScanCommand
func ContrastExecuteScanCommand() *cobra.Command {
	const STEP_NAME = "contrastExecuteScan"

	metadata := contrastExecuteScanMetadata()
	var stepConfig contrastExecuteScanOptions
	var startTime time.Time
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createContrastExecuteScanCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "",
		Long:  `This step evaluates if the audit requirements for Contrast Assess have been fulfilled after the execution of security tests by Contrast Assess. For further information on the tool, please consult the [documentation](https://github.wdf.sap.corp/pages/Security-Testing/doc/contrast/introduction/).`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.UserAPIKey)
			log.RegisterSecret(stepConfig.ServiceKey)
			log.RegisterSecret(stepConfig.Username)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 || len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if err = log.RegisterANSHookIfConfigured(GeneralConfig.CorrelationID); err != nil {
				log.Entry().WithError(err).Warn("failed to set up SAP Alert Notification Service log hook")
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.Dsn,
						GeneralConfig.HookConfig.SplunkConfig.Token,
						GeneralConfig.HookConfig.SplunkConfig.Index,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
				if len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblToken,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblIndex,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			contrastExecuteScan(stepConfig, &stepTelemetryData)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addContrastExecuteScanFlags(createContrastExecuteScanCmd, &stepConfig)
	return createContrastExecuteScanCmd
}

func addContrastExecuteScanFlags(cmd *cobra.Command, stepConfig *contrastExecuteScanOptions) {
	cmd.Flags().StringVar(&stepConfig.UserAPIKey, "userApiKey", os.Getenv("PIPER_userApiKey"), "User API Key for authorizing access to Contrast.")
	cmd.Flags().StringVar(&stepConfig.ServiceKey, "serviceKey", os.Getenv("PIPER_serviceKey"), "Service Key for authorization access to Contrast.")
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "Username or email to use for authorization access to Contrast.")
	cmd.Flags().StringVar(&stepConfig.Server, "server", os.Getenv("PIPER_server"), "Contrast server url.")
	cmd.Flags().StringVar(&stepConfig.OrganizationID, "organizationId", os.Getenv("PIPER_organizationId"), "Organization Id in Contrast.")
	cmd.Flags().StringVar(&stepConfig.ApplicationID, "applicationId", os.Getenv("PIPER_applicationId"), "Application Id in Contrast.")
	cmd.Flags().IntVar(&stepConfig.VulnerabilityThresholdTotal, "vulnerabilityThresholdTotal", 0, "Threshold for maximum number of allowed vulnerabilities.")
	cmd.Flags().BoolVar(&stepConfig.CheckForCompliance, "checkForCompliance", false, "If set to true, the piper step checks for compliance based on vulnerability thresholds. Example - If total vulnerabilities are 10 and vulnerabilityThresholdTotal is set as 0, then the steps throws an compliance error.")

}

// retrieve step metadata
func contrastExecuteScanMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "contrastExecuteScan",
			Aliases:     []config.Alias{},
			Description: "",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "userCredentialsId", Description: "Jenkins 'Username with password' credentials ID containing username and user API Key to communicate with the Contrast server.", Type: "jenkins"},
					{Name: "serviceKeyCredentialsId", Description: "Jenkins 'Secret text' credentials ID containing service key to communicate with the Contrast server.", Type: "jenkins"},
				},
				Parameters: []config.StepParameters{
					{
						Name: "userApiKey",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "userCredentialsId",
								Param: "userApiKey",
								Type:  "secret",
							},

							{
								Name:    "userApiKey",
								Type:    "vaultSecret",
								Default: "contrast",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_userApiKey"),
					},
					{
						Name: "serviceKey",
						ResourceRef: []config.ResourceReference{
							{
								Name: "serviceKeyCredentialsId",
								Type: "secret",
							},

							{
								Name:    "serviceKey",
								Type:    "vaultSecret",
								Default: "contrast",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "service_key"}},
						Default:   os.Getenv("PIPER_serviceKey"),
					},
					{
						Name: "username",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "userCredentialsId",
								Param: "username",
								Type:  "secret",
							},

							{
								Name:    "username",
								Type:    "vaultSecret",
								Default: "contrast",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_username"),
					},
					{
						Name:        "server",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_server"),
					},
					{
						Name:        "organizationId",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_organizationId"),
					},
					{
						Name:        "applicationId",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_applicationId"),
					},
					{
						Name:        "vulnerabilityThresholdTotal",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "int",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     0,
					},
					{
						Name:        "checkForCompliance",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
				},
			},
		},
	}
	return theMetaData
}