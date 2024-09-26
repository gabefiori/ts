package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeConfigs(t *testing.T) {
	additionalConfigs := []string{
		"--height=75%",
		"--border-label= Custom Label ",
		"--new-option=value",
	}

	// Expected result after merging
	expectedResult := []string{
		"--height=75%",                  // Updated from default
		"--border",                      // Unchanged from default
		"--border-label= Custom Label ", // Updated from default
		"--new-option=value",            // New option added
	}

	result := mergeConfigs(additionalConfigs)
	assert.ElementsMatch(t, expectedResult, result, "Merged configurations did not match expected results.")
}
