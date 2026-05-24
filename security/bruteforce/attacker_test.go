package bruteforce

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpsertIpAddressStats(t *testing.T) {
	testCases := []struct {
		name           string
		ipAddress      string
		produceAttacks int
		expectBlock    bool
	}{
		{
			name:           "Sequential attack from one address and block as a result",
			ipAddress:      "192.167.99.144",
			produceAttacks: 1000,
			expectBlock:    true,
		},
		{
			name:           "User attempts to remember his password",
			ipAddress:      "192.168.0.201",
			produceAttacks: 10,
			expectBlock:    false,
		},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			// 1. Create attackerList
			attackers := createAttackerList(context.Background(), 3600, 600)
			errNumber := 0
			// 2. Implement a set of "attacks"
			for range tCase.produceAttacks {
				err := attackers.UpsertIpAddressStats(tCase.ipAddress)
				if err != nil {
					errNumber++
				}
				time.Sleep(10 * time.Millisecond)
			}
			// 3. Check final blocked status
			assert.Equal(t, 0, errNumber)
			stats := attackers.GetAttackerStats("", tCase.ipAddress)
			assert.NotNil(t, stats)
			assert.Equal(t, tCase.expectBlock, stats.blocked)
			assert.Equal(t, tCase.produceAttacks, int(stats.attackCount))
		})
	}
}

func TestUpsertDeviceStats(t *testing.T) {
	/* browser fingerprint could be obtained here:
	 * https://scrapfly.io/web-scraping-tools/browser-fingerprint
	 */
	testCases := []struct {
		name           string
		deviceId       string
		produceAttacks int
		expectBlock    bool
	}{
		{
			name:           "Sequential attack from one address and block as a result",
			deviceId:       "1b9ee5b8d043698cd13c9f11481cd037d44b8cf3",
			produceAttacks: 1000,
			expectBlock:    true,
		},
		{
			name:           "User attempts to remember his password",
			deviceId:       "d3681338021a33149a9f8ef2f48eb8bfb46b10b3",
			produceAttacks: 10,
			expectBlock:    false,
		},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			// 1. Create attackerList
			attackers := createAttackerList(context.Background(), 3600, 600)
			errNumber := 0
			// 2. Implement a set of "attacks"
			for range tCase.produceAttacks {
				err := attackers.UpsertDeviceStats(tCase.deviceId)
				if err != nil {
					errNumber++
				}
				time.Sleep(10 * time.Millisecond)
			}
			// 3. Check final blocked status
			assert.Equal(t, 0, errNumber)
			stats := attackers.GetAttackerStats(tCase.deviceId, "")
			assert.NotNil(t, stats)
			assert.Equal(t, tCase.expectBlock, stats.blocked)
			assert.Equal(t, tCase.produceAttacks, int(stats.attackCount))
		})
	}
}

func TestGetAttackerStats(t *testing.T) {
	attackers := createAttackerList(context.Background(), 3600, 600)
	err := attackers.UpsertIpAddressStats("167.134.30.55")
	assert.NoError(t, err)
	err = attackers.UpsertIpAddressStats("55.22.90.14")
	assert.NoError(t, err)
	err = attackers.UpsertIpAddressStats("102.36.99.202")
	assert.NoError(t, err)
	err = attackers.UpsertDeviceStats("1b9ee5b8d043698cd13c9f11481cd037d44b8cf3")
	assert.NoError(t, err)
	err = attackers.UpsertDeviceStats("d3681338021a33149a9f8ef2f48eb8bfb46b10b3")
	assert.NoError(t, err)
	testCases := []struct {
		name      string
		ipAddress string
		deviceId  string
		exists    bool
	}{
		{
			name:      "non-existing-ip",
			ipAddress: "127.0.0.1",
			deviceId:  "",
			exists:    false,
		},
		{
			name:      "existing-ip",
			ipAddress: "55.22.90.14",
			deviceId:  "",
			exists:    true,
		},
		{
			name:      "existing-device-id",
			ipAddress: "",
			deviceId:  "1b9ee5b8d043698cd13c9f11481cd037d44b8cf3",
			exists:    true,
		},
		{
			name:      "non-existing-device-id",
			ipAddress: "",
			deviceId:  "atatatatatatta-no-such-device-anyway",
			exists:    false,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			stats := attackers.GetAttackerStats(tCase.deviceId, tCase.ipAddress)
			if tCase.exists {
				assert.NotNil(t, stats)
			} else {
				assert.Nil(t, stats)
			}
		})
	}
}

func TestAttackersCleanup(t *testing.T) {
	watchTime := 10
	attackers := createAttackerList(context.Background(), 3600, watchTime)
	ipAddress1 := "167.134.30.55"
	ipAddress2 := "55.22.90.14"
	ipAddress3 := "102.36.99.202"
	deviceId1 := "1b9ee5b8d043698cd13c9f11481cd037d44b8cf3"
	deviceId2 := "d3681338021a33149a9f8ef2f48eb8bfb46b10b3"
	err := attackers.UpsertIpAddressStats(ipAddress1)
	assert.NoError(t, err)
	err = attackers.UpsertIpAddressStats(ipAddress2)
	assert.NoError(t, err)
	err = attackers.UpsertIpAddressStats(ipAddress3)
	assert.NoError(t, err)
	err = attackers.UpsertDeviceStats(deviceId1)
	assert.NoError(t, err)
	err = attackers.UpsertDeviceStats(deviceId2)
	assert.NoError(t, err)
	ipAddress4 := "127.0.0.1"
	for range blockThreshold {
		err = attackers.UpsertIpAddressStats(ipAddress4)
		assert.NoError(t, err)
	}
	time.Sleep(time.Duration(watchTime+2) * time.Second)
	stats := attackers.GetAttackerStats("", ipAddress1)
	assert.Nil(t, stats)
}
