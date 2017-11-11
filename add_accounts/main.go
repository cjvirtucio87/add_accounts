package main

import (
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

type fnCreateUserWithConfig func (userConfig)

type fnCreateUserWithPort func (string) fnCreateUserWithConfig

type fnCreateUserWithHost func (string) fnCreateUserWithPort

type fnCreateUserWithAdminPassword func (string) fnCreateUserWithHost

type fnCreateUserWithAdmin func (string) fnCreateUserWithAdminPassword

func createUserWithAdmin (admin string) fnCreateUserWithAdminPassword {
	return func (password string) fnCreateUserWithHost {
		return func (host string) fnCreateUserWithPort {
			return func (port string) fnCreateUserWithConfig {
				dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", admin, password, host, port)

				return func (uconf userConfig) {
					username := uconf.Name
					password := uconf.Password
				
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
			}
		}
	}
}

func main() {
	var addAccountsConfig addAccountsConfig

	filename := "aaconfig.json"
	confj, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(confj, &addAccountsConfig)

	if err != nil {
		panic(err)
	}

	admin := addAccountsConfig.Admin
	password := addAccountsConfig.Password
	host := addAccountsConfig.Host
	port := addAccountsConfig.Port
	userConfigs := addAccountsConfig.UserConfigs

	fnCreateUserWithAdminPassword := createUserWithAdmin(admin)
	fnCreateUserWithHost := fnCreateUserWithAdminPassword(password)
	fnCreateUserWithPort := fnCreateUserWithHost(host)
	fnCreateUserWithConfig := fnCreateUserWithPort(port)

	for _, uconf := range userConfigs {
		fnCreateUserWithConfig(uconf);
	}
}