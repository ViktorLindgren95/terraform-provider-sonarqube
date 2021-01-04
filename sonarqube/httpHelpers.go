package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

// ErrorResponse struct
type ErrorResponse struct {
	Errors []ErrorMessage `json:"errors,omitempty"`
}

// ErrorMessage struct
type ErrorMessage struct {
	Message string `json:"msg,omitempty"`
}

// Paging used in /search API endpoints
type Paging struct {
	PageIndex int64 `json:"pageIndex"`
	PageSize  int64 `json:"pageSize"`
	Total     int64 `json:"total"`
}

// helper function to make api request to sonarqube
func httpRequestHelper(client *retryablehttp.Client, method string, sonarqubeURL string, expectedResponseCode int, errormsg string) (http.Response, error) {
	// Prepare request
	req, err := retryablehttp.NewRequest(method, sonarqubeURL, http.NoBody)
	if err != nil {
		return http.Response{}, fmt.Errorf("Failed to prepare http request: %v. Request: %v", err, req)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return http.Response{}, fmt.Errorf("Failed to execute http request: %v. Request: %v", err, req)
	}

	// Check response code
	if resp.StatusCode != expectedResponseCode {
		if resp.Body == http.NoBody {
			// No error message in the body
			return *resp, fmt.Errorf("StatusCode: %v does not match expectedResponseCode: %v", resp.StatusCode, expectedResponseCode)
		}

		// The response body has content, try to decode the error message
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return *resp, fmt.Errorf("Failed to decode error response json into struct: %+v", err)
		}
		return *resp, fmt.Errorf("API returned an error: %+v", errorResponse.Errors[0].Message)
	}

	return *resp, nil
}
