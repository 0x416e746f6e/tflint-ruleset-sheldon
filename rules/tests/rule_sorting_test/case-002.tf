terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.30.0"
    }
  }

  backend "gcs" {
    bucket                      = "terraform"
    impersonate_service_account = "terraform@project.iam.gserviceaccount.com"
  }
}

provider "google" {
  impersonate_service_account = "terraform@project.iam.gserviceaccount.com"
}

### Expected Issues ###

# [
#     {
#         "Message": "block `required_providers` should be placed after `backend gcs` (alphabetical sorting)",
#         "Range": {
#             "Start": { "Line": 2, "Column": 3 },
#             "End": { "Line": 7, "Column": 4 }
#         }
#     }
# ]
