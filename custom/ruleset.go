package custom

import (
	"fmt"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/config"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// RuleSet is the custom ruleset.
type RuleSet struct {
	tflint.BuiltinRuleSet
	config *config.Config
}

// ConfigSchema returns the ruleset plugin config schema.
func (r *RuleSet) ConfigSchema() *hclext.BodySchema {
	r.config = config.New()
	return hclext.ImpliedBodySchema(r.config)
}

// ApplyConfig applies the configuration to the ruleset.
func (r *RuleSet) ApplyConfig(body *hclext.BodyContent) error {
	predefinedResources := r.config.Resources
	r.config.Resources = make([]*config.Resource, 0)

	diags := hclext.DecodeBody(body, nil, r.config)
	if diags.HasErrors() {
		return diags
	}

	r.config.Resources = append(predefinedResources, r.config.Resources...)

	return nil
}

// Check runs inspection for each rule by applying Runner.
func (r *RuleSet) Check(runner tflint.Runner) error {
	rr, err := NewRunner(runner, r.config)
	if err != nil {
		return err
	}

	for _, rule := range r.Rules {
		if !rule.Enabled() {
			rr.Disabled[rule.Name()] = true
		}
	}

	for _, rule := range r.EnabledRules {
		if err := rule.Check(rr); err != nil {
			return fmt.Errorf(
				"failed to check `%s` rule, aborting: %s",
				rule.Name(),
				err,
			)
		}
	}

	return nil
}
