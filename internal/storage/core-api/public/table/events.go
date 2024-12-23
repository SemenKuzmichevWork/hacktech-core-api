//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Events = newEventsTable("public", "events", "")

type eventsTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnString
	IsSent      postgres.ColumnBool
	CreatedBy   postgres.ColumnString
	CreatedAt   postgres.ColumnTimestampz
	UpdatedAt   postgres.ColumnTimestampz
	Title       postgres.ColumnString
	Description postgres.ColumnString
	Identifier  postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type EventsTable struct {
	eventsTable

	EXCLUDED eventsTable
}

// AS creates new EventsTable with assigned alias
func (a EventsTable) AS(alias string) *EventsTable {
	return newEventsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new EventsTable with assigned schema name
func (a EventsTable) FromSchema(schemaName string) *EventsTable {
	return newEventsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new EventsTable with assigned table prefix
func (a EventsTable) WithPrefix(prefix string) *EventsTable {
	return newEventsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new EventsTable with assigned table suffix
func (a EventsTable) WithSuffix(suffix string) *EventsTable {
	return newEventsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newEventsTable(schemaName, tableName, alias string) *EventsTable {
	return &EventsTable{
		eventsTable: newEventsTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newEventsTableImpl("", "excluded", ""),
	}
}

func newEventsTableImpl(schemaName, tableName, alias string) eventsTable {
	var (
		IDColumn          = postgres.StringColumn("id")
		IsSentColumn      = postgres.BoolColumn("is_sent")
		CreatedByColumn   = postgres.StringColumn("created_by")
		CreatedAtColumn   = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn   = postgres.TimestampzColumn("updated_at")
		TitleColumn       = postgres.StringColumn("title")
		DescriptionColumn = postgres.StringColumn("description")
		IdentifierColumn  = postgres.StringColumn("identifier")
		allColumns        = postgres.ColumnList{IDColumn, IsSentColumn, CreatedByColumn, CreatedAtColumn, UpdatedAtColumn, TitleColumn, DescriptionColumn, IdentifierColumn}
		mutableColumns    = postgres.ColumnList{IsSentColumn, CreatedByColumn, CreatedAtColumn, UpdatedAtColumn, TitleColumn, DescriptionColumn, IdentifierColumn}
	)

	return eventsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		IsSent:      IsSentColumn,
		CreatedBy:   CreatedByColumn,
		CreatedAt:   CreatedAtColumn,
		UpdatedAt:   UpdatedAtColumn,
		Title:       TitleColumn,
		Description: DescriptionColumn,
		Identifier:  IdentifierColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
