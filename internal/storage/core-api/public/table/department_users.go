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

var DepartmentUsers = newDepartmentUsersTable("public", "department_users", "")

type departmentUsersTable struct {
	postgres.Table

	// Columns
	DepartmentID postgres.ColumnString
	UserID       postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type DepartmentUsersTable struct {
	departmentUsersTable

	EXCLUDED departmentUsersTable
}

// AS creates new DepartmentUsersTable with assigned alias
func (a DepartmentUsersTable) AS(alias string) *DepartmentUsersTable {
	return newDepartmentUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new DepartmentUsersTable with assigned schema name
func (a DepartmentUsersTable) FromSchema(schemaName string) *DepartmentUsersTable {
	return newDepartmentUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new DepartmentUsersTable with assigned table prefix
func (a DepartmentUsersTable) WithPrefix(prefix string) *DepartmentUsersTable {
	return newDepartmentUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new DepartmentUsersTable with assigned table suffix
func (a DepartmentUsersTable) WithSuffix(suffix string) *DepartmentUsersTable {
	return newDepartmentUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newDepartmentUsersTable(schemaName, tableName, alias string) *DepartmentUsersTable {
	return &DepartmentUsersTable{
		departmentUsersTable: newDepartmentUsersTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newDepartmentUsersTableImpl("", "excluded", ""),
	}
}

func newDepartmentUsersTableImpl(schemaName, tableName, alias string) departmentUsersTable {
	var (
		DepartmentIDColumn = postgres.StringColumn("department_id")
		UserIDColumn       = postgres.StringColumn("user_id")
		allColumns         = postgres.ColumnList{DepartmentIDColumn, UserIDColumn}
		mutableColumns     = postgres.ColumnList{}
	)

	return departmentUsersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		DepartmentID: DepartmentIDColumn,
		UserID:       UserIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}