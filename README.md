# Open Redirect Finder

A fast, modern open redirect vulnerability scanner written in Go using ChromeDP for headless browser automation.

[![Capture.png](https://s1.postimg.org/88l48isty7/Capture.png)](https://postimg.org/image/5dsg2qdn6j/)

## 🚀 Features

- **Modern Stack**: Built with Go and ChromeDP (headless Chrome)
- **Fast Scanning**: Concurrent workers for parallel testing
- **Docker Support**: Pre-configured Docker setup for easy deployment
- **Customisable**: Flexible payload lists and configurable options
- **Real Browser**: Uses actual Chrome for accurate redirect detection
- **Custom Test Domains**: Specify your own domains for redirect detection
- **Authentication Support**: Test authenticated endpoints with cookies and custom headers
- **JSON Output**: Structured output format for integration with other tools
- **Proxy Support**: Route traffic through HTTP/HTTPS/SOCKS5 proxies
- **Progress Bar**: Real-time visual progress indicator during scans
- **Simple**: Easy to use CLI interface

## 📋 About

Based on the original idea from [@ak1t4](https://github.com/ak1t4)'s [open-redirect-scanner](https://github.com/ak1t4/open-redirect-scanner).

This tool has been completely rewritten in Go with modern browser automation to replace the deprecated CasperJS/PhantomJS stack.

### How It Works

The tool:
1. Reads target URLs from a file
2. Appends payloads from a payload list to each URL
3. Uses headless Chrome to navigate and detect redirects
4. Identifies successful redirects to test domains (google.com, example.com)
5. Logs vulnerable URLs to an output file

## 🐳 Quick Start with Docker (Recommended)

### Prerequisites

- Docker
- Docker Compose (optional, for easier management)

### Using Docker Compose

1. **Clone the repository**
   ```bash
   git clone https://github.com/random-robbie/open-redirect.git
   cd open-redirect
   ```

2. **Create your URLs file**
   ```bash
   # Create the data directory and add your target URLs
   mkdir -p data
   nano data/urls.txt
   ```

3. **Run the scanner**
   ```bash
   docker-compose up --build
   ```

4. **View results**
   ```bash
   cat data/found.txt
   ```

### Using Docker Directly

```bash
# Build the image
docker build -t open-redirect .

# Run the scanner
docker run -v $(pwd)/data:/app/data open-redirect \
  -urls /app/data/urls.txt \
  -payloads /app/payloads.txt \
  -output /app/data/found.txt \
  -workers 10 \
  -verbose
```

## 💻 Local Installation

### Prerequisites

- Go 1.21 or later
- Chrome or Chromium browser

### Installation

```bash
# Clone the repository
git clone https://github.com/random-robbie/open-redirect.git
cd open-redirect

# Download dependencies
go mod download

# Build the binary
go build -o open-redirect main.go
```

### Usage

```bash
# Basic usage
./open-redirect -urls urls.txt -payloads payloads.txt

# With custom options
./open-redirect \
  -urls urls.txt \
  -payloads payloads.txt \
  -output results.txt \
  -workers 10 \
  -timeout 30 \
  -verbose
```

## 🎛️ Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-urls` | `urls.txt` | File containing target URLs to test |
| `-payloads` | `payloads.txt` | File containing redirect payloads |
| `-output` | `found.txt` | Output file for vulnerable URLs |
| `-workers` | `5` | Number of concurrent workers |
| `-timeout` | `30` | Timeout in seconds for each request |
| `-verbose` | `false` | Enable verbose output |
| `-json` | `false` | Output results in JSON format |
| `-domains` | *(default list)* | Comma-separated list of custom test domains |
| `-cookies` | *(none)* | Cookies in format 'name1=value1; name2=value2' |
| `-headers` | *(none)* | Custom headers in format 'Header1: Value1; Header2: Value2' |
| `-proxy` | *(none)* | Proxy URL (e.g., 'http://proxy.example.com:8080') |

## 📝 Input Files

### URLs File (`urls.txt`)

Create a file with target URLs (one per line):

```
https://example.com/redirect?url=
https://target.com/forward?dest=
https://site.com/goto?target=
```

### Payloads File (`payloads.txt`)

The repository includes a comprehensive payload list. You can also create your own:

```
//google.com
https://google.com
//example.com
@google.com
```

## 🎯 Detection Logic

The tool identifies successful open redirects by checking if the final URL starts with:
- `http://google.com` or `https://google.com`
- `http://example.com` or `https://example.com`

You can supply your own test domains with the `-domains` flag, or modify the `testDomains` variable in `main.go`.

## 📤 Output

### Console Output

The tool features a real-time progress bar that shows scan progress:

```
[*] ***************************************[*]
[*] Open Redirect Finder By @Random_Robbie [*]
[*]         Rewritten in Go + ChromeDP      [*]
[*] ***************************************[*]

[*] Test domains: [http://google.com https://google.com http://example.com https://example.com]
[*] Loaded 10 URLs and 504 payloads
[*] Using 5 concurrent workers
[*] Starting scan...

[*] Testing URLs... [=====>           ] 35% (1764/5040)

[*]*****Open Redirect Found*****[*]
[*] https://vulnerable.com/redirect?url=//google.com [*]
[*] Redirects to: https://google.com [*]
[*] Timestamp: 2024-01-15T10:32:15Z [*]

[*] Testing URLs... [==================>] 100% (5040/5040)

[*] Scan complete!
[*] Found 3 vulnerable URLs
[*] Results saved to: found.txt
```

The progress bar displays:
- Current completion percentage
- Number of tests completed vs total
- Visual progress indicator with color coding
- Automatically clears when printing vulnerability findings

### Output File (`found.txt`)

Vulnerable URLs are saved one per line:
```
https://vulnerable.com/redirect?url=//google.com
https://target.com/forward?dest=https://example.com
```

## ⚙️ Configuration

### Environment Variables (Docker)

The provided `docker-compose.yml` defines `WORKERS` and `TIMEOUT` environment variables:

```yaml
environment:
  - WORKERS=10
  - TIMEOUT=60
```

Note that these are not currently read by the binary — pass the `-workers` and `-timeout` flags in the compose command instead if you need to change them.

### Custom Payloads

The included `payloads.txt` contains 500+ bypass techniques. Add your own:

```bash
echo "//your-domain.com" >> payloads.txt
```

## 🔧 Advanced Usage

### Scanning Large Target Lists

```bash
# Increase workers for faster scanning
./open-redirect -urls large-list.txt -workers 20

# Increase timeout for slow targets
./open-redirect -urls urls.txt -timeout 60
```

### Custom Output Location

```bash
./open-redirect -urls urls.txt -output /path/to/results.txt
```

### Verbose Mode

```bash
# See all requests, including non-vulnerable ones
./open-redirect -urls urls.txt -verbose
```

### Custom Test Domains

By default, the tool checks for redirects to `google.com` and `example.com`. You can specify your own test domains:

```bash
# Use custom domains for detection
./open-redirect \
  -urls urls.txt \
  -domains "https://evil.com,http://attacker.com,https://test.com"
```

This is useful when:
- Testing with your own controlled domains
- Verifying specific redirect targets
- Using domains you control for bug bounty testing

### JSON Output Format

Generate structured JSON output for integration with other tools:

```bash
# Output results in JSON format
./open-redirect -urls urls.txt -json -output results.json
```

JSON output includes:
- Scan metadata (start time, end time, total tests)
- Complete results with timestamps
- Structured data for easy parsing

Example JSON output:
```json
{
  "scan_info": {
    "start_time": "2024-01-15T10:30:00Z",
    "end_time": "2024-01-15T10:35:00Z",
    "total_tests": 5040,
    "vulnerable_count": 3
  },
  "results": [
    {
      "test_url": "https://example.com/redirect?url=//google.com",
      "final_url": "https://google.com",
      "vulnerable": true,
      "timestamp": "2024-01-15T10:32:15Z"
    }
  ]
}
```

### Authentication Support

#### Using Cookies

Test authenticated endpoints by providing session cookies:

```bash
# Single cookie
./open-redirect \
  -urls urls.txt \
  -cookies "session=abc123def456"

# Multiple cookies
./open-redirect \
  -urls urls.txt \
  -cookies "session=abc123; csrf_token=xyz789; user_id=12345"
```

#### Using Custom Headers

Add custom HTTP headers for authentication or other purposes:

```bash
# Single header
./open-redirect \
  -urls urls.txt \
  -headers "Authorization: Bearer token123"

# Multiple headers
./open-redirect \
  -urls urls.txt \
  -headers "Authorization: Bearer token123; X-API-Key: key456; X-Custom: value"
```

#### Combined Authentication

Use both cookies and headers together:

```bash
./open-redirect \
  -urls urls.txt \
  -cookies "session=abc123; user=admin" \
  -headers "Authorization: Bearer token123; X-CSRF-Token: xyz789"
```

### Proxy Support

Route traffic through a proxy for:
- Corporate network requirements
- Additional anonymity
- Traffic inspection/debugging

```bash
# HTTP proxy
./open-redirect \
  -urls urls.txt \
  -proxy "http://proxy.company.com:8080"

# HTTPS proxy
./open-redirect \
  -urls urls.txt \
  -proxy "https://secure-proxy.com:443"

# SOCKS5 proxy
./open-redirect \
  -urls urls.txt \
  -proxy "socks5://127.0.0.1:1080"
```

### Combined Advanced Example

Using all features together:

```bash
./open-redirect \
  -urls urls.txt \
  -payloads custom-payloads.txt \
  -output results.json \
  -workers 15 \
  -timeout 45 \
  -json \
  -verbose \
  -domains "https://attacker.com,http://evil.com" \
  -cookies "session=abc123; user=admin" \
  -headers "Authorization: Bearer token123; X-API-Key: key456" \
  -proxy "http://proxy.company.com:8080"
```

## 🛡️ Security & Legal Notice

**⚠️ IMPORTANT:** This tool is designed for authorised security testing only.

- ✅ Only test applications you own or have explicit written permission to test
- ✅ Use for bug bounty programmes with proper authorisation
- ✅ Use for penetration testing engagements
- ❌ Unauthorised testing may be illegal in your jurisdiction
- ❌ The authors are not responsible for misuse of this tool

Always ensure you have proper authorisation before testing any application.

## 🐛 Troubleshooting

### Chrome/Chromium Not Found

If running locally and Chrome is not in the system PATH:

```bash
# Linux
export CHROME_BIN=/usr/bin/chromium

# macOS
export CHROME_BIN="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
```

### Docker Permission Issues

```bash
# Fix permissions on data directory
chmod -R 777 data/
```

### Connection Timeouts

Increase the timeout value:
```bash
./open-redirect -urls urls.txt -timeout 60
```

## 🔄 Migration from Python Version

If you're upgrading from the old Python/CasperJS version:

1. Your existing `payloads.txt` file will work as-is
2. Create a new `urls.txt` with your target URLs
3. Use Docker for the easiest setup (no need to install dependencies)
4. The output format remains compatible

## 📊 Performance

- **Python + CasperJS**: ~5-10 URLs/minute (single-threaded)
- **Go + ChromeDP**: ~50-100+ URLs/minute (with 10 workers)

Actual performance depends on network conditions and target response times.

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
# Clone the repo
git clone https://github.com/random-robbie/open-redirect.git
cd open-redirect

# Install dependencies
go mod download

# Run tests (if available)
go test ./...

# Build
go build -o open-redirect main.go
```

## 📝 To Do

- [x] Add support for custom test domains via CLI flag ✅
- [x] Implement authentication support (cookies, headers) ✅
- [x] Add JSON output format ✅
- [x] Support for proxy configuration ✅
- [x] Add progress bar for long scans ✅
- [ ] Create comprehensive test suite
- [ ] Add rate limiting options
- [ ] Add CI/CD pipeline
- [ ] Performance benchmarking
- [ ] Support for loading cookies/headers from file
- [ ] Implement retry logic for failed requests
- [ ] Add option to disable progress bar for CI/CD

## 📜 License

See [LICENSE](LICENSE) file for details.

## 🙏 Credits

- **Original Concept**: [@ak1t4](https://github.com/ak1t4)
- **Original Python Version**: [@Random_Robbie](https://github.com/random-robbie)
- **Go Rewrite**: [@Random_Robbie](https://github.com/random-robbie)
- **Browser Automation**: [ChromeDP](https://github.com/chromedp/chromedp)

## 📚 Resources

- [OWASP Open Redirect](https://cheatsheetseries.owasp.org/cheatsheets/Unvalidated_Redirects_and_Forwards_Cheat_Sheet.html)
- [ChromeDP Documentation](https://pkg.go.dev/github.com/chromedp/chromedp)
- [Open Redirect Testing Guide](https://portswigger.net/kb/issues/00500100_open-redirection-reflected)

## Disclaimer / Authorised Use

This tool is provided for authorised security testing and responsible disclosure only. Do not use it against systems you do not own or have explicit permission to test.

## Author

Robbie Wiggins — [What-Security](https://www.what-security.co.uk) · [GitHub](https://github.com/random-robbie)

---

**Star ⭐ this repository if you find it useful!**
