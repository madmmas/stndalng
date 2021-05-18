# Auth on Rale and Password Policy Base

## Appling Roles

## Appling Password policy

`config/config.json` file contains the password policy configuration as following:

```json
{
    ...
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

| Folder | Details                             |
| ------ | ----------------------------------- |
| api    | Holds the api endpoints             |
| conf   | Holds the config files              |
| db     | Database Initializer and DB manager |
| route  | router setup                        |
| model  | Models                              |
| utils  | Util functions                      |

## Run

`go run server.go`

## Test

`go test ./..`
