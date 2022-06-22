package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/venomuz/service-api-swag-gin/ApiGateway/pkg/logger"
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
func (QjwtHendler *JwtHendler) GenerateAuthJWT() (access, refresh string, err error) {
	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodES256)
	refreshToken = jwt.New(jwt.SigningMethodES256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["iss"] = QjwtHendler.Iss
	claims["sub"] = QjwtHendler.Sub
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["role"] = QjwtHendler.Role
	claims["aud"] = QjwtHendler.Aud
	access, err = accessToken.SignedString([]byte(QjwtHendler.SigninKey))
	if err != nil {
		QjwtHendler.Log.Error("error generating access token", logger.Error(err))
		return
	}
	refresh, err = refreshToken.SignedString([]byte(QjwtHendler.SigninKey))
	if err != nil {
		QjwtHendler.Log.Error("error generating refresh token", logger.Error(err))
		return
	}
	return access, refresh, nil
}
