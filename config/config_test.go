package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_test*.json")

	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	testConfig := `{
		"targets": [
			{ "path": "~/test/path1", "depth": 3 },
			{ "path": "~/test/path2", "depth": 5 }
		],
        "selector": [
            "--border",
            "--height=50%"
        ]
	}`

	if _, err := tempFile.Write([]byte(testConfig)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tempFile.Close()

	cliCfg := Cli{Path: tempFile.Name()}

	cfg, err := Load(cliCfg)
	assert.NoError(t, err)

	assert.NotNil(t, cfg)
	assert.Len(t, cfg.Targets, 2)

	assert.Equal(t, "~/test/path1", cfg.Targets[0].Path)
	assert.Equal(t, uint8(3), cfg.Targets[0].Depth)

	assert.Equal(t, "~/test/path2", cfg.Targets[1].Path)
	assert.Equal(t, uint8(5), cfg.Targets[1].Depth)

	assert.Equal(t, "--border", cfg.Selector[0])
	assert.Equal(t, "--height=50%", cfg.Selector[1])

	// Test with non-existent file
	cliCfg.Path = "non_existent_file.json"
	_, err = Load(cliCfg)
	assert.Error(t, err)

	// Test with invalid JSON
	invalidFile, _ := os.CreateTemp("", "invalid_config_test*.json")

	defer os.Remove(invalidFile.Name())
	invalidFile.Write([]byte(`{invalid json}`))
	invalidFile.Close()

	cliCfg.Path = invalidFile.Name()
	_, err = Load(cliCfg)
	assert.Error(t, err)
}
