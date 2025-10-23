package data

// ServerSettings is Ferrum server settings that contains main security settings
type ServerSettings struct {
	// Admin is a user that is using for Ferrum configure
	Admin AdminUser `json:"admin"`
	// AllowedHosts is a list of Hosts from which Ferrum could be configured, it allows to use *,
	// if * is in this list it means that all hosts are allowed
	AllowedHosts []string `json:"allowed_hosts"`
	// Prefix before admin api URL
	AdminApiUrlPrefix string `json:"url_prefix"`
}
