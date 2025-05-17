resource "google_project_service" "cloud_sql_api" {
  project = local.project
  service = "sqladmin.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_on_destroy = false
}

resource "google_project_service" "net_connectivity_api" {
  project = local.project
  service = "networkconnectivity.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_on_destroy = false
}

resource "google_project_service" "secret_manager_api" {
  project = local.project
  service = "secretmanager.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_on_destroy = false
}
