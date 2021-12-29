package sonarcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-sonarcloud/pkg/client"
)

// Provider SonarCloud
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SONAR_ORGANIZATION", nil),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SONAR_TOKEN", nil),
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sonarcloud_project_user_permissions": resourceProjectUserPermissions(),
			"sonarcloud_project_provision":        resourceProjectProvision(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	org := d.Get("organization").(string)
	token := d.Get("token").(string)

	c, err := client.NewSonarClient(org, token)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	var diags diag.Diagnostics

	return c, diags
}

type Config struct {
	Organization string
	Token        string
}
