package bruteforce

// ListBasedProtectionService is bruteforce protection service that uses lists of blocked clients/device
type ListBasedProtectionService struct {
	attackersList attackerList
}
