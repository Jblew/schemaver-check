package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
		DefinitionName:    "#/definitions/ChartSpec",
	})

	if err != nil {
		t.Errorf("Error %+v", err)
	}

	if noCalls == 0 {
		t.Errorf("CheckSchemaVerCompatibility does not call endpoint")
	}
}
