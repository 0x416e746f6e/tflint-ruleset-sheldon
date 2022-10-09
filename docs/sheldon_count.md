# Rule `sheldon_count`

Makes sure that `count` meta-attribute is placed at the top of `resource` or
`data` block.

## Example

```hcl
resource "random_password" "this" {
  length           = 16
  override_special = "!#$%&*()-_=+[]{}<>:?"
  special          = true

  count = 10
}
```

```text
Error: `count` must be the top-most attribute (sheldon_count)

  on template.tf line 6:
   6:   count = 10
```
