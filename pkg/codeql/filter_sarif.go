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
		sign, filePattern, rulePattern, err := parsePattern(pattern)
		patterns = append(patterns, &Pattern{
			sign:        sign,
			filePattern: filePattern,
			rulePattern: rulePattern,
		})
		if err != nil {
			return nil, err
		}
		s := "positive"
		if !sign {
			s = "negative"
		}
		log.Entry().Debugf("files: %s	rules: %s	(%s)", filePattern, rulePattern, s)
	}
	return patterns, nil
}

func parsePattern(line string) (bool, string, string, error) {
	sign := true
	filePattern := ""
	rulePattern := ""
	seenSeparator := false
	escChar := '\\'
	sepChar := ':'

	if strings.HasPrefix(line, "-") {
		sign = false
		line = strings.TrimPrefix(line, "-")
	} else if strings.HasPrefix(line, "+") {
		line = strings.TrimPrefix(line, "+")
	}

	for i := 0; i < len(line); i++ {
		c := rune(line[i])

		if c == sepChar {
			if seenSeparator {
				return false, "", "", fmt.Errorf("Invalid pattern: '%s'. Contains more than one separator!\n", line)
			}
			seenSeparator = true
		} else if c == escChar {
			nextC := rune(line[i+1])
			if i+1 < len(line) && (nextC == '+' || nextC == '-' || nextC == escChar || nextC == sepChar) {
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

	if rulePattern == "" {
		rulePattern = "**"
	}

	log.Entry().Debugf("rulePattern %s, filePattern %s", rulePattern, filePattern)

	return sign, filePattern, rulePattern, nil
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
		for _, result := range results {
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
				matched, err := matchPathAndRule(uri, ruleId, patterns)
				if err != nil {
					return nil, err
				}

				if uri != "" && !matched {
					log.Entry().Infof("removed %v from results", uri)
					newLocations = append(newLocations, location)
				}
			}

			if len(newLocations) == 0 {
				resultMap["locations"] = newLocations
				newResults = append(newResults, result)
			}
		}

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
		return matchComponents(patternComponents, fileNameComponents[1:]) ||
			matchComponents(patternComponents[1:], fileNameComponents)
	} else {
		return matchComponent(patternComponents[0], fileNameComponents[0]) &&
			matchComponents(patternComponents[1:], fileNameComponents[1:])
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
	patternComponents := strings.Split(pattern, "/")
	fileNameComponents := splitter.Split(fileName, -1)

	return matchComponents(patternComponents, fileNameComponents), nil

}
