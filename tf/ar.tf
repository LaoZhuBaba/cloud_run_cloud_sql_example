resource "google_artifact_registry_repository" "sandbox-repo" {
  location      = local.region
  repository_id = local.ar_repo_name
  format        = "DOCKER"
}
