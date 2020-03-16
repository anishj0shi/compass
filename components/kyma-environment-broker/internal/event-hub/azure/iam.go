package azure

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

var (
	armAuthorizer autorest.Authorizer
)

type OAuthGrantType int

const OAuthGrantTypeServicePrincipal OAuthGrantType = iota

func grantType() OAuthGrantType {
	return OAuthGrantTypeServicePrincipal
}

func GetResourceManagementAuthorizer(config *Config) (autorest.Authorizer, error) {
	if armAuthorizer != nil {
		return armAuthorizer, nil
	}

	var a autorest.Authorizer
	var err error

	environment, err := config.Environment()
	if err != nil {
		return nil, err
	}
	a, err = getAuthorizerForResource(config, grantType(), environment.ResourceManagerEndpoint)

	if err == nil {
		// cache
		armAuthorizer = a
	} else {
		// clear cache
		armAuthorizer = nil
	}
	return armAuthorizer, err
}

func getAuthorizerForResource(config *Config, grantType OAuthGrantType, resource string) (autorest.Authorizer, error) {
	var a autorest.Authorizer
	var err error

	switch grantType {
	case OAuthGrantTypeServicePrincipal:
		environment, err := config.Environment()
		if err != nil {
			return nil, err
		}
		oauthConfig, err := adal.NewOAuthConfig(environment.ActiveDirectoryEndpoint, config.tenantID)
		if err != nil {
			return nil, err
		}

		token, err := adal.NewServicePrincipalToken(*oauthConfig, config.clientID, config.clientSecret, resource)
		if err != nil {
			return nil, err
		}
		a = autorest.NewBearerAuthorizer(token)

	default:
		return a, fmt.Errorf("invalid grant type specified")
	}

	return a, err
}
