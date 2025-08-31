// package main

// import (
// 	"flag"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"runtime"
// 	"sync"
// 	"time"
// )

// type Result struct {
// 	min, max, total time.Duration
// 	count           int
// 	rps             float64
// 	avg             time.Duration
// }

// func sessionWorker(urls []string, sessions int, wg *sync.WaitGroup, resChan chan<- time.Duration) {
// 	defer wg.Done()
// 	for i := 0; i < sessions; i++ {
// 		var total time.Duration
// 		for _, url := range urls {
// 			start := time.Now()
// 			resp, err := http.Get(url)
// 			if err != nil {
// 				continue
// 			}
// 			io.Copy(io.Discard, resp.Body)
// 			resp.Body.Close()
// 			total += time.Since(start)
// 		}
// 		resChan <- total
// 	}
// }

// func benchmark(urls []string, concurrency, sessions int) Result {
// 	var wg sync.WaitGroup
// 	resChan := make(chan time.Duration, sessions)

// 	startAll := time.Now()
// 	for i := 0; i < concurrency; i++ {
// 		wg.Add(1)
// 		go sessionWorker(urls, sessions/concurrency, &wg, resChan)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(resChan)
// 	}()

// 	r := Result{min: time.Hour, max: 0, total: 0, count: 0}
// 	for t := range resChan {
// 		if t < r.min {
// 			r.min = t
// 		}
// 		if t > r.max {
// 			r.max = t
// 		}
// 		r.total += t
// 		r.count++
// 	}
// 	elapsed := time.Since(startAll)

// 	if r.count > 0 {
// 		r.avg = r.total / time.Duration(r.count)
// 		r.rps = float64(r.count) / elapsed.Seconds()
// 	}

// 	return r
// }

// func printResult(urls []string, c, sessions int, res Result, warn bool) {
// 	if warn {
// 		fmt.Print("\033[31m") // red
// 	}
// 	fmt.Printf("Concurrency: %d | Sessions: %d | RPS: %.2f | Lat min/avg/max: %v / %v / %v\n",
// 		c, res.count, res.rps, res.min, res.avg, res.max)
// 	if warn {
// 		fmt.Print("\033[0m") // reset
// 	}
// }

// func rampUp(urls []string, start, step, sessions int) {
// 	prevRPS := 0.0
// 	for c := start; ; c += step {
// 		res := benchmark(urls, c, sessions)
// 		// detect instability
// 		warn := false
// 		if res.min > 0 && res.avg > 5*res.min {
// 			warn = true
// 		}
// 		if prevRPS > 0 && res.rps < prevRPS*0.8 {
// 			warn = true
// 		}
// 		printResult(urls, c, sessions, res, warn)
// 		prevRPS = res.rps
// 		time.Sleep(2 * time.Second) // cooldown
// 	}
// }

// func main() {
// 	runtime.GOMAXPROCS(runtime.NumCPU())

// 	staggered := flag.Bool("staggered-mode", false, "run /item then /api/item per session")
// 	ramp := flag.Bool("ramp", false, "run ramp-up test until stopped (Ctrl+C)")
// 	sessions := flag.Int("sessions", 10000, "total number of sessions per test")
// 	concurrency := flag.Int("concurrency", 1000, "number of concurrent workers (fixed mode)")
// 	flag.Parse()

// 	itemURL := "http://localhost:8888/item?id=44"
// 	apiURL := "http://localhost:8888/api/item?id=44"

// 	urls := []string{itemURL}
// 	if *staggered {
// 		urls = []string{itemURL, apiURL}
// 	}

// 	if *ramp {
// 		rampUp(urls, 10, 100, *sessions) // start=10 workers, +100 each loop
// 	} else {
// 		res := benchmark(urls, *concurrency, *sessions)
// 		printResult(urls, *concurrency, *sessions, res, false)
// 	}
// }

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Result struct {
	min, max, total time.Duration
	count           int
	rps             float64
	avg             time.Duration
}

// sessionWorker simulates user sessions
// In staggered mode: fetches /item then /api/item (simulating page load + API call)
// In normal mode: only fetches /item (simulating just page load)
func sessionWorker(urls []string, sessions int, wg *sync.WaitGroup, resChan chan<- time.Duration) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for i := 0; i < sessions; i++ {
		sessionStart := time.Now()

		// Execute all URLs in sequence for this session
		for _, url := range urls {
			resp, err := client.Get(url)
			if err != nil {
				// Skip this session on error
				break
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		sessionDuration := time.Since(sessionStart)
		resChan <- sessionDuration
	}
}

func benchmark(urls []string, concurrency, sessions int) Result {
	var wg sync.WaitGroup
	resChan := make(chan time.Duration, sessions)

	startAll := time.Now()

	// Distribute sessions across workers
	sessionsPerWorker := sessions / concurrency
	remainingSessions := sessions % concurrency

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		workerSessions := sessionsPerWorker
		if i < remainingSessions {
			workerSessions++ // distribute remainder
		}
		go sessionWorker(urls, workerSessions, &wg, resChan)
	}

	// Close channel when all workers complete
	go func() {
		wg.Wait()
		close(resChan)
	}()

	// Collect results
	r := Result{min: time.Hour, max: 0, total: 0, count: 0}
	for sessionTime := range resChan {
		if sessionTime < r.min {
			r.min = sessionTime
		}
		if sessionTime > r.max {
			r.max = sessionTime
		}
		r.total += sessionTime
		r.count++
	}

	elapsed := time.Since(startAll)
	if r.count > 0 {
		r.avg = r.total / time.Duration(r.count)
		r.rps = float64(r.count) / elapsed.Seconds()
	}

	return r
}

func printResult(urls []string, concurrency, sessions int, res Result, warn bool) {
	if warn {
		fmt.Print("\033[31m") // red text for warnings
	}

	fmt.Printf("Mode: %s\n", getModeDescription(urls))
	fmt.Printf("Concurrency: %d | Sessions: %d | Sessions/s: %.2f\n",
		concurrency, res.count, res.rps)
	fmt.Printf("Session latency min/avg/max: %v / %v / %v\n",
		res.min, res.avg, res.max)

	if warn {
		fmt.Print("\033[0m") // reset color
	}
	fmt.Println()
}

func getModeDescription(urls []string) string {
	if len(urls) == 1 {
		return "Normal (page load only)"
	}
	return "Staggered (page load + API call)"
}

func rampUp(urls []string, start, step, sessions int) {
	fmt.Printf("Starting ramp-up test: %s\n", getModeDescription(urls))
	fmt.Printf("Starting concurrency: %d, Step: %d, Sessions per test: %d\n\n", start, step, sessions)

	prevRPS := 0.0
	for c := start; ; c += step {
		fmt.Printf("Testing concurrency %d...\n", c)
		res := benchmark(urls, c, sessions)

		// Detect performance degradation
		warn := false
		if res.min > 0 && res.avg > 5*res.min {
			warn = true // High latency variance
		}
		if prevRPS > 0 && res.rps < prevRPS*0.8 {
			warn = true // Significant RPS drop
		}

		printResult(urls, c, sessions, res, warn)
		prevRPS = res.rps

		time.Sleep(2 * time.Second) // Cooldown between tests
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Command line flags
	staggered := flag.Bool("staggered-mode", false,
		"simulate user workflow: fetch /item then /api/item per session")
	ramp := flag.Bool("ramp", false,
		"run ramp-up test increasing concurrency until stopped (Ctrl+C)")
	sessions := flag.Int("sessions", 10000,
		"total number of sessions per test")
	concurrency := flag.Int("concurrency", 1000,
		"number of concurrent workers (fixed mode only)")
	flag.Parse()

	// Define URLs
	itemURL := "http://localhost:8888/item?id=44"
	apiURL := "http://localhost:8888/api/item?id=44"

	// Configure session workflow
	urls := []string{itemURL}
	if *staggered {
		urls = []string{itemURL, apiURL} // Sequential: page load, then API call
	}

	fmt.Printf("HTTP Load Testing Tool\n")
	fmt.Printf("======================\n")

	if *ramp {
		rampUp(urls, 10, 100, *sessions) // Start at 10, increment by 100
	} else {
		fmt.Printf("Running single test: %s\n\n", getModeDescription(urls))
		res := benchmark(urls, *concurrency, *sessions)
		printResult(urls, *concurrency, *sessions, res, false)
	}
}
