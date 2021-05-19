# Auth Based On Multiple Roles and Password Policy

## Appling Roles

Roles can be applied on uri configured in the `config/handler_roles.json` file.

```json
{
  "Login": [],
  "Logout": [],
  "NewRole": ["ROOT", "ADMIN"],
  "UpdateRole": ["ROOT", "ADMIN"],
  "GetRoles": ["ROOT", "ADMIN"],
  "GetRole": ["*"],
  "GetRolesKV": ["*"],
  "UserInfo": ["*"],
  "GetUsers": ["ROOT", "ADMIN"],
  "GetUser": ["*"],
  "NewUser": ["ROOT", "ADMIN"],
  "UpdateUser": ["ROOT", "ADMIN"],
  "GetDeUsers": ["ROOT", "ADMIN"],
  "DeleteUser": ["ROOT", "ADMIN"],
  "ChangePassword": ["*"],
  "ChangeUserDeactiveFlag": ["ROOT"],
  "ChangeUserActiveFlag": ["ROOT"]
}
```

## Appling Password policy

`config/pass_policy.json` file contains the password policy configuration as following:

```json
{
  "PASS_SIZE": 6,
  "PASS_UPPER": false,
  "PASS_LETTER": true,
  "PASS_NUMBER": true,
  "PASS_HISTORY": 3,
  "PASS_SPECIAL": false,
  "IF_PASS_EXPIRE": true,
  "LOCKOUT_DURATION": 1,
  "DAYS_TOBE_EXPIRED": 45,
  "LOCKOUT_THRESHOLD": 3,
  "TOKEN_TOBE_EXPIRED": 24,
  "TOKEN_CRYPTO_KEY": "F6F61L8L7CUCF61L8L7CUCGN0NKF61L8L7CUCGN0NKF61L8L7CUCGN0NK6336I8TFP9Y2ZOS43OS43"
}
```

## Architecture

| Folder    | Details                             |
| --------- | ----------------------------------- |
| api       | Holds the api endpoints             |
| conf      | Holds the config files              |
| repo      | Database Initializer and DB manager |
| migration | Sql dump                            |
| route     | router setup                        |
| model     | Models                              |
| utils     | Util functions                      |

## Run

`go run server.go`
