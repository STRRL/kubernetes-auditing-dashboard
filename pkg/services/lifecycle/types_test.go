package lifecycle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strrl/kubernetes-auditing-dashboard/pkg/services/lifecycle"
)

func TestResourceIdentifierParseFromURL(t *testing.T) {
	t.Run("should parse GVK from URL format apps-v1-Deployment", func(t *testing.T) {
		ri, err := lifecycle.ParseFromURL("apps-v1-Deployment", "default", "webapp")
		require.NoError(t, err)
		assert.Equal(t, "apps", ri.APIGroup)
		assert.Equal(t, "v1", ri.Version)
		assert.Equal(t, "Deployment", ri.Kind)
		assert.Equal(t, "default", ri.Namespace)
		assert.Equal(t, "webapp", ri.Name)
	})

	t.Run("should handle core resources with empty apiGroup", func(t *testing.T) {
		// Test core resource format
		ri, err := lifecycle.ParseFromURL("v1-ConfigMap", "default", "my-config")
		require.NoError(t, err)
		assert.Equal(t, "", ri.APIGroup)
		assert.Equal(t, "v1", ri.Version)
		assert.Equal(t, "ConfigMap", ri.Kind)
		assert.Equal(t, "default", ri.Namespace)
		assert.Equal(t, "my-config", ri.Name)

		// Test core-v1 format
		ri2, err := lifecycle.ParseFromURL("core-v1-ConfigMap", "default", "my-config")
		require.NoError(t, err)
		assert.Equal(t, "", ri2.APIGroup)
		assert.Equal(t, "v1", ri2.Version)
		assert.Equal(t, "ConfigMap", ri2.Kind)
	})

	t.Run("should URL decode special characters in resource names", func(t *testing.T) {
		// Test URL encoded name with dot
		ri, err := lifecycle.ParseFromURL("apps-v1-Deployment", "default", "test-app-v2%2E1")
		require.NoError(t, err)
		assert.Equal(t, "test-app-v2.1", ri.Name)

		// Test URL encoded name with space
		ri2, err := lifecycle.ParseFromURL("apps-v1-Deployment", "default", "test%20app")
		require.NoError(t, err)
		assert.Equal(t, "test app", ri2.Name)
	})

	t.Run("should validate required fields", func(t *testing.T) {
		// Test empty Kind
		_, err := lifecycle.ParseFromURL("apps-v1-", "default", "webapp")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "kind cannot be empty")

		// Test empty Version
		_, err = lifecycle.ParseFromURL("apps--Deployment", "default", "webapp")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "version cannot be empty")

		// Test empty Name
		_, err = lifecycle.ParseFromURL("apps-v1-Deployment", "default", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("should handle cluster-scoped resources without namespace", func(t *testing.T) {
		ri, err := lifecycle.ParseFromURL("v1-Namespace", "", "production")
		require.NoError(t, err)
		assert.Equal(t, "", ri.APIGroup)
		assert.Equal(t, "v1", ri.Version)
		assert.Equal(t, "Namespace", ri.Kind)
		assert.Equal(t, "", ri.Namespace)
		assert.Equal(t, "production", ri.Name)
	})

	t.Run("should parse complex GVK formats", func(t *testing.T) {
		// Test batch API
		ri, err := lifecycle.ParseFromURL("batch-v1-Job", "default", "my-job")
		require.NoError(t, err)
		assert.Equal(t, "batch", ri.APIGroup)
		assert.Equal(t, "v1", ri.Version)
		assert.Equal(t, "Job", ri.Kind)

		// Test networking.k8s.io API
		ri2, err := lifecycle.ParseFromURL("networking.k8s.io-v1-NetworkPolicy", "default", "my-policy")
		require.NoError(t, err)
		assert.Equal(t, "networking.k8s.io", ri2.APIGroup)
		assert.Equal(t, "v1", ri2.Version)
		assert.Equal(t, "NetworkPolicy", ri2.Kind)

		// Test rbac.authorization.k8s.io API
		ri3, err := lifecycle.ParseFromURL("rbac.authorization.k8s.io-v1-Role", "default", "my-role")
		require.NoError(t, err)
		assert.Equal(t, "rbac.authorization.k8s.io", ri3.APIGroup)
		assert.Equal(t, "v1", ri3.Version)
		assert.Equal(t, "Role", ri3.Kind)
	})
}

func TestResourceIdentifierToEntQuery(t *testing.T) {
	t.Run("should convert ResourceIdentifier to Ent query parameters", func(t *testing.T) {
		ri := &lifecycle.ResourceIdentifier{
			APIGroup:  "apps",
			Version:   "v1",
			Kind:      "Deployment",
			Namespace: "default",
			Name:      "webapp",
		}

		apiGroup, apiVersion, resource, namespace, name := ri.ToEntQuery()

		assert.Equal(t, "apps", apiGroup)
		assert.Equal(t, "v1", apiVersion)
		assert.Equal(t, "deployments", resource)
		assert.Equal(t, "default", namespace)
		assert.Equal(t, "webapp", name)
	})

	t.Run("should handle resource type conversion", func(t *testing.T) {
		testCases := []struct {
			kind     string
			expected string
		}{
			{"Deployment", "deployments"},
			{"ConfigMap", "configmaps"},
			{"Service", "services"},
			{"Pod", "pods"},
			{"StatefulSet", "statefulsets"},
			{"DaemonSet", "daemonsets"},
			{"ReplicaSet", "replicasets"},
			{"Namespace", "namespaces"},
			{"Node", "nodes"},
			{"PersistentVolume", "persistentvolumes"},
			{"PersistentVolumeClaim", "persistentvolumeclaims"},
			{"StorageClass", "storageclasses"},
			{"Ingress", "ingresses"},
			{"NetworkPolicy", "networkpolicies"},
			{"Role", "roles"},
			{"RoleBinding", "rolebindings"},
			{"ClusterRole", "clusterroles"},
			{"ClusterRoleBinding", "clusterrolebindings"},
			{"ServiceAccount", "serviceaccounts"},
			{"CustomResourceDefinition", "customresourcedefinitions"},
			{"PodDisruptionBudget", "poddisruptionbudgets"},
			{"Secret", "secrets"},
			{"CustomResource", "customresources"}, // Default pluralization
		}

		for _, tc := range testCases {
			t.Run(tc.kind, func(t *testing.T) {
				ri := &lifecycle.ResourceIdentifier{
					APIGroup: "test",
					Version:  "v1",
					Kind:     tc.kind,
					Name:     "test",
				}
				_, _, resource, _, _ := ri.ToEntQuery()
				assert.Equal(t, tc.expected, resource)
			})
		}
	})

	t.Run("should return only version in apiVersion field", func(t *testing.T) {
		// Test with API group - should return only version, not apiGroup/version
		ri1 := &lifecycle.ResourceIdentifier{
			APIGroup: "apps",
			Version:  "v1",
			Kind:     "Deployment",
			Name:     "test",
		}
		apiGroup1, apiVersion1, _, _, _ := ri1.ToEntQuery()
		assert.Equal(t, "apps", apiGroup1)
		assert.Equal(t, "v1", apiVersion1)

		// Test without API group (core)
		ri2 := &lifecycle.ResourceIdentifier{
			APIGroup: "",
			Version:  "v1",
			Kind:     "ConfigMap",
			Name:     "test",
		}
		apiGroup2, apiVersion2, _, _, _ := ri2.ToEntQuery()
		assert.Equal(t, "", apiGroup2)
		assert.Equal(t, "v1", apiVersion2)

		// Test beta version - should return only version
		ri3 := &lifecycle.ResourceIdentifier{
			APIGroup: "batch",
			Version:  "v1beta1",
			Kind:     "CronJob",
			Name:     "test",
		}
		apiGroup3, apiVersion3, _, _, _ := ri3.ToEntQuery()
		assert.Equal(t, "batch", apiGroup3)
		assert.Equal(t, "v1beta1", apiVersion3)
	})
}

func TestResourceIdentifierValidation(t *testing.T) {
	testCases := []struct {
		name         string
		gvk          string
		namespace    string
		resourceName string
		shouldError  bool
		errorMsg     string
	}{
		{
			name:         "valid namespaced resource",
			gvk:          "apps-v1-Deployment",
			namespace:    "default",
			resourceName: "webapp",
			shouldError:  false,
		},
		{
			name:         "valid cluster-scoped resource",
			gvk:          "v1-Namespace",
			namespace:    "",
			resourceName: "production",
			shouldError:  false,
		},
		{
			name:         "empty kind",
			gvk:          "apps-v1-",
			namespace:    "default",
			resourceName: "webapp",
			shouldError:  true,
			errorMsg:     "kind cannot be empty",
		},
		{
			name:         "empty version",
			gvk:          "apps--Deployment",
			namespace:    "default",
			resourceName: "webapp",
			shouldError:  true,
			errorMsg:     "version cannot be empty",
		},
		{
			name:         "empty name",
			gvk:          "apps-v1-Deployment",
			namespace:    "default",
			resourceName: "",
			shouldError:  true,
			errorMsg:     "name cannot be empty",
		},
		{
			name:         "malformed GVK",
			gvk:          "invalid",
			namespace:    "default",
			resourceName: "webapp",
			shouldError:  true,
			errorMsg:     "invalid GVK format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ri, err := lifecycle.ParseFromURL(tc.gvk, tc.namespace, tc.resourceName)

			if tc.shouldError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, ri)
			}
		})
	}
}

func TestResourceIdentifierSpecialCases(t *testing.T) {
	t.Run("should handle resources with dots in API group", func(t *testing.T) {
		ri, err := lifecycle.ParseFromURL("networking.k8s.io-v1-NetworkPolicy", "default", "my-policy")
		require.NoError(t, err)
		assert.Equal(t, "networking.k8s.io", ri.APIGroup)
		assert.Equal(t, "v1", ri.Version)
		assert.Equal(t, "NetworkPolicy", ri.Kind)
	})

	t.Run("should handle resources with hyphens in kind", func(t *testing.T) {
		// Most Kubernetes kinds don't have hyphens, but CRDs might
		// The parser should handle them correctly
		ri, err := lifecycle.ParseFromURL("custom.io-v1-MyCustomResource", "default", "my-resource")
		require.NoError(t, err)
		assert.Equal(t, "custom.io", ri.APIGroup)
		assert.Equal(t, "v1", ri.Version)
		assert.Equal(t, "MyCustomResource", ri.Kind)
	})

	t.Run("should handle beta and alpha versions", func(t *testing.T) {
		// Test beta version
		ri1, err := lifecycle.ParseFromURL("apps-v1beta1-Deployment", "default", "webapp")
		require.NoError(t, err)
		assert.Equal(t, "apps", ri1.APIGroup)
		assert.Equal(t, "v1beta1", ri1.Version)
		assert.Equal(t, "Deployment", ri1.Kind)

		// Test alpha version
		ri2, err := lifecycle.ParseFromURL("storage-v1alpha1-VolumeSnapshot", "default", "snapshot")
		require.NoError(t, err)
		assert.Equal(t, "storage", ri2.APIGroup)
		assert.Equal(t, "v1alpha1", ri2.Version)
		assert.Equal(t, "VolumeSnapshot", ri2.Kind)

		// Test CronJob with batch/v1beta1
		ri3, err := lifecycle.ParseFromURL("batch-v1beta1-CronJob", "default", "my-cron")
		require.NoError(t, err)
		assert.Equal(t, "batch", ri3.APIGroup)
		assert.Equal(t, "v1beta1", ri3.Version)
		assert.Equal(t, "CronJob", ri3.Kind)
	})
}

func TestMapVerbToEventType(t *testing.T) {
	testCases := []struct {
		verb     string
		expected lifecycle.EventType
	}{
		{"create", lifecycle.EventTypeCreate},
		{"CREATE", lifecycle.EventTypeCreate},
		{"update", lifecycle.EventTypeUpdate},
		{"UPDATE", lifecycle.EventTypeUpdate},
		{"patch", lifecycle.EventTypeUpdate},
		{"PATCH", lifecycle.EventTypeUpdate},
		{"delete", lifecycle.EventTypeDelete},
		{"DELETE", lifecycle.EventTypeDelete},
		{"get", lifecycle.EventTypeGet},
		{"watch", lifecycle.EventTypeUpdate},
		{"list", lifecycle.EventTypeUpdate},
		{"unknown", lifecycle.EventTypeUpdate},
	}

	for _, tc := range testCases {
		t.Run(tc.verb, func(t *testing.T) {
			result := lifecycle.MapVerbToEventType(tc.verb)
			assert.Equal(t, tc.expected, result)
		})
	}
}