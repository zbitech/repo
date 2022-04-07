package helper

import (
	"github.com/zbitech/common/pkg/id"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/model/object"
)

var (
	SECRET_KEY = "KaPdSgVkYp3s6v9y$B&E)H@McQeThWmZ"
)

func GenerateJwtToken(user entity.User) (*string, error) {

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, object.ZBIBasicClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  "ZBI",
			ExpiresAt: now.Add(time.Hour * 24).Unix(),
			Id:        id.GenerateRequestID(),
			IssuedAt:  now.Unix(),
			Issuer:    "ZBI",
			NotBefore: now.Unix(),
			Subject:   user.UserId,
		},
		Role:  user.Role,
		Email: user.Email,
	}) //TODO - add teams to claims?

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}
