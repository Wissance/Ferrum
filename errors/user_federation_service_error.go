package errors

import sf "github.com/wissance/stringFormatter"

type FederatedUserNotFoundError struct {
	FederationType string
	Name           string
	Url            string
	Username       string
}

type MultipleUserResultError struct {
	Name     string
	Username string
}

func (e FederatedUserNotFoundError) Error() string {
	return sf.Format("User: \"{0}\" was n't found in service {1} of type {2} by url: \"{3}\"",
		e.Username, e.Name, e.FederationType, e.Url)
}

func NewFederatedUserNotFound(federationType string, name string, url string, username string) FederatedUserNotFoundError {
	return FederatedUserNotFoundError{
		FederationType: federationType,
		Name:           name,
		Url:            url,
		Username:       username,
	}
}

func (e MultipleUserResultError) Error() string {
	return sf.Format("Multiple federated user with name: \"{0}\" for federation service: {1}", e.Username, e.Name)
}

func NewMultipleUserResultError(name string, username string) MultipleUserResultError {
	return MultipleUserResultError{
		Name:     name,
		Username: username,
	}
}
