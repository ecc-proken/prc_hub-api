package users

import (
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/users"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(c echo.Context) (err error) {
	// リクエストボディをバインド
	p := new(users.VerifyPostBody)
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

	// emailをもとにuserを取得
	u, notFound, err := users.GetByEmail(p.Email)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 未登録のメールアドレス
		// 404: Not found
		c.Logger().Debug("404: user not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}
	verify, err := u.Verify(p.Password)
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if !verify {
		// 不正なパスワード
		// 403: Forbidden
		c.Logger().Debug("403: incorrect password")
		return c.JSONPretty(http.StatusForbidden, map[string]string{"message": "incorrect password"}, "	")
	}

	// トークンを生成
	t, err := jwt.GenerateToken(
		u,
		*flags.Get().JwtIssuer,
		*flags.Get().JwtSecret,
	)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// Cookieを追加
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    t,
		HttpOnly: true,
	})

	// 200: Success
	c.Logger().Debug("200: signin successful")
	return c.JSONPretty(
		http.StatusOK,
		map[string]string{"token": t},
		"	",
	)
}
