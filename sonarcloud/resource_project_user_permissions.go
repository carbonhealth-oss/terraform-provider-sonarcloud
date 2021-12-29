package sonarcloud

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"terraform-provider-sonarcloud/pkg/client"
	"terraform-provider-sonarcloud/pkg/collection"
)

// Returns the resource represented by this file.
func resourceProjectUserPermissions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectUserPermissionsCreate,
		ReadContext:   resourceProjectUserPermissionsRead,
		UpdateContext: resourceProjectUserPermissionsUpdate,
		DeleteContext: resourceProjectUserPermissionsDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						var diags diag.Diagnostics

						allowed := []string{
							// Project permissions
							"admin",
							"codeviewer",
							"issueadmin",
							"securityhotspotadmin",
							"scan",
							"user",
						}
						ok := false
						for _, v := range allowed {
							if v == i.(string) {
								ok = true
								break
							}
						}
						if !ok {
							return diag.Errorf("unsupported permission '%s'", i.(string))
						}

						return diags
					},
				},
			},
		},
	}
}

func resourceProjectUserPermissionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get all resource values
	login := d.Get("login").(string)
	permissions := expandPermissions(d)
	projectKey := d.Get("project_key").(string)

	sc := m.(*client.SonarClient)

	// loop through all permissions that should be applied
	for _, permission := range *&permissions {
		permission := permission
		if err := sc.AddPermissionForUser(
			login,
			projectKey,
			permission,
		); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set resource id
	if diags == nil || !diags.HasError() {
		d.SetId(uuid.New().String())
		//d.SetId(projectKey + login)
		return resourceProjectUserPermissionsRead(ctx, d, m)
	}
	return diags
}

func resourceProjectUserPermissionsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	projectKey := d.Get("project_key").(string)
	login := d.Get("login").(string)
	sc := m.(*client.SonarClient)

	users, err := sc.GetUsersForProject(login, projectKey)
	if err != nil {
		return diag.FromErr(err)
	}
	// Loop over all users to see if the user login we need exists.
	found := false
	for _, value := range users {
		if strings.EqualFold(value.Login, login) {
			found = true
			if err := d.Set("permissions", value.Permissions); err != nil {
				return diag.FromErr(err)
			}
			break
		}
	}

	// Unset the id if the resource has not been found
	if !found {
		d.SetId("")
	}

	return diags
}
func resourceProjectUserPermissionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	login := d.Get("login").(string)
	projectKey := d.Get("project_key").(string)

	if d.HasChange("permissions") {
		// Find which permissions to add and which to remove
		o, n := d.GetChange("permissions")
		added, removed := collection.Diff(o.([]interface{}), n.([]interface{}))

		// Get client and prepare synchronization
		sc := m.(*client.SonarClient)

		// Remove permissions
		for _, permission := range removed {
			if err := sc.RemovePermissionForUser(login, projectKey, permission.(string)); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		// Add permissions
		for _, permission := range added {
			if err := sc.AddPermissionForUser(
				login,
				projectKey,
				permission.(string),
			); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		if diags == nil || !diags.HasError() {
			return resourceProjectUserPermissionsRead(ctx, d, m)
		}
	}
	return diags
}

func resourceProjectUserPermissionsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	//// Get resource value that's needed to read the remote resource
	login := d.Get("login").(string)
	permissions := expandPermissions(d)
	projectKey := d.Get("project_key").(string)

	// Cast m to SonarClient and create POST request for URI with encoded values
	sc := m.(*client.SonarClient)
	// loop through all permissions that should be applied
	for _, permission := range *&permissions {
		err := sc.RemovePermissionForUser(login, projectKey, permission)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// Unset the id
	d.SetId("")

	return diags
}

func expandPermissions(d *schema.ResourceData) []string {
	expandedPermissions := make([]string, 0)
	flatPermissions := d.Get("permissions").([]interface{})
	for _, permission := range flatPermissions {
		expandedPermissions = append(expandedPermissions, permission.(string))
	}

	return expandedPermissions
}
