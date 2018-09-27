package digitaloceanSDK

import (
	"github.com/digitalocean/godo"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"golang.org/x/oauth2"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

//SDK encompasses all authorization for working with digital ocean api
type SDK struct {
	accessToken string
}

//New is a constructor function for SDK
func New(accessToken string) *SDK {
	return &SDK{
		accessToken: accessToken,
	}
}

//NewFromAccount extracts credentials from accounts and returns SDK ready to be used
func NewFromAccount(account *model.CloudAccount) (*SDK, error) {
	token, ok := account.Credentials[clouds.DigitalOceanAccessToken]
	if !ok {
		return nil, sgerrors.ErrInvalidCredentials
	}
	return New(token), nil
}

//GetClient return digital ocean client
func (s *SDK) GetClient() *godo.Client {
	token := &TokenSource{
		AccessToken: s.accessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, token)
	return godo.NewClient(oauthClient)
}
