# Auth Based On Multiple Roles and Password Policy

## Appling Roles

Roles can be applied on uri configured in the `config/handler_roles.json` file.

```json
{
  "Login": [], // anonymous
  "Logout": [], // anonymous
  "NewRole": ["ROOT", "ADMIN"],
  "UpdateRole": ["ROOT", "ADMIN"],
  "GetRoles": ["ROOT", "ADMIN"],
  "GetRole": ["*"], // authenticated users/owner
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

`config/config.json` file contains the password policy configuration as following:

```json
{
  "PASS_POLICY": {
    "pass_size": 6,
    "pass_upper": false,
    "pass_letter": true,
    "pass_number": true,
    "pass_history": 3,
    "pass_special": false,
    "if_pass_expire": true,
    "lockout_duration": 1,
    "days_tobe_expired": 45,
    "lockout_threshold": 3
  }
}
```

## Architecture

| Folder    | Details                             |
| --------- | ----------------------------------- |
| api       | Holds the api endpoints             |
| conf      | Holds the config files              |
| db        | Database Initializer and DB manager |
| migration | Sql dump                            |
| route     | router setup                        |
| model     | Models                              |
| utils     | Util functions                      |

## Run

`go run server.go`

## Test

`go test ./..`
