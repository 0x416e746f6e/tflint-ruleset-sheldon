module "website_s3_bucket" {
  bucket_name = "<UNIQUE BUCKET NAME>"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }

  source = "./modules/aws-s3-static-website-bucket"
}

### Expected Issues ###

# [
#     {
#         "Message": "`source` must be the top-most attribute",
#         "Range": {
#             "Start": { "Line": 9, "Column": 3 },
#             "End": { "Line": 9, "Column": 52 }
#         }
#     }
# ]
