// Code generated by ent, DO NOT EDIT.

package task

import (
	"fmt"
)

const (
	// Label holds the string label denoting the task type in the database.
	Label = "task"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldPriority holds the string denoting the priority field in the database.
	FieldPriority = "priority"
	// FieldDueDate holds the string denoting the due_date field in the database.
	FieldDueDate = "due_date"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// Table holds the table name of the task in the database.
	Table = "tasks"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "tasks"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_id"
)

// Columns holds all SQL columns for task fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldDescription,
	FieldStatus,
	FieldPriority,
	FieldDueDate,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "tasks"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_id",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// TitleValidator is a validator for the "title" field. It is called by the builders before save.
	TitleValidator func(string) error
)

// Status defines the type for the "status" enum field.
type Status string

// Status values.
const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted:
		return nil
	default:
		return fmt.Errorf("task: invalid enum value for status field: %q", s)
	}
}

// Priority defines the type for the "priority" enum field.
type Priority string

// Priority values.
const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (pr Priority) String() string {
	return string(pr)
}

// PriorityValidator is a validator for the "priority" field enum values. It is called by the builders before save.
func PriorityValidator(pr Priority) error {
	switch pr {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	default:
		return fmt.Errorf("task: invalid enum value for priority field: %q", pr)
	}
}
