terraform {
  required_providers {
    sonarcloud = {
      source  = "carbon/platform/sonarcloud"
      version = "0.0.1"
    }
  }
}

provider "sonarcloud" {
  organization = var.organization
  token        = var.token
}


resource "sonarcloud_project_provision" "this" {
  repo_name          = var.repo
  automatic_analysis = false
  quality_gate_name  = "Sonar way"
}

resource "sonarcloud_project_user_permissions" "user1" {
  login       = var.test_user
  project_key = sonarcloud_project_provision.this.project_key
  permissions = ["admin"]
}
