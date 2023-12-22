package config

import "crypto/tls"

// todo (UMV): probably we don't need RedisConfig, because we are using map 4 this

// RedisConfig is a simplified redis.Options config
type RedisConfig struct {
	Address   string `json:"address" example:"localhost:6379"`
	Password  string `json:"password"`
	DbNumber  int    `json:"db_number"`
	Namespace string `json:"namespace"`
	// MaxRetries is a number of attempts to
	MaxRetries int `json:"max_retries"`
	// MinRetryBackoff is a backoff in milliseconds
	MinRetryBackoff int `json:"min_retry_backoff"`
	MaxRetryBackoff int `json:"max_retry_backoff"`
	// go-redis dial timeout option is time, here we simplify config we assume here Seconds as a time value, 0 means no timeout
	DialTimeout uint `json:"dial_timeout"`
	// go-redis read timeout option is time, here we simplify config we assume here Seconds as a time value, 0 means no timeout
	ReadTimeout uint `json:"read_timeout"`
	// go-redis write timeout option is time, here we simplify config we assume here Seconds as a time value, 0 means no timeout
	WriteTimeout uint `json:"write_timeout"`
	PoolSize     uint `json:"pool_size"`
	// go-redis pool timeout option is time, here we simplify config we assume here Seconds as a time value, 0 means 1 sec to pool timeout
	PoolTimeout uint `json:"pool_timeout"`
	MinIdleConn int  `json:"min_idle_conn"`
	MaxIdleConn int  `json:"max_idle_conn"`
	// default is 30 min (go-redis), 0 disable max timeout (pass -1 to go-redis)
	ConnIdleTimeout uint `json:"conn_idle_timeout"`
	// by default go-redis do not close idle connections
	ConnMaxLifetimeTimeout int         `json:"conn_max_lifetime_timeout"`
	ReadOnlySlaveEnabled   bool        `json:"read_only_slave_enabled"`
	TlsCfg                 *tls.Config `json:"tls_cfg"`
}
