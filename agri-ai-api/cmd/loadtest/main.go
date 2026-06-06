package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Defina a URL base. Para o teste usaremos o /ping que é aberto
	// Ou o login (se tiver as credenciais mockadas)
	url := "http://localhost:8080/api/v1/ping"
	numRequests := 1000 // Total de requisições
	concurrency := 50   // Requisições simultâneas

	fmt.Printf("Iniciando Load Test em %s\n", url)
	fmt.Printf("Total de Requisições: %d | Concorrência: %d\n", numRequests, concurrency)

	var wg sync.WaitGroup
	requestsChan := make(chan int, numRequests)
	
	start := time.Now()

	// Inicia os workers
	var successCount, errCount int
	var mu sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestsChan {
				resp, err := http.Get(url)
				mu.Lock()
				if err != nil {
					errCount++
				} else {
					if resp.StatusCode == 200 {
						successCount++
					} else {
						errCount++
					}
					resp.Body.Close()
				}
				mu.Unlock()
			}
		}()
	}

	// Envia as requisições
	for i := 0; i < numRequests; i++ {
		requestsChan <- i
	}
	close(requestsChan)

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("-----------------------------------------------------")
	fmt.Println("Resultados do Teste:")
	fmt.Printf("Tempo Total: %v\n", duration)
	fmt.Printf("Sucessos: %d\n", successCount)
	fmt.Printf("Erros: %d\n", errCount)
	fmt.Printf("Requisições por segundo: %.2f\n", float64(numRequests)/duration.Seconds())
	fmt.Println("-----------------------------------------------------")
}
