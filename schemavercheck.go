package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type SchemaVerCompatibilityArgs struct {
	SchemaPath        string
	EndpointURLFormat string
	DefinitionName    string
}

type SchemaVerCompatibilityResult struct {
	IsValid  bool
	ErrorMsg string
}

type compatibilityEndpointResponse struct {
	Ok       bool   `json:"ok"`
	ErrorMsg string `json:"error"`
}

func CheckSchemaVerCompatibility(args SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
	url := fmt.Sprintf(args.EndpointURLFormat, args.DefinitionName)
	schemaBytes, err := ioutil.ReadFile(args.SchemaPath)
	if err != nil {
		return SchemaVerCompatibilityResult{IsValid: false}, err
	}
	status, response, err := makePostRequest(url, schemaBytes)
	if err != nil {
		return SchemaVerCompatibilityResult{IsValid: false}, err
	}
	if status == 200 {
		return SchemaVerCompatibilityResult{IsValid: true}, nil
	}
	if status == 409 {
		return SchemaVerCompatibilityResult{IsValid: false, ErrorMsg: response.ErrorMsg}, nil
	}
	return SchemaVerCompatibilityResult{IsValid: false}, fmt.Errorf("Received error response(%d): %s", status, response.ErrorMsg)
}

func makePostRequest(url string, requestBody []byte) (int, compatibilityEndpointResponse, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, compatibilityEndpointResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, compatibilityEndpointResponse{}, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return status, compatibilityEndpointResponse{}, err
	}
	responseBodyTrimmed := strings.Trim(string(responseBody), " \n\t")
	log.Printf("Response from compatibility endpoint: %s (code=%d)", responseBodyTrimmed, status)
	var response compatibilityEndpointResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return status, compatibilityEndpointResponse{}, err
	}
	return status, response, nil
}
