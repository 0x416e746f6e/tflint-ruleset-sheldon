# Rule `sheldon_lifecycle`

Makes sure that `lifecycle` meta-block is placed at the bottom of
`resource` or `data` block (but before `depends_on` if it's present).

## Example

```hcl
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"

  depends_on = [
    aws_iam_role_policy.example
  ]

  iam_instance_profile = aws_iam_instance_profile.example
}
```

```text
Error: `lifecycle` block must be at the end of the definition (but before `depends_on`) (sheldon_lifecycle)

  on template.tf line 9:
   9:   lifecycle {
  10:     create_before_destroy = true
  11:   }
```
