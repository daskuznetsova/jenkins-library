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
	log.Entry().Debugf("start parsePattern %s", line)

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
	//re1 := regexp.MustCompile(`[^/\\]\*\*`)
	re1 := regexp.MustCompile(`[^\x2f\x5c]\*\*`)
	//re2 := regexp.MustCompile(`^\\*\\*[^/]`)
	re2 := regexp.MustCompile(`^\*\*[^/]`)
	//re3 := regexp.MustCompile(`[^\\]\*\*[^/]`)
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

func FilterSarif(input string, output string, patterns []*Pattern) error {
	outputSarif := Sarif{}

	file, err := os.Open(input)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to open sarif file: %s", err)
	}

	decoder := json.NewDecoder(file)
	sarif := Sarif{}

	err = decoder.Decode(&sarif)
	if err != nil {
		return fmt.Errorf("failed to Decode the JSON: %s", err)
	}

	for _, run := range sarif.Runs {
		if run.Results != nil {
			newResults := []Result{}
			for _, result := range run.Results {
				if result.Locations != nil {
					newLocations := []Location{}
					for _, location := range result.Locations {
						uri := location.PhysicalLocation.ArtifactLocation.Uri
						ruleId := result.RuleID
						matched, err := matchPathAndRule(uri, ruleId, patterns)
						if err != nil {
							return err
						}
						if uri == "" || matched {
							newLocations = append(newLocations, location)
						}
					}
					if len(newLocations) > 0 {
						result.Locations = newLocations
						newResults = append(newResults, result)
					}
				} else {
					newResults = append(newResults, result)
				}
			}
			run.Results = newResults
			outputSarif.Runs = append(outputSarif.Runs, run)
		}
	}

	file, err = os.Create(output)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to create filtered sarif file: %s", err)
	}

	writer := json.NewEncoder(file)
	writer.SetIndent("", "    ")
	err = writer.Encode(outputSarif)
	if err != nil {
		return fmt.Errorf("failed to encode filtered sarif file: %s", err)
	}

	log.Entry().Infof("Successfully written the JSON log to %s", output)
	return nil
}
