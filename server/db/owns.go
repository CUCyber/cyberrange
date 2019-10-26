package db

import "time"

type MachineRootOwn struct {
	Id       uint64
	SolvedAt time.Time
	OwnId    uint64
}

type MachineUserOwn struct {
	Id       uint64
	SolvedAt time.Time
	OwnId    uint64
}

type MachineOwn struct {
	Id        uint64
	UserId    uint64
	MachineId uint64
}

func GetOwns() (*[]MachineOwn, error) {
	var owns []MachineOwn

	_, err := db.Select("*").From("machine_owns").
		Load(&owns)
	if err != nil {
		return nil, err
	}

	return &owns, nil
}

func GetRootOwns(user_id uint64) (*[]Machine, error) {
	var machines []Machine

	_, err := db.Select("*").From("machines").
		Join("machine_owns", "machine_owns.machine_id = machines.id").
		Join("machine_root_owns", "machine_root_owns.own_id = machine_owns.id").
		Join("users", "users.id = machine_owns.user_id").
		Where("users.id = ?", user_id).
		Load(&machines)
	if err != nil {
		return nil, err
	}

	return &machines, nil
}

func GetUserOwns(user_id uint64) (*[]Machine, error) {
	var machines []Machine

	_, err := db.Select("*").From("machines").
		Join("machine_owns", "machine_owns.machine_id = machines.id").
		Join("machine_user_owns", "machine_user_owns.own_id = machine_owns.id").
		Join("users", "users.id = machine_owns.user_id").
		Where("users.id = ?", user_id).
		Load(&machines)
	if err != nil {
		return nil, err
	}

	return &machines, nil
}

func HasSubmittedUser(user *User, machine *Machine) (bool, error) {
	userOwns, err := GetUserOwns(user.Id)
	if err != nil {
		return false, err
	}

	for _, owned := range *userOwns {
		if machine.Id == owned.Id {
			return true, nil
		}
	}

	return false, nil
}

func HasSubmittedRoot(user *User, machine *Machine) (bool, error) {
	rootOwns, err := GetRootOwns(user.Id)
	if err != nil {
		return false, err
	}

	for _, owned := range *rootOwns {
		if machine.Id == owned.Id {
			return true, nil
		}
	}

	return false, nil
}

func UserOwnMachine(user *User, machine *Machine) error {
	var own MachineOwn

	/* User has already submitted a correct user flag */
	hasSubmitted, err := HasSubmittedUser(user, machine)
	if hasSubmitted || err != nil {
		return nil
	}

	/* Find User-Machine MachineOwn row */
	_, err = db.Select("*").From("machine_owns").
		Where("user_id = ? AND machine_id = ?", user.Id, machine.Id).
		Load(&own)
	if err != nil {
		return err
	}

	/* Update internal machine root counter */
	_, err = db.Update("machines").
		Set("user_owns", machine.UserOwns+1).
		Where("id = ?", machine.Id).
		Exec()
	if err != nil {
		return err
	}

	/* Create MachineRootOwn record */
	_, err = db.InsertInto("machine_user_owns").
		Pair("own_id", own.Id).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func RootOwnMachine(user *User, machine *Machine) error {
	var own MachineOwn

	/* User has already submitted a correct root flag */
	hasSubmitted, err := HasSubmittedRoot(user, machine)
	if hasSubmitted || err != nil {
		return nil
	}

	/* Find User-Machine MachineOwn row */
	_, err = db.Select("*").From("machine_owns").
		Where("user_id = ? AND machine_id = ?", user.Id, machine.Id).
		Load(&own)
	if err != nil {
		return err
	}

	/* Update internal machine root counter */
	_, err = db.Update("machines").
		Set("root_owns", machine.RootOwns+1).
		Where("id = ?", machine.Id).
		Exec()
	if err != nil {
		return err
	}

	/* Create MachineRootOwn record */
	_, err = db.InsertInto("machine_root_owns").
		Pair("own_id", own.Id).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func MachineCreateMachineOwns(machine_id uint64) error {
	users, err := GetUsers()
	if err != nil {
		panic(err.Error())
	}

	for _, user := range *users {
		_, err := db.InsertInto("machine_owns").
			Pair("user_id", user.Id).
			Pair("machine_id", machine_id).
			Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func UserCreateMachineOwns(user_id uint64) error {
	machines, err := GetMachines()
	if err != nil {
		panic(err.Error())
	}

	for _, machine := range *machines {
		_, err := db.InsertInto("machine_owns").
			Pair("user_id", user_id).
			Pair("machine_id", machine.Id).
			Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
