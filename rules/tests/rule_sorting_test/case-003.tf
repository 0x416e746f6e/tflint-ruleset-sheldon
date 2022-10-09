resource "google_storage_bucket" "cloudbuild" {
  for_each = local.cfg.gcp.cloudbuild

  project = local.cfg.gcp.project_id
  name    = "${local.cfg.gcp.cloudbuild_bucket_prefix}-${lower(each.key)}"

  location = upper(each.key)

  force_destroy               = true
  uniform_bucket_level_access = true

  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age = 1
    }
  }

  versioning {
    enabled = false
  }
}
