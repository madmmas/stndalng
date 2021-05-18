package api

import (
	"fmt"
	"net/http"
	"stndalng/config"
	"stndalng/model"
	"stndalng/repo"
	"time"
	"unicode"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
)

func bCheckUserExists(email string) bool {
	db := repo.DbManager()
	dt := model.User{}
	if db.Where(&model.User{Email: email}).First(&dt).RecordNotFound() {
		return false
	}
	return true
}
func bCheckUsernameExists(username string) bool {
	db := repo.DbManager()
	dt := model.User{}
	if db.Where(&model.User{Username: username}).First(&dt).RecordNotFound() {
		return false
	}
	return true
}

func GetUsers(c echo.Context) error {
	// check if root
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()
	Users := []model.User{}
	db.Where("is_sys = 0").Where("active=1").Select("id, username, email, active, roles, is_root, updated").Find(&Users)
	return c.JSON(http.StatusOK, Users)
}

func GetDeUsers(c echo.Context) error {
	// check if root
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()
	Users := []model.User{}
	db.Where("is_sys = 0").Where("active=0").Select("id, username, email, active, roles, is_root, updated").Find(&Users)
	return c.JSON(http.StatusOK, Users)
}

func DeleteUser(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()
	dt := model.User{}
	id := c.Param("id")
	fmt.Println(id)
	if id == "empty" {
		fmt.Println(id)
		return c.JSON(http.StatusOK, map[string]interface{}{})
	} else {
		db.Where("id = ?", id).Delete(&dt)
		return c.JSON(http.StatusOK, id)
	}
}

func GetUser(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()
	dt := model.User{}
	db.Where("id = ?", c.Param("id")).First(&dt)
	return c.JSON(http.StatusOK, dt)
}

func NewUser(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	dt := new(model.User)
	if err := c.Bind(dt); err != nil {
		return err
	}

	fmt.Println("USER DATA::", dt)

	username := dt.Username

	email := dt.Email

	roles := dt.Roles

	defaultRole := dt.Default_role

	profile := []byte(`{"avatar": "", "designation": "` + defaultRole + `", "display_name": "` + email + `"}`)

	fmt.Println("Profile::", profile)

	password := dt.Password

	confirmPassword := dt.ConfirmPassword

	if password == "" || confirmPassword == "" || username == "" || email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please enter username, email and password")
	}
	if roles == "" || defaultRole == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please enter Roles and Default Roles")
	}
	if password != confirmPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "Confirm password is not same to password provided.")
	}

	if bCheckUserExists(email) == true {
		return echo.NewHTTPError(http.StatusBadRequest, "Email provided already exists")
	}
	if bCheckUsernameExists(username) == true {
		return echo.NewHTTPError(http.StatusBadRequest, "Email provided already exists")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// add to the password history and also check the password policy
	// need to make a generic password validator
	dt_f := model.User{Username: username, Email: email, Password: string(hashedPassword), Roles: roles, Profile: profile, Default_role: defaultRole, Status: 1, Active: true}

	db := repo.DbManager()

	db.Create(&dt_f)

	return c.JSON(http.StatusCreated, map[string]string{
		"ID": dt.ID.String(),
	})
}

func UpdateUser(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	dt := new(model.User)
	if err := c.Bind(dt); err != nil {
		return err
	}
	fmt.Println(dt.ID)
	db := repo.DbManager()

	dt_f := model.User{}

	if db.Where("id = ?", dt.ID).First(&dt_f).RecordNotFound() {
		// record not found
		return c.JSON(http.StatusCreated, map[string]string{
			"Msg": "Not Found",
			"ID":  dt.ID.String(),
		})
	} else {
		// only update
		fmt.Println(dt)
		username := dt.Username
		email := dt.Email
		roles := dt.Roles
		if username == "" || email == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Please enter username, email and password")
		}
		// db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
		if dt.Email != dt_f.Email {
			if bCheckUserExists(email) == true {
				return echo.NewHTTPError(http.StatusBadRequest, "Email provided already exists")
			}
		}
		if dt.Username != dt_f.Username {
			if bCheckUsernameExists(username) == true {
				return echo.NewHTTPError(http.StatusBadRequest, "Username provided already exists")
			}
		}

		db.Model(&dt).Updates(model.User{Email: email, Roles: roles})

		return c.JSON(http.StatusCreated, map[string]string{
			"ID": dt.ID.String(),
		})
	}

}

func verifyPassword(s string) (letter, number, upper, special bool) {

	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letter = true
		default:
			//return false, false, false, false
		}
	}
	return
}

func VarifyPasswordPolicy(userid, password string) (bool, string, bool) {

	db := repo.DbManager()
	letter, number, upper, special := verifyPassword(password)
	fmt.Println("Password::", letter, number, upper, special)

	// get policy
	configuration := config.GetConfig()
	passwordPolicy := configuration.PASS_POLICY

	pass_size := passwordPolicy.PASS_SIZE       // fastjson.GetInt(passwordPolicy, "pass_size")
	pass_letter := passwordPolicy.PASS_LETTER   //fastjson.GetBool(passwordPolicy, "pass_letter")
	pass_number := passwordPolicy.PASS_NUMBER   //fastjson.GetBool(passwordPolicy, "pass_number")
	pass_upper := passwordPolicy.PASS_UPPER     //fastjson.GetBool(passwordPolicy, "pass_upper")
	pass_special := passwordPolicy.PASS_SPECIAL //fastjson.GetBool(passwordPolicy, "pass_special")
	fmt.Println("Password Policy::", pass_size, pass_letter, pass_number, pass_upper, pass_special)
	if len(password) < pass_size || letter != pass_letter || number != pass_number || upper != pass_upper || special != pass_special {
		return false, "Password policy does not match. Please provide a valid password.", false
	}

	// check if password history is enabled

	i_pass_history := passwordPolicy.PASS_HISTORY //fastjson.GetInt(passwordPolicy, "pass_history")
	fmt.Println("PassHistory", i_pass_history)
	if i_pass_history > 0 {
		// get all the password
		dt_cps := []model.Changepass{}
		if db.Table("_changepasses").Limit(i_pass_history).Offset(0).Where("userid = ?", userid).Select("password").Order("created desc").Find(&dt_cps).RecordNotFound() {
			// return echo.ErrUnauthorized
		} else {
			fmt.Println("PassHistory::", dt_cps)
			for _, dt_cp := range dt_cps {
				tt := bcrypt.CompareHashAndPassword([]byte(dt_cp.Password), []byte(password))
				fmt.Println(tt)
				if tt == nil {
					return false, "Password already used. Please input a different password.", false
				}
			}
		}
	}
	return true, "", i_pass_history > 0
	// return true, "", false
}

func ChangePassword(c echo.Context) error {

	userid := GetClaimUserid(c)
	fmt.Println("Userid: ", userid.String())

	dt := new(model.ChangePass)
	if err := c.Bind(dt); err != nil {
		return err
	}
	fmt.Println(dt.Password)
	fmt.Println(dt.Repassword)

	db := repo.DbManager()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	dt_f := model.User{}

	if tx.Where("id = ?", userid).First(&dt_f).RecordNotFound() {
		// record not found
		return c.JSON(http.StatusNotFound, echo.Map{
			"code":    20000,
			"message": "User not found.",
		})

	} else {
		// only update

		password := dt.Password
		confirmPassword := dt.Repassword

		if password == "" || confirmPassword == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Please enter password")
		}

		if password != confirmPassword {
			return echo.NewHTTPError(http.StatusBadRequest, "Confirm password is not same to password provided")
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

		// varify password with the policy
		flg, msg, is_pass_history := VarifyPasswordPolicy(userid.String(), password)
		if flg == false {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"code":    20000,
				"message": msg,
			})
		}
		fmt.Println("PassHisotry", is_pass_history)
		if is_pass_history {
			// save the model in the changepass table
			dt_cp := model.Changepass{Userid: userid, Password: string(hashedPassword)}

			if err := tx.Create(&dt_cp).Error; err != nil {
				fmt.Println("ERROR::", err)
				tx.Rollback()
				return err
			}
		}
		if err := tx.Model(&dt_f).Update(model.User{Password: string(hashedPassword), Updated: time.Now()}).Error; err != nil {
			tx.Rollback()
			return err
		}

	}
	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"code":    "50000",
			"message": "Password failed to changed.",
		})

	} else {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"code":    "20000",
			"message": "Password changed.",
		})
	}
}

func ChangeUserDeactiveFlag(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	id := c.Param("id")
	db := repo.DbManager()

	dt_f := model.User{}

	if db.Where("id = ?", id).First(&dt_f).RecordNotFound() {
		// record not found
		return c.JSON(http.StatusCreated, map[string]string{
			"Msg": "Not Found",
			"ID":  id,
		})
	} else {
		// db.Model(&dt_f).Update(model.User{Active: false})
		if !dt_f.IsRoot {
			db.Model(&dt_f).Update(map[string]interface{}{
				"Active": false,
			})
			return c.JSON(http.StatusCreated, map[string]string{
				"code":    "20000",
				"message": "User Deactivate.",
			})
		} else {
			return c.JSON(http.StatusCreated, map[string]string{
				"code":    "50000",
				"message": "Root Cannot Deactivate.",
			})
		}
	}
}
func ChangeUserActiveFlag(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	id := c.Param("id")
	db := repo.DbManager()

	dt_f := model.User{}

	if db.Where("id = ?", id).First(&dt_f).RecordNotFound() {
		// record not found
		return c.JSON(http.StatusCreated, map[string]string{
			"Msg": "Not Found",
			"ID":  id,
		})
	} else {
		db.Model(&dt_f).Update(map[string]interface{}{
			"Active": true,
		})
		return c.JSON(http.StatusCreated, map[string]string{
			"ID": id,
		})
	}
}
