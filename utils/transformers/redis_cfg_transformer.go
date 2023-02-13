package transformers

import (
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"github.com/wissance/Ferrum/config"
)

// TransformRedisConfig functions that transforms internal config to redis.Options that is required to establish connection
func TransformRedisConfig(redisCfg *config.RedisConfig) (*redis.Options, error) {
	// 1. Creation of minimal required to connect options (address, db, password)
	opts := redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.DbNumber,
	}

	// 2. Configure pool
	opts.PoolFIFO = true
	opts.PoolSize = int(redisCfg.PoolSize)

	// 3. Configure timeouts && retry && limitation

	// 4. TLS Configuration (further)
	opts.TLSConfig = &tls.Config{}

	return &opts, nil
}
