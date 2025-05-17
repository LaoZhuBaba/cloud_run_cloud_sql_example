resource "google_secret_manager_secret" "db_password_secret" {
  secret_id = "example-db-password"

  replication {
    user_managed {
      replicas {
        location = local.region
      }
    }
  }
}

# Password is generated here but not stored in the state file
ephemeral "random_password" "db_password" {
  length = 32
  # Avoid special characters because some may be unsupported.
  special = false
}

# Define an ephemeral resource to fetch the latest version of the secret
resource "google_secret_manager_secret_version" "example_db_password_secret_version" {
  secret                 = google_secret_manager_secret.db_password_secret.id
  secret_data_wo         = ephemeral.random_password.db_password.result
  secret_data_wo_version = local.password_version_number
}
