package bruteforce

import (
	"context"
	"github.com/wissance/Ferrum/logging"
)

// ListBasedProtectionService is bruteforce protection service that uses lists of blocked ip addresses/devices
type ListBasedProtectionService struct {
	attackersList *attackerList
	ctx           context.Context
	logger        *logging.AppLogger
}

func CreateListBasedProtectionService(ctx context.Context, watchTime int, logger *logging.AppLogger) *ListBasedProtectionService {
	// block time actually not using, we are using permanent block until app is restarted
	attackersList := createAttackerList(ctx, 86400, watchTime)
	return &ListBasedProtectionService{
		attackersList: attackersList,
		ctx:           ctx,
		logger:        logger,
	}
}

func (service *ListBasedProtectionService) RegisterIpAddressAttempt(ipAddress string) bool {
	err := service.attackersList.UpsertIpAddressStats(ipAddress)
	if err != nil {
		// todo(UMV): wrap error in the future
		service.logger.Error(err.Error())
		return false
	}

	stats := service.attackersList.GetAttackerStats("", ipAddress)
	// no stats, wasn't blocked
	if stats == nil {
		return true
	}
	return stats.blocked
}

func (service *ListBasedProtectionService) RegisterDeviceAttempt(deviceId string) bool {
	err := service.attackersList.UpsertDeviceStats(deviceId)
	if err != nil {
		// todo(UMV): wrap error in the future
		service.logger.Error(err.Error())
		return false
	}

	stats := service.attackersList.GetAttackerStats(deviceId, "")
	// no stats, wasn't blocked
	if stats == nil {
		return true
	}
	return stats.blocked
}

func (service *ListBasedProtectionService) BlockDevice(deviceId string) {
	// this err is just more informational neither real error
	err := service.attackersList.setDeviceIdBlockedStatus(deviceId, true)
	if err != nil {
		service.logger.Warn(err.Error())
	}
}

func (service *ListBasedProtectionService) BlockIpAddress(ipAddress string) {
	// this err is just more informational neither real error
	err := service.attackersList.setIpAddressBlockedStatus(ipAddress, true)
	if err != nil {
		service.logger.Warn(err.Error())
	}
}

func (service *ListBasedProtectionService) IsIpAddressBlocked(ipAddress string) bool {
	stats := service.attackersList.GetAttackerStats("", ipAddress)
	if stats == nil {
		return true
	}
	return stats.blocked
}

func (service *ListBasedProtectionService) IsDeviceBlocked(deviceId string) bool {
	stats := service.attackersList.GetAttackerStats(deviceId, "")
	if stats == nil {
		return true
	}
	return stats.blocked
}

func (service *ListBasedProtectionService) UnblockIpAddress(ipAddress string) {
	// this err is just more informational neither real error
	err := service.attackersList.setIpAddressBlockedStatus(ipAddress, false)
	if err != nil {
		service.logger.Warn(err.Error())
	}
}

func (service *ListBasedProtectionService) UnblockDevice(deviceId string) {
	// this err is just more informational neither real error
	err := service.attackersList.setDeviceIdBlockedStatus(deviceId, true)
	if err != nil {
		service.logger.Warn(err.Error())
	}
}
