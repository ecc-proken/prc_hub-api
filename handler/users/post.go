package users

import (
	"encoding/json"
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/users"

	"github.com/labstack/echo/v4"
)

func Post(c echo.Context) (err error) {
	// リクエストボディをバインド
	p := new(users.PostBody)
	if err = c.Bind(p); err != nil {
		// 400: Bad request
		c.Logger().Debug("400: " + err.Error())
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// リクエストボディを検証
	if err = c.Validate(p); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug("422: " + err.Error())
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// 書込
	u, invalidEmail, usedEmail, err := users.Post(*p)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if invalidEmail {
		// 422: Unprocessable entity
		c.Logger().Debug("422: invalid email")
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": "invalid email"}, "	")
	}
	if usedEmail {
		// 400: Bad request
		c.Logger().Debug("400: email already used")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "email already used"}, "	")
	}

	// トークンを生成
	t, err := jwt.GenerateToken(u, *flags.Get().JwtIssuer, *flags.Get().JwtSecret)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// jsonにjwtトークンを追加
	b, err := json.Marshal(u)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	m := map[string]interface{}{"token": t}
	err = json.Unmarshal(b, &m)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	c.Logger().Debug("200: create user successful")
	return c.JSONPretty(http.StatusOK, m, "	")
}
