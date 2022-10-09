# Rule `key_attributes`

Makes sure that the key-attributes (those that uniquely identify a resource) are
placed on the top of the definition in the prioritised order.

## Example

```hcl
resource "kubernetes_service_account" "this" {
  metadata {
    name      = "service-account"
    namespace = kubernetes_namespace.this.metadata[0].name
  }
}
```

```text
Error: higher-priority key-attribute `namespace` should be defined before `name` (sheldon_key_attributes)

  on template.tf line 4:
   4:     namespace = kubernetes_namespace.this.metadata[0].name
```
