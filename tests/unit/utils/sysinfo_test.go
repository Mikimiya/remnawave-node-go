package utils_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hteppl/remnawave-node-go/internal/utils"
)

func TestGetCPUCores(t *testing.T) {
	cores := utils.GetCPUCores()
	assert.Equal(t, runtime.NumCPU(), cores)
	assert.Greater(t, cores, 0)
}

func TestGetCPUModel(t *testing.T) {
	model := utils.GetCPUModel()
	assert.NotEmpty(t, model)

	if runtime.GOOS == "linux" {
		assert.NotEqual(t, "unknown", model)
		assert.Contains(t, model, "GHz")
	}
}

func TestGetTotalMemory(t *testing.T) {
	mem := utils.GetTotalMemory()
	assert.NotEmpty(t, mem)

	if runtime.GOOS == "linux" {
		assert.NotEqual(t, "unknown", mem)
	}
}

func TestPrettyBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"zero", 0, "0 B"},
		{"bytes", 500, "500 B"},
		{"kilobytes", 1000, "1 kB"},
		{"kilobytes_fractional", 1500, "1.50 kB"},
		{"megabytes", 1000000, "1 MB"},
		{"megabytes_fractional", 1500000, "1.50 MB"},
		{"gigabytes", 1000000000, "1 GB"},
		{"gigabytes_fractional", 8590000000, "8.59 GB"},
		{"sixteen_gb", 16000000000, "16 GB"},
		{"terabytes", 1000000000000, "1 TB"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.PrettyBytes(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
