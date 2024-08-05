package federation

import "github.com/wissance/Ferrum/data"

// UserFederation is interface to external User Storage systems (AD, LDAP or FreeIPA)
/* UserFederation instances are classes that have config to connect external providers
 * and Authenticate in system using this provider
 */
type UserFederation interface {
	// GetUser searches for User in external Provider and return data.User mapped with mask (jsonpath)
	GetUser(userName string, mask string) (data.User, error)
	// GetUsers searches for Users in external Provider and return []data.User mapped with mask (jsonpath)
	GetUsers(mask string) []data.User
}
