package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"terraform-provider-sonarcloud/pkg/models"
)

func (sc *SonarClient) getQualityGates() (*models.ListQualityGates, error) {

	resp, err := sc.HttpReqWithParams(http.MethodGet, fmt.Sprintf("%s/qualitygates/list", API))
	if err != nil {
		return nil, err
	}

	qualityGates := &models.ListQualityGates{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGates)
	if err != nil {
		return nil, fmt.Errorf("decode error: %+v", err)
	}

	return qualityGates, nil
}
func (sc *SonarClient) GetQualityGate(gateName string) (*models.QualityGate, error) {

	resp, err := sc.getQualityGates()
	if err != nil {
		return nil, err
	}
	// loop through all permissions that should be applied
	for _, qualityGate := range *&resp.QualityGate {
		if qualityGate.Name == gateName {
			return &qualityGate, nil
		}
	}

	return nil, fmt.Errorf("invalid quality gate name: %+v", gateName)
}

func (sc *SonarClient) SelectQualityGate(gateId int, projectKey string) error {
	var params []string
	params = []string{"gateId", strconv.Itoa(gateId), "projectKey", projectKey}
	_, err := sc.HttpReqWithParams(http.MethodPost, fmt.Sprintf("%s/qualitygates/select", API), params...)
	if err != nil {
		return err
	}
	return nil
}
