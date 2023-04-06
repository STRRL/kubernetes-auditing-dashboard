// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AuditEventsColumns holds the columns for the "audit_events" table.
	AuditEventsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "raw", Type: field.TypeString, Size: 2147483647},
		{Name: "level", Type: field.TypeString},
		{Name: "audit_id", Type: field.TypeString},
		{Name: "verb", Type: field.TypeString},
		{Name: "user_agent", Type: field.TypeString},
		{Name: "request_timestamp", Type: field.TypeTime},
		{Name: "stage_timestamp", Type: field.TypeTime},
		{Name: "namespace", Type: field.TypeString, Default: ""},
		{Name: "name", Type: field.TypeString, Default: ""},
		{Name: "api_version", Type: field.TypeString, Default: ""},
		{Name: "api_group", Type: field.TypeString, Default: ""},
		{Name: "resource", Type: field.TypeString, Default: ""},
		{Name: "sub_resource", Type: field.TypeString, Default: ""},
		{Name: "stage", Type: field.TypeString},
	}
	// AuditEventsTable holds the schema information for the "audit_events" table.
	AuditEventsTable = &schema.Table{
		Name:       "audit_events",
		Columns:    AuditEventsColumns,
		PrimaryKey: []*schema.Column{AuditEventsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "auditevent_level_verb",
				Unique:  false,
				Columns: []*schema.Column{AuditEventsColumns[2], AuditEventsColumns[4]},
			},
			{
				Name:    "auditevent_verb",
				Unique:  false,
				Columns: []*schema.Column{AuditEventsColumns[4]},
			},
			{
				Name:    "auditevent_audit_id",
				Unique:  false,
				Columns: []*schema.Column{AuditEventsColumns[3]},
			},
			{
				Name:    "auditevent_user_agent",
				Unique:  false,
				Columns: []*schema.Column{AuditEventsColumns[5]},
			},
			{
				Name:    "auditevent_request_timestamp",
				Unique:  false,
				Columns: []*schema.Column{AuditEventsColumns[6]},
			},
			{
				Name:    "auditevent_stage_timestamp",
				Unique:  false,
				Columns: []*schema.Column{AuditEventsColumns[7]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AuditEventsTable,
	}
)

func init() {
}
