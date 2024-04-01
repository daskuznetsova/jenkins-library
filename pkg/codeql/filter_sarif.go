package codeql

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/SAP/jenkins-library/pkg/log"
)

type Pattern struct {
	sign        bool
	filePattern string
	rulePattern string
}

type ArtifactLocation struct {
	Uri string `json:"uri"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
}

type Location struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type Result struct {
	RuleID    string     `json:"ruleId"`
	Locations []Location `json:"locations"`
}

type Run struct {
	Results []Result `json:"results"`
}

type Sarif struct {
	Runs []Run `json:"runs"`
}

func ParsePatterns(filterPattern string) ([]*Pattern, error) {
	patterns := []*Pattern{}
	patternsSplit := strings.Split(filterPattern, " ")
	for _, pattern := range patternsSplit {
		parsedPattern, err := parsePattern(pattern)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, parsedPattern)
		log.Entry().Infof("files: %s, rules: %s (include: %t)", parsedPattern.filePattern, parsedPattern.rulePattern, parsedPattern.sign)
	}
	return patterns, nil
}

// Helper function to get sign and trim pattern
func getSignAndTrimPattern(pattern string) (bool, string) {
	sign := true

	if strings.HasPrefix(pattern, "-") {
		sign = false
		pattern = strings.TrimPrefix(pattern, "-")
	} else if strings.HasPrefix(pattern, "+") {
		pattern = strings.TrimPrefix(pattern, "+")
	}

	return sign, pattern
}

// Helper function to separate file and rule pattern
func separateFileAndRulePattern(pattern string) (string, string, error) {
	escChar := '\\'
	sepChar := ':'
	filePattern := ""
	rulePattern := ""
	seenSeparator := false

	for i := 0; i < len(pattern); i++ {
		c := rune(pattern[i])

		if c == sepChar {
			if seenSeparator {
				return "", "", fmt.Errorf("Invalid pattern: '%s'. Contains more than one separator!\n", pattern)
			}
			seenSeparator = true
			continue
		} else if c == escChar {
			// If we find an escape character and the current position
			// is less than total length - 1, we get the next character instead
			var nextC rune
			if i+1 < len(pattern) {
				nextC = rune(pattern[i+1])
			}
			if nextC == '+' || nextC == '-' || nextC == escChar || nextC == sepChar {
				i++
				c = nextC
			}
		}

		if seenSeparator {
			rulePattern += string(c)
		} else {
			filePattern += string(c)
		}
	}

	return filePattern, rulePattern, nil
}

func parsePattern(line string) (*Pattern, error) {
	sign, pattern := getSignAndTrimPattern(line)
	filePattern, rulePattern, err := separateFileAndRulePattern(pattern)
	if err != nil {
		return nil, err
	}
	if rulePattern == "" {
		rulePattern = "**"
	}

	log.Entry().Infof("rulePattern %s, filePattern %s", rulePattern, filePattern)

	return &Pattern{
		sign:        sign,
		filePattern: filePattern,
		rulePattern: rulePattern,
	}, nil
}

func ReadSarifFile(input string) (map[string]interface{}, error) {
	var sarif map[string]interface{}
	file, err := os.Open(input)
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to open sarif file: %s", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sarif)
	if err != nil {
		return nil, fmt.Errorf("failed to Decode the JSON: %s", err)
	}
	return sarif, nil
}

func ProcessSarif(sarif map[string]interface{}, patterns []*Pattern) (map[string]interface{}, error) {
	runs, ok := sarif["runs"].([]interface{})
	if !ok {
		return sarif, nil
	}

	for _, run := range runs {
		runMap, ok := run.(map[string]interface{})
		if !ok {
			continue
		}

		results, ok := runMap["results"].([]interface{})
		if !ok {
			continue
		}

		newResults := []interface{}{}
		for i, result := range results {
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				continue
			}

			locations, ok := resultMap["locations"].([]interface{})
			if !ok {
				continue
			}

			newLocations := []interface{}{}
			for _, location := range locations {
				locationMap, ok := location.(map[string]interface{})
				if !ok {
					continue
				}

				uri, ok := locationMap["physicalLocation"].(map[string]interface{})["artifactLocation"].(map[string]interface{})["uri"].(string)
				if !ok {
					continue
				}
				ruleId := resultMap["ruleId"].(string)
				log.Entry().Infof("checking location: %s", uri)
				matched, err := matchPathAndRule(uri, ruleId, patterns)
				if err != nil {
					return nil, err
				}

				if uri != "" && !matched {
					log.Entry().Infof("added location to results: %s", uri)
					newLocations = append(newLocations, location)
				} else {
					log.Entry().Infof("removed %v from results", uri)
				}
			}

			if len(newLocations) > 0 {
				log.Entry().Infof("added result %d with locations (len=%d)", i, len(newLocations))
				resultMap["locations"] = newLocations
				newResults = append(newResults, result)
				log.Entry().Info(result)
			} else {
				log.Entry().Infof("empty locations, removed result %d", i)
			}
		}
		log.Entry().Infof("len of new results: %d", len(newResults))
		runMap["results"] = newResults
	}

	return sarif, nil
}

func WriteSarifFile(output string, sarif map[string]interface{}) error {
	file, err := os.Create(output)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to create filtered sarif file: %s", err)
	}
	writer := json.NewEncoder(file)
	writer.SetIndent("", "    ")
	err = writer.Encode(sarif)
	if err != nil {
		return fmt.Errorf("failed to encode filtered sarif file: %s", err)
	}
	log.Entry().Infof("Successfully written the JSON log to %s", output)
	return nil
}

// implements glob matching
func matchPathAndRule(uri string, ruleId string, patterns []*Pattern) (bool, error) {
	result := true
	for _, p := range patterns {
		matchedRule, err := match(p.rulePattern, ruleId)
		if err != nil {
			return false, err
		}
		matchedFile, err := match(p.filePattern, uri)
		if err != nil {
			return false, err
		}
		if matchedRule && matchedFile {
			result = p.sign
		}
	}
	return result, nil
}

func matchComponent(patternComponent string, fileNameComponent string) bool {
	if len(patternComponent) == 0 && len(fileNameComponent) == 0 {
		return true
	}
	if len(patternComponent) == 0 {
		return false
	}
	if len(fileNameComponent) == 0 {
		return patternComponent == "*"
	}
	if string(patternComponent[0]) == "*" {
		return matchComponent(patternComponent, fileNameComponent[1:]) ||
			matchComponent(patternComponent[1:], fileNameComponent)
	}
	if string(patternComponent[0]) == "?" {
		return matchComponent(patternComponent[1:], fileNameComponent[1:])
	}
	if string(patternComponent[0]) == "\\" {
		return len(patternComponent) >= 2 && patternComponent[1] == fileNameComponent[0] &&
			matchComponent(patternComponent[2:], fileNameComponent[1:])
	}
	if patternComponent[0] != fileNameComponent[0] {
		return false
	}

	return matchComponent(patternComponent[1:], fileNameComponent[1:])
}

func matchComponents(patternComponents []string, fileNameComponents []string) bool {
	if len(patternComponents) == 0 && len(fileNameComponents) == 0 {
		return true
	}
	if len(patternComponents) == 0 {
		return false
	}
	if len(fileNameComponents) == 0 {
		return len(patternComponents) == 1 && patternComponents[0] == "**"
	}
	if patternComponents[0] == "**" {
		matchFiles := matchComponents(patternComponents, fileNameComponents[1:])
		matchPattern := matchComponents(patternComponents[1:], fileNameComponents)
		return matchFiles || matchPattern
	} else {
		matchStart := matchComponent(patternComponents[0], fileNameComponents[0])
		matchEnd := matchComponents(patternComponents[1:], fileNameComponents[1:])
		return matchStart && matchEnd
	}
}

func match(pattern string, fileName string) (bool, error) {
	re1 := regexp.MustCompile(`[^\x2f\x5c]\*\*`)
	re2 := regexp.MustCompile(`^\*\*[^/]`)
	re3 := regexp.MustCompile(`[^\x5c]\*\*[^/]`)

	if re1.MatchString(pattern) || re2.MatchString(pattern) || re3.MatchString(pattern) {
		return false, fmt.Errorf("`**` in %v not alone between path separators \n", pattern)
	}

	pattern = strings.TrimSuffix(pattern, "/")
	fileName = strings.TrimSuffix(fileName, "/")
	for strings.Contains(pattern, "**/**") {
		pattern = strings.Replace(pattern, "**/**", "**", -1)
	}

	splitter := regexp.MustCompile(`[\\/]+`)
	fileNameComponents := splitter.Split(fileName, -1)
	patternComponents := strings.Split(pattern, "/")

	return matchComponents(patternComponents, fileNameComponents), nil

}
