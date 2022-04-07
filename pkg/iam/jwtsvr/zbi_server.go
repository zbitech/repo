package jwtsvr

import (
	"github.com/golang-jwt/jwt"
	"github.com/zbitech/common/interfaces"
	"github.com/zbitech/common/pkg/model/object"
	"github.com/zbitech/common/pkg/model/ztypes"
	"github.com/zbitech/repo/internal/helper"
)

type ZBIJwtServer struct {
	secretKey string
}

func NewZBIJwtServer() interfaces.JwtServerIF {
	return &ZBIJwtServer{
		secretKey: helper.SECRET_KEY,
	}
}

func (s *ZBIJwtServer) GetKey() (interface{}, error) {
	return []byte(s.secretKey), nil
}

func (s *ZBIJwtServer) ValidateToken(token *jwt.Token) error {

	_, ok := token.Method.(*jwt.SigningMethodHMAC)

	if !ok {
		//		log.Printf("Method %s is not valid", method)
		return jwt.ErrInvalidKey
	}

	return nil
}

func (s *ZBIJwtServer) GetPayload() jwt.Claims {
	return &object.ZBIBasicClaims{}
}

func (s *ZBIJwtServer) GetUserId(claim jwt.Claims) string {
	zbiClaims := claim.(*object.ZBIBasicClaims)
	return zbiClaims.Subject
}

func (s *ZBIJwtServer) GetEmail(claim jwt.Claims) string {
	zbiClaims := claim.(*object.ZBIBasicClaims)
	return zbiClaims.Email
}

func (s *ZBIJwtServer) GetRole(claim jwt.Claims) ztypes.Role {
	zbiClaims := claim.(*object.ZBIBasicClaims)
	return zbiClaims.Role
}
