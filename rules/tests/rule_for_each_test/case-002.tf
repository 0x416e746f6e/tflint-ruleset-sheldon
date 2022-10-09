resource "google_container_registry" "this" {
  project  = "foo"
  location = each.key
  for_each = ["a", "b"]
}

### Expected Issues ###

# [
#     {
#         "Message": "`for_each` must be the top-most attribute",
#         "Range": {
#             "Start": { "Line": 4, "Column": 3 },
#             "End": { "Line": 4, "Column": 24 }
#         }
#     }
# ]
