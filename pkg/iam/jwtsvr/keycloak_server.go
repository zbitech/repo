package jwtsvr

import (
	"github.com/golang-jwt/jwt"
)

type KeyCloakJwt struct {
	publicKey string
}

type KeyCloakClaims struct {
	Roles []string
	Orgs  []string
	jwt.StandardClaims
}

func NewKeyCloakJwt(publicKey string) KeyCloakJwt {
	return KeyCloakJwt{
		publicKey: publicKey,
	}
}

func (s *KeyCloakJwt) GetKey() (interface{}, error) {
	return jwt.ParseRSAPublicKeyFromPEM([]byte(s.publicKey))
}

func (s *KeyCloakJwt) ValidateToken(token *jwt.Token) error {

	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return jwt.ErrInvalidKey
	}

	return nil
}

func (s *KeyCloakJwt) GetPayload() jwt.Claims {
	return KeyCloakClaims{}
}
