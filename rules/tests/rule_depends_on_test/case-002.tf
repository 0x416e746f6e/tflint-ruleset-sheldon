resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"

  depends_on = [
    aws_iam_role_policy.example
  ]

  iam_instance_profile = aws_iam_instance_profile.example
}

### Expected Issues ###

# [
#     {
#         "Message": "`depends_on` clause must be the last one in the definition",
#         "Range": {
#             "Start": { "Line": 5, "Column": 3 },
#             "End": { "Line": 7, "Column": 4 }
#         }
#     }
# ]
