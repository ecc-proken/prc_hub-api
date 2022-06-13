package events

import (
	"fmt"
	"net/http"
	"prc_hub-api/events"
	"prc_hub-api/flags"
	"strconv"
	"strings"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetById(c echo.Context) (err error) {
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

	// id
	idStr := c.Param("id")
	// string -> uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		c.Logger().Debug("404: cannot parse `:id`")
		return echo.ErrNotFound
	}

	// userを取得
	e, notFound, err := events.GetById(id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: event not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event not found"}, "	")
	}
	if !admin && !e.Published && (userId == nil || e.UserId != *userId) {
		// 403: Forbidden
		c.Logger().Debug("403: you cannot access this event")
		return c.JSONPretty(http.StatusForbidden, map[string]string{"message": "you cannot access this event"}, "	")
	}

	// 200: Success
	c.Logger().Debug("200: get event successful")
	return c.JSONPretty(http.StatusOK, e, "	")
}
