package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Rules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Endpoint string         `json:"endpoint"`
	Method   string         `json:"method"`
	Items    []EndpointRule `json:"items"`
}

type EndpointRule struct {
	QueryString string           `json:"queryString"`
	Body        string           `json:"body"`
	Counter     *int             `json:"counter,omitempty"`
	Response    EndpointResponse `json:"response"`
}

type EndpointResponse struct {
	Status  int      `json:"status"`
	Delay   int      `json:"delay"`
	Headers []string `json:"headers"`
	Body    string   `json:"body"`
}

func LoadRulesFromFile(fileName string) (Rules, error) {
	var rules Rules

	if len(fileName) == 0 {
		return rules, nil
	}

	jsonFile, err := os.Open(fileName)

	if err != nil {
		return rules, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(byteValue, &rules); err != nil {
		return rules, err
	}

	return rules, nil
}

func IsQueryStringMatchRule(requestQueryString string, ruleQueryString string) bool {
	if len(ruleQueryString) > 0 {
		r := regexp.MustCompile(fmt.Sprintf("%s%s", `(?m)`, ruleQueryString))
		return r.MatchString(requestQueryString)
	}
	return true
}

func IsBodyMatchRule(bRequest string, bRule string) bool {
	if len(bRule) > 0 {
		sampleRegexp := regexp.MustCompile(bRule)
		return sampleRegexp.MatchString(bRequest)
	}

	return true
}
