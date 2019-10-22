package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var db *gorm.DB

func CloseDatabase() {
	db.Close()
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
		mysqlOptions["username"], mysqlOptions["password"],
		mysqlOptions["host"], mysqlOptions["port"], mysqlOptions["db"],
	)
}

func SeedDatabase() {
	FirstOrCreateMachine(
		&Machine{
			Name:       "Ellingson",
			Difficulty: "Hard",
			Flag:       "HelloWorld",
		},
	)

	FirstOrCreateMachine(
		&Machine{
			Name:       "Smasher2",
			Difficulty: "Insane",
			Flag:       "HelloWorld",
		},
	)

	for i := 0; i < 30; i++ {
		UserOwnMachine(&Machine{Name: "Ellingson"})
	}

	for i := 0; i < 24; i++ {
		RootOwnMachine(&Machine{Name: "Ellingson"})
	}

	for i := 0; i < 6; i++ {
		UserOwnMachine(&Machine{Name: "Smasher2"})
	}

	for i := 0; i < 4; i++ {
		RootOwnMachine(&Machine{Name: "Smasher2"})
	}
}

func InitializeDatabase() {
	var err error

	connStr := ParseDBOptions()

	db, err = gorm.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}

	db.DropTableIfExists(&User{}, &Machine{})
	db.AutoMigrate(&User{}, &Machine{})

	SeedDatabase()
}
