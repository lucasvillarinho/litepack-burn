package helpers

import (
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func GetFakerDataVegetaMetrics() vegeta.Metrics {
	return vegeta.Metrics{
		Requests:   1000,
		Throughput: 95.0,
		Success:    0.99,
		Latencies: vegeta.LatencyMetrics{
			Mean: time.Millisecond * 200,
			P50:  time.Millisecond * 190,
			P99:  time.Millisecond * 250,
			Max:  time.Millisecond * 300,
		},
		BytesIn: vegeta.ByteMetrics{
			Total: 1024000,
			Mean:  1024,
		},
		BytesOut: vegeta.ByteMetrics{
			Total: 512000,
			Mean:  512,
		},
		Errors: []string{}, // No errors for default
	}
}
