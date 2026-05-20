package bruteforce

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestUpsertIpAddressStats(t *testing.T) {
	testCases := []struct {
		name           string
		ipList         []string
		produceAttacks int
		expectBlock    bool
	}{
		{
			name:           "Sequential attack from one address and block as a result",
			ipList:         []string{"192.167.99.144"},
			produceAttacks: 1000,
			expectBlock:    true,
		},
		{
			name:           "User attempts to remember his password",
			ipList:         []string{"192.168.0.201"},
			produceAttacks: 10,
			expectBlock:    false,
		},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			// 1. Create attackerList
			attackers := createAttackerList(3600)
			// 2. take random IP from the list
			ipIndex := rand.Intn(len(tCase.ipList))
			selectedIP := tCase.ipList[ipIndex]
			errNumber := 0
			for range tCase.produceAttacks {
				err := attackers.UpsertIpAddressStats(selectedIP)
				if err != nil {
					errNumber++
				}
				time.Sleep(10 * time.Millisecond)
			}
			assert.Equal(t, 0, errNumber)
			stats := attackers.GetAttackerStats("", selectedIP)
			assert.NotNil(t, stats)
			assert.Equal(t, tCase.expectBlock, stats.blocked)
		})
	}
}

func TestUpsertDeviceStats(t *testing.T) {

}

func TestGetAttackerStats(t *testing.T) {

}
