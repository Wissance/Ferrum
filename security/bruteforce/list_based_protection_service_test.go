package bruteforce

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/wissance/Ferrum/logging"
	"testing"
	"time"
)

func TestRegisterAttemptsAndCheckIsBlocked(t *testing.T) {
	logger := logging.AppLogger{}
	cfg := ProtectionServiceConfig{
		WatchTimeSec: 60,
	}
	protectionService := CreateProtectionService(context.Background(), &cfg, &logger)
	testCases := []struct {
		name      string
		attackers []struct {
			ipAddress     string
			deviceId      string
			attacks       int
			expectedBlock bool
		}
	}{
		{
			name: "attackers both from an IP addresses and devices",
			attackers: []struct {
				ipAddress     string
				deviceId      string
				attacks       int
				expectedBlock bool
			}{
				{
					ipAddress:     "199.90.178.54",
					deviceId:      "",
					attacks:       109,
					expectedBlock: true,
				},
				{
					ipAddress:     "153.108.162.18",
					deviceId:      "",
					attacks:       27,
					expectedBlock: false,
				},
				{
					ipAddress:     "122.201.104.32",
					deviceId:      "",
					attacks:       95,
					expectedBlock: false,
				},
				{
					ipAddress:     "",
					deviceId:      "somedevice-1234567890",
					attacks:       133,
					expectedBlock: true,
				},
			},
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			for _, a := range tCase.attackers {
				for range a.attacks {
					if a.deviceId == "" {
						go protectionService.RegisterIpAddressAttempt(a.ipAddress)
					} else {
						go protectionService.RegisterDeviceAttempt(a.deviceId)
					}
				}
			}
			// Wait pause until all goroutines are completed
			time.Sleep(time.Duration(5) * time.Second)

			for _, a := range tCase.attackers {
				var isBlocked bool
				if a.deviceId == "" {
					isBlocked = protectionService.IsIpAddressBlocked(a.ipAddress)
				} else {
					isBlocked = protectionService.IsDeviceBlocked(a.deviceId)
				}
				assert.Equal(t, a.expectedBlock, isBlocked)
			}
		})
	}

}

func TestBlockAttackersManuallyAndAutoRemoveSomeOfThem(t *testing.T) {
	ipAddress1 := "122.130.205.221"
	ipAddress2 := "192.100.205.221"
	deviceId1 := "device_one"
	deviceId2 := "device_two"
	logger := logging.AppLogger{}
	cfg := ProtectionServiceConfig{
		WatchTimeSec: 10,
	}
	protectionService := CreateProtectionService(context.Background(), &cfg, &logger)
	protectionService.RegisterIpAddressAttempt(ipAddress1)
	protectionService.RegisterIpAddressAttempt(ipAddress2)
	protectionService.RegisterDeviceAttempt(deviceId1)
	protectionService.RegisterDeviceAttempt(deviceId2)
	assert.Equal(t, 4, protectionService.GetWatchingAttackersCount())
	// do manual block after that they can't be removed
	protectionService.BlockIpAddress(ipAddress1)
	protectionService.BlockDevice(deviceId2)
	// wait period to skip auto-remove
	time.Sleep(time.Duration(cfg.WatchTimeSec+2) * time.Second)
	assert.Equal(t, true, protectionService.IsIpAddressBlocked(ipAddress1))
	assert.Equal(t, false, protectionService.IsIpAddressBlocked(ipAddress2))
	assert.Equal(t, false, protectionService.IsDeviceBlocked(deviceId1))
	assert.Equal(t, true, protectionService.IsDeviceBlocked(deviceId2))
	assert.Equal(t, 2, protectionService.GetWatchingAttackersCount())
	// unblock and wait (however it is hard to check
	protectionService.UnblockIpAddress(ipAddress1)
	protectionService.UnblockDevice(deviceId2)
	time.Sleep(time.Duration(cfg.WatchTimeSec+2) * time.Second)
	assert.Equal(t, false, protectionService.IsIpAddressBlocked(ipAddress1))
	assert.Equal(t, false, protectionService.IsIpAddressBlocked(ipAddress2))
	assert.Equal(t, false, protectionService.IsDeviceBlocked(deviceId1))
	assert.Equal(t, false, protectionService.IsDeviceBlocked(deviceId2))
	assert.Equal(t, 0, protectionService.GetWatchingAttackersCount())
}
