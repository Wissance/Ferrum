package bruteforce

type attacker struct {
	IPAddress string
	DeviceId  *string
}

func (a attacker) Equals(other attacker) bool {
	// attacker use browser with specific deviceId
	if a.DeviceId != nil && other.DeviceId != nil {
		return *a.DeviceId == *other.DeviceId
	}
	// attacker did not used browser, but IP could be equal
	return a.IPAddress == other.IPAddress
}
