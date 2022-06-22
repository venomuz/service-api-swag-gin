package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/logger"
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
}
