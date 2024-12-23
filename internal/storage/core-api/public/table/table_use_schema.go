//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	CompanyPositions = CompanyPositions.FromSchema(schema)
	CompanyRoles = CompanyRoles.FromSchema(schema)
	DepartmentUsers = DepartmentUsers.FromSchema(schema)
	Departments = Departments.FromSchema(schema)
	Events = Events.FromSchema(schema)
	GooseDbVersion = GooseDbVersion.FromSchema(schema)
	Kpi = Kpi.FromSchema(schema)
	UserReports = UserReports.FromSchema(schema)
	UserReportsEvents = UserReportsEvents.FromSchema(schema)
	Users = Users.FromSchema(schema)
}
