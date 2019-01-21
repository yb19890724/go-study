package databases

import (
	"fmt"
	"gopkg.in/ini.v1"
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
)

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

func MYSQLConfig() (device string, config string) {

	conf, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("try load config file[%s] error[%s]\n", "config.ini", err.Error())
	}

	device = conf.Section("DB_CONNECTION").Key("DEVICE").String()

	host := conf.Section("MYSQL").Key("HOST").String()
	username := conf.Section("MYSQL").Key( "USERNAME").String()
	password := conf.Section("MYSQL").Key( "PASSWORD").String()
	port := conf.Section("MYSQL").Key( "PORT").String()
	charset := conf.Section("MYSQL").Key( "CHARSET").String()
	database := conf.Section("MYSQL").Key( "DATABASE").String()

	config = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", username, password, host, port, database,charset)

	return device, config
}
