# Rule `sheldon_source`

Makes sure that `source` meta-attribute is placed at the top of `module` block.

## Example

```hcl
module "website_s3_bucket" {
  bucket_name = "<UNIQUE BUCKET NAME>"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }

  source = "./modules/aws-s3-static-website-bucket"
}
```

```text
Error: `source` must be the top-most attribute (sheldon_source)

  on template.tf line 9:
   9:   source = "./modules/aws-s3-static-website-bucket"
```
