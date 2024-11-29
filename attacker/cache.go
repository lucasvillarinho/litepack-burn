package attacker

import (
	"log"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type cacheAttacker struct {
	vegeta *vegeta.Attacker
}

type CacheAttacker interface {
	AttackCacheSet() error
}

func NewCacheAttacker() CacheAttacker {
	return &cacheAttacker{
		vegeta: vegeta.NewAttacker(),
	}
}

func (ca *cacheAttacker) AttackCacheSet() error {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 10 * time.Second

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8080/cache/set",
		Body:   []byte(`{"key":"test-key","value":"test-value","ttl":5}`),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	})

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Set Benchmark") {
		metrics.Add(res)
	}
	metrics.Close()

	log.Printf("Requests: %d", metrics.Requests)
	log.Printf("Throughput: %.7f req/s", metrics.Throughput)
	log.Printf("Latency: min=%s, avg=%s, max=%s", metrics.Latencies, metrics.Latencies.Mean, metrics.Latencies.Max)
	log.Printf("Errors: %d", len(metrics.Errors))

	return nil
}
