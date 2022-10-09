resource "random_password" "this" {
  length           = 16
  override_special = "!#$%&*()-_=+[]{}<>:?"
  special          = true

  count = 10
}

### Expected Issues ###

# [
#     {
#         "Message": "`count` must be the top-most attribute",
#         "Range": {
#             "Start": { "Line": 6, "Column": 3 },
#             "End": { "Line": 6, "Column": 13 }
#         }
#     }
# ]
