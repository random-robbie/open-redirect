package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/schollz/progressbar/v3"
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
	urlsFile      string
	payloadsFile  string
	outputFile    string
	workers       int
	timeout       int
	verbose       bool
	jsonOutput    bool
	customDomains string
	cookiesFlag   string
	headersFlag   string
	proxyURL      string
	testDomains   []string
	cookies       []*network.CookieParam
	headers       map[string]string
)

type TestCase struct {
	URL     string
	Payload string
}

type Result struct {
	TestURL      string `json:"test_url"`
	FinalURL     string `json:"final_url"`
	IsVulnerable bool   `json:"vulnerable"`
	Timestamp    string `json:"timestamp"`
}

type JSONOutput struct {
	ScanInfo struct {
		StartTime       string `json:"start_time"`
		EndTime         string `json:"end_time"`
		TotalTests      int    `json:"total_tests"`
		VulnerableCount int    `json:"vulnerable_count"`
	} `json:"scan_info"`
	Results []Result `json:"results"`
}

func init() {
	flag.StringVar(&urlsFile, "urls", "urls.txt", "File containing URLs to test")
	flag.StringVar(&payloadsFile, "payloads", "payloads.txt", "File containing payloads")
	flag.StringVar(&outputFile, "output", "found.txt", "Output file for vulnerable URLs")
	flag.IntVar(&workers, "workers", 5, "Number of concurrent workers")
	flag.IntVar(&timeout, "timeout", 30, "Timeout in seconds for each request")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	flag.StringVar(&customDomains, "domains", "", "Comma-separated list of custom test domains (e.g., 'https://evil.com,http://test.com')")
	flag.StringVar(&cookiesFlag, "cookies", "", "Cookies in format 'name1=value1; name2=value2'")
	flag.StringVar(&headersFlag, "headers", "", "Custom headers in format 'Header1: Value1; Header2: Value2'")
	flag.StringVar(&proxyURL, "proxy", "", "Proxy URL (e.g., 'http://proxy.example.com:8080')")
}

func main() {
	flag.Parse()

	fmt.Println(banner)

	// Parse custom test domains if provided
	if customDomains != "" {
		testDomains = strings.Split(customDomains, ",")
		for i := range testDomains {
			testDomains[i] = strings.TrimSpace(testDomains[i])
		}
	} else {
		// Default test domains
		testDomains = []string{
			"http://google.com",
			"https://google.com",
			"http://example.com",
			"https://example.com",
		}
	}

	// Parse cookies if provided
	if cookiesFlag != "" {
		cookies = parseCookies(cookiesFlag)
	}

	// Parse headers if provided
	if headersFlag != "" {
		headers = parseHeaders(headersFlag)
	}

	// Display configuration
	fmt.Printf("[*] Test domains: %v\n", testDomains)
	if len(cookies) > 0 {
		fmt.Printf("[*] Using %d cookie(s)\n", len(cookies))
	}
	if len(headers) > 0 {
		fmt.Printf("[*] Using %d custom header(s)\n", len(headers))
	}
	if proxyURL != "" {
		fmt.Printf("[*] Using proxy: %s\n", proxyURL)
	}

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

	startTime := time.Now()
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

	// Create progress bar
	totalTests := len(urls) * len(payloads)
	bar := progressbar.NewOptions(totalTests,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("[cyan][*][reset] Testing URLs..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	// Collect results
	var allResults []Result
	vulnerableCount := 0
	for result := range results {
		bar.Add(1)

		if result.IsVulnerable {
			vulnerableCount++
			allResults = append(allResults, result)
			if !jsonOutput {
				// Clear progress bar line, print result, redraw progress bar
				fmt.Print("\r\033[K")
				logVulnerability(result)
			}
		} else if verbose {
			if !jsonOutput {
				fmt.Print("\r\033[K")
				fmt.Printf("[*] Not vulnerable: %s\n", result.TestURL)
			}
		}
	}
	bar.Finish()

	endTime := time.Now()

	// Output results
	if jsonOutput {
		outputJSON(allResults, startTime, endTime, len(urls)*len(payloads), vulnerableCount)
	} else {
		fmt.Printf("\n[*] Scan complete!\n")
		fmt.Printf("[*] Found %d vulnerable URLs\n", vulnerableCount)
		if vulnerableCount > 0 {
			fmt.Printf("[*] Results saved to: %s\n", outputFile)
		}
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

	// Setup Chrome options
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}

	// Add proxy if configured
	if proxyURL != "" {
		opts = append(opts, chromedp.ProxyServer(proxyURL))
	}

	// Create context with custom options
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	var finalURL string
	tasks := []chromedp.Action{}

	// Set cookies if provided
	if len(cookies) > 0 {
		tasks = append(tasks, network.SetCookies(cookies))
	}

	// Set custom headers if provided
	if len(headers) > 0 {
		tasks = append(tasks, network.Enable())
		tasks = append(tasks, network.SetExtraHTTPHeaders(network.Headers(headers)))
	}

	// Navigate and get final URL
	tasks = append(tasks,
		chromedp.Navigate(testURL),
		chromedp.Sleep(2*time.Second), // Wait for redirects
		chromedp.Location(&finalURL),
	)

	err := chromedp.Run(ctx, tasks...)

	if err != nil {
		if verbose {
			fmt.Printf("[!] Error testing %s: %v\n", testURL, err)
		}
		return Result{
			TestURL:      testURL,
			FinalURL:     "",
			IsVulnerable: false,
			Timestamp:    time.Now().Format(time.RFC3339),
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
		Timestamp:    time.Now().Format(time.RFC3339),
	}
}

func logVulnerability(result Result) {
	// Print to console
	fmt.Println("\n\n[*]*****Open Redirect Found*****[*]")
	fmt.Printf("[*] %s [*]\n", result.TestURL)
	fmt.Printf("[*] Redirects to: %s [*]\n", result.FinalURL)
	fmt.Printf("[*] Timestamp: %s [*]\n", result.Timestamp)
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

func outputJSON(results []Result, startTime, endTime time.Time, totalTests, vulnerableCount int) {
	output := JSONOutput{}
	output.ScanInfo.StartTime = startTime.Format(time.RFC3339)
	output.ScanInfo.EndTime = endTime.Format(time.RFC3339)
	output.ScanInfo.TotalTests = totalTests
	output.ScanInfo.VulnerableCount = vulnerableCount
	output.Results = results

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("[!] Error marshaling JSON: %v", err)
	}

	// Write to file
	if outputFile != "" && outputFile != "found.txt" {
		err = os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("[!] Error writing JSON to file: %v", err)
		}
		fmt.Printf("[*] JSON output saved to: %s\n", outputFile)
	} else {
		// Write to default JSON file
		jsonFile := "results.json"
		err = os.WriteFile(jsonFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("[!] Error writing JSON to file: %v", err)
		}
		fmt.Printf("[*] JSON output saved to: %s\n", jsonFile)
	}

	// Also print to stdout
	fmt.Println(string(jsonData))
}

func parseCookies(cookieStr string) []*network.CookieParam {
	var cookies []*network.CookieParam
	pairs := strings.Split(cookieStr, ";")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			cookie := &network.CookieParam{
				Name:  strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			}
			cookies = append(cookies, cookie)
		}
	}

	return cookies
}

func parseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	pairs := strings.Split(headerStr, ";")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	return headers
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
