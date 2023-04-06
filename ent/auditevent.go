// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/auditevent"
)

// AuditEvent is the model entity for the AuditEvent schema.
type AuditEvent struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Raw holds the value of the "raw" field.
	Raw string `json:"raw,omitempty"`
	// Level holds the value of the "level" field.
	Level string `json:"level,omitempty"`
	// AuditID holds the value of the "auditID" field.
	AuditID string `json:"auditID,omitempty"`
	// Verb holds the value of the "verb" field.
	Verb string `json:"verb,omitempty"`
	// UserAgent holds the value of the "userAgent" field.
	UserAgent string `json:"userAgent,omitempty"`
	// RequestTimestamp holds the value of the "requestTimestamp" field.
	RequestTimestamp time.Time `json:"requestTimestamp,omitempty"`
	// StageTimestamp holds the value of the "stageTimestamp" field.
	StageTimestamp time.Time `json:"stageTimestamp,omitempty"`
	// Namespace holds the value of the "namespace" field.
	Namespace string `json:"namespace,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// ApiVersion holds the value of the "apiVersion" field.
	ApiVersion string `json:"apiVersion,omitempty"`
	// ApiGroup holds the value of the "apiGroup" field.
	ApiGroup string `json:"apiGroup,omitempty"`
	// Resource holds the value of the "resource" field.
	Resource string `json:"resource,omitempty"`
	// SubResource holds the value of the "subResource" field.
	SubResource string `json:"subResource,omitempty"`
	// Stage holds the value of the "stage" field.
	Stage string `json:"stage,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AuditEvent) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case auditevent.FieldID:
			values[i] = new(sql.NullInt64)
		case auditevent.FieldRaw, auditevent.FieldLevel, auditevent.FieldAuditID, auditevent.FieldVerb, auditevent.FieldUserAgent, auditevent.FieldNamespace, auditevent.FieldName, auditevent.FieldApiVersion, auditevent.FieldApiGroup, auditevent.FieldResource, auditevent.FieldSubResource, auditevent.FieldStage:
			values[i] = new(sql.NullString)
		case auditevent.FieldRequestTimestamp, auditevent.FieldStageTimestamp:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type AuditEvent", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AuditEvent fields.
func (ae *AuditEvent) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case auditevent.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ae.ID = int(value.Int64)
		case auditevent.FieldRaw:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field raw", values[i])
			} else if value.Valid {
				ae.Raw = value.String
			}
		case auditevent.FieldLevel:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field level", values[i])
			} else if value.Valid {
				ae.Level = value.String
			}
		case auditevent.FieldAuditID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field auditID", values[i])
			} else if value.Valid {
				ae.AuditID = value.String
			}
		case auditevent.FieldVerb:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field verb", values[i])
			} else if value.Valid {
				ae.Verb = value.String
			}
		case auditevent.FieldUserAgent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field userAgent", values[i])
			} else if value.Valid {
				ae.UserAgent = value.String
			}
		case auditevent.FieldRequestTimestamp:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field requestTimestamp", values[i])
			} else if value.Valid {
				ae.RequestTimestamp = value.Time
			}
		case auditevent.FieldStageTimestamp:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field stageTimestamp", values[i])
			} else if value.Valid {
				ae.StageTimestamp = value.Time
			}
		case auditevent.FieldNamespace:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field namespace", values[i])
			} else if value.Valid {
				ae.Namespace = value.String
			}
		case auditevent.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				ae.Name = value.String
			}
		case auditevent.FieldApiVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field apiVersion", values[i])
			} else if value.Valid {
				ae.ApiVersion = value.String
			}
		case auditevent.FieldApiGroup:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field apiGroup", values[i])
			} else if value.Valid {
				ae.ApiGroup = value.String
			}
		case auditevent.FieldResource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field resource", values[i])
			} else if value.Valid {
				ae.Resource = value.String
			}
		case auditevent.FieldSubResource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field subResource", values[i])
			} else if value.Valid {
				ae.SubResource = value.String
			}
		case auditevent.FieldStage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field stage", values[i])
			} else if value.Valid {
				ae.Stage = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this AuditEvent.
// Note that you need to call AuditEvent.Unwrap() before calling this method if this AuditEvent
// was returned from a transaction, and the transaction was committed or rolled back.
func (ae *AuditEvent) Update() *AuditEventUpdateOne {
	return NewAuditEventClient(ae.config).UpdateOne(ae)
}

// Unwrap unwraps the AuditEvent entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ae *AuditEvent) Unwrap() *AuditEvent {
	_tx, ok := ae.config.driver.(*txDriver)
	if !ok {
		panic("ent: AuditEvent is not a transactional entity")
	}
	ae.config.driver = _tx.drv
	return ae
}

// String implements the fmt.Stringer.
func (ae *AuditEvent) String() string {
	var builder strings.Builder
	builder.WriteString("AuditEvent(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ae.ID))
	builder.WriteString("raw=")
	builder.WriteString(ae.Raw)
	builder.WriteString(", ")
	builder.WriteString("level=")
	builder.WriteString(ae.Level)
	builder.WriteString(", ")
	builder.WriteString("auditID=")
	builder.WriteString(ae.AuditID)
	builder.WriteString(", ")
	builder.WriteString("verb=")
	builder.WriteString(ae.Verb)
	builder.WriteString(", ")
	builder.WriteString("userAgent=")
	builder.WriteString(ae.UserAgent)
	builder.WriteString(", ")
	builder.WriteString("requestTimestamp=")
	builder.WriteString(ae.RequestTimestamp.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("stageTimestamp=")
	builder.WriteString(ae.StageTimestamp.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("namespace=")
	builder.WriteString(ae.Namespace)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(ae.Name)
	builder.WriteString(", ")
	builder.WriteString("apiVersion=")
	builder.WriteString(ae.ApiVersion)
	builder.WriteString(", ")
	builder.WriteString("apiGroup=")
	builder.WriteString(ae.ApiGroup)
	builder.WriteString(", ")
	builder.WriteString("resource=")
	builder.WriteString(ae.Resource)
	builder.WriteString(", ")
	builder.WriteString("subResource=")
	builder.WriteString(ae.SubResource)
	builder.WriteString(", ")
	builder.WriteString("stage=")
	builder.WriteString(ae.Stage)
	builder.WriteByte(')')
	return builder.String()
}

// AuditEvents is a parsable slice of AuditEvent.
type AuditEvents []*AuditEvent
