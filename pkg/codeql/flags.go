package codeql

import (
	"strings"
)

var DatabaseCreateFlags = map[string]bool{
	"--no-db-cluster":           true,
	"--db-cluster":              true,
	"--language":                true,
	"-l":                        true,
	"--command":                 true,
	"-c":                        true,
	"--source-root":             true,
	"-s":                        true,
	"--github-url":              true,
	"-g":                        true,
	"--mode":                    true,
	"-m":                        true,
	"--cleanup-upgrade-backups": true,
	"--extractor-option":        true,
	"-O":                        true,
	"--extractor-options-file":  true,
	"--registries-auth-stdin":   true,
	"--github-auth-stdin":       true,
	"-a":                        true,
	"--threads":                 true,
	"--ram":                     true,
	"-j":                        true,
	"-M":                        true,
	"--search-path":             true,
	"--max-disk-cache":          true,
}

var DatabaseAnalyzeFlags = map[string]bool{
	"--no-rerun":                     true,
	"--rerun":                        true,
	"--no-print-diagnostics-summary": true,
	"--no-print-metrics-summary":     true,
	"--max-paths":                    true,
	"--sarif-add-file-contents":      true,
	"--sarif-add-snippets":           true,
	"--sarif-add-query-help":         true,
	"--sarif-group-rules-by-pack":    true,
	"--sarif-multicause-markdown":    true,
	"--no-sarif-add-file-contents":   true,
	"--no-sarif-add-snippets":        true,
	"--no-sarif-add-query-help":      true,
	"--no-sarif-group-rules-by-pack": true,
	"--no-sarif-multicause-markdown": true,
	"--no-group-results":             true,
	"--csv-location-format":          true,
	"--dot-location-url-format":      true,
	"--sarif-category":               true,
	"--no-download":                  true,
	"--download":                     true,
	"--external":                     true,
	"--warnings":                     true,
	"--no-debug-info":                true,
	//"--no-fast-compilation":          true,	// deprecated
	"--no-local-checking": true,
	//"--fast-compilation":             true,	// deprecated
	"--local-checking":           true,
	"--no-metadata-verification": true,
	"--additional-packs":         true,
	"--registries-auth-stdin":    true,
	"--github-auth-stdin":        true,
	"--threads":                  true,
	"--ram":                      true,
	"-j":                         true,
	"-M":                         true,
	"--search-path":              true,
	"--max-disk-cache":           true,
}

var longShortFlagsMap = map[string]string{
	"--language":          "-l",
	"--command":           "-c",
	"--source-root":       "-s",
	"--github-url":        "-g",
	"--mode":              "-m",
	"--extractor-option":  "-O",
	"--github-auth-stdin": "-a",
	"--threads":           "-j",
	"--ram":               "-M",
}

func ValidateFlags(input map[string]string, validFlags map[string]bool) ([]string, error) {
	params := []string{}

	for flag, value := range input {
		if _, exists := validFlags[flag]; exists {
			appendFlag := flag
			if value != "" {
				appendFlag = appendFlag + "=" + value
			}
			params = append(params, appendFlag)
		}
	}

	return params, nil
}

func CheckIfFlagSetByUser(customFlags map[string]string, flagsToCheck []string) bool {
	for _, flag := range flagsToCheck {
		if _, exists := customFlags[flag]; exists {
			return true
		}
	}
	return false
}

func AppendFlagIfNotPresent(cmd []string, flagToCheck []string, appendFlag []string, customFlags map[string]string) []string {
	if !CheckIfFlagSetByUser(customFlags, flagToCheck) {
		cmd = append(cmd, appendFlag...)
	}
	return cmd
}

func ParseCustomFlags(databaseCreateFlagsStr, databaseAnalyzeFlagsStr string) map[string]string {
	flagStrings := []string{databaseCreateFlagsStr, databaseAnalyzeFlagsStr}
	jointFlags := make(map[string]string)

	for _, flagString := range flagStrings {
		individualFlags := strings.Fields(flagString)
		for _, flag := range individualFlags {
			flagName := strings.Split(flag, "=")[0]
			jointFlags[flagName] = flag
		}
	}

	removeDuplicateFlags(jointFlags, longShortFlagsMap)
	return jointFlags
}

func removeDuplicateFlags(customFlags map[string]string, shortFlags map[string]string) {
	for longFlag, correspondingShortFlag := range shortFlags {
		if _, exists := customFlags[longFlag]; exists {
			delete(customFlags, correspondingShortFlag)
		}
	}
}

func GetRamAndThreadsFromConfig(threads, ram string, customFlags map[string]string) []string {
	params := make([]string, 0, 2)
	if len(threads) > 0 && !CheckIfFlagSetByUser(customFlags, []string{"--threads", "-j"}) {
		params = append(params, "--threads="+threads)
	}
	if len(ram) > 0 && !CheckIfFlagSetByUser(customFlags, []string{"--ram", "-M"}) {
		params = append(params, "--ram="+ram)
	}
	return params
}
