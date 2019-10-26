package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/spf13/viper"
)

var db *dbr.Session

func CloseDatabase() {
	db.Connection.DB.Close()
}

func ParseDBOptions() string {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./db/")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error in db config file: %s\n", err))
	}

	mysqlOptions := viper.GetStringMap("mysql")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mysqlOptions["user"], mysqlOptions["pass"],
		mysqlOptions["host"], mysqlOptions["port"],
		mysqlOptions["db"],
	)
}

func SeedDatabase() {
	machine1, err := FindOrCreateMachine(
		&Machine{
			Name:       "Ellingson",
			Flag:       "FlagValue",
			Difficulty: "Hard",
			Points:     30,
		},
	)
	if err != nil {
		panic(err.Error())
	}

	machine2, err := FindOrCreateMachine(
		&Machine{
			Name:       "Smasher2",
			Flag:       "HelloWorld",
			Difficulty: "Insane",
			Points:     50,
		},
	)
	if err != nil {
		panic(err.Error())
	}

	user, err := FindOrCreateUser(
		&User{
			Username: "nbulisc",
		},
	)
	if err != nil {
		panic(err.Error())
	}

	UserOwnMachine(user, machine1)

	RootOwnMachine(user, machine1)

	UserOwnMachine(user, machine2)
}

func InitializeDatabase() {
	conn, err := dbr.Open("mysql", ParseDBOptions(), nil)
	if err != nil {
		panic(err.Error())
	}
	conn.SetMaxOpenConns(10)

	db = conn.NewSession(nil)
	db.Begin()

	SeedDatabase()
}
