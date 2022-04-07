package iam

import (
	"context"
	"github.com/zbitech/repo/pkg/iam/auth"
	"github.com/zbitech/repo/pkg/iam/basic"
	"github.com/zbitech/repo/pkg/iam/jwtsvr"

	"github.com/zbitech/common/interfaces"
)

type BasicAuthorizationFactory struct {
	jwtServer        interfaces.JwtServerIF
	accessAuthorizer interfaces.AccessAuthorizerIF
	iamService       interfaces.IAMServiceIF
}

func NewBasicAuthorizationFactory() interfaces.AuthorizationFactoryIF {
	return &BasicAuthorizationFactory{}
}

func (j *BasicAuthorizationFactory) Init(ctx context.Context) error {
	//TODO - introduce logic to determine which JWT Server to use
	j.jwtServer = jwtsvr.NewZBIJwtServer()

	//TODO - introduce logic to determine which IAM service to use
	j.iamService = basic.NewBasicIAMService(j.jwtServer)
	j.accessAuthorizer = auth.NewAccessAuthorizer(j.iamService)

	return nil
}

func (j *BasicAuthorizationFactory) GetJwtServer() interfaces.JwtServerIF {
	return j.jwtServer
}

//func (j *BasicAuthorizationFactory) GetJwtTokenVerifierIF() interfaces.JwtTokenVerifierIF {
//	return j.tokenVerifier
//}

func (j *BasicAuthorizationFactory) GetAccessAuthorizer() interfaces.AccessAuthorizerIF {
	return j.accessAuthorizer
}

func (j *BasicAuthorizationFactory) GetIAMService() interfaces.IAMServiceIF {
	return j.iamService
}
