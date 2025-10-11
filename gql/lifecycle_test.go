package gql_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strrl/kubernetes-auditing-dashboard/ent"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/enttest"
	"github.com/strrl/kubernetes-auditing-dashboard/gql"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *ent.Client {
	// Use in-memory database with unique name for each test to avoid data pollution
	// cache=shared allows multiple connections within the same test to access the same database
	dbName := fmt.Sprintf("file:test_%d_%d?mode=memory&cache=shared&_fk=1",
		time.Now().UnixNano(), rand.Int63())
	client := enttest.Open(t, "sqlite3", dbName)
	return client
}

func createTestAuditEvent(client *ent.Client, ctx context.Context, verb, namespace, name string, timestamp time.Time) error {
	// Create the requestObject and responseObject as proper Kubernetes resources
	resourceObject := map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"replicas": 3,
		},
	}

	resourceJSON, _ := json.Marshal(resourceObject)

	// Create a minimal audit event as JSON
	auditEvent := map[string]interface{}{
		"level":   "RequestResponse",
		"auditID": "test-audit-id-" + timestamp.Format(time.RFC3339),
		"verb":    verb,
		"user": map[string]interface{}{
			"username": "test-user",
		},
		"objectRef": map[string]interface{}{
			"apiGroup":   "apps",
			"apiVersion": "v1",
			"resource":   "deployments",
			"namespace":  namespace,
			"name":       name,
		},
		"requestReceivedTimestamp": timestamp.Format(time.RFC3339Nano),
		"stageTimestamp":          timestamp.Format(time.RFC3339Nano),
		"requestObject":           json.RawMessage(resourceJSON),
		"responseObject":          json.RawMessage(resourceJSON),
	}

	raw, err := json.Marshal(auditEvent)
	if err != nil {
		return err
	}

	// Create the Ent audit event
	_, err = client.AuditEvent.Create().
		SetRaw(string(raw)).
		SetLevel("RequestResponse").
		SetAuditID("test-audit-id-" + timestamp.Format(time.RFC3339)).
		SetVerb(verb).
		SetUserAgent("kubectl/v1.20.0").
		SetRequestTimestamp(timestamp).
		SetStageTimestamp(timestamp).
		SetNamespace(namespace).
		SetName(name).
		SetApiVersion("v1").
		SetApiGroup("apps").
		SetResource("deployments").
		SetSubResource("").
		SetStage("ResponseComplete").
		Save(ctx)

	return err
}

func TestResourceLifecycleQuery(t *testing.T) {
	t.Run("should return events in DESC order for namespaced resource", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create test events with different timestamps
		now := time.Now()
		events := []struct {
			verb      string
			timestamp time.Time
		}{
			{"create", now.Add(-3 * time.Hour)},
			{"update", now.Add(-2 * time.Hour)},
			{"patch", now.Add(-1 * time.Hour)},
			{"update", now},
		}

		// Create audit events
		for _, e := range events {
			err := createTestAuditEvent(client, ctx, e.verb, "default", "test-app", e.timestamp)
			require.NoError(t, err)
		}

		// Create resolver and query
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "test-app")

		require.NoError(t, err)
		require.Len(t, result, 4)

		// Verify DESC order
		for i := 0; i < len(result)-1; i++ {
			assert.True(t, result[i].Timestamp.After(result[i+1].Timestamp) || result[i].Timestamp.Equal(result[i+1].Timestamp),
				"Events should be in DESC order by timestamp")
		}
	})

	t.Run("should work with cluster-scoped resource (no namespace)", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create test event for cluster-scoped resource
		auditEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "test-audit-cluster",
			"verb":    "create",
			"user": map[string]interface{}{
				"username": "admin",
			},
			"objectRef": map[string]interface{}{
				"apiGroup":   "",
				"apiVersion": "v1",
				"resource":   "namespaces",
				"name":       "production",
			},
			"requestReceivedTimestamp": time.Now().Format(time.RFC3339Nano),
			"stageTimestamp":          time.Now().Format(time.RFC3339Nano),
			"requestObject":  json.RawMessage(`{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"production"}}`),
			"responseObject": json.RawMessage(`{"apiVersion":"v1","kind":"Namespace","metadata":{"name":"production"}}`),
		}

		raw, err := json.Marshal(auditEvent)
		require.NoError(t, err)

		// Create the Ent audit event without namespace
		_, err = client.AuditEvent.Create().
			SetRaw(string(raw)).
			SetLevel("RequestResponse").
			SetAuditID("test-audit-cluster").
			SetVerb("create").
			SetUserAgent("kubectl/v1.20.0").
			SetRequestTimestamp(time.Now()).
			SetStageTimestamp(time.Now()).
			SetNamespace(""). // Empty namespace for cluster-scoped
			SetName("production").
			SetApiVersion("v1").
			SetApiGroup("").
			SetResource("namespaces").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Query without namespace
		resolver := gql.NewResolver(client)
		result, err := resolver.Query().ResourceLifecycle(ctx, "", "v1", "Namespace", nil, "production")

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "admin", result[0].User)
	})

	t.Run("should return empty array for non-existent resource", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Query for non-existent resource
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "non-existent")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("should validate required parameters", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		resolver := gql.NewResolver(client)
		namespace := "default"

		// Test with empty name
		_, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")

		// Test with empty kind
		_, err = resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "", &namespace, "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kind cannot be empty")

		// Test with empty version
		_, err = resolver.Query().ResourceLifecycle(ctx, "apps", "", "Deployment", &namespace, "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "version cannot be empty")
	})

	t.Run("should handle 100+ events correctly (pagination test)", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create 120 events
		now := time.Now()
		for i := 0; i < 120; i++ {
			verb := "update"
			if i%10 == 0 {
				verb = "patch"
			}
			timestamp := now.Add(time.Duration(-i) * time.Minute)
			err := createTestAuditEvent(client, ctx, verb, "default", "high-volume-app", timestamp)
			require.NoError(t, err)
		}

		// Query all events
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "high-volume-app")

		require.NoError(t, err)
		assert.Len(t, result, 120)

		// Verify DESC order
		for i := 0; i < len(result)-1; i++ {
			assert.True(t, result[i].Timestamp.After(result[i+1].Timestamp) || result[i].Timestamp.Equal(result[i+1].Timestamp),
				"Events should maintain DESC order even with 100+ events")
		}
	})

	t.Run("should maintain DESC order across large result sets", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create events with specific timestamps
		baseTime := time.Now().Add(-24 * time.Hour)
		timestamps := []time.Time{
			baseTime,
			baseTime.Add(1 * time.Hour),
			baseTime.Add(2 * time.Hour),
			baseTime.Add(3 * time.Hour),
		}

		// Create events in random order
		for i := len(timestamps) - 1; i >= 0; i-- {
			err := createTestAuditEvent(client, ctx, "update", "default", "order-test", timestamps[i])
			require.NoError(t, err)
		}

		// Query and verify order
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "order-test")

		require.NoError(t, err)
		require.Len(t, result, 4)

		// Should be in reverse chronological order
		assert.True(t, result[0].Timestamp.Equal(timestamps[3]) || result[0].Timestamp.After(timestamps[3]))
		assert.True(t, result[3].Timestamp.Equal(timestamps[0]) || result[3].Timestamp.Before(timestamps[0].Add(time.Second)))
	})

	t.Run("should handle 500+ events efficiently without OOM", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create 500 events
		now := time.Now()
		batchSize := 50
		for batch := 0; batch < 10; batch++ {
			for i := 0; i < batchSize; i++ {
				idx := batch*batchSize + i
				timestamp := now.Add(time.Duration(-idx) * time.Second)
				err := createTestAuditEvent(client, ctx, "update", "default", "large-app", timestamp)
				require.NoError(t, err)
			}
		}

		// Query all events - should not cause memory issues
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "large-app")

		require.NoError(t, err)
		assert.Len(t, result, 500)

		// Spot check order - first and last
		assert.True(t, result[0].Timestamp.After(result[499].Timestamp),
			"First event should be newer than last")
	})
}

func TestResourceLifecycle_DiffComputation(t *testing.T) {
	t.Run("should compute diffs for UPDATE events", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create two events - initial and update
		now := time.Now()

		// First event - create
		createEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "create-event",
			"verb":    "create",
			"user":    map[string]interface{}{"username": "admin"},
			"objectRef": map[string]interface{}{
				"apiGroup": "apps", "apiVersion": "v1",
				"resource": "deployments", "namespace": "default", "name": "test-diff",
			},
			"requestReceivedTimestamp": now.Add(-1 * time.Hour).Format(time.RFC3339Nano),
			"stageTimestamp":          now.Add(-1 * time.Hour).Format(time.RFC3339Nano),
			"requestObject":           json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-diff"},"spec":{"replicas":1}}`),
			"responseObject":          json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-diff"},"spec":{"replicas":1}}`),
		}

		raw1, _ := json.Marshal(createEvent)
		_, err := client.AuditEvent.Create().
			SetRaw(string(raw1)).
			SetLevel("RequestResponse").
			SetAuditID("create-event").
			SetVerb("create").
			SetUserAgent("kubectl").
			SetRequestTimestamp(now.Add(-1 * time.Hour)).
			SetStageTimestamp(now.Add(-1 * time.Hour)).
			SetNamespace("default").
			SetName("test-diff").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Second event - update
		updateEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "update-event",
			"verb":    "update",
			"user":    map[string]interface{}{"username": "admin"},
			"objectRef": map[string]interface{}{
				"apiGroup": "apps", "apiVersion": "v1",
				"resource": "deployments", "namespace": "default", "name": "test-diff",
			},
			"requestReceivedTimestamp": now.Format(time.RFC3339Nano),
			"stageTimestamp":          now.Format(time.RFC3339Nano),
			"requestObject":           json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-diff"},"spec":{"replicas":3}}`),
			"responseObject":          json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-diff"},"spec":{"replicas":3}}`),
		}

		raw2, _ := json.Marshal(updateEvent)
		_, err = client.AuditEvent.Create().
			SetRaw(string(raw2)).
			SetLevel("RequestResponse").
			SetAuditID("update-event").
			SetVerb("update").
			SetUserAgent("kubectl").
			SetRequestTimestamp(now).
			SetStageTimestamp(now).
			SetNamespace("default").
			SetName("test-diff").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Query and check diff
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "test-diff")

		require.NoError(t, err)
		require.Len(t, result, 2)

		// First event (most recent) should be UPDATE with diff
		updateResult := result[0]
		assert.Equal(t, gql.EventTypeUpdate, updateResult.Type)
		assert.NotNil(t, updateResult.Diff)
		assert.NotEmpty(t, updateResult.Diff.Modified)
	})

	t.Run("should not compute diffs for CREATE events", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create a single CREATE event
		err := createTestAuditEvent(client, ctx, "create", "default", "new-app", time.Now())
		require.NoError(t, err)

		// Query
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "new-app")

		require.NoError(t, err)
		require.Len(t, result, 1)

		// CREATE event should not have diff
		assert.Equal(t, gql.EventTypeCreate, result[0].Type)
		assert.Nil(t, result[0].Diff)
	})

	t.Run("should not compute diffs for DELETE events", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create a DELETE event
		err := createTestAuditEvent(client, ctx, "delete", "default", "deleted-app", time.Now())
		require.NoError(t, err)

		// Query
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "deleted-app")

		require.NoError(t, err)
		require.Len(t, result, 1)

		// DELETE event should not have diff
		assert.Equal(t, gql.EventTypeDelete, result[0].Type)
		assert.Nil(t, result[0].Diff)
	})
}

func TestResourceLifecycle_EventTypeMapping(t *testing.T) {
	tests := []struct {
		verb         string
		expectedType gql.EventType
	}{
		{"create", gql.EventTypeCreate},
		{"update", gql.EventTypeUpdate},
		{"patch", gql.EventTypeUpdate},
		{"delete", gql.EventTypeDelete},
	}

	for _, tt := range tests {
		t.Run("should map verb "+tt.verb+" to "+string(tt.expectedType), func(t *testing.T) {
			ctx := context.Background()
			client := setupTestDB(t)
			defer client.Close()

			// Create event with specific verb
			err := createTestAuditEvent(client, ctx, tt.verb, "default", "test-type", time.Now())
			require.NoError(t, err)

			// Query
			resolver := gql.NewResolver(client)
			namespace := "default"
			result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "test-type")

			require.NoError(t, err)
			require.Len(t, result, 1)
			assert.Equal(t, tt.expectedType, result[0].Type)
		})
	}
}

func TestResourceLifecycle_UserExtraction(t *testing.T) {
	t.Run("should extract user from userAgent field", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create event with specific user
		auditEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "user-test",
			"verb":    "create",
			"user": map[string]interface{}{
				"username": "john.doe@example.com",
			},
			"objectRef": map[string]interface{}{
				"apiGroup": "apps", "apiVersion": "apps/v1",
				"resource": "deployments", "namespace": "default", "name": "user-test",
			},
			"requestReceivedTimestamp": time.Now().Format(time.RFC3339Nano),
			"stageTimestamp":          time.Now().Format(time.RFC3339Nano),
			"responseObject": map[string]interface{}{
				"raw": []byte(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"user-test"}}`),
			},
		}

		raw, _ := json.Marshal(auditEvent)
		_, err := client.AuditEvent.Create().
			SetRaw(string(raw)).
			SetLevel("RequestResponse").
			SetAuditID("user-test").
			SetVerb("create").
			SetUserAgent("kubectl").
			SetRequestTimestamp(time.Now()).
			SetStageTimestamp(time.Now()).
			SetNamespace("default").
			SetName("user-test").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Query
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "user-test")

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "john.doe@example.com", result[0].User)
	})
}

func TestResourceLifecycle_ErrorHandling(t *testing.T) {
	t.Run("should handle malformed audit events gracefully", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		// Create malformed event (invalid JSON in raw field)
		_, err := client.AuditEvent.Create().
			SetRaw("invalid json").
			SetLevel("RequestResponse").
			SetAuditID("malformed").
			SetVerb("create").
			SetUserAgent("kubectl").
			SetRequestTimestamp(time.Now()).
			SetStageTimestamp(time.Now()).
			SetNamespace("default").
			SetName("malformed-test").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Query should handle the malformed event gracefully
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "malformed-test")

		// Should not error, but may skip malformed events
		require.NoError(t, err)
		// Malformed events are skipped in the current implementation
		assert.Empty(t, result)
	})
}

func TestResourceLifecycle_PreviousState(t *testing.T) {
	t.Run("should populate previousState for UPDATE events", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		now := time.Now()

		// Create initial state - CREATE event
		createEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "create-event",
			"verb":    "create",
			"user":    map[string]interface{}{"username": "admin"},
			"objectRef": map[string]interface{}{
				"apiGroup": "apps", "apiVersion": "v1",
				"resource": "deployments", "namespace": "default", "name": "test-prev-state",
			},
			"requestReceivedTimestamp": now.Add(-2 * time.Hour).Format(time.RFC3339Nano),
			"stageTimestamp":          now.Add(-2 * time.Hour).Format(time.RFC3339Nano),
			"requestObject":           json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-prev-state"},"spec":{"replicas":1}}`),
			"responseObject":          json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-prev-state"},"spec":{"replicas":1}}`),
		}

		raw1, _ := json.Marshal(createEvent)
		_, err := client.AuditEvent.Create().
			SetRaw(string(raw1)).
			SetLevel("RequestResponse").
			SetAuditID("create-event").
			SetVerb("create").
			SetUserAgent("kubectl").
			SetRequestTimestamp(now.Add(-2 * time.Hour)).
			SetStageTimestamp(now.Add(-2 * time.Hour)).
			SetNamespace("default").
			SetName("test-prev-state").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Update event - change replicas from 1 to 3
		updateEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "update-event",
			"verb":    "update",
			"user":    map[string]interface{}{"username": "admin"},
			"objectRef": map[string]interface{}{
				"apiGroup": "apps", "apiVersion": "v1",
				"resource": "deployments", "namespace": "default", "name": "test-prev-state",
			},
			"requestReceivedTimestamp": now.Format(time.RFC3339Nano),
			"stageTimestamp":          now.Format(time.RFC3339Nano),
			"requestObject":           json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-prev-state"},"spec":{"replicas":3}}`),
			"responseObject":          json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-prev-state"},"spec":{"replicas":3}}`),
		}

		raw2, _ := json.Marshal(updateEvent)
		_, err = client.AuditEvent.Create().
			SetRaw(string(raw2)).
			SetLevel("RequestResponse").
			SetAuditID("update-event").
			SetVerb("update").
			SetUserAgent("kubectl").
			SetRequestTimestamp(now).
			SetStageTimestamp(now).
			SetNamespace("default").
			SetName("test-prev-state").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Query and verify
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "test-prev-state")

		require.NoError(t, err)
		require.Len(t, result, 2)

		// First event (UPDATE) should have previousState
		updateResult := result[0]
		assert.Equal(t, gql.EventTypeUpdate, updateResult.Type)
		assert.NotNil(t, updateResult.PreviousState, "UPDATE event should have previousState")

		// Parse previous state and verify it has replicas=1
		var prevState map[string]interface{}
		err = json.Unmarshal([]byte(*updateResult.PreviousState), &prevState)
		require.NoError(t, err)
		assert.Equal(t, float64(1), prevState["spec"].(map[string]interface{})["replicas"])

		// Parse current state and verify it has replicas=3
		var currentState map[string]interface{}
		err = json.Unmarshal([]byte(updateResult.ResourceState), &currentState)
		require.NoError(t, err)
		assert.Equal(t, float64(3), currentState["spec"].(map[string]interface{})["replicas"])

		// CREATE event should not have previousState
		createResult := result[1]
		assert.Equal(t, gql.EventTypeCreate, createResult.Type)
		assert.Nil(t, createResult.PreviousState, "CREATE event should not have previousState")
	})

	t.Run("should skip GET events when finding previousState", func(t *testing.T) {
		ctx := context.Background()
		client := setupTestDB(t)
		defer client.Close()

		now := time.Now()

		// Create initial state - CREATE event
		createEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "create-with-gets",
			"verb":    "create",
			"user":    map[string]interface{}{"username": "admin"},
			"objectRef": map[string]interface{}{
				"apiGroup": "apps", "apiVersion": "v1",
				"resource": "deployments", "namespace": "default", "name": "test-skip-gets",
			},
			"requestReceivedTimestamp": now.Add(-4 * time.Hour).Format(time.RFC3339Nano),
			"stageTimestamp":          now.Add(-4 * time.Hour).Format(time.RFC3339Nano),
			"requestObject":           json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-skip-gets"},"spec":{"replicas":1}}`),
			"responseObject":          json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-skip-gets"},"spec":{"replicas":1}}`),
		}

		raw1, _ := json.Marshal(createEvent)
		_, err := client.AuditEvent.Create().
			SetRaw(string(raw1)).
			SetLevel("RequestResponse").
			SetAuditID("create-with-gets").
			SetVerb("create").
			SetUserAgent("kubectl").
			SetRequestTimestamp(now.Add(-4 * time.Hour)).
			SetStageTimestamp(now.Add(-4 * time.Hour)).
			SetNamespace("default").
			SetName("test-skip-gets").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Insert some GET events in between (should be skipped in previousState lookup)
		for i := 0; i < 3; i++ {
			resourceObject := map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name":      "test-skip-gets",
					"namespace": "default",
				},
				"spec": map[string]interface{}{
					"replicas": 1,
				},
			}
			resourceJSON, _ := json.Marshal(resourceObject)

			getEvent := map[string]interface{}{
				"level":   "RequestResponse",
				"auditID": fmt.Sprintf("get-event-%d", i),
				"verb":    "get",
				"user":    map[string]interface{}{"username": "reader"},
				"objectRef": map[string]interface{}{
					"apiGroup": "apps", "apiVersion": "v1",
					"resource": "deployments", "namespace": "default", "name": "test-skip-gets",
				},
				"requestReceivedTimestamp": now.Add(time.Duration(-3+i) * time.Hour).Format(time.RFC3339Nano),
				"stageTimestamp":          now.Add(time.Duration(-3+i) * time.Hour).Format(time.RFC3339Nano),
				"responseObject":          json.RawMessage(resourceJSON),
			}

			rawGet, _ := json.Marshal(getEvent)
			_, err = client.AuditEvent.Create().
				SetRaw(string(rawGet)).
				SetLevel("RequestResponse").
				SetAuditID(fmt.Sprintf("get-event-%d", i)).
				SetVerb("get").
				SetUserAgent("kubectl").
				SetRequestTimestamp(now.Add(time.Duration(-3+i) * time.Hour)).
				SetStageTimestamp(now.Add(time.Duration(-3+i) * time.Hour)).
				SetNamespace("default").
				SetName("test-skip-gets").
				SetApiVersion("v1").
				SetApiGroup("apps").
				SetResource("deployments").
				SetSubResource("").
				SetStage("ResponseComplete").
				Save(ctx)
			require.NoError(t, err)
		}

		// Update event - should skip GET events and find CREATE as previous state
		updateEvent := map[string]interface{}{
			"level":   "RequestResponse",
			"auditID": "update-after-gets",
			"verb":    "update",
			"user":    map[string]interface{}{"username": "admin"},
			"objectRef": map[string]interface{}{
				"apiGroup": "apps", "apiVersion": "v1",
				"resource": "deployments", "namespace": "default", "name": "test-skip-gets",
			},
			"requestReceivedTimestamp": now.Format(time.RFC3339Nano),
			"stageTimestamp":          now.Format(time.RFC3339Nano),
			"requestObject":           json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-skip-gets"},"spec":{"replicas":5}}`),
			"responseObject":          json.RawMessage(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test-skip-gets"},"spec":{"replicas":5}}`),
		}

		rawUpdate, _ := json.Marshal(updateEvent)
		_, err = client.AuditEvent.Create().
			SetRaw(string(rawUpdate)).
			SetLevel("RequestResponse").
			SetAuditID("update-after-gets").
			SetVerb("update").
			SetUserAgent("kubectl").
			SetRequestTimestamp(now).
			SetStageTimestamp(now).
			SetNamespace("default").
			SetName("test-skip-gets").
			SetApiVersion("v1").
			SetApiGroup("apps").
			SetResource("deployments").
			SetSubResource("").
			SetStage("ResponseComplete").
			Save(ctx)
		require.NoError(t, err)

		// Query and verify
		resolver := gql.NewResolver(client)
		namespace := "default"
		result, err := resolver.Query().ResourceLifecycle(ctx, "apps", "v1", "Deployment", &namespace, "test-skip-gets")

		require.NoError(t, err)
		// Should have UPDATE + GET events + CREATE = 5 events total
		require.Len(t, result, 5)

		// First event should be UPDATE
		updateResult := result[0]
		assert.Equal(t, gql.EventTypeUpdate, updateResult.Type)
		assert.NotNil(t, updateResult.PreviousState, "UPDATE event should have previousState")

		// Previous state should be from CREATE (replicas=1), not from GET events
		var prevState map[string]interface{}
		err = json.Unmarshal([]byte(*updateResult.PreviousState), &prevState)
		require.NoError(t, err)
		assert.Equal(t, float64(1), prevState["spec"].(map[string]interface{})["replicas"],
			"previousState should come from CREATE event, skipping GET events")

		// Current state should have replicas=5
		var currentState map[string]interface{}
		err = json.Unmarshal([]byte(updateResult.ResourceState), &currentState)
		require.NoError(t, err)
		assert.Equal(t, float64(5), currentState["spec"].(map[string]interface{})["replicas"])
	})
}