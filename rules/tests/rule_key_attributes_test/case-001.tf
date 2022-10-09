resource "google_compute_subnetwork" "network-with-private-secondary-ip-ranges" {
  region        = "us-central1"
  name          = "test-subnetwork"
  ip_cidr_range = "10.2.0.0/16"
  network       = google_compute_network.custom-test.id
  secondary_ip_range {
    range_name    = "tf-test-secondary-range-update1"
    ip_cidr_range = "192.168.10.0/24"
  }
}
