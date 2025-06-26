package roster

import (
	"github.com/quincy/scoutbook-tools/date"
)

type AdultUser struct {
	Name        string
	BsaId       int64
	Email       string
	Gender      string
	UnitNumber  string
	Training    []UserStatusRecord
	HealthForms []UserStatusRecord
	SwimClass   UserStatusRecord
	Positions   []string
}

type YouthUser struct {
	Name        string
	BsaId       int64
	Email       string
	Gender      string
	DateOfBirth date.Date
	Age         int
	Patrol      string
	Training    []UserStatusRecord
	HealthForms []UserStatusRecord
	SwimClass   UserStatusRecord
	Positions   []string
}
