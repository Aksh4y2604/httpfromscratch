## httpfromscratch

Build-and-learn project: a minimal HTTP/1.1 server implemented from raw TCP sockets in Go. It parses requests and writes responses without `net/http` handlers, including manual status lines, headers, bodies, and a chunked transfer encoding example with trailers.

### What this implements (HTTP/1.1 from scratch)
- **TCP server**: accepts connections on a port and reads raw bytes.
- **Request line parsing**: `METHOD SP request-target SP HTTP/1.1 CRLF`.
- **Header parsing**: validates header field-name and aggregates duplicate keys per RFC semantics.
- **Body handling**: respects `Content-Length` for request bodies.
- **Response writing**:
  - Status line: `HTTP/1.1 <code> <reason> CRLF` (200/400/500 supported)
  - Headers: explicit `Content-Length`, `Content-Type`, `Connection`
  - Body: raw write for fixed-length responses
  - **Chunked Transfer-Encoding**: streaming body with hexadecimal chunk sizes, and a final zero-length chunk
  - **Trailers**: emits trailer headers after chunked body (e.g., `X-Content-SHA256`, `X-Content-Length`)

Key files:
- `internal/request`: incremental request parser for start-line, headers, and body
- `internal/headers`: header map with parser and helpers (get/set/replace/delete)
- `internal/response`: response writer with status line, headers, body, chunked encoding, and trailers
- `internal/server`: tiny TCP server that wires the parser and writer
- `cmd/httpserver`: example HTTP server with a few routes, including chunked streaming from httpbin
- `cmd/tcplistener`: prints parsed HTTP requests received on TCP
- `cmd/udpsender`: simple UDP line sender for experimentation

### Requirements
- Go 1.20+ (module declares 1.25). Install from `https://go.dev/dl/`.

### Getting started
Clone and build:

```bash
cd /workspace
go mod tidy
go build ./...
```

### Running the HTTP server
The server listens on port 42069.

```bash
go run ./cmd/httpserver
```

Test it:

```bash
# Default route
curl -i http://localhost:42069/

# Returns 400 with an HTML body
curl -i http://localhost:42069/yourproblem

# Returns 500 with an HTML body
curl -i http://localhost:42069/myproblem

# Chunked streaming proxy from httpbin, with trailers
curl -i http://localhost:42069/httpbin/5
```

What happens for `/httpbin/<n>`:
- The server streams `http://httpbin.org/stream/<n>`
- It switches to `Transfer-Encoding: chunked`
- Each chunk is written with size and CRLF, finishing with `0\r\n` and trailers
- Trailers include a SHA-256 digest and the content length of the streamed body

### Running the raw TCP listener (request inspector)
Prints the parsed HTTP request line, headers, and body. Useful for seeing how the parser behaves.

```bash
go run ./cmd/tcplistener
```

From another terminal, send a raw request:

```bash
printf "GET /hello HTTP/1.1\r\nHost: example\r\nContent-Length: 5\r\n\r\nhello" | nc localhost 42069
```

### Running the UDP sender (optional)
Simple REPL that sends lines over UDP to localhost:42069.

```bash
go run ./cmd/udpsender
```

### Project layout
```text
cmd/
  httpserver/     # sample HTTP server using the primitives
  tcplistener/    # raw TCP inspector for parsed HTTP requests
  udpsender/      # simple UDP line sender
internal/
  headers/        # header parsing and utilities
  request/        # incremental HTTP/1.1 request parser
  response/       # response writer, chunked encoding, trailers
  server/         # thin TCP accept loop and handler wiring
```

### Notes and limitations
- Only HTTP/1.1 is supported.
- No TLS; run behind a reverse proxy or use stunnel if needed.
- Limited status codes and content negotiation; this is for learning purposes.

### License
MIT

