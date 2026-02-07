# Capro 🧠⚡

**A Simple Disk-Backed Caching Proxy CLI written in Go**

Capro is a lightweight HTTP caching proxy built using Go and Cobra.
It forwards incoming HTTP requests to an origin server, caches responses on disk, and serves cached responses on subsequent requests — even after restarts.

This project focuses on **understanding proxy servers, caching fundamentals, and persistent storage**, not just making API calls.

---

## ✨ Features

- 🚀 HTTP proxy server
- 📦 Disk-backed cache (`cache.json`)
- 🔁 Cache persists across restarts
- 🧠 Cache HIT / MISS detection
- 🛠 CLI powered by Cobra
- 🔄 Configurable origin server
- 🧹 Cache clearing support

---

## 📦 Tech Stack

- **Language**: Go
- **CLI Framework**: Cobra
- **HTTP**: `net/http`
- **Storage**: JSON file (disk persistence)

---

## 📁 Project Structure

```
capro/
├── cmd/
│   └── root.go        # CLI + proxy server logic
├── cache.json         # Disk cache (auto-created)
├── main.go            # Entry point
└── README.md
```

---

## 🚀 Getting Started

### 1️⃣ Clone the repository

```bash
git clone https://github.com/your-username/capro.git
cd capro
```

---

### 2️⃣ Build the CLI

```bash
go build -o capro
```

---

### 3️⃣ Run the proxy server

```bash
./capro --origin https://api.github.com --port 3000
```

Output:

```
Starting Proxy Server on port 3000
Forwarding requests to: https://api.github.com
```

---

## 🔄 How It Works

1. Client sends a request to Capro

   ```
   http://localhost:3000/users
   ```

2. Capro checks the cache:

   - ✅ **HIT** → returns cached response
   - ❌ **MISS** → forwards request to origin

3. On MISS:

   - Fetches response from origin
   - Saves response to `cache.json`
   - Returns response to client

4. On restart:

   - Cache is loaded from disk
   - Previous responses are still available

---

## 🧪 Example Usage

```bash
curl http://localhost:3000/users
```

First request:

```
Cache-MISS : /users
```

Second request:

```
Cache-HIT : /users
```

Response headers include:

```
X-Cache: HIT
```

---

## 🧹 Clear Cache

```bash
./capro --clear-cache
```

This:

- Clears `cache.json`
- Resets in-memory cache

---

## 🧠 Cache Key Strategy

Cache entries are keyed by:

```
<origin> + <request URI>
```

Example:

```
https://api.github.com/users
```

This ensures:

- Different paths are cached separately
- Query parameters are respected

---

## 🧑‍💻 Author

**Harshith Rajesh**
Backend Developer | Go | Systems & Infrastructure Enthusiast

---

Project was inspired from Roadmap.sh backend projects
<https://roadmap.sh/projects/caching-server>
