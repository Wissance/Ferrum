package config

// RedisOptionConfig is a simplified redis.Options config
type RedisOptionConfig struct {
	DbNumber int `json:"db_number"`
	// MaxRetries is a number of attempts to
	MaxRetries int `json:"max_retries"`
}
