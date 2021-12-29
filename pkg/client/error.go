package client

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
)

type ErrorResponse struct {
	Errors []struct {
		Message string `json:"msg"`
	}
}

func DiagErrorResponse(resp *http.Response, diags diag.Diagnostics) diag.Diagnostics {
	errResponse := &ErrorResponse{}
	err := json.NewDecoder(resp.Body).Decode(&errResponse)
	if err != nil {
		return diag.Errorf("API returned a code %d but the error response could not be parsed: %+v", resp.StatusCode, err)
	}

	for _, e := range errResponse.Errors {
		diags = appendDiagErrorFromStr(diags, fmt.Sprintf("API (%d): %s", resp.StatusCode, e.Message))
	}

	return diags
}

func errorResponse(resp *http.Response) error {
	errResponse := ErrorResponse{}
	err := json.NewDecoder(resp.Body).Decode(&errResponse)
	if err != nil {
		return fmt.Errorf("API returned a code %d but the error response could not be parsed: %+v", resp.StatusCode, err)
	}
	return fmt.Errorf("API returned an error: %+v", errResponse.Errors[0].Message)
}

func appendDiagErrorFromStr(diags diag.Diagnostics, s string) diag.Diagnostics {
	return append(diags, diag.Diagnostic{
		Severity: 0,
		Summary:  s,
	})
}
