resource "kubernetes_service_account" "this" {
  metadata {
    namespace = kubernetes_namespace.this.metadata[0].name
    name      = "service-account"
  }
}
