package db

type User struct {
	Id       uint64
	Username string
	Points   uint64
	UserOwns uint64
	RootOwns uint64
}

func GetRank(user *User) (uint64, error) {
	var rank uint64

	_, err := db.Select("COUNT(*) AS rank").From("users").
		Where("points >= (SELECT points FROM users WHERE id = ?)", user.Id).
		Load(&rank)
	if err != nil {
		return 0, err
	}

	return rank, nil
}

func Scoreboard() (*[]User, error) {
	var users []User

	_, err := db.Select("*").From("users").
		OrderDesc("points").
		Load(&users)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func GetUsers() (*[]User, error) {
	var users []User

	_, err := db.Select("*").From("users").
		Load(&users)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func FindUserById(user *User) (*User, error) {
	var queryUser User

	_, err := db.Select("*").From("users").
		Where("id = ?", user.Id).
		Load(&queryUser)
	if err != nil {
		return nil, err
	}

	if queryUser.Username == "" {
		return nil, ErrUserNotFound
	}

	return &queryUser, nil
}

func FindUserByUsername(user *User) (*User, error) {
	var queryUser User

	_, err := db.Select("*").From("users").
		Where("username = ?", user.Username).
		Load(&queryUser)
	if err != nil {
		return nil, err
	}

	if queryUser.Username == "" {
		return nil, ErrUserNotFound
	}

	return &queryUser, nil
}

func CreateUser(user *User) (*User, error) {
	_, err := db.InsertInto("users").Ignore().
		Columns("username", "points").
		Record(user).
		Exec()
	if err != nil {
		return nil, err
	}

	/*
	   While insert...returning isn't available, we need
	   to get the updated user id through a second query

	   https://mariadb.com/kb/en/library/insertreturning/
	*/

	user, err = FindUserByUsername(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindOrCreateUser(user *User) (*User, error) {
	queryUser, err := FindUserByUsername(user)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}

	if queryUser == nil {
		queryUser, err := CreateUser(user)
		if err != nil {
			return nil, err
		}

		err = UserCreateMachineOwns(queryUser.Id)
		if err != nil {
			return nil, err
		}
	}

	return queryUser, nil
}
