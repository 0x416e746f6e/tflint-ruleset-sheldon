resource "kubernetes_config_map" "this" {
  count = var.create_map ? 1 : 0

  metadata {
    namespace = kubernetes_namespace.this.metadata[0].name
  }


  metadata {
    name = "config-map"
  }
}

### Expected Issues ###

# [
#     {
#         "Message": "1 redundant blank line in front",
#         "Range": {
#             "Start": { "Line": 9, "Column": 3 },
#             "End": { "Line": 11, "Column": 4 }
#         }
#     }
# ]
