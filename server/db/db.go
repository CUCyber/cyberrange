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
	var err error

	_, err = FindOrCreateMachine(
		&Machine{
			Name:       "Ellingson",
			Points:     40,
			Difficulty: "Hard",
			UserFlag:   "EllingsonUser",
			RootFlag:   "EllingsonRoot",
		},
	)
	if err != nil {
		panic(err.Error())
	}

	_, err = FindOrCreateMachine(
		&Machine{
			Name:       "Smasher2",
			Points:     50,
			Difficulty: "Insane",
			UserFlag:   "Smasher2User",
			RootFlag:   "Smasher2Root",
		},
	)
	if err != nil {
		panic(err.Error())
	}

	_, err = FindOrCreateUser(
		&User{
			Username: "nbulisc",
		},
	)
	if err != nil {
		panic(err.Error())
	}
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
