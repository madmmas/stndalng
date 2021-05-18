package api

import (
	"fmt"
	"net/http"
	"stndalng/model"
	"stndalng/repo"

	"github.com/labstack/echo"
)

func GetRoles(c echo.Context) error {
	// check if root
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()
	Role := []model.Role{}
	db.Select("id, name, description").Find(&Role)
	return c.JSON(http.StatusOK, Role)
}

func GetRole(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()
	dt := model.Role{}
	db.Where("id = ?", c.Param("id")).First(&dt)
	return c.JSON(http.StatusOK, dt)
}

func GetRolesKV(c echo.Context) error {
	// check if root
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()
	Role := []model.Role{}
	db.Select("id, name, description").Find(&Role)
	m := make(map[string]string)

	for _, role := range Role {
		m[role.Name] = role.Name
	}

	return c.JSON(http.StatusOK, m)
}

func NewRole(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	db := repo.DbManager()

	dt := new(model.Role)
	if err := c.Bind(dt); err != nil {
		return err
	}
	db.Create(&dt)

	return c.JSON(http.StatusCreated, map[string]string{
		"ID": dt.ID.String(),
	})
}

func UpdateRole(c echo.Context) error {
	if !IsRoot(c) {
		return echo.ErrUnauthorized
	}

	dt := new(model.Role)
	if err := c.Bind(dt); err != nil {
		return err
	}
	fmt.Println(dt.ID)
	db := repo.DbManager()

	dt_f := model.Role{}
	if db.Where("id = ?", dt.ID).First(&dt_f).RecordNotFound() {
		// record not found
		return c.JSON(http.StatusCreated, map[string]string{
			"Msg": "Not Found",
			"ID":  dt.ID.String(),
		})
	} else {
		fmt.Println(dt)
		db.Model(&dt).Updates(dt)

		return c.JSON(http.StatusCreated, map[string]string{
			"ID": dt.ID.String(),
		})
	}

}
