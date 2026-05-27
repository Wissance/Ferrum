package bruteforce

import (
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
	"golang.org/x/net/context"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// run this test : go test ./security/bruteforce -v -bench=^Benchmark -run=^$ -benchmem -cpu 2
func BenchmarkRegisterAndCheckIsBlocked(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const clientsNumber = 100000
		const validClients = 10000
		const potentialAttackers = clientsNumber - validClients
		/* among the client some of them are attackers, at least 10K are valid users, others could be
		 * an attackers
		 */
		realAttackers := rand.Intn(potentialAttackers)
		realAttackerCount := 0
		personForgetPasswordCount := 0
		cfg := ProtectionServiceConfig{
			WatchTimeSec: 300,
		}
		logger := logging.AppLogger{}
		protectionService := CreateProtectionService(context.Background(), &cfg, &logger)

		wg := sync.WaitGroup{}
		wg.Add(potentialAttackers)

		for range potentialAttackers {
			// 1. Get randomly real attacker or person who forget the password
			isRealAttacker := false
			attemptsNumber := 0
			if realAttackerCount < realAttackers {
				isRealAttacker = rand.Intn(2) > 0
			}
			if isRealAttacker {
				realAttackerCount++
				attemptsNumber = 101 + rand.Intn(50)
			} else {
				personForgetPasswordCount++
				attemptsNumber = 10 + rand.Intn(10)
			}
			// 2. Perform actions
			go func() {
				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
				ipAddr := getRandomIp()
				for range attemptsNumber {
					protectionService.RegisterIpAddressAttempt(ipAddr)
				}
				wg.Done()
			}()
		}

		wg.Wait()
	}
}

func getRandomIp() string {
	d1 := rand.Intn(250)
	if d1 < 100 {
		d1 += 100 - d1
	}
	d2 := rand.Intn(250)
	if d2 < 100 {
		d2 += 100 - d2
	}
	d3 := rand.Intn(250)
	if d3 < 100 {
		d3 += 100 - d3
	}
	d4 := rand.Intn(250)
	if d4 < 100 {
		d4 += 100 - d4
	}
	return sf.Format("{0}.{1}.{2}.{3}", d1, d2, d3, d4)
}
