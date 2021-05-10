package api

import (
	"fmt"
	"net/http"
	"stndalng/config"
	"stndalng/db"
	"stndalng/model"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/valyala/fastjson"

	"golang.org/x/crypto/bcrypt"
)

// jwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	UserID uuid.UUID `json:"id"`
	Roles  string    `json:"roles"`
	Role   string    `json:"role"`
	jwt.StandardClaims
}

func GetClaimUserid(c echo.Context) uuid.UUID {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims.UserID
}

func GetClaimRoles(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims.Roles
}

func GetClaimRole(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims.Role
}

func IsRoot(c echo.Context) bool {
	role := GetClaimRole(c)
	return role == "ROOT"
}

func Login(c echo.Context) error {

	user := new(model.UserLogin)
	if err := c.Bind(user); err != nil {
		fmt.Println("Error in there:")
		return err
	}
	fmt.Println("ROLE::", user.Role)

	// check if there are appcode
	appcode := c.Param("appcode")
	if appcode == "" {
		return echo.ErrUnauthorized
	}

	// Throws unauthorized error
	db := db.DbManager()
	dt := model.User{}
	role := user.Role

	if user.Role == "" {
		// get user default role
		fmt.Println("Role is empty.")
		fmt.Println("USER::", user)
		sql := "select * from _users where appcode=? and username=?"
		if db.Raw(sql, appcode, user.Username).Scan(&dt).RecordNotFound() {
			return echo.ErrUnauthorized
		}
		role = dt.Default_role
	} else {
		sql := "select * from _users where appcode=? and username=? and find_in_set(?, roles)"
		if user.Role == "ROOT" {
			sql = "select * from _users where appcode=? and username=? and is_root=1 and 'ROOT'=?"
		}
		if db.Raw(sql, appcode, user.Username, user.Role).Scan(&dt).RecordNotFound() {
			return echo.ErrUnauthorized
		}
	}

	if !dt.Active {
		return echo.ErrUnauthorized
	}
	fmt.Println(dt)

	// Of_pass_policy := false
	lockout_threshold := 0
	lockout_duration := 0

	configuration := config.GetConfig()
	passwordPolicy := configuration.PASS_POLICY

	lockout_threshold = fastjson.GetInt(passwordPolicy, "lockout_threshold")
	lockout_duration = fastjson.GetInt(passwordPolicy, "lockout_duration")

	is_lockout := (lockout_threshold > 0) && dt.IsLockout
	failpass_count := dt.Failpasscount
	fmt.Println("failpass_count::", failpass_count)
	if is_lockout {
		// check if the threshold is over then reset
		now := time.Now()
		elapsed := now.Sub(dt.LockoutStart)
		fmt.Println("Lockout Time elapsed::", elapsed.Minutes())
		if elapsed.Minutes() <= float64(lockout_duration) {
			fmt.Println("Lock time is not over.")
			return echo.ErrUnauthorized
		} else {
			fmt.Println("Updated as lockout over.")
			failpass_count = 0
		}
	}
	fmt.Println("2: failpass_count::", failpass_count)
	if !ValidateUsassword(dt.Password, user.Password) {
		// if password does not match than update
		if lockout_threshold > 0 {
			fmt.Println(lockout_duration)
			failpass_count = failpass_count + 1
			fmt.Println("Fail pass:", failpass_count)
			if failpass_count == lockout_threshold {
				fmt.Println("Updated as lockout set to true.")
				db.Model(&dt).Update(map[string]interface{}{
					"Failpasscount": failpass_count,
					"Lastfailpass":  time.Now(),
					"IsLockout":     true,
					"LockoutStart":  time.Now(),
				})
			} else {
				fmt.Println("Updated as lockout set to false.")
				db.Model(&dt).Update(map[string]interface{}{
					"Failpasscount": failpass_count,
					"Lastfailpass":  time.Now(),
					"IsLockout":     false,
					"LockoutStart":  time.Time{},
				})
			}
		}
		return echo.ErrUnauthorized
	} else {
		fmt.Println("Updated as password matched.")
		db.Model(&dt).Update(map[string]interface{}{
			"Failpasscount": 0,
			"IsLockout":     false,
			"LockoutStart":  time.Time{},
		})
	}

	claims := &JwtCustomClaims{
		dt.ID,
		dt.Roles,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	if_pass_expire := fastjson.GetBool(passwordPolicy, "if_pass_expire")
	if if_pass_expire {
		days_tobe_expired := fastjson.GetInt(passwordPolicy, "days_tobe_expired")
		fmt.Println("days_tobe_expired::", days_tobe_expired)
		fmt.Println("user updated::", dt.Updated)

		now := time.Now()
		elapsed := now.Sub(dt.Updated)
		fmt.Println("Time elapsed::", elapsed.Hours()/24)
		if int(elapsed.Hours()) >= (days_tobe_expired * 24) {
			fmt.Println("Time Expired!!!")
			return c.JSON(http.StatusForbidden, echo.Map{
				"code":    20000,
				"message": "Password expired. Please change your password.",
				"data":    t,
			})
			// return echo.ErrForbidden
		} else {
			fmt.Println("Still have time.")
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code": 20000,
		"data": t,
	})
}

func ValidateUsassword(hashpassword, password string) bool {
	tt := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
	fmt.Println(tt)
	return tt == nil
}

func UserInfo(c echo.Context) error {

	selected_role := GetClaimRole(c)
	fmt.Println("Selected_role: ", selected_role)

	userid := GetClaimUserid(c)
	fmt.Println("UserID: ", userid)

	// get the user info
	db := db.DbManager()
	dt := new(model.User)
	if db.Where("id = ?", userid).Find(&dt).RecordNotFound() {
		return echo.NewHTTPError(http.StatusBadRequest, "Badrequest")
	}
	roles := strings.Split(dt.Roles, ",")
	return c.JSON(http.StatusOK, echo.Map{
		"code": 20000,
		"data": echo.Map{
			"roles":      roles,
			"role":       selected_role,
			"name":       dt.Username,
			"email":      dt.Email,
			"profile_id": dt.ProfileId,
			"profile":    dt.Profile,
		},
	})
}
