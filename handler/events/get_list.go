package events

import (
	"fmt"
	"net/http"
	"prc_hub-api/events"
	"prc_hub-api/flags"
	"strings"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Get(c echo.Context) (err error) {
	// jwtトークン確認
	var (
		userId *uint64 = nil
		admin  bool    = false
	)
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		var token *jwtGo.Token
		token, err = jwtGo.Parse(tokenString, func(token *jwtGo.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(*flags.Get().JwtSecret), nil
		})
		if err == nil {
			if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid && claims["id"].(float64) >= 1 {
				tmpUserId := uint64(claims["id"].(float64))
				userId = &tmpUserId
				admin = claims["admin"].(bool)
			}
		}
	}

	// リクエストボディをバインド
	q := new(events.GetQuery)
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

	// eventを取得
	events, err := events.Get(*q, userId, admin)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	if events == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, events, "	")
}
