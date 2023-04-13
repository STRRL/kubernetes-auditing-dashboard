// Code generated by ent, DO NOT EDIT.

package view

const (
	// Label holds the string label denoting the view type in the database.
	Label = "view"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// Table holds the table name of the view in the database.
	Table = "views"
)

// Columns holds all SQL columns for view fields.
var Columns = []string{
	FieldID,
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
