resource "google_cloud_run_v2_service" "example-app" {
  name     = local.app_name
  location = local.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "${local.region}-docker.pkg.dev/${local.project}/${local.ar_repo_name}/${local.app_name}"
      #  env {
      #     name = "INSTANCE_CONNECTION_NAME"
      #     //        value = "qoria-sandbox:australia-southeast1:pgsql-instance"
      #     value = "${local.project}:${local.region}:${local.pg_instance_name}"
      #   }
      env {
        name  = "DB_USER"
        value = "david"
      }
      env {
        name = "DB_PASS"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_password_secret.secret_id
            version = "latest"
          }
        }
      }
      env {
        name  = "DB_NAME"
        value = google_sql_database.example_database.name
      }
      env {
        name  = "DB_IP_ADDRESS"
        value = local.sql_internal_ip
      }
      env {
        name  = "DB_IP_PORT"
        value = local.sql_internal_port
      }
    }
    vpc_access {
      network_interfaces {
        network    = google_compute_network.example_vpc.id
        subnetwork = google_compute_subnetwork.cloud_run_subnet.id
      }
    }
    service_account = google_service_account.cloud_run_sa.email
  }
  depends_on = [google_secret_manager_secret_version.example_db_password_secret_version]
}
