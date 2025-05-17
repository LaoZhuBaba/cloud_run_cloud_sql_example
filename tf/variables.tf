locals {
  project                 = "qoria-sandbox"
  region                  = "australia-southeast1"
  zone                    = "australia-southeast1-a"
  ar_repo_name            = "sandbox-repo"
  app_name                = "example-app"
  pg_instance_name        = "pgsql-instance"
  sql_subnet_cidr         = "10.1.0.0/24"
  cloud_run_subnet_cidr   = "10.2.0.0/24"
  vpc_name                = "example-vpc"
  sql_internal_ip         = "10.1.0.250"
  sql_internal_port       = "5432"
  sql_db_user_name        = "david"
  password_version_number = 2
}
