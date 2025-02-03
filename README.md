# TFLint Ruleset Template

<div style="text-align: right"><em>„This is my spot!“</em></div>

<br>

TFLint ruleset plugin that will check your terraform formatting from sorting and
spacing perspective.

## Requirements

- TFLint v0.41+
- Go v1.19

## Installation

You can install the plugin with `tflint --init`. Declare a config in
`.tflint.hcl` as follows:

```hcl
plugin "sheldon" {
  enabled = true

  version = "0.0.6"
  source  = "github.com/0x416e746f6e/tflint-ruleset-sheldon"

  signing_key = <<-KEY
    -----BEGIN PGP PUBLIC KEY BLOCK-----

    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
    SgvumlA=
    =ys3x
    -----END PGP PUBLIC KEY BLOCK-----
  KEY
}
```

## Building the plugin

Clone the repository locally and run the following command:

```bash
make
```

You can easily install the built plugin with the following:

```bash
make install
```

You can run the built plugin like the following:

```bash
cat << EOF > .tflint.hcl
config {
  plugin_dir = "~/.tflint.d/plugins"
}

plugin "sheldon" {
  enabled = true
}
EOF

tflint
```

Some of the resources come with their key-attributes pre-defined.  However,
their set is far from being exhaustive (and will never be). To define the
key-attributes for a resource/data blocks (so that `key_attributes`
rule picks them up) add them to in `.tflint.hcl` like follows:

```hcl
plugin "sheldon" {
  enabled = true

  resource "kubernetes_deployment" {
    key_attributes = ["metadata.namespace", "metadata.name"]
  }
}
```
