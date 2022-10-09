package project

import "fmt"

// Name is the name of the plugin.
const Name string = "sheldon"

// Version is ruleset version.
const Version string = "0.0.1"

// ReferenceLink returns the rule reference link.
func ReferenceLink(name string) string {
	return fmt.Sprintf(
		"https://github.com/0x416e746f6e/tflint-ruleset-sheldon/blob/v%s/docs/%s_%s.md",
		Version,
		Name,
		name,
	)
}

// RuleName returns the name of the rule.
func RuleName(id string) string {
	return fmt.Sprintf(
		"%s_%s",
		Name,
		id,
	)
}
