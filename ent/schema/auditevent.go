package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AuditEvent holds the schema definition for the AuditEvent entity.
type AuditEvent struct {
	ent.Schema
}

// Fields of the AuditEvent.
func (AuditEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Text("raw").NotEmpty().Immutable(),
		field.String("level").NotEmpty().Immutable(),
		field.String("auditID").NotEmpty().Immutable(),
		field.String("verb").NotEmpty().Immutable(),
		field.String("userAgent").NotEmpty().Immutable(),
		field.Time("requestTimestamp").Immutable(),
		field.Time("stageTimestamp").Immutable(),
		field.String("namespace").Immutable().Default(""),
		field.String("name").Immutable().Default(""),
		field.String("apiVersion").Immutable().Default(""),
		field.String("apiGroup").Immutable().Default(""),
		field.String("resource").Immutable().Default(""),
		field.String("subResource").Immutable().Default(""),
		field.String("stage").Immutable(),
	}
}

// Edges of the AuditEvent.
func (AuditEvent) Edges() []ent.Edge {
	return nil
}

func (AuditEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("level", "verb"),
		index.Fields("verb"),
		index.Fields("auditID"),
		index.Fields("userAgent"),
		index.Fields("requestTimestamp"),
		index.Fields("stageTimestamp"),
	}
}
