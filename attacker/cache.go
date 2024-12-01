package attacker

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	vegeta "github.com/tsenart/vegeta/lib"

	"github.com/lucasvillarinho/litepack-burn/helpers"
	"github.com/lucasvillarinho/litepack-burn/table"
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
	// rate := vegeta.Rate{Freq: 100, Per: time.Second}
	// duration := 10 * time.Second

	// targeter := vegeta.NewStaticTargeter(vegeta.Target{
	// 	Method: "POST",
	// 	URL:    "http://localhost:8080/cache/set",
	// 	Body:   []byte(`{"key":"test-key","value":"test-value","ttl":5}`),
	// 	Header: map[string][]string{
	// 		"Content-Type": {"application/json"},
	// 	},
	// })

	// attacker := vegeta.NewAttacker()
	// var metrics vegeta.Metrics
	// for res := range attacker.Attack(targeter, rate, duration, "Set Benchmark") {
	// 	metrics.Add(res)
	// }
	// metrics.Close()

	metrics := helpers.GetFakerDataVegetaMetrics()

	renderCacheMetrics("SET", metrics)

	return nil
}

func createCacheTable(rows []string) *table.Table {
	HeaderStyle := lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center)

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF69B4")).
		Bold(true)

	table := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(borderStyle).
		StyleFunc(table.StyleFunc(func(row, col int) lipgloss.Style {
			return HeaderStyle
		})).
		Headers(headerCacheMetrics()...).
		Rows(rows)

	return table

}

func renderCacheMetrics(method string, metrics vegeta.Metrics) {
	table := createCacheTable(rowCacheMetrics(method, metrics))

	fmt.Println(table.Render())
}

func headerCacheMetrics() []string {
	return []string{
		"Method",
		"Requests",
		"Throughput",
		"Success",
		"Mean Latency",
		"P50 Latency",
		"P99 Latency",
		"Max Latency",
		"Errors"}
}

func rowCacheMetrics(method string, metrics vegeta.Metrics) []string {
	return []string{
		method,
		fmt.Sprintf("%d", metrics.Requests),
		fmt.Sprintf("%.4f", metrics.Throughput),
		fmt.Sprintf("%.4f", metrics.Success),
		metrics.Latencies.Mean.String(),
		metrics.Latencies.P50.String(),
		metrics.Latencies.P99.String(),
		metrics.Latencies.Max.String(),
		fmt.Sprintf("%v", metrics.Errors),
	}
}
