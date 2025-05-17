resource "google_service_account" "cloud_run_sa" {
  account_id   = "cloud-run-sa"
  display_name = "Service Account for Cloud Run example"
}

resource "google_project_iam_binding" "sa_bindings" {
  for_each = toset([
    "roles/cloudsql.client",
    "roles/cloudsql.instanceUser",
    "roles/logging.logWriter",
    "roles/secretmanager.secretAccessor"
  ])
  project = local.project
  role    = each.value

  members = [
    google_service_account.cloud_run_sa.member,
  ]
}
