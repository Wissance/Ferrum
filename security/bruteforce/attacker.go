package bruteforce

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type attackerList struct {
	mutex                *sync.RWMutex
	attackersIPAddresses map[string]uuid.UUID
	attackersDevices     map[string]uuid.UUID
	attackersStats       map[uuid.UUID]AttackerStats
}

type AttackerStats struct {
	firstAttackDetection time.Time
	attackCount          int64
	blockedAt            time.Time
	blockTill            time.Time
}

/*func (a attacker) Equals(other attacker) bool {
	// attacker use browser with specific deviceId
	if a.DeviceId != nil && other.DeviceId != nil {
		return *a.DeviceId == *other.DeviceId
	}
	// attacker did not used browser, but IP could be equal
	return a.IPAddress == other.IPAddress
}*/

func createAttackerList() *attackerList {
	return &attackerList{
		mutex:                &sync.RWMutex{},
		attackersDevices:     map[string]uuid.UUID{},
		attackersIPAddresses: map[string]uuid.UUID{},
		attackersStats:       map[uuid.UUID]AttackerStats{},
	}
}

// GetAttackerStats is a function for searching the stats
func (attackers *attackerList) getAttackerStats(deviceId string, ipAddress string) *AttackerStats {
	if deviceId == "" && ipAddress == "" {
		return nil
	}
	var id uuid.UUID
	var ok bool
	var stats AttackerStats
	attackers.mutex.RLock()
	if deviceId != "" {
		id, ok = attackers.attackersDevices[deviceId]
	}
	if !ok {
		id, ok = attackers.attackersDevices[deviceId]
	}

	if ok {
		stats = attackers.attackersStats[id]
	}
	attackers.mutex.RUnlock()
	return &stats
}
