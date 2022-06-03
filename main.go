package main

import (
	"fmt"
	"prc_hub-api/flags"
	"prc_hub-api/migration"
	"prc_hub-api/mysql"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func main() {
	f := flags.Get()

	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: int(*f.GzipLevel),
	}))
	e.Logger.SetLevel(log.Lvl(*f.LogLevel))
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup db client instance
	e.Logger.Info(mysql.SetDSNTCP(*f.MysqlUser, *f.MysqlPasswd, *f.MysqlHost, int(*f.MysqlPORT), *f.MysqlDB))

	// Migrate
	adminFound, invalidEmail, usedEmail, err := migration.MigrateAdminUser(*f.AdminEmail, *f.AdminPasswd)
	if err != nil {
		e.Logger.Fatal(err.Error())
		return
	}
	if !adminFound && invalidEmail && usedEmail {
		e.Logger.Fatalf("Admin email already used or invalid. %s", *f.AdminEmail)
	}
	if !adminFound && !invalidEmail && !usedEmail {
		e.Logger.Info("Migrate admin user successful.")
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *f.Port)))
}
