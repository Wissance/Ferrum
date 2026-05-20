package bruteforce

import "testing"

func TestUpsertIpAddressStats(t *testing.T) {

}

func TestUpsertDeviceStats(t *testing.T) {

}

func TestGetAttackerStats(t *testing.T) {
	
}

/*import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAttackerEquals(t *testing.T) {
	device1Id := "8e59868c-55b4-42e3-88e2-a709ba403ecc"
	device2Id := "8ad11bba-a072-48d7-86d6-63a33cc00e1c"

	testCases := []struct {
		name         string
		items        []attacker
		itemToLookup attacker
		lookupResult bool
	}{
		{
			name: "Search attacker by IP among IP and Devices, present in table",
			items: []attacker{
				{IPAddress: "192.168.122.20", DeviceId: nil},
				{IPAddress: "192.168.123.44", DeviceId: nil},
				{IPAddress: "167.165.190.239", DeviceId: nil},
				{IPAddress: "167.165.190.239", DeviceId: &device1Id},
				{IPAddress: "10.116.200.15", DeviceId: &device2Id},
			},
			itemToLookup: attacker{IPAddress: "192.168.123.44", DeviceId: nil},
			lookupResult: true,
		},
		{
			name: "Search attacker by Device among IP and Devices, present in table",
			items: []attacker{
				{IPAddress: "192.168.122.20", DeviceId: nil},
				{IPAddress: "192.168.123.44", DeviceId: nil},
				{IPAddress: "167.165.190.239", DeviceId: nil},
				{IPAddress: "167.165.190.239", DeviceId: &device1Id},
				{IPAddress: "10.116.200.15", DeviceId: &device2Id},
			},
			itemToLookup: attacker{IPAddress: "167.165.190.178", DeviceId: &device1Id},
			lookupResult: true,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			attackersMap := map[attacker]struct{}{}
			for _, hacker := range tCase.items {
				attackersMap[hacker] = struct{}{}
			}
			_, ok := attackersMap[tCase.itemToLookup]
			assert.Equal(t, tCase.lookupResult, ok)
		})
	}
}
*/
