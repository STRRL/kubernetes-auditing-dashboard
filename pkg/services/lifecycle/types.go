package lifecycle

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// EventType represents the type of lifecycle event
type EventType string

const (
	EventTypeCreate EventType = "CREATE"
	EventTypeUpdate EventType = "UPDATE"
	EventTypeDelete EventType = "DELETE"
	EventTypeGet    EventType = "GET"
)

// ResourceIdentifier uniquely identifies a Kubernetes resource
type ResourceIdentifier struct {
	APIGroup  string
	Version   string
	Kind      string
	Namespace string
	Name      string
}

// ParseFromURL parses resource identifier from URL segments
func ParseFromURL(gvk, namespace, name string) (*ResourceIdentifier, error) {
	// Parse GVK format: "apiGroup-version-kind" or "version-kind" for core resources
	parts := strings.Split(gvk, "-")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GVK format: %s", gvk)
	}

	ri := &ResourceIdentifier{}

	// Handle different GVK formats
	if len(parts) == 2 {
		// Core resource: "v1-ConfigMap"
		ri.APIGroup = ""
		ri.Version = parts[0]
		ri.Kind = parts[1]
	} else if len(parts) == 3 {
		// Regular: "apps-v1-Deployment" or "core-v1-ConfigMap"
		if parts[0] == "core" {
			ri.APIGroup = ""
		} else {
			ri.APIGroup = parts[0]
		}
		ri.Version = parts[1]
		ri.Kind = parts[2]
	} else {
		// Complex API group with dots: "networking.k8s.io-v1-NetworkPolicy"
		// Join all parts except the last two (version and kind)
		apiGroupParts := parts[:len(parts)-2]
		ri.APIGroup = strings.Join(apiGroupParts, ".")
		ri.Version = parts[len(parts)-2]
		ri.Kind = parts[len(parts)-1]
	}

	// Validate required fields
	if ri.Kind == "" {
		return nil, fmt.Errorf("kind cannot be empty")
	}
	if ri.Version == "" {
		return nil, fmt.Errorf("version cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	// URL decode the name to handle special characters
	decodedName, err := url.QueryUnescape(name)
	if err != nil {
		return nil, fmt.Errorf("failed to decode resource name: %w", err)
	}
	ri.Name = decodedName

	// Handle namespace: empty for cluster-scoped, or convert "_cluster" sentinel to empty
	if namespace != "" && namespace != "_cluster" {
		decodedNamespace, err := url.QueryUnescape(namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to decode namespace: %w", err)
		}
		ri.Namespace = decodedNamespace
	}
	// If namespace is "_cluster" or empty, ri.Namespace remains empty for cluster-scoped resources

	return ri, nil
}

// ToEntQuery converts ResourceIdentifier to Ent query parameters
func (ri *ResourceIdentifier) ToEntQuery() (apiGroup, apiVersion, resource, namespace, name string) {
	apiGroup = ri.APIGroup

	// Construct full apiVersion (e.g., "apps/v1" or "v1" for core resources)
	if ri.APIGroup != "" {
		apiVersion = ri.APIGroup + "/" + ri.Version
	} else {
		apiVersion = ri.Version
	}

	// Convert Kind to resource (lowercase plural)
	resource = kindToResource(ri.Kind)

	namespace = ri.Namespace
	name = ri.Name

	return
}

// kindToResource converts a Kind to a resource name (lowercase plural)
func kindToResource(kind string) string {
	// Simple pluralization rules
	// TODO: Use a proper pluralization library or mapping table
	lower := strings.ToLower(kind)

	// Handle common irregular plurals
	switch lower {
	case "configmap":
		return "configmaps"
	case "secret":
		return "secrets"
	case "service":
		return "services"
	case "deployment":
		return "deployments"
	case "statefulset":
		return "statefulsets"
	case "daemonset":
		return "daemonsets"
	case "replicaset":
		return "replicasets"
	case "pod":
		return "pods"
	case "namespace":
		return "namespaces"
	case "node":
		return "nodes"
	case "persistentvolume":
		return "persistentvolumes"
	case "persistentvolumeclaim":
		return "persistentvolumeclaims"
	case "storageclass":
		return "storageclasses"
	case "ingress":
		return "ingresses"
	case "networkpolicy":
		return "networkpolicies"
	case "poddisruptionbudget":
		return "poddisruptionbudgets"
	case "role":
		return "roles"
	case "rolebinding":
		return "rolebindings"
	case "clusterrole":
		return "clusterroles"
	case "clusterrolebinding":
		return "clusterrolebindings"
	case "serviceaccount":
		return "serviceaccounts"
	case "customresourcedefinition":
		return "customresourcedefinitions"
	default:
		// Default: just add 's'
		return lower + "s"
	}
}

// LifecycleEvent represents a processed lifecycle event for API responses
type LifecycleEvent struct {
	ID            int
	Type          EventType
	Timestamp     time.Time
	User          string
	ResourceState map[string]interface{}
	Diff          *ResourceDiff
}

// ResourceDiff represents changes between resource versions
type ResourceDiff struct {
	Added    map[string]interface{}
	Removed  map[string]interface{}
	Modified map[string]DiffEntry
}

// DiffEntry represents a single field change
type DiffEntry struct {
	OldValue interface{}
	NewValue interface{}
	Path     string
}

// MapVerbToEventType maps Kubernetes audit verb to EventType
func MapVerbToEventType(verb string) EventType {
	switch strings.ToLower(verb) {
	case "create":
		return EventTypeCreate
	case "update", "patch":
		return EventTypeUpdate
	case "delete":
		return EventTypeDelete
	case "get":
		return EventTypeGet
	default:
		return EventTypeUpdate
	}
}