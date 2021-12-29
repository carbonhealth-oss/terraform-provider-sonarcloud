package models


// PermissionsAddUser Add permission to a user.<br /> This service defaults to global permissions, but can be limited to project permissions by providing project id or project key.<br />Requires the permission 'Administer' on the specified project.
type PermissionsAddUser struct {
	Login        string `form:"login,omitempty"`        // User login
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning</li><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	ProjectId    string `form:"projectId,omitempty"`    // Project id
	ProjectKey   string `form:"projectKey,omitempty"`   // Project key
}

// PermissionsRemoveUser Remove permission from a user.<br /> This service defaults to global permissions, but can be limited to project permissions by providing project id or project key.<br /> Requires the permission 'Administer' on the specified project.
type PermissionsRemoveUser struct {
	Login        string `form:"login,omitempty"`        // User login
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning</li><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	ProjectId    string `form:"projectId,omitempty"`    // Project id
	ProjectKey   string `form:"projectKey,omitempty"`   // Project key
}
