package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type userConfig struct {
	Name string
	Password string
}

type addAccountsConfig struct {
	Admin string
	Password string
	Host string
	Port string
	UserConfigs []userConfig
}

type fnCreateUserWithUConf func (userConfig)

func createUser (dsn string, username string, password string) {
	db, err := sql.Open("mysql", dsn)
	
	if err != nil {
		panic(err)
	}

	sqlCreateUser := fmt.Sprintf(`
		CREATE USER 
			'%s' 
		IDENTIFIED BY
			'%s'
	`, username, password)

	_, err = db.Exec(sqlCreateUser)

	if err != nil {
		panic(err)
	}

	defer db.Close()	
}

func createUserWithAddAccountsConf (addAccountsConfig addAccountsConfig) fnCreateUserWithUConf {
	admin := addAccountsConfig.Admin
	password := addAccountsConfig.Password
	host := addAccountsConfig.Host
	port := addAccountsConfig.Port
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", admin, password, host, port)

	return func (uconf userConfig) {
		username := uconf.Name
		password := uconf.Password

		createUser(dsn, username, password);
	}
}

func createAddAccountsConfig (filename string) addAccountsConfig {
	var addAccountsConfig addAccountsConfig
	
	tmp := os.TempDir()
	fileDir := fmt.Sprintf("%s\\%s", tmp, filename)
	confj, err := ioutil.ReadFile(fileDir)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(confj, &addAccountsConfig)

	if err != nil {
		panic(err)
	}	

	return addAccountsConfig
}

func main() {
	filename := "aaconfig.json"

	addAccountsConfig := createAddAccountsConfig(filename)
	userConfigs := addAccountsConfig.UserConfigs

	fnCreateUserWithUConf := createUserWithAddAccountsConf(addAccountsConfig)

	for _, uconf := range userConfigs {
		fnCreateUserWithUConf(uconf);
	}
}