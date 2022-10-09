package custom

import (
	"testing"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/config"
	"github.com/stretchr/testify/assert"
)

func TestParseConfigResource(t *testing.T) {
	res, err := parseConfigResource(&config.Resource{
		Kind: "kubernetes_manifest",
		Keys: []string{"manifest.metadata.namespace", "manifest.metadata.name"},
	})
	assert.NoError(t, err)
	assert.Equal(t, res.KeyBlocks, []string{"manifest", "metadata"})
	assert.Equal(t, res.KeyAttributes, []string{"namespace", "name"})
}
