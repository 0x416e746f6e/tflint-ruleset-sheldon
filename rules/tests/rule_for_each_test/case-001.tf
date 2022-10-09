resource "google_container_registry" "this" {
  for_each = local.cfg.gcp.gcr

  project  = local.cfg.gcp.project_id
  location = upper(each.key)
}
