resource "google_sql_database_instance" "pgsql" {
  name             = local.pg_instance_name
  database_version = "POSTGRES_15"
  region           = local.region
  project          = local.project

  settings {
    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = "db-f1-micro"
    ip_configuration {
      psc_config {
        psc_enabled               = true
        allowed_consumer_projects = [local.project]
        # psc_auto_connections {
        #   consumer_network            = google_compute_network.example_vpc.id
        #   consumer_service_project_id = local.project
        # }
      }
      ipv4_enabled = false
    }
    backup_configuration {
      enabled = false
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
