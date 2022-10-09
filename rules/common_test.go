package rules

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/0x416e746f6e/tflint-ruleset-sheldon/config"
	"github.com/0x416e746f6e/tflint-ruleset-sheldon/custom"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const marker = "### Expected Issues ###"

func runTests(t *testing.T, r tflint.Rule) {
	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("Could not get the caller of `runTestCases` from runtime")
	}

	dir := path.Join(
		path.Dir(testFilename),
		"tests",
		strings.TrimSuffix(path.Base(testFilename), path.Ext(testFilename)),
	)

	if err := filepath.Walk(dir, func(tfFilename string, tfInfo fs.FileInfo, _err error) error {
		// Skip directories and non terraform files
		if _err != nil {
			return _err
		}
		if tfInfo.IsDir() {
			return nil
		}
		ext := path.Ext(tfFilename)
		if ext != ".tf" {
			return nil
		}

		// Get terraform content
		tfBytes, err := os.ReadFile(tfFilename)
		if err != nil {
			return fmt.Errorf("%s: %s", tfFilename, err)
		}
		tfContent := string(tfBytes)

		// Learn the expected results
		issContent := "[]"
		for _, l := range strings.Split(tfContent, "\n") {
			if strings.TrimSpace(l) == marker {
				issContent = ""
				continue
			}
			if issContent == "[]" {
				continue
			}
			issContent = issContent + strings.TrimSpace(strings.TrimPrefix(l, "# "))
		}
		iss := helper.Issues{}
		err = json.Unmarshal([]byte(issContent), &iss)
		if err != nil {
			return fmt.Errorf("%s: %s: %s", tfFilename, err, issContent)
		}

		for _, i := range iss {
			i.Rule = r
			i.Range.Filename = tfFilename
		}

		// Create a runner
		helperRunner := helper.TestRunner(
			t,
			map[string]string{tfFilename: tfContent},
		)
		runner, err := custom.NewRunner(helperRunner, config.New())
		if err != nil {
			return fmt.Errorf("%s: %s", tfFilename, err)
		}

		// Run the test
		t.Logf("Testing `%s` against `%s` rule", tfFilename, r.Name())
		if err := r.Check(runner); err != nil {
			return fmt.Errorf("%s: %s", tfFilename, err)
		}
		helper.AssertIssues(t, iss, helperRunner.Issues)

		return nil
	}); err != nil {
		t.Fatalf("Failed to run tests: %s", err)
	}
}
