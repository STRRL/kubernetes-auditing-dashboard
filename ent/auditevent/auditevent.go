// Code generated by ent, DO NOT EDIT.

package auditevent

const (
	// Label holds the string label denoting the auditevent type in the database.
	Label = "audit_event"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldRaw holds the string denoting the raw field in the database.
	FieldRaw = "raw"
	// FieldLevel holds the string denoting the level field in the database.
	FieldLevel = "level"
	// FieldAuditID holds the string denoting the auditid field in the database.
	FieldAuditID = "audit_id"
	// FieldVerb holds the string denoting the verb field in the database.
	FieldVerb = "verb"
	// FieldUserAgent holds the string denoting the useragent field in the database.
	FieldUserAgent = "user_agent"
	// FieldRequestTimestamp holds the string denoting the requesttimestamp field in the database.
	FieldRequestTimestamp = "request_timestamp"
	// FieldStageTimestamp holds the string denoting the stagetimestamp field in the database.
	FieldStageTimestamp = "stage_timestamp"
	// FieldNamespace holds the string denoting the namespace field in the database.
	FieldNamespace = "namespace"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldApiVersion holds the string denoting the apiversion field in the database.
	FieldApiVersion = "api_version"
	// FieldApiGroup holds the string denoting the apigroup field in the database.
	FieldApiGroup = "api_group"
	// FieldResource holds the string denoting the resource field in the database.
	FieldResource = "resource"
	// FieldSubResource holds the string denoting the subresource field in the database.
	FieldSubResource = "sub_resource"
	// FieldStage holds the string denoting the stage field in the database.
	FieldStage = "stage"
	// Table holds the table name of the auditevent in the database.
	Table = "audit_events"
)

// Columns holds all SQL columns for auditevent fields.
var Columns = []string{
	FieldID,
	FieldRaw,
	FieldLevel,
	FieldAuditID,
	FieldVerb,
	FieldUserAgent,
	FieldRequestTimestamp,
	FieldStageTimestamp,
	FieldNamespace,
	FieldName,
	FieldApiVersion,
	FieldApiGroup,
	FieldResource,
	FieldSubResource,
	FieldStage,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// RawValidator is a validator for the "raw" field. It is called by the builders before save.
	RawValidator func(string) error
	// LevelValidator is a validator for the "level" field. It is called by the builders before save.
	LevelValidator func(string) error
	// AuditIDValidator is a validator for the "auditID" field. It is called by the builders before save.
	AuditIDValidator func(string) error
	// VerbValidator is a validator for the "verb" field. It is called by the builders before save.
	VerbValidator func(string) error
	// UserAgentValidator is a validator for the "userAgent" field. It is called by the builders before save.
	UserAgentValidator func(string) error
	// DefaultNamespace holds the default value on creation for the "namespace" field.
	DefaultNamespace string
	// DefaultName holds the default value on creation for the "name" field.
	DefaultName string
	// DefaultApiVersion holds the default value on creation for the "apiVersion" field.
	DefaultApiVersion string
	// DefaultApiGroup holds the default value on creation for the "apiGroup" field.
	DefaultApiGroup string
	// DefaultResource holds the default value on creation for the "resource" field.
	DefaultResource string
	// DefaultSubResource holds the default value on creation for the "subResource" field.
	DefaultSubResource string
)
