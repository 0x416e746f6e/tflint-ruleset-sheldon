# Rule `sheldon_sorting`

Makes sure that blocks and dictionary keys are sorted alphabetically.

- Signle-line elements are defined before the multi-line ones.
- Each consecutive chunk of single-liners is sorted individually.
- The position of special cases like `for_each`/`count`, `depends_on`, and so on
  is respected so that there is no contradictions.

## Example

```hcl
resource "something" "this" {
  template {
    metadata {
      annotations = {
        "def" = "blah-blah"
        "abc" = "yada-yada"
      }
    }
  }
}
```

```text
Error: key `abc` is out of order (should not follow alphabetically greater `def`) (sheldon_sorting)

  on case-001.tf line 6:
   6:         "abc" = "yada-yada"
```
