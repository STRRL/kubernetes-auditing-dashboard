package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ResourceKind holds the schema definition for the ResourceKind entity.
type ResourceKind struct {
	ent.Schema
}

// Fields of the ResourceKind.
func (ResourceKind) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("apiVersion").NotEmpty().Unique(),
		field.Bool("namespaced").Default(true),
		field.String("kind").NotEmpty().Unique(),
	}
}

// Edges of the ResourceKind.
func (ResourceKind) Edges() []ent.Edge {
	return nil
}

func (ResourceKind) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("apiVersion", "name"),
		index.Fields("name"),
	}
}

func (ResourceKind) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
	}
}
