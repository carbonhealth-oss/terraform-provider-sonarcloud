package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"terraform-provider-sonarcloud/pkg/models"
)

func (sc *SonarClient) ProvisionProject(installationKey string) (models.ProvisionRepositoryResponse, error) {
	var params []string
	params = []string{"installationKeys", installationKey}

	resp, err := sc.HttpReqWithParams(http.MethodPost, fmt.Sprintf("%s/alm_integration/provision_projects", API), params...)

	emptyResponse := models.ProvisionRepositoryResponse{}
	if err != nil {
		return emptyResponse, err
	}

	provisionRepo := models.ProvisionRepositoryResponse{}
	err = json.NewDecoder(resp.Body).Decode(&provisionRepo)
	if err != nil {
		return emptyResponse, fmt.Errorf("decode error: %+v", err)
	}
	return provisionRepo, nil
}

func (sc *SonarClient) DeleteProject(projectKey string) error {
	var params []string
	params = []string{"project", projectKey}
	_, err := sc.HttpReqWithParams(http.MethodPost, fmt.Sprintf("%s/projects/delete", API), params...)
	return err
}

func (sc *SonarClient) GetRepository(repoName string) (*models.Repository, error) {
	resp, err := sc.HttpReqWithParams(http.MethodGet, fmt.Sprintf("%s/alm_integration/list_repositories", API))
	if err != nil {
		return nil, err
	}

	repos := &models.GetRepositories{}
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, fmt.Errorf("decode error: %+v", err)
	}

	// Loop over all repos to see if the user login we need exists.
	for _, repository := range repos.Repositories {
		if strings.EqualFold(repository.Label, repoName) {
			return &repository, nil
		}
	}
	return nil, fmt.Errorf("unable to find given repo %s, please check if it's a valid git repo and associated with sonarcloud", repoName)
}

func (sc *SonarClient) UpdateAutoAnalysis(enableAutoAnalysis bool, projectKey string) error {
	var params []string
	params = []string{"enable", strconv.FormatBool(enableAutoAnalysis), "projectKey", projectKey}
	_, err := sc.HttpReqWithParams(http.MethodPost, fmt.Sprintf("%s/autoscan/activation", API), params...)
	return err
}
