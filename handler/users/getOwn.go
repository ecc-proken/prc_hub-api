package users

import (
	"errors"
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/users"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Get(c echo.Context) (err error) {
	// jwtトークン確認
	t := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, t)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}
	id := claims.Id

	if claims.Admin {
		// リクエストボディをバインド
		q := new(users.GetQuery)
		if err = c.Bind(q); err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}

		// リクエストボディを検証
		if err = c.Validate(q); err != nil {
			// 422: Unprocessable entity
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
		}

		// userを取得
		users, err := users.Get(*q)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// 200: Success
		return c.JSONPretty(http.StatusOK, users, "	")
	} else {
		// userを取得
		u, notFound, err := users.GetById(id)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if notFound {
			// 404: Not found
			c.Logger().Debug(errors.New("user not found"))
			return echo.ErrNotFound
		}

		// 200: Success
		return c.JSONPretty(http.StatusOK, []interface{}{u}, "	")
	}
}
