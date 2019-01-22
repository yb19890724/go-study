package databases

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
)

var Config = struct {
	DB struct {
		Host     	string	`default:"127.0.0.1"`
		Username    string 	`default:"root"`
		Password    string  `default:"root"`
		Port 		string  `default:"3306"`
		Database 	string  `default:"test"`
		Charset 	string  `default:"utf8mb4"`
		Device      string 	`default:"mysql"`
	}
}{}

var Eloquent *gorm.DB

func init() {

	var err error

	device, config := MYSQLConfig()
	Eloquent, err = gorm.Open(device, config)

	if err != nil {
		fmt.Printf("MYSQL connect error %v", err)
	}

	if Eloquent.Error != nil {
		fmt.Printf("database error %v", Eloquent.Error)
	}

}

func MYSQLConfig() (device string, dbConfig string) {


	configor.Load(&Config, "src/config.yml")

	device 	 = Config.DB.Device

	host 	 := Config.DB.Host
	username := Config.DB.Username
	password := Config.DB.Password
	port 	 :=	Config.DB.Port
	charset  := Config.DB.Charset
	database := Config.DB.Database

	dbConfig = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s",
		username,
		password,
		host,
		port,
		database,
		charset)
	return device,dbConfig
}
