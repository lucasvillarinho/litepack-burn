package attacker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/go-faker/faker/v4"
	vegeta "github.com/tsenart/vegeta/lib"

	"github.com/lucasvillarinho/litepack-burn/table"
)

type cacheAttacker struct {
	vegeta *vegeta.Attacker
}

type CacheAttacker interface {
	Attack() error
}

var (
	rate     = vegeta.Rate{Freq: 100, Per: time.Second}
	duration = 10 * time.Second
)

type FakeCacheData struct {
	ID               string `json:"id" faker:"uuid_hyphenated"`
	Name             string `json:"name" faker:"name"`
	Email            string `json:"email" faker:"email"`
	CreditCardNumber string `json:"cc_number" faker:"cc_number"`
	PaymentMethod    string `json:"payment_method" faker:"oneof: cc, paypal, check, money order"`
	Age              int    `json:"age" faker:"boundary_start=18, boundary_end=99"`
}

func NewCacheAttacker() CacheAttacker {
	return &cacheAttacker{
		vegeta: vegeta.NewAttacker(),
	}
}

func (ca *cacheAttacker) Attack() error {

	data, err := GenerateCacheFakeData(1000)
	if err != nil {
		return fmt.Errorf("❌ Error generating fake data: %v", err)
	}
	fmt.Printf("✅ Generated %d fake data\n", len(data))

	setMetrics, err := ca.AttackCacheSet(data)
	if err != nil {
		return fmt.Errorf("❌ Error running set benchmark: %v", err)
	}
	fmt.Println("✅ Set benchmark finished")

	upsertMetrics, err := ca.AttackCacheSet(data)
	if err != nil {
		return fmt.Errorf("❌ Error running upsert benchmark: %v", err)
	}
	fmt.Println("✅ Upsert benchmark finished")

	getMetrics, err := ca.AttackCacheGet(data)
	if err != nil {
		return fmt.Errorf("❌ Error running get benchmark: %v", err)
	}
	fmt.Println("✅ Get benchmark finished")

	deleteMetrics, err := ca.AttackCacheDelete(data)
	if err != nil {
		return fmt.Errorf("❌ Error running delete benchmark: %v", err)
	}

	fmt.Println("🪖 Results:")
	var rows [][]string
	rows = append(rows, rowCacheMetrics("SET", setMetrics)...)
	rows = append(rows, rowCacheMetrics("UPSERT", upsertMetrics)...)
	rows = append(rows, rowCacheMetrics("GET", getMetrics)...)
	rows = append(rows, rowCacheMetrics("DELETE", deleteMetrics)...)

	renderCacheMetrics(rows)
	return nil
}

func (ca *cacheAttacker) AttackCacheSet(data map[string]string) (vegeta.Metrics, error) {
	targets := setupSetTarget(data)
	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Set Benchmark") {
		metrics.Add(res)
	}
	metrics.Close()

	return metrics, nil
}

func setupSetTarget(data map[string]string) []vegeta.Target {
	targets := make([]vegeta.Target, 0, len(data))
	for key, value := range data {
		body := fmt.Sprintf(`{"key":"%s","value":"%s","ttl":5}`, key, value)
		targets = append(targets, vegeta.Target{
			Method: "POST",
			URL:    "http://localhost:8080/cache/set",
			Body:   []byte(body),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		})
	}

	return targets
}

func (ca *cacheAttacker) AttackCacheGet(data map[string]string) (vegeta.Metrics, error) {
	targets := setupGetTarget(data)
	targeter := vegeta.NewStaticTargeter(targets...)

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Get Benchmark") {
		metrics.Add(res)
	}
	metrics.Close()

	return metrics, nil
}

func setupGetTarget(data map[string]string) []vegeta.Target {
	targets := make([]vegeta.Target, 0, len(data))
	for key := range data {
		targets = append(targets, vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("http://localhost:8080/cache/get/%s", key),
		})
	}

	return targets
}

func (ca *cacheAttacker) AttackCacheDelete(data map[string]string) (vegeta.Metrics, error) {
	targets := setupDeleteTarget(data)

	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Delete Benchmark") {
		metrics.Add(res)
	}
	metrics.Close()

	return metrics, nil
}

func setupDeleteTarget(data map[string]string) []vegeta.Target {
	targets := make([]vegeta.Target, 0, len(data))
	for key := range data {
		targets = append(targets, vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("http://localhost:8080/cache/delete/%s", key),
		})
	}

	return targets
}

func createCacheTable(rows [][]string) *table.Table {
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
		Rows(rows...)

	return table

}

func renderCacheMetrics(rows [][]string) {
	table := createCacheTable(rows)

	fmt.Println(table.Render())
}

func headerCacheMetrics() []string {
	return []string{
		"Method",
		"Requests",
		"Success",
		"Mean Latency",
		"Max Latency",
		"P50 Latency",
		"P99 Latency",
		"Errors"}
}

func rowCacheMetrics(method string, metrics vegeta.Metrics) [][]string {
	return [][]string{
		{
			method,
			fmt.Sprintf("%d", metrics.Requests),
			fmt.Sprintf("%.2f%%", metrics.Success*100),
			metrics.Latencies.Mean.String(),
			metrics.Latencies.Max.String(),
			metrics.Latencies.P50.String(),
			metrics.Latencies.P99.String(),
			fmt.Sprintf("%v", metrics.Errors),
		},
	}
}

func GenerateCacheFakeData(size int) (map[string]string, error) {
	data := make(map[string]string)

	for i := 0; i < size; i++ {
		var fake FakeCacheData
		if err := faker.FakeData(&fake); err != nil {
			return nil, fmt.Errorf("error generating fake data: %w", err)
		}

		jsonData, err := json.Marshal(fake)
		if err != nil {
			return nil, fmt.Errorf("error marshalling JSON: %w", err)
		}

		data[fake.ID] = escapeJSON(string(jsonData))

	}

	return data, nil
}

func escapeJSON(input string) string {
	escaped, _ := json.Marshal(input)
	return string(escaped[1 : len(escaped)-1])
}
