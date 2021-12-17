package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
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
	Param       string           `json:"param"`
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

func IsParamMatchRule(pRequest []string, qRule string) bool {
	return (len(pRequest) > 0 && qRule == pRequest[0]) || qRule == ""
}

func IsQueryStringMatchRule(qString map[string][]string, queryString string) bool {
	pUrl, err := url.ParseQuery(queryString)
	if err != nil {
		return false
	}
	ruleQueryString := map[string][]string(pUrl)
	return reflect.DeepEqual(qString, ruleQueryString)
}

func IsBodyMatchRule(bRequest string, bRule string) bool {
	if len(bRule) > 0 {
		sampleRegexp := regexp.MustCompile(bRule)
		return sampleRegexp.MatchString(bRequest)
	}

	return true
}
