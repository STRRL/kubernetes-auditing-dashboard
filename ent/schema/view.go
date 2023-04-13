package schema

import "entgo.io/ent"

// View holds the schema definition for the View entity.
type View struct {
	ent.Schema
}

// Fields of the View.
func (View) Fields() []ent.Field {
	return nil
}

// Edges of the View.
func (View) Edges() []ent.Edge {
	return nil
}
