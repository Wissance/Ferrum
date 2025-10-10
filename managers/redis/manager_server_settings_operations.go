package redis

import "github.com/wissance/Ferrum/data"

// GetServerSettings function that returns ServerSettings
func (mn *RedisDataManager) GetServerSettings() (data.ServerSettings, error) {
	return data.ServerSettings{}, nil
}

// SetServerSettings function that updates ServerSettings by full new settings replace
func (mn *RedisDataManager) SetServerSettings(settings data.ServerSettings) error {

	return nil
}
