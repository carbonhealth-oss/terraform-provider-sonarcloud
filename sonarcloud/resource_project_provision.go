package sonarcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-sonarcloud/pkg/client"
)

// Returns the resource represented by this file.
func resourceProjectProvision() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectProvisionCreate,
		ReadContext:   resourceProjectProvisionRead,
		UpdateContext: resourceProjectProvisionUpdate,
		DeleteContext: resourceProjectProvisionDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"repo_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"installation_key": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: true,
				Computed: true,
			},
			"project_key": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: true,
				Computed: true,
			},
			"automatic_analysis": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				//Default:  true,
			},
			"quality_gate_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

func resourceProjectProvisionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var _ diag.Diagnostics
	repoName := d.Get("repo_name").(string)
	// Cast m to SonarClient and create POST request for URI with encoded values
	sc := m.(*client.SonarClient)

	repository, err := sc.GetRepository(repoName)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = sc.ProvisionProject(repository.InstallationKey)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set resource id
	d.SetId(repository.InstallationKey)

	return resourceProjectProvisionRead(ctx, d, m)
}

func resourceProjectProvisionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	sc := m.(*client.SonarClient)
	projectKey := d.Get("project_key").(string)

	if d.HasChange("automatic_analysis") {
		_, n := d.GetChange("automatic_analysis")
		err := sc.UpdateAutoAnalysis(n.(bool), projectKey)
		if err != nil {
			return diag.FromErr(err)
		}
		return diags
	}

	if d.HasChange("quality_gate_name") {
		_, n := d.GetChange("quality_gate_name")
		qualityGate, err := sc.GetQualityGate(n.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		err = sc.SelectQualityGate(qualityGate.Id, projectKey)
		if err != nil {
			return diag.FromErr(err)
		}
		return diags
	}
	return resourceProjectProvisionRead(ctx, d, m)
}

func resourceProjectProvisionRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	repoName := d.Get("repo_name").(string)

	sc := m.(*client.SonarClient)
	repository, err := sc.GetRepository(repoName)
	if err != nil {
		return diag.FromErr(err)
	}
	if repository != nil {
		err := d.Set("installation_key", repository.InstallationKey)
		if err != nil {
			return nil
		}
		if len(repository.LinkedProjects) > 0 {
			projectKey := repository.LinkedProjects[0].Key
			err := d.Set("project_key", projectKey)
			if err != nil {
				return diag.FromErr(err)
			}
			component, err := sc.GetComponent(projectKey)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set("automatic_analysis", component.AutoscanEnabled)
			if err != nil {
				return diag.FromErr(err)
			}
			err = d.Set("quality_gate_name", component.QualityGate.Name)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	} else {
		d.SetId("")
	}
	return diags
}

func resourceProjectProvisionDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	sc := m.(*client.SonarClient)
	err := sc.DeleteProject(d.Get("project_key").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Unset the id
	d.SetId("")

	return diags
}
