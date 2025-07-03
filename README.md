

# Satsukashii

<img src="assets/banner.jpg" alt="banner with text satsukashi" width="50%">

A Go-based project for historical Bitcoin price analysis. Combines "satoshi" and "natsukashii" (nostalgia) to showcase what you might have earned if you had invested in Bitcoin earlier. Includes comparisons with other assets like Big Mac prices.

---

## Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Run it
 ```
 make run
 ```
---

## Tools

- **Go**: Primary programming language.
- **golangci-lint**: For code quality checks.

Install the linter:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## License

DAYW License. See [LICENSE](LICENSE) for details.

## Screenshot

<img src="asset/screenshot.png" alt="screenshot of page">
