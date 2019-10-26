package db

import "github.com/gocraft/dbr/v2"

type Machine struct {
	Id         uint64
	Name       string
	Flag       string
	Difficulty string
	Points     uint
	UserOwns   uint64
	RootOwns   uint64
}

func GetMachines() (*[]Machine, error) {
	var machines []Machine

	_, err := db.Select("*").From("machines").
		Load(&machines)
	if err != nil {
		return nil, err
	}

	return &machines, nil
}

func FindMachineById(machine *Machine) (*Machine, error) {
	var queryMachine Machine

	_, err := db.Select("*").From("machines").
		Where("id = ?", machine.Id).
		Load(&queryMachine)
	if err != nil {
		return nil, err
	}

	if queryMachine.Name == "" {
		return nil, dbr.ErrNotFound
	}

	return &queryMachine, nil
}

func FindMachineByName(machine *Machine) (*Machine, error) {
	var queryMachine Machine

	_, err := db.Select("*").From("machines").
		Where("name = ?", machine.Name).
		Load(&queryMachine)
	if err != nil {
		return nil, err
	}

	if queryMachine.Name == "" {
		return nil, dbr.ErrNotFound
	}

	return &queryMachine, nil
}

func FindOrCreateMachine(machine *Machine) (*Machine, error) {
	/*
	   While insert...returning isn't available, we need
	   to get the updated machine id through a second query

	   https://mariadb.com/kb/en/library/insertreturning/
	*/

	queryMachine, err := FindMachineByName(machine)
	if err != nil {
		/* Unexpected error */
		if err != dbr.ErrNotFound {
			return nil, err
		}

		/* User not found, create user */
		_, err := db.InsertInto("machines").
			Columns("name", "flag", "difficulty", "points").
			Record(machine).
			Exec()
		if err != nil {
			return nil, err
		}

		/* Get new user data */
		queryMachine, err = FindMachineByName(machine)
		if err != nil {
			if err != dbr.ErrNotFound {
				return nil, err
			}
		}

		/* Pre-populate machine owns */
		err = MachineCreateMachineOwns(queryMachine.Id)
		if err != nil {
			panic(err.Error())
		}
	}

	return queryMachine, nil
}
