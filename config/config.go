package config

import (
	"encoding/json"
)

// Config is the configuration for the ruleset.
type Config struct {
	Resources []*Resource `hclext:"resource,block"`
}

// Resource is the custom configuration of the resource-specific behaviour.
type Resource struct {
	Kind string   `hclext:"name,label"`
	Keys []string `hclext:"key_attributes"`
}

// New creates a new configuration stucture that is pre-filled with defaults.
func New() *Config {
	return &Config{
		Resources: []*Resource{
			{Kind: "external", Keys: []string{}},
			{Kind: "google_cloud_run_service_iam_policy", Keys: []string{"project", "location", "service"}},
			{Kind: "google_cloud_run_service", Keys: []string{"project", "location", "name"}},
			{Kind: "google_compute_address", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_backend_service", Keys: []string{"project", "name"}},
			{Kind: "google_compute_firewall", Keys: []string{"project", "name"}},
			{Kind: "google_compute_forwarding_rule", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_global_address", Keys: []string{"project", "name"}},
			{Kind: "google_compute_global_forwarding_rule", Keys: []string{"project", "name"}},
			{Kind: "google_compute_health_check", Keys: []string{"project", "name"}},
			{Kind: "google_compute_instance_group_manager", Keys: []string{"project", "zone", "name"}},
			{Kind: "google_compute_instance_template", Keys: []string{"project", "name"}},
			{Kind: "google_compute_instance", Keys: []string{"project", "zone", "name"}},
			{Kind: "google_compute_managed_ssl_certificate", Keys: []string{"project", "name"}},
			{Kind: "google_compute_network_endpoint_group", Keys: []string{"project", "zone", "name"}},
			{Kind: "google_compute_network", Keys: []string{"project", "name"}},
			{Kind: "google_compute_region_instance_group_manager", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_region_network_endpoint_group", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_route", Keys: []string{"project", "name"}},
			{Kind: "google_compute_router_nat", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_router", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_ssl_policy", Keys: []string{"project", "name"}},
			{Kind: "google_compute_subnetwork", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_target_http_proxy", Keys: []string{"project", "name"}},
			{Kind: "google_compute_target_https_proxy", Keys: []string{"project", "name"}},
			{Kind: "google_compute_url_map", Keys: []string{"project", "name"}},
			{Kind: "google_compute_vpn_gateway", Keys: []string{"project", "region", "name"}},
			{Kind: "google_compute_vpn_tunnel", Keys: []string{"project", "region", "name"}},
			{Kind: "google_container_cluster", Keys: []string{"project", "location", "name"}},
			{Kind: "google_container_node_pool", Keys: []string{"project", "locaion", "cluster", "name"}},
			{Kind: "google_container_registry", Keys: []string{"project", "location"}},
			{Kind: "google_dns_managed_zone", Keys: []string{"project", "name"}},
			{Kind: "google_dns_record_set", Keys: []string{"project", "managed_zone", "name", "type"}},
			{Kind: "google_iam_policy", Keys: []string{}},
			{Kind: "google_iap_brand", Keys: []string{"project", "application_title"}},
			{Kind: "google_iap_client", Keys: []string{"brand"}},
			{Kind: "google_iap_web_backend_service_iam_policy", Keys: []string{"project", "web_backend_service"}},
			{Kind: "google_logging_project_sink", Keys: []string{"project", "name"}},
			{Kind: "google_project_iam_audit_config", Keys: []string{"project", "service"}},
			{Kind: "google_project_iam_binding", Keys: []string{"project", "role"}}, // TODO: condition.title
			{Kind: "google_project_iam_member", Keys: []string{"project", "role", "member"}},
			{Kind: "google_project", Keys: []string{"project_id"}},
			{Kind: "google_secret_manager_secret_iam_member", Keys: []string{"project", "secret_id", "role", "member"}},
			{Kind: "google_secret_manager_secret_iam_policy", Keys: []string{"project", "secret_id"}},
			{Kind: "google_secret_manager_secret_version", Keys: []string{"project", "secret"}},
			{Kind: "google_secret_manager_secret", Keys: []string{"project", "secret_id"}},
			{Kind: "google_service_account_access_token", Keys: []string{}},
			{Kind: "google_service_account_iam_policy", Keys: []string{"service_account_id"}},
			{Kind: "google_service_account_key", Keys: []string{"project", "service_account_id"}},
			{Kind: "google_service_account", Keys: []string{"project", "account_id"}},
			{Kind: "google_service_networking_connection", Keys: []string{"project", "network", "service"}},
			{Kind: "google_sql_ca_certs", Keys: []string{"project", "instance"}},
			{Kind: "google_sql_database_instance", Keys: []string{"project", "name"}},
			{Kind: "google_sql_ssl_cert", Keys: []string{"project", "instance", "common_name"}},
			{Kind: "google_sql_user", Keys: []string{"project", "instance", "host", "name"}},
			{Kind: "google_storage_bucket_iam_binding", Keys: []string{"bucket", "role"}},
			{Kind: "google_storage_bucket_iam_policy", Keys: []string{"bucket", "role"}},
			{Kind: "google_storage_bucket", Keys: []string{"project", "name"}},
			{Kind: "google_vpc_access_connector", Keys: []string{"project", "region", "name"}},
			{Kind: "helm_release", Keys: []string{"namespace", "name"}},
			{Kind: "http", Keys: []string{"url"}},
			{Kind: "kubernetes_cluster_role_binding", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_cluster_role", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_config_map_v1_data", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_config_map", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_cron_job", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_deployment", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_endpoints", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_manifest", Keys: []string{"manifest.metadata.namespace", "manifest.metadata.name"}}, // TODO: Search for keys in attributes too
			{Kind: "kubernetes_namespace", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_persistent_volume_claim", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_resource_quota", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_role_binding", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_role", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_secret", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_service_account", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_service", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "kubernetes_stateful_set", Keys: []string{"metadata.namespace", "metadata.name"}},
			{Kind: "null_resource", Keys: []string{}},
			{Kind: "postgresql_database", Keys: []string{"name"}},
			{Kind: "postgresql_grant_role", Keys: []string{"role", "grant_role"}},
			{Kind: "postgresql_role", Keys: []string{"name"}},
			{Kind: "random_password", Keys: []string{}},
			{Kind: "template_file", Keys: []string{}},
			{Kind: "vault_auth_backend", Keys: []string{"path"}},
			{Kind: "vault_database_secret_backend_connection", Keys: []string{"backend", "name"}},
			{Kind: "vault_database_secret_backend_role", Keys: []string{"backend", "name"}},
			{Kind: "vault_gcp_auth_backend_role", Keys: []string{"backend", "role"}},
			{Kind: "vault_gcp_secret_backend", Keys: []string{"path"}},
			{Kind: "vault_gcp_secret_static_account", Keys: []string{"backend", "static_account"}},
			{Kind: "vault_generic_secret", Keys: []string{"path"}},
			{Kind: "vault_identity_group_alias", Keys: []string{"name"}},
			{Kind: "vault_identity_group", Keys: []string{}},
			{Kind: "vault_kubernetes_auth_backend_config", Keys: []string{"backend"}},
			{Kind: "vault_kubernetes_auth_backend_role", Keys: []string{"backend", "role_name"}},
			{Kind: "vault_ldap_auth_backend", Keys: []string{"path"}},
			{Kind: "vault_mount", Keys: []string{"path"}},
			{Kind: "vault_pki_secret_backend_config_urls", Keys: []string{"backend"}},
			{Kind: "vault_pki_secret_backend_intermediate_cert_request", Keys: []string{"backend"}},
			{Kind: "vault_pki_secret_backend_intermediate_set_signed", Keys: []string{"backend"}},
			{Kind: "vault_pki_secret_backend_role", Keys: []string{"backend", "name"}},
			{Kind: "vault_pki_secret_backend_root_cert", Keys: []string{"backend", "name"}},
			{Kind: "vault_pki_secret_backend_root_sign_intermediate", Keys: []string{"backend"}},
			{Kind: "vault_policy", Keys: []string{"name"}},
		},
	}
}

func (r *Resource) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
