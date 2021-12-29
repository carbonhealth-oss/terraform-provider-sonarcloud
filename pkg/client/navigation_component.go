package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-sonarcloud/pkg/models"
)

func (sc *SonarClient) GetComponent(projectKey string) (*models.GetComponentResponse, error) {
	var params []string
	params = []string{"component", projectKey}
	resp, err := sc.HttpReqWithParams(http.MethodGet, fmt.Sprintf("%s/navigation/component", API),params...)
	if err != nil {
		return nil, err
	}

	component := &models.GetComponentResponse{}
	err = json.NewDecoder(resp.Body).Decode(&component)
	if err != nil {
		return nil, fmt.Errorf("decode error: %+v", err)
	}

	return component, nil
}
