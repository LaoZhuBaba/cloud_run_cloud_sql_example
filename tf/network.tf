resource "google_compute_network" "example_vpc" {
  name                    = local.vpc_name
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "sql_subnet" {
  name          = "sql-subnet"
  ip_cidr_range = local.sql_subnet_cidr
  region        = local.region
  network       = google_compute_network.example_vpc.id
}

resource "google_compute_subnetwork" "cloud_run_subnet" {
  name          = "cloud-run-subnet"
  ip_cidr_range = local.cloud_run_subnet_cidr
  region        = local.region
  network       = google_compute_network.example_vpc.id
}

resource "google_compute_address" "default" {
  name         = "psc-compute-address"
  region       = local.region
  address_type = "INTERNAL"
  subnetwork   = google_compute_subnetwork.sql_subnet.name
  address      = local.sql_internal_ip
}

data "google_sql_database_instance" "pgsql" {
  name = resource.google_sql_database_instance.pgsql.name
}

resource "google_compute_forwarding_rule" "default" {
  name                  = "psc-forwarding-rule-${google_sql_database_instance.pgsql.name}"
  region                = local.region
  network               = google_compute_network.example_vpc.id
  ip_address            = google_compute_address.default.self_link
  load_balancing_scheme = ""
  target                = data.google_sql_database_instance.pgsql.psc_service_attachment_link
}
