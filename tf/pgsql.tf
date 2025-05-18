resource "google_sql_database_instance" "pgsql" {
  name             = local.pg_instance_name
  database_version = "POSTGRES_15"
  region           = local.region
  project          = local.project

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      psc_config {
        psc_enabled               = true
        allowed_consumer_projects = [local.project]
      }
      ipv4_enabled = false
    }
    backup_configuration {
      enabled = false
    }
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = true
    }
  }
}

resource "google_sql_database" "example_database" {
  name     = "example-db"
  instance = google_sql_database_instance.pgsql.name
}

resource "google_sql_user" "sql_db_user" {
  name                = "david"
  instance            = google_sql_database_instance.pgsql.name
  type                = "BUILT_IN"
  password_wo         = ephemeral.random_password.db_password.result
  password_wo_version = local.password_version_number
}

resource "google_sql_user" "sql_db_sa_user" {
  name     = trimsuffix(google_service_account.cloud_run_sa.email, ".gserviceaccount.com")
  instance = google_sql_database_instance.pgsql.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

