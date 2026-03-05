package utils

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetCPUCores() int {
	return runtime.NumCPU()
}

func GetCPUModel() string {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	defer f.Close()

	var model string
	var mhz string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if model == "" && strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				model = strings.TrimSpace(parts[1])
			}
		}
		if mhz == "" && strings.HasPrefix(line, "cpu MHz") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				mhz = strings.TrimSpace(parts[1])
			}
		}
		if model != "" && mhz != "" {
			break
		}
	}

	if model == "" {
		return "unknown"
	}

	if mhz != "" {
		if freq, err := strconv.ParseFloat(mhz, 64); err == nil {
			return fmt.Sprintf("%s/%.2f GHz", model, freq/1000)
		}
	}

	return model
}

func GetTotalMemory() string {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return "unknown"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				kb, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return parts[1] + " kB"
				}
				return PrettyBytes(kb * 1024)
			}
		}
	}
	return "unknown"
}

func PrettyBytes(bytes uint64) string {
	units := []string{"B", "kB", "MB", "GB", "TB"}
	value := float64(bytes)
	i := 0
	for value >= 1000 && i < len(units)-1 {
		value /= 1000
		i++
	}
	if value == float64(uint64(value)) {
		return fmt.Sprintf("%d %s", uint64(value), units[i])
	}
	return fmt.Sprintf("%.2f %s", value, units[i])
}
