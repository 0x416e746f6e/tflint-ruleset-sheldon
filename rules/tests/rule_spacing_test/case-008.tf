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
#         "Message": "attribute `count` must be separated from the rest of the definition by an extra line",
#         "Range": {
#             "Start": { "Line": 2, "Column": 3 },
#             "End": { "Line": 2, "Column": 33 }
#         }
#     },
#     {
#         "Message": "multi-line element must be separated from the previous one by an extra line",
#         "Range": {
#             "Start": { "Line": 3, "Column": 3 },
#             "End": { "Line": 6, "Column": 4 }
#         }
#     },
#     {
#         "Message": "1 redundant blank line in front",
#         "Range": {
#             "Start": { "Line": 12, "Column": 1 },
#             "End": { "Line": 12, "Column": 2 }
#         }
#     }
# ]
