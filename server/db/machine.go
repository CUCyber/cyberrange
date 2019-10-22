package db

import (
	"github.com/jinzhu/gorm"
)

type Machine struct {
	gorm.Model
	Name       string `gorm:"unique;not null"`
	Flag       string
	Difficulty string
	UserOwns   int
	RootOwns   int
}

func UserOwnMachine(machine *Machine) {
	result := db.Find(&machine, machine)
	if result.Error != nil {
		return
	}
	machine.UserOwns += 1
	db.Save(machine)
}

func RootOwnMachine(machine *Machine) {
	result := db.Find(&machine, machine)
	if result.Error != nil {
		return
	}
	machine.RootOwns += 1
	db.Save(machine)
}

func FindMachine(machine *Machine) (*Machine, error) {
	result := db.Find(&machine, machine)
	if result.Error != nil {
		return nil, result.Error
	}
	return machine, nil
}

func GetMachineList() (*[]Machine, error) {
	var machines []Machine
	result := db.Find(&machines)
	if result.Error != nil {
		return nil, result.Error
	}
	return &machines, nil
}

func FirstOrCreateMachine(machine *Machine) *Machine {
	db.FirstOrCreate(&machine, machine)
	return machine
}
