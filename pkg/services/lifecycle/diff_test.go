package lifecycle_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strrl/kubernetes-auditing-dashboard/pkg/services/lifecycle"
)

func TestComputeDiff(t *testing.T) {
	t.Run("should compute diff between two YAML states with added fields", func(t *testing.T) {
		oldYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  key1: value1`

		newYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  key1: value1
  key2: value2`

		diff, err := lifecycle.ComputeDiff(oldYAML, newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)
		assert.NotEmpty(t, diff.Added)
		assert.Empty(t, diff.Removed)
		assert.Contains(t, diff.Added, "data.key2")
		assert.Equal(t, "value2", diff.Added["data.key2"])
	})

	t.Run("should compute diff between two YAML states with removed fields", func(t *testing.T) {
		oldYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  key1: value1
  key2: value2`

		newYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  key1: value1`

		diff, err := lifecycle.ComputeDiff(oldYAML, newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)
		assert.Empty(t, diff.Added)
		assert.NotEmpty(t, diff.Removed)
		assert.Contains(t, diff.Removed, "data.key2")
		assert.Equal(t, "value2", diff.Removed["data.key2"])
	})

	t.Run("should compute diff between two YAML states with modified fields", func(t *testing.T) {
		oldYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  key1: value1`

		newYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  key1: value2`

		diff, err := lifecycle.ComputeDiff(oldYAML, newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)
		assert.Empty(t, diff.Added)
		assert.Empty(t, diff.Removed)
		assert.NotEmpty(t, diff.Modified)
		assert.Contains(t, diff.Modified, "data.key1")
		assert.Equal(t, "value1", diff.Modified["data.key1"].OldValue)
		assert.Equal(t, "value2", diff.Modified["data.key1"].NewValue)
	})

	t.Run("should handle identical YAML states (no diff)", func(t *testing.T) {
		yaml := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  key1: value1`

		diff, err := lifecycle.ComputeDiff(yaml, yaml)
		require.NoError(t, err)
		assert.NotNil(t, diff)
		assert.Empty(t, diff.Added)
		assert.Empty(t, diff.Removed)
		assert.Empty(t, diff.Modified)
	})

	t.Run("should handle malformed YAML gracefully", func(t *testing.T) {
		malformed := `
apiVersion: v1
kind: ConfigMap
metadata
  name: test`

		valid := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test`

		_, err := lifecycle.ComputeDiff(malformed, valid)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse")
	})

	t.Run("should compute diff for deeply nested structures", func(t *testing.T) {
		oldYAML := `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        image: nginx:1.14`

		newYAML := `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        image: nginx:1.15`

		diff, err := lifecycle.ComputeDiff(oldYAML, newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)
		// The array diff is handled as a whole modification
		assert.NotEmpty(t, diff.Modified)
	})
}

func TestDiffYAMLExamples(t *testing.T) {
	t.Run("should handle Kubernetes Deployment YAML diff", func(t *testing.T) {
		oldDeployment := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webapp
  namespace: default
  resourceVersion: "12345"
  generation: 1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: webapp
  template:
    metadata:
      labels:
        app: webapp
    spec:
      containers:
      - name: webapp
        image: nginx:1.14
        ports:
        - containerPort: 80`

		newDeployment := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webapp
  namespace: default
  resourceVersion: "12346"
  generation: 2
spec:
  replicas: 5
  selector:
    matchLabels:
      app: webapp
  template:
    metadata:
      labels:
        app: webapp
        version: v2
    spec:
      containers:
      - name: webapp
        image: nginx:1.15
        ports:
        - containerPort: 80
        env:
        - name: DEBUG
          value: "true"`

		diff, err := lifecycle.ComputeDiff(oldDeployment, newDeployment)
		require.NoError(t, err)
		assert.NotNil(t, diff)

		// Should detect replicas change
		assert.Contains(t, diff.Modified, "spec.replicas")
		assert.Equal(t, float64(3), diff.Modified["spec.replicas"].OldValue)
		assert.Equal(t, float64(5), diff.Modified["spec.replicas"].NewValue)

		// Should detect new label
		assert.Contains(t, diff.Added, "spec.template.metadata.labels.version")
		assert.Equal(t, "v2", diff.Added["spec.template.metadata.labels.version"])

		// Metadata fields should be filtered out
		assert.NotContains(t, diff.Modified, "metadata.resourceVersion")
		assert.NotContains(t, diff.Modified, "metadata.generation")
	})

	t.Run("should handle ConfigMap YAML diff", func(t *testing.T) {
		oldConfigMap := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
  uid: "abc-123"
  creationTimestamp: "2023-01-01T00:00:00Z"
data:
  database.url: "postgres://localhost/mydb"
  app.debug: "false"
  app.timeout: "30"`

		newConfigMap := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
  uid: "abc-123"
  creationTimestamp: "2023-01-01T00:00:00Z"
data:
  database.url: "postgres://prod.example.com/mydb"
  app.debug: "true"
  app.timeout: "60"
  app.newfeature: "enabled"`

		diff, err := lifecycle.ComputeDiff(oldConfigMap, newConfigMap)
		require.NoError(t, err)
		assert.NotNil(t, diff)

		// Should detect modified values
		assert.Contains(t, diff.Modified, "data.database.url")
		assert.Contains(t, diff.Modified, "data.app.debug")
		assert.Contains(t, diff.Modified, "data.app.timeout")

		// Should detect added field
		assert.Contains(t, diff.Added, "data.app.newfeature")
		assert.Equal(t, "enabled", diff.Added["data.app.newfeature"])

		// UID and creationTimestamp should be filtered out
		assert.NotContains(t, diff.Modified, "metadata.uid")
		assert.NotContains(t, diff.Modified, "metadata.creationTimestamp")
	})
}

func TestDiffArrayHandling(t *testing.T) {
	t.Run("should detect changes in YAML arrays", func(t *testing.T) {
		oldYAML := `
spec:
  containers:
  - name: app
    image: nginx:1.14
  - name: sidecar
    image: proxy:1.0`

		newYAML := `
spec:
  containers:
  - name: app
    image: nginx:1.15
  - name: sidecar
    image: proxy:1.0
  - name: logger
    image: fluentd:latest`

		diff, err := lifecycle.ComputeDiff(oldYAML, newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)

		// Arrays are compared as a whole, so change should be in Modified
		assert.NotEmpty(t, diff.Modified)
	})

	t.Run("should handle array reordering", func(t *testing.T) {
		oldYAML := `
spec:
  ports:
  - port: 80
    name: http
  - port: 443
    name: https`

		newYAML := `
spec:
  ports:
  - port: 443
    name: https
  - port: 80
    name: http`

		diff, err := lifecycle.ComputeDiff(oldYAML, newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)

		// Reordered arrays are considered modified
		assert.NotEmpty(t, diff.Modified)
	})
}

func TestPartialDiffOnError(t *testing.T) {
	t.Run("should return error when YAML is malformed", func(t *testing.T) {
		validYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test`

		malformedYAML := `
apiVersion: v1
kind: ConfigMap
metadata
  name: test
  this is not valid YAML`

		// Should error on malformed old YAML
		_, err := lifecycle.ComputeDiff(malformedYAML, validYAML)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse")

		// Should error on malformed new YAML
		_, err = lifecycle.ComputeDiff(validYAML, malformedYAML)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse")
	})
}

func TestDiffPerformance(t *testing.T) {
	t.Run("should handle large YAML files efficiently", func(t *testing.T) {
		// Create a large YAML with many fields
		largeYAML1 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: large-config
data:`

		largeYAML2 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: large-config
data:`

		// Add 100 fields to each
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d", i)
			largeYAML1 += fmt.Sprintf("\n  %s: %s", key, value)
			if i < 50 {
				// First 50 unchanged, last 50 modified
				largeYAML2 += fmt.Sprintf("\n  %s: %s", key, value)
			} else {
				largeYAML2 += fmt.Sprintf("\n  %s: modified%d", key, i)
			}
		}

		diff, err := lifecycle.ComputeDiff(largeYAML1, largeYAML2)
		require.NoError(t, err)
		assert.NotNil(t, diff)

		// Should have 50 modifications
		assert.Len(t, diff.Modified, 50)
	})
}

func TestDiffEmptyStates(t *testing.T) {
	t.Run("should handle CREATE event (empty old state)", func(t *testing.T) {
		newYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: new-config
data:
  key1: value1`

		diff, err := lifecycle.ComputeDiff("", newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)
		assert.NotEmpty(t, diff.Added)
		assert.Empty(t, diff.Removed)
		assert.Empty(t, diff.Modified)

		// All fields should be in Added
		assert.Contains(t, diff.Added, "apiVersion")
		assert.Contains(t, diff.Added, "kind")
		assert.Contains(t, diff.Added, "metadata.name")
		assert.Contains(t, diff.Added, "data.key1")
	})

	t.Run("should handle DELETE event (empty new state)", func(t *testing.T) {
		oldYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: old-config
data:
  key1: value1`

		diff, err := lifecycle.ComputeDiff(oldYAML, "")
		require.NoError(t, err)
		assert.NotNil(t, diff)
		assert.Empty(t, diff.Added)
		assert.NotEmpty(t, diff.Removed)
		assert.Empty(t, diff.Modified)

		// All fields should be in Removed
		assert.Contains(t, diff.Removed, "apiVersion")
		assert.Contains(t, diff.Removed, "kind")
		assert.Contains(t, diff.Removed, "metadata.name")
		assert.Contains(t, diff.Removed, "data.key1")
	})

	t.Run("should handle both empty states", func(t *testing.T) {
		diff, err := lifecycle.ComputeDiff("", "")
		require.NoError(t, err)
		assert.NotNil(t, diff)
		assert.Empty(t, diff.Added)
		assert.Empty(t, diff.Removed)
		assert.Empty(t, diff.Modified)
	})
}

func TestDiffMetadataFiltering(t *testing.T) {
	t.Run("should filter out volatile metadata fields", func(t *testing.T) {
		oldYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
  namespace: default
  resourceVersion: "100"
  generation: 1
  uid: "abc-123"
  creationTimestamp: "2023-01-01T00:00:00Z"
  selfLink: "/api/v1/configmaps/test"
  managedFields:
  - manager: kubectl
    operation: Update
data:
  key: value1`

		newYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
  namespace: default
  resourceVersion: "101"
  generation: 2
  uid: "abc-123"
  creationTimestamp: "2023-01-01T00:00:00Z"
  selfLink: "/api/v1/configmaps/test"
  managedFields:
  - manager: kubectl
    operation: Update
    time: "2023-01-02T00:00:00Z"
data:
  key: value2`

		diff, err := lifecycle.ComputeDiff(oldYAML, newYAML)
		require.NoError(t, err)
		assert.NotNil(t, diff)

		// Should only detect the actual data change
		assert.Contains(t, diff.Modified, "data.key")
		assert.Equal(t, "value1", diff.Modified["data.key"].OldValue)
		assert.Equal(t, "value2", diff.Modified["data.key"].NewValue)

		// Volatile metadata fields should be filtered out
		assert.NotContains(t, diff.Modified, "metadata.resourceVersion")
		assert.NotContains(t, diff.Modified, "metadata.generation")
		assert.NotContains(t, diff.Modified, "metadata.uid")
		assert.NotContains(t, diff.Modified, "metadata.creationTimestamp")
		assert.NotContains(t, diff.Modified, "metadata.selfLink")
		assert.NotContains(t, diff.Modified, "metadata.managedFields")
	})
}