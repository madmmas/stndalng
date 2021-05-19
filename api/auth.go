package api

import (
	"fmt"
	"net/http"
	"stndalng/config"
	"stndalng/model"
	"stndalng/repo"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"golang.org/x/crypto/bcrypt"
)

type (
	s_conf struct {
		DB         *gorm.DB
		PassPolicy *config.PassPolicy
	}
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

func (c *s_conf) getUserData(user *model.UserLogin, dt_user *model.User) (string, bool) {

	db := c.DB
	role := user.Role
	if user.Role == "" {
		// get user default role
		fmt.Println("Role is empty.")
		fmt.Println("USER::", user)
		sql := "select * from users where username=? and active=1"
		if db.Raw(sql, user.Username).Scan(&dt_user).RecordNotFound() {
			return "", false
		}
		role = dt_user.Default_role
	} else {
		sql := "select * from users where username=? and active=1 and find_in_set(?, roles)"
		if user.Role == "ROOT" {
			sql = "select * from users where username=? and active=1 and is_root=1 and 'ROOT'=?"
		}
		if db.Raw(sql, user.Username, user.Role).Scan(&dt_user).RecordNotFound() {
			return "", false
		}
	}
	return role, true
}

func (c *s_conf) isInThrottle(dt_user *model.User) (int, bool) {
	passwordPolicy := c.PassPolicy
	is_lockout := (passwordPolicy.LOCKOUT_THRESHOLD > 0) && dt_user.IsLockout
	failpass_count := dt_user.Failpasscount
	fmt.Println("failpass_count::", failpass_count)
	if is_lockout {
		// check if the threshold is over then reset
		now := time.Now()
		elapsed := now.Sub(dt_user.LockoutStart)
		fmt.Println("Lockout Time elapsed::", elapsed.Minutes())
		if elapsed.Minutes() <= float64(passwordPolicy.LOCKOUT_DURATION) {
			fmt.Println("Lock time is not over.")
			return failpass_count, true
		} else {
			fmt.Println("Updated as lockout over.")
			failpass_count = 0
		}
	}
	fmt.Println("2: failpass_count::", failpass_count)
	return failpass_count, false
}

func (c *s_conf) isPasswordValid(failpass_count int, user *model.UserLogin, dt_user *model.User) bool {
	db := c.DB
	passwordPolicy := c.PassPolicy
	if !ValidateUsassword(dt_user.Password, user.Password) {
		// if password does not match than update
		if passwordPolicy.LOCKOUT_THRESHOLD > 0 {
			// fmt.Println(lockout_duration)
			failpass_count = failpass_count + 1
			fmt.Println("Fail pass:", failpass_count)
			if failpass_count == passwordPolicy.LOCKOUT_THRESHOLD {
				fmt.Println("Updated as lockout set to true.")
				db.Model(&dt_user).Update(map[string]interface{}{
					"Failpasscount": failpass_count,
					"Lastfailpass":  time.Now(),
					"IsLockout":     true,
					"LockoutStart":  time.Now(),
				})
			} else {
				fmt.Println("Updated as lockout set to false.")
				db.Model(&dt_user).Update(map[string]interface{}{
					"Failpasscount": failpass_count,
					"Lastfailpass":  time.Now(),
					"IsLockout":     false,
					"LockoutStart":  time.Time{},
				})
			}
		}
		return false
	} else {
		fmt.Println("Updated as password matched.")
		db.Model(&dt_user).Update(map[string]interface{}{
			"Failpasscount": 0,
			"IsLockout":     false,
			"LockoutStart":  time.Time{},
		})
	}
	return true
}

func (c *s_conf) isPasswordExpired(dt_user *model.User) bool {
	passwordPolicy := c.PassPolicy
	// if_pass_expire := fastjson.GetBool(passwordPolicy, "if_pass_expire")
	if passwordPolicy.IF_PASS_EXPIRE {
		// days_tobe_expired := fastjson.GetInt(passwordPolicy, "days_tobe_expired")
		fmt.Println("days_tobe_expired::", passwordPolicy.DAYS_TOBE_EXPIRED)
		fmt.Println("user updated::", dt_user.Updated)

		now := time.Now()
		elapsed := now.Sub(dt_user.Updated)
		fmt.Println("Time elapsed::", elapsed.Hours()/24)
		if int(elapsed.Hours()) >= (passwordPolicy.DAYS_TOBE_EXPIRED * 24) {
			fmt.Println("Time Expired!!!")
			return true

		} else {
			fmt.Println("Still have time.")
		}
	}
	return false
}

func (c *s_conf) generateToken(role string, dt_user *model.User) (string, error) {
	claims := &JwtCustomClaims{
		dt_user.ID,
		dt_user.Roles,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(c.PassPolicy.TOKEN_TOBE_EXPIRED)).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(c.PassPolicy.TOKEN_CRYPTO_KEY))
	return t, err
}

func Login(c echo.Context) error {

	user := new(model.UserLogin)
	if err := c.Bind(user); err != nil {
		fmt.Println("Error in there:")
		return err
	}
	fmt.Println("ROLE::", user.Role)
	passPolicy := config.GetPassPolicy()

	h_conf := s_conf{DB: repo.DbManager(), PassPolicy: &passPolicy}

	dt_user := model.User{}
	// role := user.Role

	// IP throtolling needs to be added here
	// get client ip check if number of tries
	// apply same logic here
	// check the ip is in throtolling mode

	// check if the user wants to login with specific role
	role, usernotfound := h_conf.getUserData(user, &dt_user)
	if !usernotfound {
		return echo.ErrUnauthorized
	}

	fmt.Println(dt_user)

	// check if the user is in lookout moode
	failpass_count, inThrottle := h_conf.isInThrottle(&dt_user)
	if inThrottle {
		return echo.ErrTooManyRequests
	}

	// validate password
	if !h_conf.isPasswordValid(failpass_count, user, &dt_user) {
		return echo.ErrUnauthorized
	}

	// generate claim
	token, err := h_conf.generateToken(role, &dt_user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	// check password expired
	if h_conf.isPasswordExpired(&dt_user) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"code":    20000,
			"message": "Password expired. Please change your password.",
			"token":   token,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":  20000,
		"token": token,
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
	db := repo.DbManager()
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
