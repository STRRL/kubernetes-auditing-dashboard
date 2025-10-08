package lifecycle

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"sigs.k8s.io/yaml"
)

// ComputeDiff computes the difference between two YAML representations of a resource
func ComputeDiff(oldYAML, newYAML string) (*ResourceDiff, error) {
	// Handle empty states
	if oldYAML == "" && newYAML == "" {
		return &ResourceDiff{
			Added:    make(map[string]interface{}),
			Removed:  make(map[string]interface{}),
			Modified: make(map[string]DiffEntry),
		}, nil
	}

	// Parse YAML into maps
	var oldObj map[string]interface{}
	var newObj map[string]interface{}

	// Parse old YAML if not empty
	if oldYAML != "" {
		if err := yaml.Unmarshal([]byte(oldYAML), &oldObj); err != nil {
			return nil, fmt.Errorf("failed to parse old YAML: %w", err)
		}
	}

	// Parse new YAML if not empty
	if newYAML != "" {
		if err := yaml.Unmarshal([]byte(newYAML), &newObj); err != nil {
			return nil, fmt.Errorf("failed to parse new YAML: %w", err)
		}
	}

	// Filter metadata fields
	if oldObj != nil {
		filterMetadata(oldObj)
	}
	if newObj != nil {
		filterMetadata(newObj)
	}

	// Compute the diff
	diff := &ResourceDiff{
		Added:    make(map[string]interface{}),
		Removed:  make(map[string]interface{}),
		Modified: make(map[string]DiffEntry),
	}

	// Handle CREATE case (old is empty)
	if oldObj == nil && newObj != nil {
		diff.Added = flattenMap("", newObj)
		return diff, nil
	}

	// Handle DELETE case (new is empty)
	if oldObj != nil && newObj == nil {
		diff.Removed = flattenMap("", oldObj)
		return diff, nil
	}

	// Compute differences between two non-empty states
	computeMapDiff("", oldObj, newObj, diff)

	return diff, nil
}

// filterMetadata removes volatile metadata fields that should not be included in diffs
func filterMetadata(obj map[string]interface{}) {
	// Remove volatile metadata fields
	if metadata, ok := obj["metadata"].(map[string]interface{}); ok {
		delete(metadata, "resourceVersion")
		delete(metadata, "generation")
		delete(metadata, "uid")
		delete(metadata, "creationTimestamp")
		delete(metadata, "selfLink")
		delete(metadata, "managedFields")

		// Clean up status fields that are often server-managed
		if len(metadata) == 0 {
			delete(obj, "metadata")
		}
	}

	// Remove status for most resources as it's server-managed
	delete(obj, "status")
}

// computeMapDiff recursively computes differences between two maps
func computeMapDiff(path string, oldMap, newMap map[string]interface{}, diff *ResourceDiff) {
	// Track which keys we've processed
	processedKeys := make(map[string]bool)

	// Check for removed and modified fields
	for key, oldVal := range oldMap {
		processedKeys[key] = true
		fieldPath := buildPath(path, key)

		if newVal, exists := newMap[key]; exists {
			// Field exists in both - check if modified
			if !deepEqual(oldVal, newVal) {
				// Check if both values are maps for recursive diff
				if oldMapVal, oldIsMap := oldVal.(map[string]interface{}); oldIsMap {
					if newMapVal, newIsMap := newVal.(map[string]interface{}); newIsMap {
						// Recursively compute diff for nested maps
						computeMapDiff(fieldPath, oldMapVal, newMapVal, diff)
						continue
					}
				}

				// Check if both values are arrays
				if oldArr, oldIsArr := oldVal.([]interface{}); oldIsArr {
					if newArr, newIsArr := newVal.([]interface{}); newIsArr {
						// Handle array differences
						if !arraysEqual(oldArr, newArr) {
							diff.Modified[fieldPath] = DiffEntry{
								Path:     fieldPath,
								OldValue: oldVal,
								NewValue: newVal,
							}
						}
						continue
					}
				}

				// Simple value change
				diff.Modified[fieldPath] = DiffEntry{
					Path:     fieldPath,
					OldValue: oldVal,
					NewValue: newVal,
				}
			}
		} else {
			// Field was removed
			diff.Removed[fieldPath] = oldVal
		}
	}

	// Check for added fields
	for key, newVal := range newMap {
		if !processedKeys[key] {
			fieldPath := buildPath(path, key)
			diff.Added[fieldPath] = newVal
		}
	}
}

// buildPath constructs a dot-separated path
func buildPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}

// flattenMap flattens a nested map into a single-level map with dot-separated keys
func flattenMap(prefix string, m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range m {
		fullKey := buildPath(prefix, key)

		// Check if value is a nested map
		if nestedMap, isMap := value.(map[string]interface{}); isMap {
			// Recursively flatten nested maps
			flattened := flattenMap(fullKey, nestedMap)
			for k, v := range flattened {
				result[k] = v
			}
		} else {
			// Store the value directly
			result[fullKey] = value
		}
	}

	return result
}

// deepEqual performs a deep comparison of two values
func deepEqual(a, b interface{}) bool {
	// Use json marshaling for consistent comparison
	// This handles different numeric types that YAML might produce
	aJSON, errA := json.Marshal(a)
	bJSON, errB := json.Marshal(b)

	if errA != nil || errB != nil {
		// Fallback to reflect.DeepEqual if JSON marshaling fails
		return reflect.DeepEqual(a, b)
	}

	return string(aJSON) == string(bJSON)
}

// arraysEqual compares two arrays for equality
func arraysEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !deepEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}

// ParseYAML parses YAML string into a map
func ParseYAML(yamlStr string) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Handle empty string
	if strings.TrimSpace(yamlStr) == "" {
		return make(map[string]interface{}), nil
	}

	if err := yaml.Unmarshal([]byte(yamlStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return result, nil
}