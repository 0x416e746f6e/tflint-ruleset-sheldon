terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.30.0"
    }
  }

  # Comments restart sorting

  backend "gcs" {
    bucket                      = "terraform"
    impersonate_service_account = "terraform@project.iam.gserviceaccount.com"
  }
}

provider "google" {
  impersonate_service_account = "terraform@project.iam.gserviceaccount.com"
}
