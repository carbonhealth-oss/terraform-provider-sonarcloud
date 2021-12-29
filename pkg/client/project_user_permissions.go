package client

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"strings"
	"terraform-provider-sonarcloud/pkg/models"
)

func (sc *SonarClient) GetUsersForProject(login string, projectKey string) ([]models.User, error) {
	var params []string
	params = []string{"projectKey", projectKey, "q", login, "ps", "100"}

	resp, err := sc.HttpReqWithParams(http.MethodGet, fmt.Sprintf("%s/permissions/users", API), params...)
	if err != nil {
		return nil, err
	}

	users := &models.GetUser{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("decode error: %+v", err)
	}
	return users.Users, nil
}
func (sc *SonarClient) AddPermissionForUser(login string, projectKey string, permission string) error {
	// Use name because the organization is always set and using an id will then throw an error...
	create := models.PermissionsAddUser{
		Login:        login,
		Permission:   permission,
		ProjectKey:   projectKey,
	}
	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&create)
	if err != nil {
		return err
	}

	_, err = sc.HttpReq(http.MethodPost, fmt.Sprintf("%s/permissions/add_user", API), strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	return nil
}

func (sc *SonarClient) RemovePermissionForUser(login string, projectKey string, permission string) error {

	del := models.PermissionsRemoveUser{
		Login:        login,
		Permission:   permission,
		ProjectKey:   projectKey,
	}

	// Encode the values
	encoder := form.NewEncoder()
	values, err := encoder.Encode(&del)
	if err != nil {
		return err
	}

	_, err = sc.NewRequest(http.MethodPost, fmt.Sprintf("%s/permissions/remove_user", API), strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	return nil
}
