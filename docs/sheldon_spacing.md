# Rule `sheldon_spacing`

Normalises blank-lines in the sources.

- No multiple consecutive blank-lines.
- No unnecessary single blank-lines.
- Comments are ignored.

## Example

```hcl
resource "kubernetes_config_map" "this" {
  count = var.create_map ? 1 : 0
  metadata {
    namespace = kubernetes_namespace.this.metadata[0].name
    name      = "config-map"
  }

  data = {
    "foo" = "bar"
  }

}
```

```text
Error: attribute `count` must be separated from the rest of the definition by an extra line (sheldon_spacing)

  on case-008.tf line 2:
   2:   count = var.create_map ? 1 : 0

Error: multi-line element must be separated from the previous one by an extra line (sheldon_spacing)

  on case-008.tf line 3:
   3:   metadata {
   4:     namespace = kubernetes_namespace.this.metadata[0].name
   5:     name      = "config-map"
   6:   }

Error: 1 redundant empty line in front (sheldon_spacing)

  on case-008.tf line 12:
  12: }
```
