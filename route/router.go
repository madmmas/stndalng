package route

import (
	"stndalng/api"
	"stndalng/repo"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	repo.InitDB()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// public access
	p := e.Group("/api")

	p.POST("/login", api.Login)
	// p.POST("/logout", api.Logout)

	// restricted access
	r := e.Group("/api")

	config := middleware.JWTConfig{
		Claims:     &api.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))

	r.POST("/role", api.NewRole)
	r.PUT("/role", api.UpdateRole)
	r.GET("/roles", api.GetRoles)
	r.GET("/role/:id", api.GetRole)
	r.GET("/getroles", api.GetRolesKV)

	r.GET("/userinfo", api.UserInfo)
	r.GET("/users", api.GetUsers)
	r.GET("/user/:id", api.GetUser)
	r.POST("/user", api.NewUser)
	r.PUT("/user", api.UpdateUser)
	r.GET("/deusers", api.GetDeUsers)
	r.DELETE("/user/:id", api.DeleteUser)
	r.PUT("/change_password", api.ChangePassword)
	r.PUT("/deactive_user/:id", api.ChangeUserDeactiveFlag)
	r.PUT("/active_user/:id", api.ChangeUserActiveFlag)

	return e
}
