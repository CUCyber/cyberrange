package db

import (
	"sort"
	"time"
)

type OwnMetadata struct {
	Id       uint64
	Type     string
	Name     string
	Points   uint64
	SolvedAt time.Time
}

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

func ResetOwns(user *User) error {
	userOwns, err := GetUserOwns(user.Id)
	if err != nil {
		return err
	}

	rootOwns, err := GetRootOwns(user.Id)
	if err != nil {
		return err
	}

	user, err = FindUserById(user)
	if err != nil {
		return err
	}

	var points uint64
	for i := range *userOwns {
		points += (*userOwns)[i].Points / 2
	}
	for i := range *rootOwns {
		points += (*rootOwns)[i].Points
	}

	_, err = db.Update("users").
		Set("points", points).
		Set("user_owns", len(*userOwns)).
		Set("root_owns", len(*rootOwns)).
		Where("id = ?", user.Id).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserOwn(user *User, machine *Machine) error {
	var err error
	var own MachineOwn

	owned, err := HasSubmittedUser(user, machine)
	if err != nil && !owned {
		return err
	} else if !owned {
		return ErrUserNotOwned
	}

	_, err = db.Select("*").From("machine_owns").
		Where("user_id = ? AND machine_id = ?", user.Id, machine.Id).
		Load(&own)
	if err != nil {
		return err
	}

	_, err = db.DeleteFrom("machine_user_owns").
		Where("own_id = ?", own.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.Update("machines").
		Set("user_owns", machine.UserOwns-1).
		Where("id = ?", machine.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.Update("users").
		Set("points", user.Points-machine.Points/2).
		Set("user_owns", user.UserOwns-1).
		Where("id = ?", user.Id).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func DeleteRootOwn(user *User, machine *Machine) error {
	var err error
	var own MachineOwn

	owned, err := HasSubmittedRoot(user, machine)
	if err != nil && !owned {
		return err
	} else if !owned {
		return ErrRootNotOwned
	}

	_, err = db.Select("*").From("machine_owns").
		Where("user_id = ? AND machine_id = ?", user.Id, machine.Id).
		Load(&own)
	if err != nil {
		return err
	}

	_, err = db.DeleteFrom("machine_root_owns").
		Where("own_id = ?", own.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.Update("machines").
		Set("root_owns", machine.RootOwns-1).
		Where("id = ?", machine.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.Update("users").
		Set("points", user.Points-machine.Points).
		Set("root_owns", user.RootOwns-1).
		Where("id = ?", user.Id).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func GetOwnsTimeline(user_id uint64) (*[]OwnMetadata, error) {
	userOwns, err := GetUserOwns(user_id)
	if err != nil {
		return nil, err
	}

	rootOwns, err := GetRootOwns(user_id)
	if err != nil {
		return nil, err
	}

	totalOwns := append(*userOwns, *rootOwns...)

	sort.Slice(totalOwns, func(i, j int) bool {
		return totalOwns[i].SolvedAt.After(totalOwns[j].SolvedAt)
	})

	return &totalOwns, nil
}

func GetRootOwns(user_id uint64) (*[]OwnMetadata, error) {
	var owns []OwnMetadata

	_, err := db.Select("'root' AS type, name, machines.points, solved_at").From("machines").
		Join("machine_owns", "machine_owns.machine_id = machines.id").
		Join("machine_root_owns", "machine_root_owns.own_id = machine_owns.id").
		Join("users", "users.id = machine_owns.user_id").
		Where("users.id = ?", user_id).
		Load(&owns)
	if err != nil {
		return nil, err
	}

	return &owns, nil
}

func GetUserOwns(user_id uint64) (*[]OwnMetadata, error) {
	var owns []OwnMetadata

	_, err := db.Select("'user' AS type, name, machines.points, solved_at").From("machines").
		Join("machine_owns", "machine_owns.machine_id = machines.id").
		Join("machine_user_owns", "machine_user_owns.own_id = machine_owns.id").
		Join("users", "users.id = machine_owns.user_id").
		Where("users.id = ?", user_id).
		Load(&owns)
	if err != nil {
		return nil, err
	}

	return &owns, nil
}

func HasSubmittedRoot(user *User, machine *Machine) (bool, error) {
	var rootOwns []MachineOwn

	_, err := db.Select("machine_id").From("machine_owns").
		Join("machine_root_owns", "machine_root_owns.own_id = machine_owns.id").
		Where("machine_owns.user_id = ?", user.Id).
		Load(&rootOwns)
	if err != nil {
		return false, err
	}

	for i := range rootOwns {
		if machine.Id == rootOwns[i].MachineId {
			return true, ErrRootOwned
		}
	}

	return false, nil
}

func HasSubmittedUser(user *User, machine *Machine) (bool, error) {
	var userOwns []MachineOwn

	_, err := db.Select("machine_id").From("machine_owns").
		Join("machine_user_owns", "machine_user_owns.own_id = machine_owns.id").
		Where("machine_owns.user_id = ?", user.Id).
		Load(&userOwns)
	if err != nil {
		return false, err
	}

	for i := range userOwns {
		if machine.Id == userOwns[i].MachineId {
			return true, ErrUserOwned
		}
	}

	return false, nil
}

func RootOwnMachine(user *User, machine *Machine) error {
	var err error
	var own MachineOwn

	_, err = HasSubmittedRoot(user, machine)
	if err != nil {
		return err
	}

	_, err = db.Select("*").From("machine_owns").
		Where("user_id = ? AND machine_id = ?", user.Id, machine.Id).
		Load(&own)
	if err != nil {
		return err
	}

	_, err = db.Update("machines").
		Set("root_owns", machine.RootOwns+1).
		Where("id = ?", machine.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.InsertInto("machine_root_owns").
		Pair("own_id", own.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.Update("users").
		Set("points", user.Points+machine.Points).
		Set("root_owns", user.RootOwns+1).
		Where("id = ?", user.Id).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func UserOwnMachine(user *User, machine *Machine) error {
	var err error
	var own MachineOwn

	_, err = HasSubmittedUser(user, machine)
	if err != nil {
		return err
	}

	_, err = db.Select("*").From("machine_owns").
		Where("user_id = ? AND machine_id = ?", user.Id, machine.Id).
		Load(&own)
	if err != nil {
		return err
	}

	_, err = db.Update("machines").
		Set("user_owns", machine.UserOwns+1).
		Where("id = ?", machine.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.InsertInto("machine_user_owns").
		Pair("own_id", own.Id).
		Exec()
	if err != nil {
		return err
	}

	_, err = db.Update("users").
		Set("points", user.Points+machine.Points/2).
		Set("user_owns", user.UserOwns+1).
		Where("id = ?", user.Id).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func OwnMachine(flag Flag, user *User, machine *Machine) error {
	switch flag {
	case machine.UserFlag:
		err := UserOwnMachine(user, machine)
		if err != nil {
			return err
		}
	case machine.RootFlag:
		err := RootOwnMachine(user, machine)
		if err != nil {
			return err
		}
	default:
		return ErrWrongFlag
	}

	return nil
}

func MachineCreateMachineOwns(machine_id uint64) error {
	users, err := GetUsers()
	if err != nil {
		return err
	}

	for i := range *users {
		_, err := db.InsertInto("machine_owns").
			Pair("user_id", (*users)[i].Id).
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
		return err
	}

	for i := range *machines {
		_, err := db.InsertInto("machine_owns").
			Pair("user_id", user_id).
			Pair("machine_id", (*machines)[i].Id).
			Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
