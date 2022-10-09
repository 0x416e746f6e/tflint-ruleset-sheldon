resource "kubernetes_config_map" "this" {

  count = var.create_map ? 1 : 0

  metadata {
    namespace = kubernetes_namespace.this.metadata[0].name
    name      = "config-map"
  }

  data = {
    "foo" = "bar"
  }
}

### Expected Issues ###

# [
#     {
#         "Message": "1 redundant blank line in front",
#         "Range": {
#             "Start": { "Line": 3, "Column": 3 },
#             "End": { "Line": 3, "Column": 33 }
#         }
#     }
# ]
