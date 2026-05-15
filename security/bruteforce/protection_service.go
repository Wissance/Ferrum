package bruteforce

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
	// BlockSender functions that adds record for blocking access to the list
	BlockSender(ipAddress string, deviceId string)
	// IsSenderBlocked checks that sender was blocked
	IsSenderBlocked(ipAddress string) bool
	// IsSenderDeviceBlocked checks that sender device was blocked
	IsSenderDeviceBlocked(deviceId string) bool
	// UnblockSender is a function for direct Sender unblock
	UnblockSender(ipAddress string)
	// UnblockDevice is a function for direct Device unblock
	UnblockDevice(ipAddress string)
}
