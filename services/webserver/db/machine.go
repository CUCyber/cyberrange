package db

type Machine struct {
	Id         uint64 `json:"-"`
	Name       string `json:"Name"`
	Points     uint64 `json:"-"`
	Difficulty string `json:"Difficulty"`
	UserFlag   string `json:"-"`
	RootFlag   string `json:"-"`
	UserOwns   uint64 `json:"UserOwns"`
	RootOwns   uint64 `json:"RootOwns"`
	IpAddress  string `json:"IpAddress"`
}

var MachineDifficulty = [...]string{"Easy", "Medium", "Hard", "Insane"}

func GetMachines() (*[]Machine, error) {
	var machines []Machine

	_, err := db.Select("*").From("machines").
		Load(&machines)
	if err != nil {
		return nil, err
	}

	return &machines, nil
}

func SetMachineIp(machine *Machine) error {
	_, err := db.Update("machines").
		Set("ip_address", machine.IpAddress).
		Where("name = ?", machine.Name).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func MachineExists(machine *Machine) (bool, error) {
	_, err := FindMachineByName(machine)
	if err != nil {
		if err == ErrMachineNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
		Columns("name", "points", "difficulty", "user_flag", "root_flag", "ip_address").
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
