package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const API = "https://sonarcloud.io/api"

type SonarClient struct {
	client *http.Client
	org    string
	token  string
}

func NewSonarClient(org string, token string) (*SonarClient, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	return &SonarClient{
		client: client,
		org:    org,
		token:  token,
	}, nil
}

func (sc *SonarClient) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("organization", sc.org)
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(sc.token, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (sc *SonarClient) NewRequestWithParameters(method string, url string, params ...string) (*http.Request, error) {
	if l := len(params); l%2 != 0 {
		return nil, fmt.Errorf("params must be an even number, %d given", l)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("organization", sc.org)

	for i := 0; i < len(params); i++ {
		q.Add(params[i], params[i+1])
		i++
	}
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(sc.token, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (sc *SonarClient) DoReq(req *http.Request) (http.Response, error) {
	// Perform the request
	resp, err := sc.Do(req)
	if err != nil {
		return http.Response{}, fmt.Errorf("failed to execute http request: %v. Request: %v", err, req)
	}

	// Check status code and return diagnostics from ErrorResponse if needed
	if resp.StatusCode >= 300 {

		return *resp,errorResponse(resp)
	}

	return *resp,nil
}

func closeHttpResp(resp http.Response){
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
}

func (sc *SonarClient) HttpReqWithParams(method string, url string, params ...string) (http.Response, error){
	req, err := sc.NewRequestWithParameters(method, url, params...)
	if err != nil {
		return http.Response{},err
	}
	return sc.DoReq(req)
}

func (sc *SonarClient) HttpReq(method string, url string, body io.Reader) (http.Response, error){
	req, err := sc.NewRequest(method, url, body)
	if err != nil {
		return http.Response{},err
	}
	return sc.DoReq(req)
}

func (sc *SonarClient) Do(req *http.Request) (*http.Response, error) {
	return sc.client.Do(req)
}
