package client

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CloudClient interface {
	Login(string, string) (*types.AuthenticationResultType, error)
	DBInterface
}
