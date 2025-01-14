//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type UserReports struct {
	ID         uuid.UUID `sql:"primary_key"`
	ReportedBy uuid.UUID
	Kind       ReportKind
	Rating     int32
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
