package federation

import (
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
)

// UserFederation is interface to external User Storage systems (AD, LDAP or FreeIPA)
/* UserFederation instances are classes that have config to connect external providers
 * and Authenticate in system using this provider
 */
type UserFederation interface {
	// GetUser searches for User in external Provider and return data.User mapped with mask (jsonpath)
	GetUser(userName string, mask string) (data.User, error)
	// GetUsers searches for Users in external Provider and return []data.User mapped with mask (jsonpath)
	GetUsers(mask string) []data.User
	// Authenticate method for Authenticate in external Provider
	Authenticate(userName string, password string) (bool, error)
}

// CreateUserFederationService is a factory method that creates
func CreateUserFederationService(config *data.UserFederationServiceConfig, logger *logging.AppLogger) (UserFederation, error) {
	if config.Type == data.LDAP {
		s, err := CreateLdapUserFederationService(config, logger)
		if err != nil {
			logger.Error(sf.Format("An error occurred during Ldap User Federation service creation: {0}", err.Error()))
			return nil, err
		}
		return s, nil
	}
	return nil, nil
}
