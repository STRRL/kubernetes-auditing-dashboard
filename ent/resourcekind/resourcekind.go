// Code generated by ent, DO NOT EDIT.

package resourcekind

const (
	// Label holds the string label denoting the resourcekind type in the database.
	Label = "resource_kind"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldApiVersion holds the string denoting the apiversion field in the database.
	FieldApiVersion = "api_version"
	// FieldNamespaced holds the string denoting the namespaced field in the database.
	FieldNamespaced = "namespaced"
	// FieldKind holds the string denoting the kind field in the database.
	FieldKind = "kind"
	// Table holds the table name of the resourcekind in the database.
	Table = "resource_kinds"
)

// Columns holds all SQL columns for resourcekind fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldApiVersion,
	FieldNamespaced,
	FieldKind,
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// ApiVersionValidator is a validator for the "apiVersion" field. It is called by the builders before save.
	ApiVersionValidator func(string) error
	// DefaultNamespaced holds the default value on creation for the "namespaced" field.
	DefaultNamespaced bool
	// KindValidator is a validator for the "kind" field. It is called by the builders before save.
	KindValidator func(string) error
)