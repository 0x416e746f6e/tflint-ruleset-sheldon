resource "kubernetes_service_account" "this" {
  metadata {
    name      = "service-account"
    namespace = kubernetes_namespace.this.metadata[0].name
  }
}

### Expected Issues ###

# [
#     {
#         "Message": "higher-priority key-attribute `namespace` should be defined before `name`",
#         "Range": {
#             "Start": { "Line": 4, "Column": 5 },
#             "End": { "Line": 4, "Column": 59 }
#         }
#     }
# ]
