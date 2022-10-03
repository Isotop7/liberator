package models

import (
	"database/sql/driver"
	"errors"
)

type AssignmentStatus struct {
	definition string
}

func (as AssignmentStatus) String() string {
	return as.definition
}

var (
	Unknown  = AssignmentStatus{""}
	Assigned = AssignmentStatus{"assigned"}
	Active   = AssignmentStatus{"active"}
	Inactive = AssignmentStatus{"inactive"}
)

func (as AssignmentStatus) Value() (driver.Value, error) {
	return as.String(), nil
}

func (as *AssignmentStatus) Scan(value interface{}) error {
	if value == nil {
		*as = Unknown
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			*as = AssignmentStatus{
				definition: v,
			}
			return nil
		}
	}

	return errors.New("failed to scan assignmentStatus")
}
