variable "create_map" {}

resource "kubernetes_config_map" "this" {

  # Optional resource

  count = var.create_map ? 1 : 0

  metadata {
    namespace = kubernetes_namespace.this.metadata[0].name
    name      = "config-map"
  }

  data = {
    "foo" = "bar"
  }
}
