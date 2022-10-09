resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"

  depends_on = [
    aws_iam_role_policy.example
  ]

  lifecycle {
    create_before_destroy = true
  }

  iam_instance_profile = aws_iam_instance_profile.example
}

### Expected Issues ###

# [
#     {
#         "Message": "`lifecycle` block must be at the end of the definition (but before `depends_on`)",
#         "Range": {
#             "Start": { "Line": 9, "Column": 13 },
#             "End": { "Line": 11, "Column": 4 }
#         }
#     }
# ]
