# Open Redirect Finder

An automated tool to detect open redirect vulnerabilities in web applications.

[![Capture.png](https://s1.postimg.org/88l48isty7/Capture.png)](https://postimg.org/image/5dsg2qdn6j/)

## ⚠️ Important Notice

**This tool uses deprecated dependencies (CasperJS and PhantomJS) that are no longer actively maintained.** While the tool still works, consider using modern alternatives:

- [Puppeteer](https://github.com/puppeteer/puppeteer) - Headless Chrome/Chromium
- [Playwright](https://github.com/microsoft/playwright) - Cross-browser automation
- [Selenium](https://www.selenium.dev/) - Browser automation framework

## About

Based on the idea from [@ak1t4](https://github.com/ak1t4) and his script [open-redirect-scanner](https://github.com/ak1t4/open-redirect-scanner).

This tool takes two input files:
1. **URLs file** - Contains target URLs to test
2. **Payloads file** - Contains open redirect payloads (supplied)

The tool uses CasperJS headless browser to check if sites redirect to the supplied payloads. Successful redirects are logged to `found.txt`.







## Requirements

### System Requirements

- **Python 3.x** (Python 2 is no longer supported)
- **Node.js** (LTS version recommended)
- **PhantomJS** and **xvfb** (for headless operation)
- **CasperJS**

### Installation

#### Install Node.js (Ubuntu/Debian)

```bash
# Install Node.js LTS version
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs
```

#### Install Dependencies

```bash
# Install system dependencies
sudo apt-get install phantomjs xvfb -y

# Install CasperJS globally
sudo npm install -g casperjs

# Install Python requirements (optional, currently minimal)
pip3 install -r requirements.txt
```

## Usage

### Basic Usage

```bash
python3 redirect.py urls.txt payloads.txt
```

### Input Files

1. Create a `urls.txt` file with target URLs (one per line):
```
https://example.com/redirect?url=
https://target.com/forward?dest=
```

2. Use the provided `payloads.txt` or create your own with redirect payloads

### Output

- Successful redirects are logged to `found.txt`
- Console output shows real-time testing progress





## Docker Support

The tool supports Docker environments. When running in Docker, set the `DOCKER` environment variable:

```bash
export DOCKER=1
python3 redirect.py urls.txt payloads.txt
```

## Security & Legal Notice

**⚠️ IMPORTANT:** This tool is designed for authorized security testing only.

- Only test applications you own or have explicit permission to test
- Unauthorized testing may be illegal in your jurisdiction
- The authors are not responsible for misuse of this tool

## Features

- Automated testing of multiple URLs against multiple payloads
- Headless browser testing for accurate redirect detection
- Customizable payload lists
- Results logging for further analysis
- Docker environment support

## Troubleshooting

### Common Issues

1. **CasperJS not found**: Ensure CasperJS is installed globally with `npm install -g casperjs`
2. **PhantomJS errors**: PhantomJS may have compatibility issues on newer systems. Consider using the modern alternatives mentioned above.
3. **Permission errors**: Run with appropriate permissions or use `sudo` for system-level installations

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## To Do

- [ ] Migrate to modern browser automation (Puppeteer/Playwright)
- [ ] Add support for authenticated endpoints
- [ ] Implement multi-threading for faster scanning
- [ ] Auto-prepend http:// or https:// to URLs when missing
- [ ] Add JSON output format
- [ ] Implement rate limiting to avoid overwhelming targets

## License

See [LICENSE](LICENSE) file for details.

## Credits

- Original concept by [@ak1t4](https://github.com/ak1t4)
- Developed by [@Random_Robbie](https://github.com/random-robbie)
