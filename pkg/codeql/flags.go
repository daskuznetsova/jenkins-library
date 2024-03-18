package codeql

import (
	"strings"
)

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

func AppendCustomFlags(input map[string]string) ([]string, error) {
	params := []string{}

	for _, value := range input {
		params = append(params, value)
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

func ParseCustomFlags(flagsStr string) map[string]string {
	flagsMap := make(map[string]string)

	for _, flag := range strings.Fields(flagsStr) {
		if strings.Contains(flag, "=") {
			split := strings.SplitN(flag, "=", 2)
			flagsMap[split[0]] = flag
		} else {
			flagsMap[flag] = flag
		}
	}

	removeDuplicateFlags(flagsMap, longShortFlagsMap)
	return flagsMap
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
