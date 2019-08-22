package main

import (
	"fmt"
	"github.com/Shelnutt2/db2struct"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"
)

func main() {
	gp := os.Getenv("GOPATH")
	host := "192.168.72.188"
	pass := "secret"
	user := "default"
	database := "default"
	port := 3306
	table := "all_data_types"
	structName := "testUser"
	packageName := "testUser"

	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(user, pass, host, port, database, table)

	struc, err := db2struct.Generate(*columnDataTypes, table, structName, packageName, true, true, true)

	cf := filepath.Join(gp, "src", "./test.go")
	f, err := os.Create(cf)
	f.Write([]byte(string(struc)))
	f.Close()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("create gorm struct path:%s", cf)

}
