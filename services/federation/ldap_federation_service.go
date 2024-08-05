package federation

import "github.com/wissance/Ferrum/data"

type LdapUserFederation struct {
}

func (s *LdapUserFederation) GetUser(userName string, mask string) (data.User, error) {
	return nil, nil
}

func (s *LdapUserFederation) Init() {

}
