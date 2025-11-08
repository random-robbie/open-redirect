package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	banner = `
[*] ***************************************[*]
[*] Open Redirect Finder By @Random_Robbie [*]
[*]         Rewritten in Go + ChromeDP      [*]
[*] ***************************************[*]
`
)

var (
	urlsFile     string
	payloadsFile string
	outputFile   string
	workers      int
	timeout      int
	verbose      bool
	testDomains  = []string{
		"http://google.com",
		"https://google.com",
		"http://example.com",
		"https://example.com",
	}
)

type TestCase struct {
	URL     string
	Payload string
}

type Result struct {
	TestURL     string
	FinalURL    string
	IsVulnerable bool
}

func init() {
	flag.StringVar(&urlsFile, "urls", "urls.txt", "File containing URLs to test")
	flag.StringVar(&payloadsFile, "payloads", "payloads.txt", "File containing payloads")
	flag.StringVar(&outputFile, "output", "found.txt", "Output file for vulnerable URLs")
	flag.IntVar(&workers, "workers", 5, "Number of concurrent workers")
	flag.IntVar(&timeout, "timeout", 30, "Timeout in seconds for each request")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
}

func main() {
	flag.Parse()

	fmt.Println(banner)

	// Read URLs from file
	urls, err := readLines(urlsFile)
	if err != nil {
		log.Fatalf("[!] Error reading URLs file: %v", err)
	}

	// Read payloads from file
	payloads, err := readLines(payloadsFile)
	if err != nil {
		log.Fatalf("[!] Error reading payloads file: %v", err)
	}

	if len(urls) == 0 {
		log.Fatal("[!] No URLs to test. Please provide a URLs file.")
	}

	if len(payloads) == 0 {
		log.Fatal("[!] No payloads to test. Please provide a payloads file.")
	}

	fmt.Printf("[*] Loaded %d URLs and %d payloads\n", len(urls), len(payloads))
	fmt.Printf("[*] Using %d concurrent workers\n", workers)
	fmt.Println("[*] Starting scan...\n")

	// Create test cases
	testCases := make(chan TestCase, len(urls)*len(payloads))
	for _, url := range urls {
		for _, payload := range payloads {
			testCases <- TestCase{URL: url, Payload: payload}
		}
	}
	close(testCases)

	// Process test cases with workers
	results := make(chan Result, len(urls)*len(payloads))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i+1, testCases, results, &wg)
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	vulnerableCount := 0
	for result := range results {
		if result.IsVulnerable {
			vulnerableCount++
			logVulnerability(result)
		} else if verbose {
			fmt.Printf("[*] Not vulnerable: %s\n", result.TestURL)
		}
	}

	fmt.Printf("\n[*] Scan complete!\n")
	fmt.Printf("[*] Found %d vulnerable URLs\n", vulnerableCount)
	if vulnerableCount > 0 {
		fmt.Printf("[*] Results saved to: %s\n", outputFile)
	}
}

func worker(id int, testCases <-chan TestCase, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for testCase := range testCases {
		result := testRedirect(testCase)
		results <- result
	}
}

func testRedirect(tc TestCase) Result {
	testURL := tc.URL + tc.Payload

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	var finalURL string
	err := chromedp.Run(ctx,
		chromedp.Navigate(testURL),
		chromedp.Sleep(2*time.Second), // Wait for redirects
		chromedp.Location(&finalURL),
	)

	if err != nil {
		if verbose {
			fmt.Printf("[!] Error testing %s: %v\n", testURL, err)
		}
		return Result{
			TestURL:      testURL,
			FinalURL:     "",
			IsVulnerable: false,
		}
	}

	// Check if redirected to a test domain
	isVulnerable := false
	for _, domain := range testDomains {
		if strings.HasPrefix(finalURL, domain) {
			isVulnerable = true
			break
		}
	}

	return Result{
		TestURL:      testURL,
		FinalURL:     finalURL,
		IsVulnerable: isVulnerable,
	}
}

func logVulnerability(result Result) {
	// Print to console
	fmt.Println("\n\n[*]*****Open Redirect Found*****[*]")
	fmt.Printf("[*] %s [*]\n", result.TestURL)
	fmt.Printf("[*] Redirects to: %s [*]\n", result.FinalURL)
	fmt.Println()

	// Write to file
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[!] Error writing to output file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(result.TestURL + "\n"); err != nil {
		log.Printf("[!] Error writing to output file: %v", err)
	}
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}
