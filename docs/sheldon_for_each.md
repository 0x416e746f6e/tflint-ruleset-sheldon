# Rule `for_each`

Makes sure that `for_each` meta-attribute is placed at the top of `resource` or
`data` block.

## Example

```hcl
resource "google_container_registry" "this" {
  project  = "foo"
  location = each.key
  for_each = var.locations
}
```

```text
Error: `for_each` must be the top-most attribute (sheldon_for_each)

  on template.tf line 4:
   4:   for_each = ["a", "b"]
```
