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

type fnGrantPrivilegesWithUConf func (userConfig)

type fnSQL func (dsn string, username string, password string)

type fnExecuteSQL func (userConfig)

func createUser (dsn string, username string, password string) {
	db, err := sql.Open("mysql", dsn)
	
	if err != nil {
		panic(err)
	}

	sqlQuery := fmt.Sprintf(`
		CREATE USER 
			'%s'
		IDENTIFIED BY
			'%s'
	`, username, password)

	_, err = db.Exec(sqlQuery)

	if err != nil {
		panic(err)
	}

	defer db.Close()	
}

func grantPrivileges (dsn string, username string, password string) {
	db, err := sql.Open("mysql", dsn)
	
	if err != nil {
		panic(err)
	}

	sqlQuery := fmt.Sprintf(`
		GRANT ALL PRIVILEGES ON 
			*.*
		TO 
			'%s'
	`, username)

	_, err = db.Exec(sqlQuery)

	if err != nil {
		panic(err)
	}

	sqlQuery = `FLUSH PRIVILEGES`

	_, err = db.Exec(sqlQuery)

	if err != nil {
		panic(err)
	}

	defer db.Close()
}

func composeSQL (conf addAccountsConfig, cb fnSQL) fnExecuteSQL {
	admin := conf.Admin
	password := conf.Password
	host := conf.Host
	port := conf.Port
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", admin, password, host, port)

	return func (uconf userConfig) {
		username := uconf.Name
		password := uconf.Password

		cb(dsn, username, password);
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

	fnCreateUser := composeSQL(addAccountsConfig, createUser)
	fnGrantPrivileges := composeSQL(addAccountsConfig, grantPrivileges)

	for _, uconf := range userConfigs {
		fnCreateUser(uconf)
		fnGrantPrivileges(uconf)
	}
}