package main

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsParamMatchRuleWhenParamEquals(t *testing.T) {
	assert.True(t, IsParamMatchRule([]string{"AnSpecificValueHere"}, "AnSpecificValueHere"))
}

func TestIsParamMatchRuleWhenParamEmpty(t *testing.T) {
	assert.True(t, IsParamMatchRule([]string{}, ""))
}

func TestIsParamMatchRuleWhenAnyParam(t *testing.T) {
	assert.True(t, IsParamMatchRule([]string{"AnyRandomValue"}, ""))
}

func TestIsParamMatchRuleWhenParamNotMatch(t *testing.T) {
	assert.False(t, IsParamMatchRule([]string{"AnyRandomValue"}, "AnotherValue"))
	assert.False(t, IsParamMatchRule([]string{}, "AnotherValue"))
}

func TestIsQueryStringMatchRule(t *testing.T) {
	type test struct {
		queryString     map[string][]string
		queryStringRule string
	}
	tests := []test{
		{queryString: map[string][]string{"foo": {"1"}, "bar": {"2"}}, queryStringRule: "foo=1&bar=2"},
		{queryString: map[string][]string{"bar": {"2"}, "foo": {"1"}}, queryStringRule: "foo=1&bar=2"},
		{queryString: map[string][]string{"foo": {"1"}}, queryStringRule: "foo=1"},
	}
	for _, tc := range tests {
		assert.True(t, IsQueryStringMatchRule(tc.queryString, tc.queryStringRule))
	}
}

func TestIsQueryStringDoesNotMatchRule(t *testing.T) {
	type test struct {
		queryString     map[string][]string
		queryStringRule string
	}
	tests := []test{
		{queryString: map[string][]string{"foo": {"1"}, "bar": {"2"}}, queryStringRule: ""},
		{queryString: map[string][]string{"bar": {"2"}, "foo": {"1"}}, queryStringRule: "foo=1"},
		{queryString: map[string][]string{"foo": {"1"}}, queryStringRule: "foo=2"},
	}
	for _, tc := range tests {
		assert.False(t, IsQueryStringMatchRule(tc.queryString, tc.queryStringRule))
	}
}

func TestIsBodyMatchRuleWhenBodyEquals(t *testing.T) {
	assert.True(t, IsBodyMatchRule("{\"sample-key\": \"sample-value\"}", "{\"sample-key\": \"sample-value\"}"))
}

func TestIsBodyMatchRuleWhenBodyEmpty(t *testing.T) {
	assert.True(t, IsBodyMatchRule("", ""))
}

func TestIsBodyMatchRuleWhenAnyBody(t *testing.T) {
	assert.True(t, IsBodyMatchRule("{\"some-random-key\": \"random-value\"}", ""))
}

func TestIsBodyMatchRuleWhenBodyNotMatch(t *testing.T) {
	assert.False(t, IsBodyMatchRule("{\"specific-key\": \"specific-value\"}", "{\"another-key\": \"another-value\"}"))
	assert.False(t, IsBodyMatchRule("", "{\"unmatchable-key\": \"unmatchable-value\"}"))
}

func TestIsBodyMatchRuleWhenRegexMultipleWords(t *testing.T) {
	assert.True(t, IsBodyMatchRule(
		"{\"name\": \"sample\", \"role\": \"director\", \"age\": \"77\"}",
		".*name.*role.*age",
	))
}

func TestRespondGetMethod(t *testing.T) {
	h := func(c echo.Context, er EndpointRule) {}
	r := Rule{Method: "GET", Items: []EndpointRule{{}}}

	e := echo.New()
	c := e.NewContext(
		httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(""))),
		httptest.NewRecorder(),
	)

	assert.True(t, respondOnRuleMatch(c, r, h))
}

func TestRespondPostMethodWhenBodyAndKeyEquals(t *testing.T) {
	h := func(c echo.Context, er EndpointRule) {}
	b := bytes.NewBuffer([]byte("{\"key\":123}"))
	r := Rule{Method: "POST", Items: []EndpointRule{{Param: "AnSpecificValueHere", Body: "{\"key\":123}"}}}

	e := echo.New()
	c := e.NewContext(
		httptest.NewRequest(http.MethodPost, "/", b),
		httptest.NewRecorder(),
	)

	c.SetParamNames("sampleParamKey")
	c.SetParamValues("AnSpecificValueHere")

	assert.True(t, respondOnRuleMatch(c, r, h))
}

func TestSolvePostMethodWhenBodyMatchExpression(t *testing.T) {
	h := func(c echo.Context, er EndpointRule) {}
	b := bytes.NewBuffer([]byte("{\"key\":123, \"hash\":\"HaSh123aBc\"}"))
	r := Rule{Method: "POST", Items: []EndpointRule{{Body: ".*key.*hash"}}}

	e := echo.New()
	c := e.NewContext(
		httptest.NewRequest(http.MethodPost, "/", b),
		httptest.NewRecorder(),
	)

	assert.True(t, respondOnRuleMatch(c, r, h))
}

func TestParseBody(t *testing.T) {
	bRequest := `{"sha256":{"sample-key-1":{"gen_data":{}},"sample-key-2":{"gen_data":{}}},"flags":1}`
	bRule1 := "-> {{sha256|@keys|0}} <-"
	bRule2 := "-> {{sha256|@keys|1}} <-"

	assert.Equal(t, "-> sample-key-1 <-", parseBody(bRule1, bRequest))
	assert.Equal(t, "-> sample-key-2 <-", parseBody(bRule2, bRequest))
}
