package db

type Machine struct {
	Id         uint64
	Name       string
	Points     uint64
	Difficulty string
	UserFlag   string
	RootFlag   string
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
		return nil, ErrMachineNotFound
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
		return nil, ErrMachineNotFound
	}

	return &queryMachine, nil
}

func CreateMachine(machine *Machine) (*Machine, error) {
	_, err := db.InsertInto("machines").Ignore().
		Columns("name", "points", "difficulty", "user_flag", "root_flag").
		Record(machine).
		Exec()
	if err != nil {
		return nil, err
	}

	/*
	   While insert...returning isn't available, we need
	   to get the updated user id through a second query

	   https://mariadb.com/kb/en/library/insertreturning/
	*/

	machine, err = FindMachineByName(machine)
	if err != nil {
		return nil, err
	}

	return machine, nil
}

func FindOrCreateMachine(machine *Machine) (*Machine, error) {
	queryMachine, err := FindMachineByName(machine)
	if err != nil && err != ErrMachineNotFound {
		return nil, err
	}

	if queryMachine == nil {
		queryMachine, err := CreateMachine(machine)
		if err != nil {
			return nil, err
		}

		err = MachineCreateMachineOwns(queryMachine.Id)
		if err != nil {
			return nil, err
		}
	}

	return queryMachine, nil
}
