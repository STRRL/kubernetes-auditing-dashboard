// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/resourcekind"
)

// ResourceKind is the model entity for the ResourceKind schema.
type ResourceKind struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// ApiVersion holds the value of the "apiVersion" field.
	ApiVersion string `json:"apiVersion,omitempty"`
	// Namespaced holds the value of the "namespaced" field.
	Namespaced bool `json:"namespaced,omitempty"`
	// Kind holds the value of the "kind" field.
	Kind string `json:"kind,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ResourceKind) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case resourcekind.FieldNamespaced:
			values[i] = new(sql.NullBool)
		case resourcekind.FieldID:
			values[i] = new(sql.NullInt64)
		case resourcekind.FieldName, resourcekind.FieldApiVersion, resourcekind.FieldKind:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ResourceKind", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ResourceKind fields.
func (rk *ResourceKind) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case resourcekind.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			rk.ID = int(value.Int64)
		case resourcekind.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				rk.Name = value.String
			}
		case resourcekind.FieldApiVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field apiVersion", values[i])
			} else if value.Valid {
				rk.ApiVersion = value.String
			}
		case resourcekind.FieldNamespaced:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field namespaced", values[i])
			} else if value.Valid {
				rk.Namespaced = value.Bool
			}
		case resourcekind.FieldKind:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kind", values[i])
			} else if value.Valid {
				rk.Kind = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this ResourceKind.
// Note that you need to call ResourceKind.Unwrap() before calling this method if this ResourceKind
// was returned from a transaction, and the transaction was committed or rolled back.
func (rk *ResourceKind) Update() *ResourceKindUpdateOne {
	return NewResourceKindClient(rk.config).UpdateOne(rk)
}

// Unwrap unwraps the ResourceKind entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (rk *ResourceKind) Unwrap() *ResourceKind {
	_tx, ok := rk.config.driver.(*txDriver)
	if !ok {
		panic("ent: ResourceKind is not a transactional entity")
	}
	rk.config.driver = _tx.drv
	return rk
}

// String implements the fmt.Stringer.
func (rk *ResourceKind) String() string {
	var builder strings.Builder
	builder.WriteString("ResourceKind(")
	builder.WriteString(fmt.Sprintf("id=%v, ", rk.ID))
	builder.WriteString("name=")
	builder.WriteString(rk.Name)
	builder.WriteString(", ")
	builder.WriteString("apiVersion=")
	builder.WriteString(rk.ApiVersion)
	builder.WriteString(", ")
	builder.WriteString("namespaced=")
	builder.WriteString(fmt.Sprintf("%v", rk.Namespaced))
	builder.WriteString(", ")
	builder.WriteString("kind=")
	builder.WriteString(rk.Kind)
	builder.WriteByte(')')
	return builder.String()
}

// ResourceKinds is a parsable slice of ResourceKind.
type ResourceKinds []*ResourceKind
