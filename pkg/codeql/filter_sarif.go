package codeql

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gobwas/glob"
)

type Pattern struct {
	glob     glob.Glob
	isIgnore bool
}

type Sarif struct {
	Runs []Run `json:"runs"`
}

type Run struct {
	Results []Result `json:"results"`
}

type Result struct {
	RuleID   string     `json:"ruleId"`
	Location []Location `json:"locations"`
}

type Location struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}

func FilterSarif(filename string, patterns []Pattern, outputFilename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var res Sarif
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}

	for _, run := range res.Runs {
		for _, result := range run.Results {
			for _, matchPattern := range patterns {
				if matchPattern.glob.Match(result.RuleID) {
				}
			}
		}
	}

	//print('Given patterns:')
	//for s, fp, rp in args.patterns:
	//print(
	//	'files: {file_pattern}    rules: {rule_pattern} ({sign})'.format(
	//		file_pattern=fp,
	//	rule_pattern=rp,
	//	sign='positive' if s else 'negative'
	//)
	//)

	data, err = json.MarshalIndent(res, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFilename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ParsePattern(pattern string) (sign bool, filePattern string, rulePattern string, err error) {
	sepchar := ':'
	escchar := '\\'
	filePattern = ""
	rulePattern = ""
	seenSeparator := false
	sign = true

	// inclusion or exclusion pattern?
	upattern := pattern
	if pattern != "" {
		if pattern[0] == '-' {
			sign = false
			upattern = pattern[1:]
		} else if pattern[0] == '+' {
			upattern = pattern[1:]
		}
	}

	for i := 0; i < len(upattern); i++ {
		c := upattern[i]
		if int32(c) == sepchar {
			if seenSeparator {
				return sign, filePattern, rulePattern, fmt.Errorf("invalid pattern: %q contains more than one separator", pattern)
			}
			seenSeparator = true
			continue
		} else if int32(c) == escchar {
			nextc := rune(0)
			if i+1 < len(upattern) {
				nextc = rune(upattern[i+1])
			}
			if nextc == '+' || nextc == '-' || nextc == escchar || nextc == sepchar {
				i++
				c = uint8(nextc)
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

	return sign, filePattern, rulePattern, nil
}
