//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type ReportKind string

const (
	ReportKind_Event                ReportKind = "event"
	ReportKind_Business             ReportKind = "business"
	ReportKind_ProjectParticipation ReportKind = "project_participation"
	ReportKind_DailyCheckups        ReportKind = "daily_checkups"
)

var ReportKindAllValues = []ReportKind{
	ReportKind_Event,
	ReportKind_Business,
	ReportKind_ProjectParticipation,
	ReportKind_DailyCheckups,
}

func (e *ReportKind) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "event":
		*e = ReportKind_Event
	case "business":
		*e = ReportKind_Business
	case "project_participation":
		*e = ReportKind_ProjectParticipation
	case "daily_checkups":
		*e = ReportKind_DailyCheckups
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for ReportKind enum")
	}

	return nil
}

func (e ReportKind) String() string {
	return string(e)
}
