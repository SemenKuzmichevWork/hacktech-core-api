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

var UserReportsEvents = newUserReportsEventsTable("public", "user_reports_events", "")

type userReportsEventsTable struct {
	postgres.Table

	// Columns
	ReportID postgres.ColumnString
	EventID  postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UserReportsEventsTable struct {
	userReportsEventsTable

	EXCLUDED userReportsEventsTable
}

// AS creates new UserReportsEventsTable with assigned alias
func (a UserReportsEventsTable) AS(alias string) *UserReportsEventsTable {
	return newUserReportsEventsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UserReportsEventsTable with assigned schema name
func (a UserReportsEventsTable) FromSchema(schemaName string) *UserReportsEventsTable {
	return newUserReportsEventsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UserReportsEventsTable with assigned table prefix
func (a UserReportsEventsTable) WithPrefix(prefix string) *UserReportsEventsTable {
	return newUserReportsEventsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UserReportsEventsTable with assigned table suffix
func (a UserReportsEventsTable) WithSuffix(suffix string) *UserReportsEventsTable {
	return newUserReportsEventsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUserReportsEventsTable(schemaName, tableName, alias string) *UserReportsEventsTable {
	return &UserReportsEventsTable{
		userReportsEventsTable: newUserReportsEventsTableImpl(schemaName, tableName, alias),
		EXCLUDED:               newUserReportsEventsTableImpl("", "excluded", ""),
	}
}

func newUserReportsEventsTableImpl(schemaName, tableName, alias string) userReportsEventsTable {
	var (
		ReportIDColumn = postgres.StringColumn("report_id")
		EventIDColumn  = postgres.StringColumn("event_id")
		allColumns     = postgres.ColumnList{ReportIDColumn, EventIDColumn}
		mutableColumns = postgres.ColumnList{}
	)

	return userReportsEventsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ReportID: ReportIDColumn,
		EventID:  EventIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
