package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/logger"
	"time"
)

type JwtHendler struct {
	Sub       string
	Iss       string
	Exp       string
	Iat       string
	Aud       []string
	Role      string
	Token     string
	SigninKey string
	Log       logger.Logger
}

//GenerateAuthJWT ...
func (JwtHendler *JwtHendler) GenerateAuthJWT() (access, refresh string, err error) {
	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodES256)
	refreshToken = jwt.New(jwt.SigningMethodES256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["iss"] = JwtHendler.Iss
	claims["sub"] = JwtHendler.Sub
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["role"] = JwtHendler.Role
	claims["aud"] = JwtHendler.Aud
	access, err = accessToken.SignedString([]byte(JwtHendler.SigninKey))
	if err != nil {
		JwtHendler.Log.Error("error generating access token", logger.Error(err))
		return
	}
	refresh, err = refreshToken.SignedString([]byte(JwtHendler.SigninKey))
	if err != nil {
		JwtHendler.Log.Error("error generating refresh token", logger.Error(err))
		return
	}
}
