package custom

import (
	"fmt"
	"strings"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/config"
)

func parseConfigResource(r *config.Resource) (*Resource, error) {
	if len(r.Keys) == 0 {
		return &Resource{}, nil
	}

	keys := strings.Split(r.Keys[0], ".")
	keyAttributes := make([]string, 0, len(r.Keys))
	keyBlocks := keys[:len(keys)-1]

	for _, k := range r.Keys {
		keys = strings.Split(k, ".")
		for i := 0; i < len(keys)-1; i++ {
			if keys[i] != keyBlocks[i] {
				return nil, fmt.Errorf(
					"invalid configuration for `%s`: all keys must have the same prefix `%s`: unexpected `%s`",
					r.Kind,
					strings.Join(keyBlocks, "."),
					strings.Join(keys[:i+1], "."),
				)
			}
		}
		keyAttributes = append(keyAttributes, keys[len(keys)-1])
	}

	return &Resource{
		KeyBlocks:     keyBlocks,
		KeyAttributes: keyAttributes,
	}, nil
}
