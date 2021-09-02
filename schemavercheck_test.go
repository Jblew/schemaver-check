package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckSchemaVerCompatibilityCallsHTTPEndpoint(t *testing.T) {
	noCalls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"ok\":true}")
		w.WriteHeader(200)
		noCalls++
	}))
	defer ts.Close()

	_, err := CheckSchemaVerCompatibility(SchemaVerCompatibilityArgs{
		EndpointURLFormat: ts.URL,
		SchemaPath:        "mock/schema.json",
		DefinitionName:    "ChartSpec",
	})

	if err != nil {
		t.Errorf("Error %+v", err)
	}

	if noCalls == 0 {
		t.Errorf("CheckSchemaVerCompatibility does not call endpoint")
	}
}

func TestCheckSchemaVerCompatibilityPostsSchema(t *testing.T) {
	bodyStr := ""
	methodStr := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		bodyStr = string(bodyBytes)
		methodStr = r.Method
		fmt.Fprintln(w, "{\"ok\":true}")
		w.WriteHeader(200)
	}))
	defer ts.Close()

	_, err := CheckSchemaVerCompatibility(SchemaVerCompatibilityArgs{
		EndpointURLFormat: ts.URL,
		SchemaPath:        "mock/schema.json",
		DefinitionName:    "ChartSpec",
	})

	if err != nil {
		t.Errorf("Error %+v", err)
	}

	if methodStr != "POST" {
		t.Errorf("CheckSchemaVerCompatibility does not call endpoint using POST")
	}

	if !strings.Contains(bodyStr, "CompanyStructureSpec") {
		t.Errorf("CheckSchemaVerCompatibility does not post JSON schema")
	}

	var bodyUnm map[string]interface{}
	err = json.Unmarshal([]byte(bodyStr), &bodyUnm)
	if err != nil {
		t.Errorf("CheckSchemaVerCompatibility does not post a valid JSON: %+v", err)
	}
}

func TestCheckSchemaVerCompatibilityReturnsValidOn200(t *testing.T) {
	ts := makeHttpEndpointWithResponse(200, "{\"ok\":true}")
	defer ts.Close()

	res, err := CheckSchemaVerCompatibility(SchemaVerCompatibilityArgs{
		EndpointURLFormat: ts.URL,
		SchemaPath:        "mock/schema.json",
		DefinitionName:    "ChartSpec",
	})

	if err != nil {
		t.Errorf("Error %+v", err)
	}

	if !res.IsValid {
		t.Errorf("CheckSchemaVerCompatibility does not return Valid on 200 server response")
	}
}

func TestCheckSchemaVerCompatibilityReturnsInvalidOn409(t *testing.T) {
	ts := makeHttpEndpointWithResponse(409, "{\"error\":\"Incompatible JSON Schema\"}")
	defer ts.Close()

	res, err := CheckSchemaVerCompatibility(SchemaVerCompatibilityArgs{
		EndpointURLFormat: ts.URL,
		SchemaPath:        "mock/schema.json",
		DefinitionName:    "ChartSpec",
	})

	if err != nil {
		t.Errorf("Error %+v", err)
	}

	if res.IsValid {
		t.Errorf("CheckSchemaVerCompatibility returns Valid on 409 server response")
	}

	if res.ErrorMsg != "Incompatible JSON Schema" {
		t.Errorf("CheckSchemaVerCompatibility does not return error message")
	}
}

func TestCheckSchemaVerCompatibilityReturnsErrorOn400(t *testing.T) {
	ts := makeHttpEndpointWithResponse(400, "{\"error\":\"Invalid Definition Name\"}")
	defer ts.Close()

	res, err := CheckSchemaVerCompatibility(SchemaVerCompatibilityArgs{
		EndpointURLFormat: ts.URL,
		SchemaPath:        "mock/schema.json",
		DefinitionName:    "ChartSpec",
	})

	if err == nil {
		t.Errorf("Should throw error on 400 server response")
	}

	if err != nil && !strings.Contains(fmt.Sprintf("%+v", err), "Invalid Definition Name") {
		t.Errorf("Error does not container server returned msg")
	}

	if res.IsValid {
		t.Errorf("CheckSchemaVerCompatibility returns Valid on 400 server response")
	}

	if res.ErrorMsg != "" {
		t.Errorf("CheckSchemaVerCompatibility has nonempty error msg on 400 response")
	}
}

func makeHttpEndpointWithResponse(statusCode int, response string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
		w.WriteHeader(statusCode)
	}))
	return ts
}
