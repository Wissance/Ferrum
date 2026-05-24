package bruteforce

import "context"

// ListBasedProtectionService is bruteforce protection service that uses lists of blocked clients/device
type ListBasedProtectionService struct {
	attackersList *attackerList
	ctx           context.Context
}

func CreateListBasedProtectionService(ctx context.Context, watchTime int) *ListBasedProtectionService {
	// block time actually not using, we are using permanent block until app is restarted
	attackersList := createAttackerList(ctx, 86400, watchTime)
	return &ListBasedProtectionService{
		attackersList: attackersList,
		ctx:           ctx,
	}
}

func
