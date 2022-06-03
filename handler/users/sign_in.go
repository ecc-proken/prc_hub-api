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
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// リクエストボディを検証
	if err = c.Validate(p); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// emailをもとにuserを取得
	u, notFound, err := users.GetByEmail(p.Email)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 未登録のメールアドレス
		// 404: Not found
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}
	verify, err := u.Verify(p.Password)
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if !verify {
		// 不正なパスワード
		// 403: Forbidden
		c.Logger().Debug("failed to sign in")
		return echo.ErrForbidden
	}

	// トークンを生成
	t, err := jwt.GenerateToken(
		u,
		*flags.Get().JwtIssuer,
		*flags.Get().JwtSecret,
	)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// Cookieを追加
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    t,
		HttpOnly: true,
	})

	// 200: Success
	return c.JSONPretty(
		http.StatusOK,
		map[string]string{"token": t},
		"	",
	)
}
