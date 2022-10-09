# Rule `sheldon_depends_on`

Makes sure that `depends_on` meta-attribute is placed at the bottom of
`resource` or `data` block.

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
Error: `depends_on` clause must be the last one in the definition (sheldon_depends_on)

  on template.tf line 5:
   5:   depends_on = [
   6:     aws_iam_role_policy.example
   7:   ]
```
