package federation

import (
	"errors"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
)

type LdapUserFederation struct {
	logger *logging.AppLogger
}

func CreateLdapUserFederationService(config *data.UserFederationServiceConfig, logger *logging.AppLogger) (*LdapUserFederation, error) {
	return &LdapUserFederation{logger: logger}, nil
}

func (s *LdapUserFederation) GetUser(userName string, mask string) (data.User, error) {
	return nil, nil
}

func (s *LdapUserFederation) GetUsers(mask string) []data.User {
	return []data.User{}
}

func (s *LdapUserFederation) Authenticate(userName string, password string) (bool, error) {
	return false, errors.New("not implemented yet")
}

func (s *LdapUserFederation) Init() {

}
