package services

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func init() {
	st := gobreaker.Settings{
		Name:        "OpenMeteo HTTP",
		MaxRequests: 3, // Máximo de requisições permitidas quando em Half-Open
		Interval:    10 * time.Second, // Tempo de expiração do estado Closed para resetar falhas
		Timeout:     30 * time.Second, // Tempo até tentar voltar de Open para Half-Open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Abre o circuito se houver mais de 5 falhas consecutivas
			return counts.ConsecutiveFailures >= 5
		},
	}
	cb = gobreaker.NewCircuitBreaker(st)
}

// FetchOpenMeteo realiza a requisição à API usando Retry e Circuit Breaker
func FetchOpenMeteo(lat, lon float64) ([]byte, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,precipitation,wind_speed_10m", lat, lon)

	// A execução do cb.Execute envolve toda a lógica de chamada e retry
	result, err := cb.Execute(func() (interface{}, error) {
		return doRequestWithRetry(url, 3)
	})

	if err != nil {
		return nil, fmt.Errorf("openmeteo request failed: %w", err)
	}

	return result.([]byte), nil
}

// doRequestWithRetry implementa um retry simples com exponential backoff
func doRequestWithRetry(url string, maxRetries int) ([]byte, error) {
	var lastErr error
	backoff := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		// #nosec G107
		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
		} else {
			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				return body, err
			}
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			_ = resp.Body.Close()
			// Não tentar novamente se for erro 4xx (erro do cliente, como bad request)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return nil, lastErr
			}
		}

		// Espera antes do próximo retry (exponential backoff)
		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, fmt.Errorf("max retries exceeded, last error: %w", lastErr)
}
