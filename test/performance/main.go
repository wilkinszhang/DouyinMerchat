package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	TotalRequests   int
	SuccessRequests int
	FailedRequests  int
	AverageLatency  time.Duration
	MinLatency      time.Duration
	MaxLatency      time.Duration
}

func TestEndpoint(url string, payload interface{}, concurrency, requests int) Metrics {
	var wg sync.WaitGroup
	var metrics Metrics
	var mu sync.Mutex

	metrics.MinLatency = time.Hour
	requestsPerWorker := requests / concurrency

	jsonPayload, _ := json.Marshal(payload)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < requestsPerWorker; j++ {
				start := time.Now()

				resp, err := http.Post(url,
					"application/json",
					bytes.NewBuffer(jsonPayload))

				latency := time.Since(start)

				mu.Lock()
				metrics.TotalRequests++
				if err == nil && resp.StatusCode == 200 {
					metrics.SuccessRequests++
					metrics.AverageLatency += latency
					if latency < metrics.MinLatency {
						metrics.MinLatency = latency
					}
					if latency > metrics.MaxLatency {
						metrics.MaxLatency = latency
					}
				} else {
					metrics.FailedRequests++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}
		}()
	}

	wg.Wait()

	if metrics.SuccessRequests > 0 {
		metrics.AverageLatency /= time.Duration(metrics.SuccessRequests)
	}

	return metrics
}

func main() {
	baseURL := "http://localhost:8888"
	concurrency := 100
	totalRequests := 10000

	// Test Register
	registerMetrics := TestEndpoint(
		baseURL+"/api/user/register",
		map[string]string{
			"email":            "test@example.com",
			"password":         "password123",
			"confirm_password": "password123",
		},
		concurrency,
		totalRequests,
	)

	fmt.Printf("\nRegister Endpoint Performance:\n")
	fmt.Printf("Total Requests: %d\n", registerMetrics.TotalRequests)
	fmt.Printf("Success Rate: %.2f%%\n", float64(registerMetrics.SuccessRequests)/float64(registerMetrics.TotalRequests)*100)
	fmt.Printf("Average Latency: %v\n", registerMetrics.AverageLatency)
	fmt.Printf("Min Latency: %v\n", registerMetrics.MinLatency)
	fmt.Printf("Max Latency: %v\n", registerMetrics.MaxLatency)

	// Test Login
	loginMetrics := TestEndpoint(
		baseURL+"/api/user/login",
		map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		},
		concurrency,
		totalRequests,
	)

	fmt.Printf("\nLogin Endpoint Performance:\n")
	fmt.Printf("Total Requests: %d\n", loginMetrics.TotalRequests)
	fmt.Printf("Success Rate: %.2f%%\n", float64(loginMetrics.SuccessRequests)/float64(loginMetrics.TotalRequests)*100)
	fmt.Printf("Average Latency: %v\n", loginMetrics.AverageLatency)
	fmt.Printf("Min Latency: %v\n", loginMetrics.MinLatency)
	fmt.Printf("Max Latency: %v\n", loginMetrics.MaxLatency)
}
