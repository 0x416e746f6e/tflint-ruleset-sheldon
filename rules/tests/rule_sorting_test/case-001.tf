resource "something" "this" {
  template {
    metadata {
      annotations = {
        "def" = "blah-blah"
        "abc" = "yada-yada"
      }
    }
  }
}

### Expected Issues ###

# [
#     {
#         "Message": "key `abc` is out of order (should not follow alphabetically greater `def`)",
#         "Range": {
#             "Start": { "Line": 6, "Column": 9 },
#             "End": { "Line": 6, "Column": 14 }
#         }
#     }
# ]
