package config

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GlobalSecurityConfig struct {
	AllowedHosts      []string  `json:"allowed_hosts"`
	AdminApiUrlPrefix string    `json:"admin_api_url_prefix"`
	Admin             AdminUser `json:"admin"`
}
