# web-reader

A Cloudflare Workers service for converting web pages to clean, LLM-friendly markdown content. Inspired by [Jina Reader](https://r.jina.ai), this service provides an easy way to convert any webpage into clean, readable markdown format.

## Features

- **Simple URL-based Conversion**: Simply prepend your deployed worker URL to any webpage URL
- **Caching with D1**: Automatically caches converted content in Cloudflare D1 for faster responses
- **Cache Control**: Use `X-No-Cache` header to bypass cache when needed
- **Multiple Input Methods**: Support for both GET and POST requests

## Deployment

### Prerequisites

- [TinyGo](https://tinygo.org/getting-started/install/)
- [Wrangler](https://developers.cloudflare.com/workers/wrangler/install-and-update/)

### Steps

1. Clone the repository:
```bash
git clone https://github.com/AdjectiveAllison/web-reader.git
cd web-reader
```

2. Deploy the worker:
```bash
wrangler deploy
```

3. Initialize the D1 cache database:
```bash
wrangler d1 migrations apply web-reader-cache --remote
```

## Usage

### Simple GET Request
Simply prepend your worker URL to the target webpage URL:

```bash
curl "https://<your-worker>.workers.dev/https://example.com"
```

### POST Request
You can also send a POST request with the URL in the request body:

```bash
curl -X POST "https://<your-worker>.workers.dev/convert" \
-H "Content-Type: application/json" \
-d '{
    "url": "https://example.com"
}'
```

### Bypass Cache
To skip cached content and force a fresh fetch, use the `X-No-Cache` header:

```bash
curl -H "X-No-Cache: true" "https://<your-worker>.workers.dev/https://example.com"
```

## Response Format

The service returns the content in a clean, standardized format:

```
Title: Example Domain

URL Source: https://example.com

Markdown Content:
[converted markdown content here]
```

## License

MIT License - See [LICENSE](LICENSE) for details
