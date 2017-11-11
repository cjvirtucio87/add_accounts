# Add Accounts script

A simple script to add user accounts to your MySQL database. Easily configurable with a JSON file. Just place it in your platform's `temp` folder.

The JSON config file must be named `aaconfig.json` and should have the following shape:

```
# linux: /tmp/aaconfig.json
# windows: %TMP%\aaconfig.json
{
  "Admin": "root",
  "Password": "123asiaslocal_lr7",
  "Host": "127.0.0.1",
  "Port": "3306",
  "UserConfigs": [
    {
      "Name": "test_user0",
      "Password": "test_pw0"
    },
    {
      "Name": "test_user1",
      "Password": "test_pw1"
    }
  ]
}
```

`Admin` represents the account that has admin privileges over the database. `Password` is that account's password. This script interfaces with the MySQL daemon over `tcp`, so you need to specify a `Host` and a `Port`. You then define the credentials for the user accounts that you want ot create as an array mapped to `UserConfigs`.
