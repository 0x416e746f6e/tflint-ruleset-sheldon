resource "kubernetes_manifest" "qux" {
  manifest {
    apiVersion = "eek"
    kind       = "ook"

    metadata {
      name      = "foo"
      namespace = "bar"
    }
  }
}

### Expected Issues ###

# [
#     {
#         "Message": "higher-priority key-attribute `namespace` should be defined before `name`",
#         "Range": {
#             "Start": { "Line": 8, "Column": 7 },
#             "End": { "Line": 8, "Column": 24 }
#         }
#     }
# ]
