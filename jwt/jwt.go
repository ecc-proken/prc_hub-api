package jwt

import (
	"errors"
	"prc_hub-api/users"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtCustumClaims struct {
	Id            uint64 `json:"id"`
	Email         string `json:"email"`
	Admin         bool   `json:"admin"`
	MigratedAdmin bool   `json:"migrated_admin"`
	jwt.StandardClaims
}

func GenerateToken(user users.User, issuer string, secret string) (token string, err error) {
	// 有効期限
	expAt := time.Now().Add(time.Hour * 72).Unix()
	if user.MigrateAdmin {
		expAt = time.Now().Add(time.Hour * 8).Unix()
	}

	// Set custom claims
	claims := &JwtCustumClaims{
		user.Id,
		user.Email,
		user.Admin,
		user.MigrateAdmin,
		jwt.StandardClaims{
			ExpiresAt: expAt,
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}

	// トークンを生成
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString([]byte(secret))
}

func CheckToken(issuer string, token *jwt.Token) (claims *JwtCustumClaims, err error) {
	claims = token.Claims.(*JwtCustumClaims)

	if !claims.VerifyIssuer(issuer, true) {
		// 不正なトークン
		err = errors.New("invalid token")
		return
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		// 期限切れのトークン
		err = errors.New("token expired")
		return
	}

	return
}
