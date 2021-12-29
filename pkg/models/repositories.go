package models

// GetRepositories for unmarshalling response body where repositories are returned
type GetRepositories struct {
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	Label           string `json:"label,omitempty"`
	InstallationKey string `json:"installationKey,omitempty"`
	LinkedProjects        []struct {
		Key  string `json:"key,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"linkedProjects,omitempty"`
	Private bool `json:"private,omitempty"`
}

type ProvisionRepositoryResponse struct {
	Projects []struct {
		ProjectKey string `json:"projectKey,omitempty"`
	} `json:"projects,omitempty"`
}
