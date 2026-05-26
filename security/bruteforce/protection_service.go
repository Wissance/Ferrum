package bruteforce

import (
	"context"
	"github.com/wissance/Ferrum/logging"
)

// ProtectionServiceConfig contains a set of properties that defines bruteforce protection behaviour
type ProtectionServiceConfig struct {
	// WatchTimeSec is set a period for inspect potential attackers
	WatchTimeSec int
}

// ProtectionService is a struct that is using for
/* Service manages Sender (someone who sends request on issue token or|and auth)
 * Sender is representing by the following fingerprint
 *     1. Combination of IP addresses via headers X-Forwarded-For or X-Real-IP (for gin could be obtained
 *        through ClientIP())
 *     2. X-Device-ID header send by frontend (usually identifies Browser)
 * Useful articles:
 *     1. https://www.sobyte.net/post/2021-09/gin-get-client-real-ip/
 */
type ProtectionService interface {
	// RegisterIpAddressAttempt function that register attempt to enter non-valid credentials from ipAddress
	RegisterIpAddressAttempt(ipAddress string) bool
	// RegisterDeviceAttempt function that register attempt to enter non-valid credentials from specific deviceId
	RegisterDeviceAttempt(deviceId string) bool
	// BlockDevice functions that adds record for blocking access for deviceId to the list
	BlockDevice(deviceId string)
	// BlockIpAddress functions that adds record for blocking access for ipAddress to the list
	BlockIpAddress(ipAddress string)
	// IsIpAddressBlocked checks that sender was blocked
	IsIpAddressBlocked(ipAddress string) bool
	// IsDeviceBlocked checks that sender device was blocked
	IsDeviceBlocked(deviceId string) bool
	// UnblockIpAddress is a function for direct IP Address unblock
	UnblockIpAddress(ipAddress string)
	// UnblockDevice is a function for direct Device unblock
	UnblockDevice(ipAddress string)
	// GetWatchingAttackersCount returns number of watching attackers
	GetWatchingAttackersCount() int
}

// CreateProtectionService is a factory function that build implementation of ProtectionService
func CreateProtectionService(ctx context.Context, config *ProtectionServiceConfig,
	logger *logging.AppLogger) ProtectionService {
	return CreateListBasedProtectionService(ctx, config.WatchTimeSec, logger)
}
