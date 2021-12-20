package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var counterMap = map[string]int{}

func main() {
	flags := ParseFlags()
	e := prepareEcho()
	rules, err := LoadRulesFromFile(*flags.Rules)

	if err != nil {
		fmt.Println("Failed to load rules. ", err)
		os.Exit(1)
	}

	defineRoutesFromRules(e.Router(), rules)

	if err := e.Start(*flags.Host + ":" + *flags.Port); err != nil {
		fmt.Println("Failed to initialize API router. ", err)
		os.Exit(1)
	}
}

func prepareEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e.DefaultHTTPErrorHandler(err, c)
		FormattedLog(c)
	}

	return e
}

// Parse the requested body against the ruled response body
func parseBody(bRule string, bRequest string) string {
	r := regexp.MustCompile(`{{[^{}]*}}`)
	dynRuleTags := r.FindAllString(bRule, -1)
	result := bRule

	for _, dynRuleTag := range dynRuleTags {
		dynRule := dynRuleTag[2 : len(dynRuleTag)-2]
		dynValue := gjson.Get(bRequest, dynRule)

		result = strings.ReplaceAll(result, dynRuleTag, dynValue.Str)
	}

	return result
}

// default rule response executor. Basically just parse the content from the rule response
func processResponseBasedOnRule(c echo.Context, er EndpointRule) {
	if er.Response.Delay > 0 {
		time.Sleep(time.Duration(er.Response.Delay) * time.Millisecond)
	}

	for _, h := range er.Response.Headers {
		header := strings.Split(h, ":")
		c.Response().Header().Add(strings.TrimSpace(header[0]), strings.TrimSpace(header[1]))
	}

	bResponse := er.Response.Body
	bRequest := bodyAsString(c.Request())

	if gjson.Valid(bRequest) {
		bResponse = parseBody(bResponse, bRequest)
	}

	c.Blob(er.Response.Status, "", []byte(bResponse))
}

func bodyAsString(c *http.Request) string {
	reqBody, _ := ioutil.ReadAll(c.Body)
	c.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	return string(reqBody)
}

func respondOnRuleMatch(
	c echo.Context,
	rule Rule,
	handlerFunc func(echo.Context, EndpointRule),
) bool {
	for _, er := range rule.Items {
		hash := createCounterHash(rule.Endpoint, rule.Method, er.Body)
		if isRuleMatch(c, er) && isCounterMatch(hash, er) {
			counterMap[hash]++
			handlerFunc(c, er)
			return true
		}
	}

	return false
}

func isRuleMatch(c echo.Context, er EndpointRule) bool {
	return IsQueryStringMatchRule(c.QueryString(), er.QueryString) &&
		IsBodyMatchRule(bodyAsString(c.Request()), er.Body)
}

func isCounterMatch(counterHash string, er EndpointRule) bool {
	if er.Counter == nil {
		return true
	}
	return counterMap[counterHash] == *er.Counter
}

func defineRoutesFromRules(router *echo.Router, rules Rules) {
	for _, rule := range rules.Rules {
		r := rule

		router.Add(rule.Method, rule.Endpoint, func(c echo.Context) error {
			if !respondOnRuleMatch(c, r, processResponseBasedOnRule) {
				c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
					"rule": r,
					"body": bodyAsString(c.Request()),
				})
			}

			FormattedLog(c)

			return nil
		})
	}
}

func createCounterHash(endpoint string, method string, body string) string {
	hashObject, _ := json.Marshal(struct {
		Endpoint string
		Method   string
		Body     string
	}{
		Endpoint: endpoint,
		Method:   method,
		Body:     body,
	})
	hasher := sha1.New()
	hasher.Write(hashObject)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
